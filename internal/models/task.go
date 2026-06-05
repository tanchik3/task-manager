package models

import "time"

// Константы статусов задач.
const (
	TaskStatusPending    = "pending"
	TaskStatusInProgress = "in_progress"
	TaskStatusDone       = "done"
)

// Task представляет задачу пользователя.
type Task struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	Status      string     `json:"status" gorm:"default:pending"`
	DueDate     *time.Time `json:"due_date"`
	UserID      uint       `json:"user_id" gorm:"index;not null"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
