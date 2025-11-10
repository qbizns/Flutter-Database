.PHONY: help build run test lint migrate clean docker-up docker-down

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

migrate-up: ## Run database migrations up
	@echo "Running migrations..."
	psql -h localhost -U postgres -d pos_saas -f ../postgres/migrations/V*.sql
	psql -h localhost -U postgres -d pos_saas -f ../accounting/migrations/V*.sql

migrate-down: ## Run database migrations down (careful!)
	@echo "Rolling back migrations not implemented - manually drop tables if needed"

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
