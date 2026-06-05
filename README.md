# Task Manager API (Go + Gin + GORM + PostgreSQL/SQLite)

REST API сервис для управления задачами с JWT-аутентификацией.

Проект поддерживает:
- PostgreSQL
- JWT авторизацию
- CRUD задач
- пагинацию и фильтрацию
- Docker
- GitHub Actions CI/CD


# Технологии

- Go 1.25+
- Gin
- GORM
- SQLite3 / PostgreSQL
- JWT (github.com/golang-jwt/jwt/v5)
- bcrypt
- godotenv
- slog
- Docker
- GitHub Actions

# Структура проекта

task-manager-api/
├── cmd/api/main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   ├── services/
│   └── routes/
├── pkg/logger/
├── migrations/
├── tests/
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── Makefile
└── README.md

# Переменные окружения

Создай `.env` файл:

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=taskmanager

JWT_SECRET=секрет-нужно-изменить
JWT_EXPIRE_HOURS=24

SERVER_PORT=8080

# Запуск проекта (SQLite)

go mod tidy
go run cmd/api/main.go

Сервер будет доступен:
http://localhost:8080

# Тесты

go test ./...

# Форматирование

go fmt ./...

# Docker

# Сборка
docker build -t task-manager-api .

# Запуск
docker run -p 8080:8080 task-manager-api

# GitHub Actions

При push в main:

- тесты
- build
- Docker build
- push в Docker Hub