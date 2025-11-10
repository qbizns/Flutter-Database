.PHONY: help build run test lint migrate migrate-up migrate-down migrate-down-all migrate-status migrate-force migrate-goto clean docker-up docker-down

# Variables
APP_NAME=pos-backend
API_CMD=./cmd/api
GRPC_CMD=./cmd/grpc
WORKER_CMD=./cmd/worker
BUILD_DIR=./bin

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Migration parameters
MIGRATE=~/go/bin/migrate
MIGRATIONS_POSTGRES=./db/migrations-postgres
MIGRATIONS_ACCOUNTING=./db/migrations-accounting
DATABASE_URL?=postgres://postgres:postgres@localhost:5432/pos_saas?sslmode=disable

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) tidy

build-api: ## Build API server
	$(GOBUILD) -o $(BUILD_DIR)/api $(API_CMD)

build-grpc: ## Build gRPC server
	$(GOBUILD) -o $(BUILD_DIR)/grpc $(GRPC_CMD)

build-worker: ## Build background worker
	$(GOBUILD) -o $(BUILD_DIR)/worker $(WORKER_CMD)

build: build-api build-grpc build-worker ## Build all binaries

run-api: ## Run API server
	$(GOCMD) run $(API_CMD)/main.go

run-grpc: ## Run gRPC server
	$(GOCMD) run $(GRPC_CMD)/main.go

run-worker: ## Run background worker
	$(GOCMD) run $(WORKER_CMD)/main.go

test: ## Run tests
	$(GOTEST) -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	$(GOCMD) tool cover -html=coverage.out

lint: ## Run linter
	golangci-lint run ./...

fmt: ## Format code
	$(GOFMT) ./...
	goimports -w .

vet: ## Run go vet
	$(GOVET) ./...

check: fmt vet lint test ## Run all checks

migrate-up: ## Run all database migrations up
	@echo "Running PostgreSQL migrations..."
	$(MIGRATE) -path $(MIGRATIONS_POSTGRES) -database "$(DATABASE_URL)" up
	@echo "Running Accounting migrations..."
	$(MIGRATE) -path $(MIGRATIONS_ACCOUNTING) -database "$(DATABASE_URL)" up
	@echo "All migrations applied successfully!"

migrate-down: ## Rollback last migration (careful!)
	@echo "Rolling back Accounting migrations (last 1 step)..."
	$(MIGRATE) -path $(MIGRATIONS_ACCOUNTING) -database "$(DATABASE_URL)" down 1
	@echo "Rolling back PostgreSQL migrations (last 1 step)..."
	$(MIGRATE) -path $(MIGRATIONS_POSTGRES) -database "$(DATABASE_URL)" down 1

migrate-down-all: ## Rollback ALL migrations (DANGEROUS!)
	@echo "WARNING: Rolling back ALL migrations..."
	@read -p "Are you sure? Type 'yes' to continue: " confirm && [ "$$confirm" = "yes" ]
	$(MIGRATE) -path $(MIGRATIONS_ACCOUNTING) -database "$(DATABASE_URL)" down -all
	$(MIGRATE) -path $(MIGRATIONS_POSTGRES) -database "$(DATABASE_URL)" down -all

migrate-status: ## Show migration status
	@echo "PostgreSQL migrations status:"
	@$(MIGRATE) -path $(MIGRATIONS_POSTGRES) -database "$(DATABASE_URL)" version || echo "No migrations applied"
	@echo "\nAccounting migrations status:"
	@$(MIGRATE) -path $(MIGRATIONS_ACCOUNTING) -database "$(DATABASE_URL)" version || echo "No migrations applied"

migrate-force: ## Force migration version (use: make migrate-force VERSION=X)
	@echo "Forcing migration version to $(VERSION)..."
	$(MIGRATE) -path $(MIGRATIONS_POSTGRES) -database "$(DATABASE_URL)" force $(VERSION)
	$(MIGRATE) -path $(MIGRATIONS_ACCOUNTING) -database "$(DATABASE_URL)" force $(VERSION)

migrate-goto: ## Go to specific migration version (use: make migrate-goto VERSION=X)
	@echo "Migrating to version $(VERSION)..."
	$(MIGRATE) -path $(MIGRATIONS_POSTGRES) -database "$(DATABASE_URL)" goto $(VERSION)
	$(MIGRATE) -path $(MIGRATIONS_ACCOUNTING) -database "$(DATABASE_URL)" goto $(VERSION)

seed: ## Seed database with test data
	psql -h localhost -U postgres -d pos_saas -f ../postgres/seed_data/*.sql
	psql -h localhost -U postgres -d pos_saas -f ../accounting/seed_data/*.sql

docker-up: ## Start Docker containers
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)
	rm -f coverage.out

install-tools: ## Install development tools
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

proto-gen: ## Generate protobuf code
	@echo "Generating protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/**/*.proto

.DEFAULT_GOAL := help
