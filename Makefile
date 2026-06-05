APP_NAME=task-manager-api

.PHONY: run
run:
	go run cmd/api/main.go

.PHONY: build
build:
	go build -o bin/$(APP_NAME) cmd/api/main.go

.PHONY: test
test:
	go test ./... -v

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: clean
clean:
	rm -rf bin
	rm -f coverage.out

.PHONY: coverage
coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out