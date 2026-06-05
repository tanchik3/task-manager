# ---------- STAGE 1: BUILD ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# зависимости для сборки
RUN apk add --no-cache git

# копируем модули
COPY go.mod go.sum ./
RUN go mod download

# копируем весь проект
COPY . .

# собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd/api

# ---------- STAGE 2: RUNTIME ----------
FROM alpine:latest

WORKDIR /app

# timezone + сертификаты (важно для JWT/HTTPS)
RUN apk add --no-cache ca-certificates

# копируем бинарник из build stage
COPY --from=builder /app/app .

# порт приложения
EXPOSE 8080

# запуск
CMD ["./app"]