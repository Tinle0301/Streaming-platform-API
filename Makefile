.PHONY: help build test run run-api run-ws docker-up docker-down migrate lint clean

# Variables
BINARY_API=bin/api-server
BINARY_WS=bin/ws-server
GO=go
GOFLAGS=-v
DOCKER_COMPOSE=docker-compose

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "StreamHub API Platform - Available Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build all binaries
	@echo "Building API server..."
	@mkdir -p bin
	$(GO) build $(GOFLAGS) -o $(BINARY_API) ./cmd/api-server
	@echo "Building WebSocket server..."
	$(GO) build $(GOFLAGS) -o $(BINARY_WS) ./cmd/ws-server
	@echo "✅ Build complete"

build-api: ## Build API server only
	@echo "Building API server..."
	@mkdir -p bin
	$(GO) build $(GOFLAGS) -o $(BINARY_API) ./cmd/api-server

build-ws: ## Build WebSocket server only
	@echo "Building WebSocket server..."
	@mkdir -p bin
	$(GO) build $(GOFLAGS) -o $(BINARY_WS) ./cmd/ws-server

run-api: ## Run API server
	@echo "Starting API server..."
	$(GO) run ./cmd/api-server/main.go

run-ws: ## Run WebSocket server
	@echo "Starting WebSocket server..."
	$(GO) run ./cmd/ws-server/main.go

run: ## Run both servers concurrently
	@echo "Starting all servers..."
	@make -j2 run-api run-ws

test: ## Run all tests
	@echo "Running tests..."
	$(GO) test -v -race -coverprofile=coverage.out ./...
	@echo "✅ Tests complete"

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	$(GO) test -v -tags=integration ./tests/integration/...

test-load: ## Run load tests
	@echo "Running load tests..."
	@if command -v k6 >/dev/null 2>&1; then \
		k6 run tests/load/api-load-test.js; \
	else \
		echo "❌ k6 not installed. Install from https://k6.io"; \
	fi

lint: ## Run linters
	@echo "Running linters..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "⚠️  golangci-lint not installed. Install from https://golangci-lint.run"; \
		$(GO) fmt ./...; \
		$(GO) vet ./...; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	$(GO) fmt ./...
	@echo "✅ Code formatted"

generate: ## Generate GraphQL code
	@echo "Generating GraphQL code..."
	@if command -v gqlgen >/dev/null 2>&1; then \
		gqlgen generate; \
	else \
		echo "⚠️  gqlgen not installed. Run: go install github.com/99designs/gqlgen@latest"; \
	fi

docker-build: ## Build Docker images
	@echo "Building Docker images..."
	docker build -f deployments/docker/Dockerfile.api -t streamhub-api:latest .
	docker build -f deployments/docker/Dockerfile.ws -t streamhub-ws:latest .
	@echo "✅ Docker images built"

docker-up: ## Start all services with Docker Compose
	@echo "Starting services with Docker Compose..."
	$(DOCKER_COMPOSE) -f deployments/docker/docker-compose.yml up -d
	@echo "✅ Services started"
	@echo "   PostgreSQL: localhost:5432"
	@echo "   Redis: localhost:6379"
	@echo "   RabbitMQ: localhost:5672"
	@echo "   RabbitMQ Management: http://localhost:15672"

docker-down: ## Stop all Docker Compose services
	@echo "Stopping services..."
	$(DOCKER_COMPOSE) -f deployments/docker/docker-compose.yml down
	@echo "✅ Services stopped"

docker-logs: ## Show Docker Compose logs
	$(DOCKER_COMPOSE) -f deployments/docker/docker-compose.yml logs -f

migrate: ## Run database migrations
	@echo "Running database migrations..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path ./migrations -database "$(DATABASE_URL)" up; \
	else \
		echo "⚠️  migrate tool not installed. Install from https://github.com/golang-migrate/migrate"; \
	fi

migrate-down: ## Rollback database migrations
	@echo "Rolling back migrations..."
	@if command -v migrate >/dev/null 2>&1; then \
		migrate -path ./migrations -database "$(DATABASE_URL)" down; \
	else \
		echo "⚠️  migrate tool not installed."; \
	fi

migrate-create: ## Create a new migration (usage: make migrate-create NAME=create_users_table)
	@if [ -z "$(NAME)" ]; then \
		echo "❌ Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	@echo "Creating migration: $(NAME)"
	@if command -v migrate >/dev/null 2>&1; then \
		migrate create -ext sql -dir ./migrations -seq $(NAME); \
	else \
		echo "⚠️  migrate tool not installed."; \
	fi

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@$(GO) clean
	@echo "✅ Clean complete"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod verify
	@echo "✅ Dependencies downloaded"

deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy
	@echo "✅ Dependencies updated"

install-tools: ## Install development tools
	@echo "Installing development tools..."
	$(GO) install github.com/99designs/gqlgen@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO) install github.com/golang/mock/mockgen@latest
	@echo "✅ Tools installed"

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	$(GO) test -bench=. -benchmem ./...

dev: docker-up ## Start development environment
	@echo "🚀 Development environment ready!"
	@echo ""
	@echo "Services running:"
	@echo "  - PostgreSQL: localhost:5432"
	@echo "  - Redis: localhost:6379"
	@echo "  - RabbitMQ: localhost:5672"
	@echo ""
	@echo "Run 'make run' to start the servers"

# Kubernetes deployment commands
k8s-deploy: ## Deploy to Kubernetes
	@echo "Deploying to Kubernetes..."
	kubectl apply -f deployments/k8s/

k8s-delete: ## Delete Kubernetes deployment
	@echo "Deleting Kubernetes deployment..."
	kubectl delete -f deployments/k8s/

# AWS deployment commands
terraform-init: ## Initialize Terraform
	@echo "Initializing Terraform..."
	cd deployments/terraform && terraform init

terraform-plan: ## Plan Terraform changes
	@echo "Planning Terraform changes..."
	cd deployments/terraform && terraform plan

terraform-apply: ## Apply Terraform changes
	@echo "Applying Terraform changes..."
	cd deployments/terraform && terraform apply

terraform-destroy: ## Destroy Terraform infrastructure
	@echo "Destroying Terraform infrastructure..."
	cd deployments/terraform && terraform destroy
