package services

import (
	"context"
	"errors"

	"task-manager-api/internal/models"
	"task-manager-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService предоставляет бизнес-логику авторизации.
type AuthService struct {
	userRepo     repository.UserRepository
	tokenService *TokenService
}

// NewAuthService создает сервис авторизации.
func NewAuthService(
	userRepo repository.UserRepository,
	tokenService *TokenService,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Register регистрирует нового пользователя.
func (s *AuthService) Register(
	ctx context.Context,
	req models.RegisterRequest,
) (*models.User, error) {
	_, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("пользователь уже существует")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login выполняет вход пользователя.
func (s *AuthService) Login(
	ctx context.Context,
	req models.LoginRequest,
) (*models.LoginResponse, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return nil, errors.New("неверный email или пароль")
	}

	token, err := s.tokenService.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: token,
		User: models.UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
	}, nil
}
