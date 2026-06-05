package repository

import (
	"context"

	"task-manager-api/internal/models"

	"gorm.io/gorm"
)

// TaskRepository определяет контракт работы с задачами.
type TaskRepository interface {
	Create(ctx context.Context, task *models.Task) error
	GetByID(ctx context.Context, id uint) (*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	Delete(ctx context.Context, task *models.Task) error
	List(ctx context.Context, userID uint, status string, offset int, limit int) ([]models.Task, int64, error)
}

// GormTaskRepository реализует TaskRepository через GORM.
type GormTaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository создает новый репозиторий задач.
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &GormTaskRepository{
		db: db,
	}
}

// Create создает задачу.
func (r *GormTaskRepository) Create(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Create(task).Error
}

// GetByID получает задачу по идентификатору.
func (r *GormTaskRepository) GetByID(ctx context.Context, id uint) (*models.Task, error) {
	var task models.Task

	err := r.db.WithContext(ctx).
		First(&task, id).Error

	if err != nil {
		return nil, err
	}

	return &task, nil
}

// Update обновляет задачу.
func (r *GormTaskRepository) Update(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Save(task).Error
}

// Delete удаляет задачу.
func (r *GormTaskRepository) Delete(ctx context.Context, task *models.Task) error {
	return r.db.WithContext(ctx).Delete(task).Error
}

// List возвращает список задач пользователя.
func (r *GormTaskRepository) List(
	ctx context.Context,
	userID uint,
	status string,
	offset int,
	limit int,
) ([]models.Task, int64, error) {
	var tasks []models.Task
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Task{}).
		Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&tasks).Error

	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}
