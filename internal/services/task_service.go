package services

import (
	"context"
	"errors"

	"task-manager-api/internal/models"
	"task-manager-api/internal/repository"

	"gorm.io/gorm"
)

var ErrForbidden = errors.New("forbidden")

// TaskService предоставляет бизнес-логику работы с задачами.
type TaskService struct {
	taskRepo repository.TaskRepository
}

// NewTaskService создает сервис задач.
func NewTaskService(taskRepo repository.TaskRepository) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
	}
}

// CreateTask создает задачу.
func (s *TaskService) CreateTask(
	ctx context.Context,
	userID uint,
	req models.CreateTaskRequest,
) (*models.Task, error) {
	status := req.Status

	if status == "" {
		status = models.TaskStatusPending
	}

	if !isValidStatus(status) {
		return nil, errors.New("некорректный статус")
	}

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Status:      status,
		DueDate:     req.DueDate,
		UserID:      userID,
	}

	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// GetTask получает задачу.
func (s *TaskService) GetTask(
	ctx context.Context,
	taskID uint,
	userID uint,
) (*models.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.UserID != userID {
		return nil, ErrForbidden
	}

	return task, nil
}

// ListTasks возвращает список задач пользователя.
func (s *TaskService) ListTasks(
	ctx context.Context,
	userID uint,
	page int,
	limit int,
	status string,
) (*models.TaskListResponse, error) {
	if status != "" && !isValidStatus(status) {
		return nil, errors.New("некорректный статус")
	}

	offset := (page - 1) * limit

	tasks, total, err := s.taskRepo.List(
		ctx,
		userID,
		status,
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}

	return &models.TaskListResponse{
		Tasks: tasks,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

// UpdateTask обновляет задачу.
func (s *TaskService) UpdateTask(
	ctx context.Context,
	taskID uint,
	userID uint,
	req models.UpdateTaskRequest,
) (*models.Task, error) {
	task, err := s.GetTask(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		task.Title = req.Title
	}

	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if req.Status != "" {
		if !isValidStatus(req.Status) {
			return nil, errors.New("некорректный статус")
		}

		task.Status = req.Status
	}

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// UpdateStatus обновляет статус задачи.
func (s *TaskService) UpdateStatus(
	ctx context.Context,
	taskID uint,
	userID uint,
	status string,
) (*models.Task, error) {
	if !isValidStatus(status) {
		return nil, errors.New("некорректный статус")
	}

	task, err := s.GetTask(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}

	task.Status = status

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// DeleteTask удаляет задачу.
func (s *TaskService) DeleteTask(
	ctx context.Context,
	taskID uint,
	userID uint,
) error {
	task, err := s.GetTask(ctx, taskID, userID)
	if err != nil {
		return err
	}

	return s.taskRepo.Delete(ctx, task)
}

// IsForbiddenError определяет ошибку доступа.
func (s *TaskService) IsForbiddenError(err error) bool {

	return errors.Is(err, ErrForbidden)
}

// IsNotFoundError определяет ошибку отсутствия записи.
func (s *TaskService) IsNotFoundError(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

// isValidStatus проверяет корректность статуса.
func isValidStatus(status string) bool {
	switch status {
	case models.TaskStatusPending,
		models.TaskStatusInProgress,
		models.TaskStatusDone:
		return true
	default:
		return false
	}
}
