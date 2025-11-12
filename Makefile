# Cross-Platform Build System for POS Backend
# Supports: Linux (amd64, arm64), macOS (amd64, arm64/M1), Windows (amd64)

.PHONY: all build clean test install uninstall help \
        build-linux build-macos build-windows build-all \
        package-linux package-macos package-windows package-all \
        install-linux install-macos install-windows \
        docker run dev deps

# Version information
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Project information
PROJECT_NAME := pos-backend
BINARY_NAME := pos-backend
WORKER_NAME := pos-worker
MIGRATE_NAME := pos-migrate

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOVET := $(GOCMD) vet
GOFMT := $(GOCMD) fmt

# Directories
BUILD_DIR := build
BIN_DIR := $(BUILD_DIR)/bin
PKG_DIR := $(BUILD_DIR)/packages
DIST_DIR := $(BUILD_DIR)/dist
CMD_DIR := backend/cmd
SRC_DIR := backend

# Build flags
LDFLAGS := -ldflags "\
	-X main.Version=$(VERSION) \
	-X main.Commit=$(COMMIT) \
	-X main.BuildDate=$(BUILD_DATE) \
	-s -w"

# Platform detection
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_S),Linux)
    DETECTED_OS := linux
endif
ifeq ($(UNAME_S),Darwin)
    DETECTED_OS := macos
endif
ifeq ($(OS),Windows_NT)
    DETECTED_OS := windows
endif

# Colors for output
COLOR_RESET := \033[0m
COLOR_BOLD := \033[1m
COLOR_GREEN := \033[32m
COLOR_YELLOW := \033[33m
COLOR_BLUE := \033[34m
COLOR_CYAN := \033[36m

# Default target
all: clean deps test build

help: ## Show this help message
	@echo "$(COLOR_BOLD)$(COLOR_CYAN)POS Backend - Cross-Platform Build System$(COLOR_RESET)"
	@echo ""
	@echo "$(COLOR_BOLD)Version:$(COLOR_RESET) $(VERSION)"
	@echo "$(COLOR_BOLD)Commit:$(COLOR_RESET)  $(COMMIT)"
	@echo "$(COLOR_BOLD)Date:$(COLOR_RESET)    $(BUILD_DATE)"
	@echo ""
	@echo "$(COLOR_BOLD)Available targets:$(COLOR_RESET)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(COLOR_GREEN)%-20s$(COLOR_RESET) %s\n", $$1, $$2}'

deps: ## Install dependencies
	@echo "$(COLOR_CYAN)Installing dependencies...$(COLOR_RESET)"
	cd $(SRC_DIR) && $(GOMOD) download
	cd $(SRC_DIR) && $(GOMOD) verify
	@echo "$(COLOR_GREEN)✓ Dependencies installed$(COLOR_RESET)"

test: ## Run tests
	@echo "$(COLOR_CYAN)Running tests...$(COLOR_RESET)"
	cd $(SRC_DIR) && $(GOTEST) -v -race -coverprofile=coverage.out ./...
	cd $(SRC_DIR) && $(GOCMD) tool cover -func=coverage.out
	@echo "$(COLOR_GREEN)✓ Tests passed$(COLOR_RESET)"

lint: ## Run linters
	@echo "$(COLOR_CYAN)Running linters...$(COLOR_RESET)"
	cd $(SRC_DIR) && $(GOVET) ./...
	cd $(SRC_DIR) && $(GOFMT) ./...
	@command -v golangci-lint >/dev/null 2>&1 && cd $(SRC_DIR) && golangci-lint run || echo "golangci-lint not installed, skipping..."
	@echo "$(COLOR_GREEN)✓ Linting complete$(COLOR_RESET)"

clean: ## Clean build artifacts
	@echo "$(COLOR_YELLOW)Cleaning build artifacts...$(COLOR_RESET)"
	rm -rf $(BUILD_DIR)
	cd $(SRC_DIR) && $(GOCLEAN)
	@echo "$(COLOR_GREEN)✓ Cleaned$(COLOR_RESET)"

#################################
# Platform-specific builds
#################################

build: build-current ## Build for current platform

build-current: ## Build for detected platform ($(DETECTED_OS))
	@echo "$(COLOR_CYAN)Building for $(DETECTED_OS)...$(COLOR_RESET)"
	@$(MAKE) build-$(DETECTED_OS)

build-linux: ## Build for Linux (amd64 and arm64)
	@echo "$(COLOR_CYAN)Building for Linux amd64...$(COLOR_RESET)"
	mkdir -p $(BIN_DIR)/linux-amd64
	cd $(SRC_DIR) && GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/linux-amd64/$(BINARY_NAME) ./cmd/api
	cd $(SRC_DIR) && GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/linux-amd64/$(WORKER_NAME) ./cmd/worker
	cd $(SRC_DIR) && GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/linux-amd64/$(MIGRATE_NAME) ./cmd/migrate

	@echo "$(COLOR_CYAN)Building for Linux arm64...$(COLOR_RESET)"
	mkdir -p $(BIN_DIR)/linux-arm64
	cd $(SRC_DIR) && GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/linux-arm64/$(BINARY_NAME) ./cmd/api
	cd $(SRC_DIR) && GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/linux-arm64/$(WORKER_NAME) ./cmd/worker
	cd $(SRC_DIR) && GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/linux-arm64/$(MIGRATE_NAME) ./cmd/migrate
	@echo "$(COLOR_GREEN)✓ Linux builds complete$(COLOR_RESET)"

build-macos: ## Build for macOS (amd64 and arm64/M1)
	@echo "$(COLOR_CYAN)Building for macOS amd64 (Intel)...$(COLOR_RESET)"
	mkdir -p $(BIN_DIR)/darwin-amd64
	cd $(SRC_DIR) && GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/darwin-amd64/$(BINARY_NAME) ./cmd/api
	cd $(SRC_DIR) && GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/darwin-amd64/$(WORKER_NAME) ./cmd/worker
	cd $(SRC_DIR) && GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/darwin-amd64/$(MIGRATE_NAME) ./cmd/migrate

	@echo "$(COLOR_CYAN)Building for macOS arm64 (Apple Silicon M1/M2)...$(COLOR_RESET)"
	mkdir -p $(BIN_DIR)/darwin-arm64
	cd $(SRC_DIR) && GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/darwin-arm64/$(BINARY_NAME) ./cmd/api
	cd $(SRC_DIR) && GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/darwin-arm64/$(WORKER_NAME) ./cmd/worker
	cd $(SRC_DIR) && GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/darwin-arm64/$(MIGRATE_NAME) ./cmd/migrate
	@echo "$(COLOR_GREEN)✓ macOS builds complete$(COLOR_RESET)"

build-windows: ## Build for Windows (amd64)
	@echo "$(COLOR_CYAN)Building for Windows amd64...$(COLOR_RESET)"
	mkdir -p $(BIN_DIR)/windows-amd64
	cd $(SRC_DIR) && GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/windows-amd64/$(BINARY_NAME).exe ./cmd/api
	cd $(SRC_DIR) && GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/windows-amd64/$(WORKER_NAME).exe ./cmd/worker
	cd $(SRC_DIR) && GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o ../$(BIN_DIR)/windows-amd64/$(MIGRATE_NAME).exe ./cmd/migrate
	@echo "$(COLOR_GREEN)✓ Windows build complete$(COLOR_RESET)"

build-all: build-linux build-macos build-windows ## Build for all platforms
	@echo "$(COLOR_GREEN)$(COLOR_BOLD)✓ All platform builds complete!$(COLOR_RESET)"

#################################
# Packaging
#################################

package: package-current ## Create package for current platform

package-current: build-current ## Package for detected platform
	@$(MAKE) package-$(DETECTED_OS)

package-linux: build-linux ## Create DEB and RPM packages for Linux
	@echo "$(COLOR_CYAN)Creating Linux packages...$(COLOR_RESET)"
	@./scripts/package/build-deb.sh $(VERSION)
	@./scripts/package/build-rpm.sh $(VERSION)
	@echo "$(COLOR_GREEN)✓ Linux packages created$(COLOR_RESET)"

package-macos: build-macos ## Create PKG installer for macOS
	@echo "$(COLOR_CYAN)Creating macOS package...$(COLOR_RESET)"
	@./scripts/package/build-macos-pkg.sh $(VERSION)
	@echo "$(COLOR_GREEN)✓ macOS package created$(COLOR_RESET)"

package-windows: build-windows ## Create MSI installer for Windows
	@echo "$(COLOR_CYAN)Creating Windows installer...$(COLOR_RESET)"
	@./scripts/package/build-windows-msi.sh $(VERSION)
	@echo "$(COLOR_GREEN)✓ Windows installer created$(COLOR_RESET)"

package-all: build-all ## Create packages for all platforms
	@echo "$(COLOR_CYAN)Creating packages for all platforms...$(COLOR_RESET)"
	@$(MAKE) package-linux
	@$(MAKE) package-macos
	@$(MAKE) package-windows
	@echo "$(COLOR_GREEN)$(COLOR_BOLD)✓ All packages created!$(COLOR_RESET)"

#################################
# Installation
#################################

install: install-current ## Install on current platform

install-current: ## Install for detected platform
	@$(MAKE) install-$(DETECTED_OS)

install-linux: build-linux ## Install on Linux
	@echo "$(COLOR_CYAN)Installing on Linux...$(COLOR_RESET)"
	@sudo ./scripts/install/install-linux.sh
	@echo "$(COLOR_GREEN)✓ Installed on Linux$(COLOR_RESET)"

install-macos: build-macos ## Install on macOS
	@echo "$(COLOR_CYAN)Installing on macOS...$(COLOR_RESET)"
	@sudo ./scripts/install/install-macos.sh
	@echo "$(COLOR_GREEN)✓ Installed on macOS$(COLOR_RESET)"

install-windows: build-windows ## Install on Windows
	@echo "$(COLOR_CYAN)Installing on Windows...$(COLOR_RESET)"
	@powershell -ExecutionPolicy Bypass -File scripts/install/install-windows.ps1
	@echo "$(COLOR_GREEN)✓ Installed on Windows$(COLOR_RESET)"

uninstall: ## Uninstall from current platform
	@$(MAKE) uninstall-$(DETECTED_OS)

uninstall-linux: ## Uninstall from Linux
	@echo "$(COLOR_YELLOW)Uninstalling from Linux...$(COLOR_RESET)"
	@sudo ./scripts/install/uninstall-linux.sh
	@echo "$(COLOR_GREEN)✓ Uninstalled$(COLOR_RESET)"

uninstall-macos: ## Uninstall from macOS
	@echo "$(COLOR_YELLOW)Uninstalling from macOS...$(COLOR_RESET)"
	@sudo ./scripts/install/uninstall-macos.sh
	@echo "$(COLOR_GREEN)✓ Uninstalled$(COLOR_RESET)"

uninstall-windows: ## Uninstall from Windows
	@echo "$(COLOR_YELLOW)Uninstalling from Windows...$(COLOR_RESET)"
	@powershell -ExecutionPolicy Bypass -File scripts/install/uninstall-windows.ps1
	@echo "$(COLOR_GREEN)✓ Uninstalled$(COLOR_RESET)"

#################################
# Development
#################################

dev: ## Run in development mode
	@echo "$(COLOR_CYAN)Starting development server...$(COLOR_RESET)"
	cd $(SRC_DIR) && $(GOCMD) run cmd/api/main.go

run: build-current ## Build and run
	@echo "$(COLOR_CYAN)Running $(BINARY_NAME)...$(COLOR_RESET)"
	@$(BIN_DIR)/$(DETECTED_OS)-*/$(BINARY_NAME)

docker-build: ## Build Docker image
	@echo "$(COLOR_CYAN)Building Docker image...$(COLOR_RESET)"
	docker build -t $(PROJECT_NAME):$(VERSION) -f backend/Dockerfile backend/
	docker tag $(PROJECT_NAME):$(VERSION) $(PROJECT_NAME):latest
	@echo "$(COLOR_GREEN)✓ Docker image built$(COLOR_RESET)"

docker-run: docker-build ## Run in Docker
	@echo "$(COLOR_CYAN)Running in Docker...$(COLOR_RESET)"
	docker-compose up

#################################
# Distribution
#################################

dist: package-all ## Create distribution archives
	@echo "$(COLOR_CYAN)Creating distribution archives...$(COLOR_RESET)"
	mkdir -p $(DIST_DIR)

	# Linux amd64
	tar -czf $(DIST_DIR)/$(PROJECT_NAME)-$(VERSION)-linux-amd64.tar.gz \
		-C $(BIN_DIR)/linux-amd64 .

	# Linux arm64
	tar -czf $(DIST_DIR)/$(PROJECT_NAME)-$(VERSION)-linux-arm64.tar.gz \
		-C $(BIN_DIR)/linux-arm64 .

	# macOS amd64
	tar -czf $(DIST_DIR)/$(PROJECT_NAME)-$(VERSION)-darwin-amd64.tar.gz \
		-C $(BIN_DIR)/darwin-amd64 .

	# macOS arm64
	tar -czf $(DIST_DIR)/$(PROJECT_NAME)-$(VERSION)-darwin-arm64.tar.gz \
		-C $(BIN_DIR)/darwin-arm64 .

	# Windows
	cd $(BIN_DIR)/windows-amd64 && zip -r ../../dist/$(PROJECT_NAME)-$(VERSION)-windows-amd64.zip .

	@echo "$(COLOR_GREEN)$(COLOR_BOLD)✓ Distribution archives created in $(DIST_DIR)$(COLOR_RESET)"
	@ls -lh $(DIST_DIR)

release: clean test build-all package-all dist ## Create full release (all platforms)
	@echo "$(COLOR_GREEN)$(COLOR_BOLD)"
	@echo "╔════════════════════════════════════════════════════════════╗"
	@echo "║                  RELEASE COMPLETE! 🎉                      ║"
	@echo "╚════════════════════════════════════════════════════════════╝"
	@echo "$(COLOR_RESET)"
	@echo "Version: $(COLOR_CYAN)$(VERSION)$(COLOR_RESET)"
	@echo "Commit:  $(COLOR_CYAN)$(COMMIT)$(COLOR_RESET)"
	@echo ""
	@echo "Binaries:   $(BIN_DIR)"
	@echo "Packages:   $(PKG_DIR)"
	@echo "Archives:   $(DIST_DIR)"
	@echo ""
	@echo "$(COLOR_GREEN)Ready for distribution!$(COLOR_RESET)"

version: ## Show version information
	@echo "$(COLOR_BOLD)Version:$(COLOR_RESET)    $(VERSION)"
	@echo "$(COLOR_BOLD)Commit:$(COLOR_RESET)     $(COMMIT)"
	@echo "$(COLOR_BOLD)BuildDate:$(COLOR_RESET)  $(BUILD_DATE)"
	@echo "$(COLOR_BOLD)Platform:$(COLOR_RESET)   $(DETECTED_OS)"

.DEFAULT_GOAL := help
