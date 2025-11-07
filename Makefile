.PHONY: help build run test clean docker-build docker-run docker-stop lint

# Variables
APP_NAME=weather-microservice
BINARY_NAME=weather-server
DOCKER_IMAGE=weather-microservice:latest
PORT=8080

help: ## Display this help message
	@echo "Weather Microservice - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "🔨 Building $(APP_NAME)..."
	@go build -o bin/$(BINARY_NAME) ./cmd/server/main.go
	@echo "✅ Build complete: bin/$(BINARY_NAME)"

run: ## Run the application locally
	@echo "🚀 Starting $(APP_NAME)..."
	@go run ./cmd/server/main.go

dev: ## Run with hot reload (requires air)
	@echo "🔥 Starting development server..."
	@air

test: ## Run tests
	@echo "🧪 Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Tests complete. Coverage report: coverage.html"

test-unit: ## Run unit tests only
	@echo "🧪 Running unit tests..."
	@go test -v ./internal/... ./pkg/...

test-integration: ## Run integration tests
	@echo "🧪 Running integration tests..."
	@go test -v ./tests/integration/...

bench: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	@go test -bench=. -benchmem ./...

lint: ## Run linter
	@echo "🔍 Running linter..."
	@golangci-lint run ./...

fmt: ## Format code
	@echo "📝 Formatting code..."
	@go fmt ./...
	@goimports -w .

clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "✅ Clean complete"

deps: ## Download dependencies
	@echo "📦 Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies downloaded"

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "✅ Docker image built: $(DOCKER_IMAGE)"

docker-run: ## Run Docker container
	@echo "🐳 Running Docker container..."
	@docker run -d \
		--name $(APP_NAME) \
		-p $(PORT):8080 \
		-v $(PWD)/.apiConfig:/app/.apiConfig:ro \
		$(DOCKER_IMAGE)
	@echo "✅ Container started on port $(PORT)"

docker-stop: ## Stop and remove Docker container
	@echo "🛑 Stopping Docker container..."
	@docker stop $(APP_NAME) || true
	@docker rm $(APP_NAME) || true
	@echo "✅ Container stopped"

docker-compose-up: ## Start with docker-compose
	@echo "🐳 Starting services with docker-compose..."
	@docker-compose up -d
	@echo "✅ Services started"

docker-compose-down: ## Stop docker-compose services
	@echo "🛑 Stopping docker-compose services..."
	@docker-compose down
	@echo "✅ Services stopped"

docker-logs: ## View Docker logs
	@docker logs -f $(APP_NAME)

install-tools: ## Install development tools
	@echo "🔧 Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/air-verse/air@latest
	@echo "✅ Tools installed"

setup: ## Setup development environment
	@echo "🔧 Setting up development environment..."
	@make deps
	@make install-tools
	@if [ ! -f .apiConfig ]; then \
		echo '{"OpenWeatherMapApiKey":"5aa85edefd94c29ea343cb21563aa912","CacheExpiryMinutes":10,"RateLimitPerMinute":100}' > .apiConfig; \
		echo "⚠️  Created .apiConfig - Please add your OpenWeatherMap API key"; \
	fi
	@echo "✅ Setup complete"

deploy: ## Build and deploy (placeholder)
	@echo "🚀 Deploying $(APP_NAME)..."
	@make build
	@echo "✅ Ready for deployment"

all: clean deps build test ## Clean, download deps, build, and test

.DEFAULT_GOAL := help