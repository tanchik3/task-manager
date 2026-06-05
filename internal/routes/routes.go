package routes

import (
	"task-manager-api/internal/handlers"
	"task-manager-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Register регистрирует маршруты приложения.
func Register(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	taskHandler *handlers.TaskHandler,
	authMiddleware *middleware.AuthMiddleware,
) {
	api := router.Group("/api")

	api.POST("/register", authHandler.Register)
	api.POST("/login", authHandler.Login)

	tasks := api.Group("/tasks")
	tasks.Use(authMiddleware.Handle())

	tasks.POST("", taskHandler.CreateTask)
	tasks.GET("", taskHandler.ListTasks)
	tasks.GET("/:id", taskHandler.GetTask)
	tasks.PUT("/:id", taskHandler.UpdateTask)
	tasks.DELETE("/:id", taskHandler.DeleteTask)
	tasks.PATCH("/:id/status", taskHandler.UpdateTaskStatus)
}
