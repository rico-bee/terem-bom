# BOM Weather Data Processor Makefile

# Build variables
BINARY_NAME=bom
BUILD_DIR=bin
MAIN_PATH=cmd/bom/main.go

# Go variables
GO=go
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)

# Build for specific platform
.PHONY: build-linux
build-linux:
	@echo "Building $(BINARY_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

.PHONY: build-darwin
build-darwin:
	@echo "Building $(BINARY_NAME) for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

.PHONY: build-windows
build-windows:
	@echo "Building $(BINARY_NAME) for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Deep clean: also clear Go build cache
.PHONY: distclean
distclean: clean
	@echo "Clearing Go build cache..."
	go clean -cache -modcache -testcache -benchcache

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	$(GO) mod download
	$(GO) mod tidy

# Run the application (requires build first)
# Usage: make run ARGS="convert -i input.csv -o output.json"
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME) $(ARGS)..."
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-linux    - Build for Linux"
	@echo "  build-darwin   - Build for macOS"
	@echo "  build-windows  - Build for Windows"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  distclean      - Deep clean: also clear Go build cache"
	@echo "  deps           - Install dependencies"
	@echo "  run            - Build and run the application (use ARGS='command options')"
	@echo "  help           - Show this help message"

# Default target
.DEFAULT_GOAL := build 