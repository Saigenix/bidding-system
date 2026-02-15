# ============================================================================
# Bidding System SDK — Makefile
# ============================================================================

# Variables
APP_NAME     = bidding-server
CMD_DIR      = cmd/server
BINARY       = $(APP_NAME).exe
MIGRATE_PATH = migrations
DB_URL       ?= postgres://postgres:password@localhost:5432/bidding?sslmode=disable

# Go
GOFLAGS      = -v
LDFLAGS      = -s -w

# ============================================================================
# Development
# ============================================================================

## run: Run the server in development mode
.PHONY: run
run:
	@echo "▶ Starting development server..."
	go run $(CMD_DIR)/main.go

## build: Build the binary
.PHONY: build
build:
	@echo "▶ Building $(BINARY)..."
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY) $(CMD_DIR)/main.go
	@echo "✓ Built $(BINARY)"

## clean: Remove build artifacts
.PHONY: clean
clean:
	@echo "▶ Cleaning..."
	@if exist $(BINARY) del /Q $(BINARY)
	go clean
	@echo "✓ Clean"

## tidy: Tidy and verify dependencies
.PHONY: tidy
tidy:
	@echo "▶ Tidying modules..."
	go mod tidy
	go mod verify
	@echo "✓ Modules verified"

## fmt: Format all Go files
.PHONY: fmt
fmt:
	@echo "▶ Formatting..."
	go fmt ./...
	@echo "✓ Formatted"

## fmt-check: Check formatting without modifying files (for CI)
.PHONY: fmt-check
fmt-check:
	@echo "▶ Checking formatting..."
	@test -z "$$(gofmt -l .)" || (echo "✗ The following files are not formatted:" && gofmt -l . && exit 1)
	@echo "✓ All files formatted"

## vet: Run go vet
.PHONY: vet
vet:
	@echo "▶ Vetting..."
	go vet ./...
	@echo "✓ Vetted"

## lint: Run fmt + vet together
.PHONY: lint
lint: fmt vet

## test: Run all tests
.PHONY: test
test:
	@echo "▶ Running tests..."
	go test ./... -v -count=1

## test-cover: Run tests with coverage
.PHONY: test-cover
test-cover:
	@echo "▶ Running tests with coverage..."
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

# ============================================================================
# Swagger
# ============================================================================

## swagger: Generate Swagger documentation
.PHONY: swagger
swagger:
	@echo "▶ Generating Swagger docs..."
	swag init -g cmd/server/main.go -o docs/swagger
	@echo "✓ Swagger docs generated at docs/swagger/"

## swagger-install: Install swag CLI tool
.PHONY: swagger-install
swagger-install:
	@echo "▶ Installing swag CLI..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✓ swag CLI installed"

# ============================================================================
# Database
# ============================================================================

## db-up: Start PostgreSQL in Docker
.PHONY: db-up
db-up:
	@echo "▶ Starting PostgreSQL..."
	docker run --name bidding-postgres \
		-e POSTGRES_PASSWORD=password \
		-e POSTGRES_DB=bidding \
		-p 5432:5432 \
		-d postgres:16
	@echo "✓ PostgreSQL started on port 5432"

## db-down: Stop and remove PostgreSQL container
.PHONY: db-down
db-down:
	@echo "▶ Stopping PostgreSQL..."
	docker stop bidding-postgres
	docker rm bidding-postgres
	@echo "✓ PostgreSQL stopped"

## migrate-up: Run database migrations
.PHONY: migrate-up
migrate-up:
	@echo "▶ Running migrations..."
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" up
	@echo "✓ Migrations applied"

## migrate-down: Rollback database migrations
.PHONY: migrate-down
migrate-down:
	@echo "▶ Rolling back migrations..."
	migrate -path $(MIGRATE_PATH) -database "$(DB_URL)" down
	@echo "✓ Migrations rolled back"

## migrate-install: Install golang-migrate CLI
.PHONY: migrate-install
migrate-install:
	@echo "▶ Installing golang-migrate..."
	go install -tags "postgres" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "✓ migrate CLI installed"

# ============================================================================
# Docker
# ============================================================================

## docker-build: Build Docker image
.PHONY: docker-build
docker-build:
	@echo "▶ Building Docker image..."
	docker build -t $(APP_NAME) .
	@echo "✓ Docker image built"

## docker-run: Run Docker container
.PHONY: docker-run
docker-run:
	@echo "▶ Running Docker container..."
	docker run -p 8080:8080 --env-file .env $(APP_NAME)

# ============================================================================
# Quick Start
# ============================================================================

## setup: Full first-time setup (install deps, start db, migrate, run)
.PHONY: setup
setup: tidy db-up
	@echo "▶ Waiting for PostgreSQL to be ready..."
	@timeout /t 5 /nobreak > nul 2>&1 || sleep 5
	$(MAKE) migrate-up
	@echo ""
	@echo "============================================"
	@echo "  ✓ Setup complete!"
	@echo "  Run 'make run' to start the server"
	@echo "============================================"

## ci: Run all CI checks (format, vet, test)
.PHONY: ci
ci: fmt-check vet test
	@echo ""
	@echo "✓ All CI checks passed!"

## all: Build and run
.PHONY: all
all: lint build run

# ============================================================================
# Help
# ============================================================================

## help: Show this help message
.PHONY: help
help:
	@echo ""
	@echo "Bidding System SDK — Available Commands"
	@echo "========================================"
	@echo ""
	@echo "Development:"
	@echo "  make run              Run the server"
	@echo "  make build            Build the binary"
	@echo "  make clean            Remove build artifacts"
	@echo "  make tidy             Tidy Go modules"
	@echo "  make fmt              Format code"
	@echo "  make vet              Run go vet"
	@echo "  make lint             Format + vet"
	@echo "  make test             Run tests"
	@echo "  make test-cover       Run tests with coverage"
	@echo "  make fmt-check        Check formatting (CI)"
	@echo "  make ci               Run all CI checks"
	@echo ""
	@echo "Database:"
	@echo "  make db-up            Start PostgreSQL (Docker)"
	@echo "  make db-down          Stop PostgreSQL"
	@echo "  make migrate-up       Apply migrations"
	@echo "  make migrate-down     Rollback migrations"
	@echo "  make migrate-install  Install migrate CLI"
	@echo ""
	@echo "Swagger:"
	@echo "  make swagger          Generate Swagger docs"
	@echo "  make swagger-install  Install swag CLI"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build     Build Docker image"
	@echo "  make docker-run       Run Docker container"
	@echo ""
	@echo "Quick Start:"
	@echo "  make setup            Full first-time setup"
	@echo "  make all              Lint, build, and run"
	@echo ""

.DEFAULT_GOAL := help
