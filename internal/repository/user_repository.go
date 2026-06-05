package repository

import (
	"context"

	"task-manager-api/internal/models"

	"gorm.io/gorm"
)

// UserRepository определяет контракт работы с пользователями.
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id uint) (*models.User, error)
}

// GormUserRepository реализует UserRepository через GORM.
type GormUserRepository struct {
	db *gorm.DB
}

// NewUserRepository создает новый репозиторий пользователей.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &GormUserRepository{
		db: db,
	}
}

// Create создает пользователя.
func (r *GormUserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByEmail получает пользователя по email.
func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByID получает пользователя по идентификатору.
func (r *GormUserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	err := r.db.WithContext(ctx).
		First(&user, id).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
