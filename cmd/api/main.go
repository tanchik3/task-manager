package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"task-manager-api/internal/config"
	"task-manager-api/internal/database"
	"task-manager-api/internal/handlers"
	"task-manager-api/internal/middleware"
	"task-manager-api/internal/repository"
	"task-manager-api/internal/routes"
	"task-manager-api/internal/services"
	"task-manager-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.New()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	tokenService := services.NewTokenService(cfg)
	authService := services.NewAuthService(
		userRepo,
		tokenService,
	)

	taskService := services.NewTaskService(taskRepo)

	authHandler := handlers.NewAuthHandler(authService)
	taskHandler := handlers.NewTaskHandler(taskService)

	authMiddleware := middleware.NewAuthMiddleware(
		tokenService,
	)

	router := gin.New()

	router.Use(logger.LoggingMiddleware(appLogger))
	router.Use(logger.RecoveryMiddleware(appLogger))

	routes.Register(
		router,
		authHandler,
		taskHandler,
		authMiddleware,
	)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		appLogger.Info(
			"server started",
			slog.String("port", cfg.ServerPort),
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	waitForShutdown(server, appLogger)
}

// waitForShutdown выполняет корректное завершение сервера.
func waitForShutdown(
	server *http.Server,
	log *slog.Logger,
) {
	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-stop

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	log.Info("shutdown started")

	if err := server.Shutdown(ctx); err != nil {
		log.Error(
			"shutdown error",
			slog.String("error", err.Error()),
		)
		return
	}

	fmt.Println("server stopped")
}
