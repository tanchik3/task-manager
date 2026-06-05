package models

import "time"

// RegisterRequest представляет запрос регистрации.
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

// LoginRequest представляет запрос авторизации.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse представляет ответ авторизации.
type LoginResponse struct {
	Token string   `json:"token"`
	User  UserInfo `json:"user"`
}

// UserInfo представляет публичную информацию о пользователе.
type UserInfo struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

// CreateTaskRequest представляет запрос создания задачи.
type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
}

// UpdateTaskRequest представляет запрос обновления задачи.
type UpdateTaskRequest struct {
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
}

// UpdateTaskStatusRequest представляет запрос обновления статуса задачи.
type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// TaskListResponse представляет список задач.
type TaskListResponse struct {
	Tasks []Task `json:"tasks"`
	Total int64  `json:"total"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
}
