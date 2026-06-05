package database

import (
	"fmt"

	"task-manager-api/internal/config"
	"task-manager-api/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Connect создает подключение к базе данных.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	var (
		db  *gorm.DB
		err error
	)

	switch cfg.DBType {
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost,
			cfg.DBPort,
			cfg.DBUser,
			cfg.DBPassword,
			cfg.DBName,
		)

		db, err = gorm.Open(
			postgres.Open(dsn),
			&gorm.Config{},
		)

	default:
		db, err = gorm.Open(
			sqlite.Open(cfg.DBPath),
			&gorm.Config{},
		)
	}

	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Task{},
	); err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
