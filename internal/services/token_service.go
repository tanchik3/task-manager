package services

import (
	"fmt"
	"time"

	"task-manager-api/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// TokenService предоставляет функциональность работы с JWT токенами.
type TokenService struct {
	config *config.Config
}

// Claims представляет набор JWT claims.
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// NewTokenService создает новый сервис токенов.
func NewTokenService(cfg *config.Config) *TokenService {
	return &TokenService{
		config: cfg,
	}
}

// GenerateToken создает JWT токен для пользователя.
func (s *TokenService) GenerateToken(userID uint) (string, error) {
	expiresAt := time.Now().Add(time.Duration(s.config.JWTExpireHours) * time.Hour)

	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.config.JWTSecret))
}

// ValidateToken валидирует JWT токен.
func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("неподдерживаемый метод подписи")
			}

			return []byte(s.config.JWTSecret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("невалидный токен")
	}

	return claims, nil
}
