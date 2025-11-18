.PHONY: help build run clean test install dev

# Variables
BINARY_NAME=twitter-backend
MAIN_PATH=main.go

help: ## Hiển thị help message
	@echo "Available commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Cài đặt dependencies
	@echo "📦 Đang cài đặt dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies đã được cài đặt"

build: ## Build binary
	@echo "🔨 Đang build..."
	go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build thành công: ./$(BINARY_NAME)"

run: ## Chạy application
	@echo "🚀 Đang chạy application..."
	go run $(MAIN_PATH)

dev: ## Chạy với hot reload (yêu cầu air)
	@echo "🔥 Đang chạy với hot reload..."
	@which air > /dev/null || (echo "❌ Chưa cài air. Chạy: go install github.com/cosmtrek/air@latest" && exit 1)
	air

test: ## Chạy tests
	@echo "🧪 Đang chạy tests..."
	go test -v ./...

test-coverage: ## Chạy tests với coverage
	@echo "📊 Đang chạy tests với coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: coverage.html"

clean: ## Xóa build artifacts
	@echo "🧹 Đang dọn dẹp..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	@echo "✅ Đã dọn dẹp xong"

lint: ## Chạy linter
	@echo "🔍 Đang chạy linter..."
	@which golangci-lint > /dev/null || (echo "❌ Chưa cài golangci-lint. Xem: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

fmt: ## Format code
	@echo "✨ Đang format code..."
	go fmt ./...
	@echo "✅ Code đã được format"

vet: ## Chạy go vet
	@echo "🔍 Đang chạy go vet..."
	go vet ./...
	@echo "✅ Vet completed"

docker-build: ## Build Docker image
	@echo "🐳 Đang build Docker image..."
	docker build -t $(BINARY_NAME):latest .
	@echo "✅ Docker image đã được build"

docker-run: ## Chạy Docker container
	@echo "🐳 Đang chạy Docker container..."
	docker run -p 8080:8080 --env-file .env $(BINARY_NAME):latest

all: clean install build ## Chạy clean, install và build

.DEFAULT_GOAL := help

