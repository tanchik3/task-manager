package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config содержит конфигурацию приложения.
type Config struct {
	DBType         string
	DBPath         string
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	JWTExpireHours int
	ServerPort     string
}

// Load загружает конфигурацию из переменных окружения.
func Load() (*Config, error) {
	_ = godotenv.Load()

	expireHours, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
	if err != nil {
		expireHours = 24
	}

	cfg := &Config{
		DBType:         getEnv("DB_TYPE", "sqlite3"),
		DBPath:         getEnv("DB_PATH", "./taskmanager.db"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "password"),
		DBName:         getEnv("DB_NAME", "taskmanager"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpireHours: expireHours,
		ServerPort:     getEnv("SERVER_PORT", "8080"),
	}

	return cfg, nil
}

// getEnv возвращает значение переменной окружения или значение по умолчанию.
func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
