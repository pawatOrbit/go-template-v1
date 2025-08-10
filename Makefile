# Go API Template Makefile
.PHONY: help build test run clean docker-build docker-run docker-stop docker-dev lint fmt vet mod-tidy mod-download sqlc-generate db-up db-down db-migrate db-seed coverage security audit

# Variables
BINARY_NAME=go-api-template
BINARY_PATH=./bin/$(BINARY_NAME)
DOCKER_IMAGE=go-api-template
DOCKER_TAG=latest
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*" -not -path "./tests/*")
TEST_PATTERN=./...
COVERAGE_FILE=coverage.out

# Default target
help: ## Show this help message
	@echo "Go API Template - Available Commands:"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-20s %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""

# Development
run: ## Run the application locally
	@echo "Running application..."
	go run main.go serve:all-api

run-dev: ## Run the application with live reload (requires air)
	@echo "Running application with live reload..."
	air -c .air.toml

build: ## Build the application
	@echo "Building application..."
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o $(BINARY_PATH) .
	@echo "Binary built: $(BINARY_PATH)"

build-local: ## Build the application for local OS
	@echo "Building application for local OS..."
	mkdir -p bin
	go build -ldflags="-w -s" -o $(BINARY_PATH) .
	@echo "Binary built: $(BINARY_PATH)"

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	go clean -cache
	go clean -modcache

# Testing
test: ## Run tests
	@echo "Running tests..."
	go test $(TEST_PATTERN) -v

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	go test ./tests/unit/... -v

test-integration: ## Run integration tests only
	@echo "Running integration tests..."
	go test ./tests/integration/... -v

test-race: ## Run tests with race condition detection
	@echo "Running tests with race detection..."
	go test $(TEST_PATTERN) -race -v

coverage: ## Run tests with coverage report
	@echo "Running tests with coverage..."
	go test $(TEST_PATTERN) -coverprofile=$(COVERAGE_FILE)
	go tool cover -html=$(COVERAGE_FILE) -o coverage.html
	@echo "Coverage report generated: coverage.html"

coverage-func: ## Show coverage by function
	@echo "Coverage by function..."
	go test $(TEST_PATTERN) -coverprofile=$(COVERAGE_FILE)
	go tool cover -func=$(COVERAGE_FILE)

benchmark: ## Run benchmarks
	@echo "Running benchmarks..."
	go test $(TEST_PATTERN) -bench=. -benchmem

# Code Quality
lint: ## Run golangci-lint
	@echo "Running golangci-lint..."
	golangci-lint run

fmt: ## Format Go code
	@echo "Formatting code..."
	go fmt ./...
	gofmt -s -w $(GO_FILES)

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

security: ## Run security scan with gosec
	@echo "Running security scan..."
	gosec ./...

audit: ## Run security audit
	@echo "Running security audit..."
	go list -json -m all | nancy sleuth

# Dependencies
mod-tidy: ## Tidy dependencies
	@echo "Tidying dependencies..."
	go mod tidy

mod-download: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

mod-verify: ## Verify dependencies
	@echo "Verifying dependencies..."
	go mod verify

# Database
sqlc-generate: ## Generate code from SQL using sqlc
	@echo "Generating code from SQL..."
	sqlc generate

db-up: ## Start database with docker-compose
	@echo "Starting database..."
	docker-compose up -d postgres redis

db-down: ## Stop database
	@echo "Stopping database..."
	docker-compose down

db-reset: ## Reset database (stop, remove volumes, start)
	@echo "Resetting database..."
	docker-compose down -v
	docker-compose up -d postgres redis
	sleep 5
	@echo "Database reset complete"

# Migration Commands
migrate: ## Run all pending migrations
	@echo "Running database migrations..."
	go run main.go migrate up

migrate-up: ## Run all pending migrations (alias for migrate)
	@echo "Running database migrations..."
	go run main.go migrate up

migrate-down: ## Rollback last migration
	@echo "Rolling back last migration..."
	go run main.go migrate down --steps=1

migrate-down-all: ## Rollback all migrations
	@echo "Rolling back all migrations..."
	go run main.go migrate down --all

migrate-status: ## Show migration status
	@echo "Checking migration status..."
	go run main.go migrate status

migrate-version: ## Show current migration version
	@echo "Current migration version:"
	go run main.go migrate version

migrate-create: ## Create a new migration (usage: make migrate-create name=create_users_table)
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=migration_name"; exit 1; fi
	@echo "Creating migration: $(name)"
	go run main.go migrate create $(name)

migrate-force: ## Force database to specific version (usage: make migrate-force version=1)
	@if [ -z "$(version)" ]; then echo "Usage: make migrate-force version=N"; exit 1; fi
	@echo "Forcing database to version $(version)..."
	go run main.go migrate force $(version)

db-seed: ## Seed database with sample data
	@echo "Seeding database..."
	docker-compose exec postgres psql -U postgres -d go_api_template -f /docker-entrypoint-initdb.d/init-db.sql

# Docker
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run application in Docker
	@echo "Running application in Docker..."
	docker-compose up -d

docker-dev: ## Run application in Docker development mode
	@echo "Running application in Docker development mode..."
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

docker-stop: ## Stop Docker containers
	@echo "Stopping Docker containers..."
	docker-compose down

docker-logs: ## Show Docker logs
	@echo "Showing Docker logs..."
	docker-compose logs -f

docker-clean: ## Clean Docker resources
	@echo "Cleaning Docker resources..."
	docker-compose down -v --remove-orphans
	docker system prune -f
	docker volume prune -f

docker-shell: ## Open shell in API container
	@echo "Opening shell in API container..."
	docker-compose exec api sh

# Tools installation
install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securecodewarrior/goat-cli/v2/cmd/goat@latest
	go install github.com/sonatype-nexus-community/nancy@latest

# CI/CD
ci: lint test ## Run CI pipeline locally
	@echo "Running CI pipeline..."

pre-commit: fmt lint test ## Run pre-commit checks
	@echo "Running pre-commit checks..."

release-dry: ## Simulate a release
	@echo "Simulating release..."
	goreleaser release --snapshot --rm-dist

# Information
info: ## Show project information
	@echo "Project Information:"
	@echo "  Binary Name: $(BINARY_NAME)"
	@echo "  Binary Path: $(BINARY_PATH)"
	@echo "  Docker Image: $(DOCKER_IMAGE):$(DOCKER_TAG)"
	@echo "  Go Version: $(shell go version)"
	@echo "  Git Commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'N/A')"
	@echo "  Git Branch: $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo 'N/A')"

version: ## Show version information
	@echo "Go API Template"
	@echo "Git Commit: $(shell git rev-parse --short HEAD 2>/dev/null || echo 'N/A')"
	@echo "Build Time: $(shell date)"

# Health checks
health-check: ## Check if application is healthy
	@echo "Checking application health..."
	curl -f http://localhost:8080/health || exit 1

# Log viewing
logs: ## Show application logs (when running with docker-compose)
	docker-compose logs -f api

logs-db: ## Show database logs
	docker-compose logs -f postgres

logs-redis: ## Show Redis logs  
	docker-compose logs -f redis