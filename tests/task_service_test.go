package tests

import (
	"context"
	"testing"

	"task-manager-api/internal/models"
	"task-manager-api/internal/repository"
	"task-manager-api/internal/services"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func createTestTaskService(t *testing.T) *services.TaskService {
	t.Helper()

	dsn := "host=localhost user=postgres password=PelmeniTop47 dbname=taskmanager_test port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("db open error: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Task{},
	); err != nil {
		t.Fatalf("migration error: %v", err)
	}

	taskRepo := repository.NewTaskRepository(db)

	return services.NewTaskService(taskRepo)
}

func TestTaskService(t *testing.T) {
	service := createTestTaskService(t)

	tests := []struct {
		name      string
		request   models.CreateTaskRequest
		wantError bool
	}{
		{
			name: "create task success",
			request: models.CreateTaskRequest{
				Title:  "Task 1",
				Status: models.TaskStatusPending,
			},
			wantError: false,
		},
		{
			name: "invalid status",
			request: models.CreateTaskRequest{
				Title:  "Task 2",
				Status: "invalid",
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CreateTask(
				context.Background(),
				1,
				tt.request,
			)

			if tt.wantError && err == nil {
				t.Fatal("expected error")
			}

			if !tt.wantError && err != nil {
				t.Fatalf(
					"unexpected error: %v",
					err,
				)
			}
		})
	}
}
