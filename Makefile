# Makefile for the Go project

# Variables
APP_NAME := "versionbump"
VERSION := "v0.6.0-alpha"
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
PKG := ./...
INTEGRATION_TEST_DIR := ./test/integration
DIST_DIR := dist
DIST_FILES := $(wildcard $(DIST_DIR)/*)
OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
BINARY_NAME := $(APP_NAME)-$(VERSION)-$(OS)-$(ARCH)
TARBALL := $(BINARY_NAME).tgz
ZIPFILE := $(BINARY_NAME).zip

# Commands
.PHONY: all run build test test-all test-integration test-integration-verbose deps clean tidy lint format init-project dist sign-dist

all: test build

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	go run main.go

clean: ## Remove previous build artifacts and test cache
	@echo "Cleaning up..."
	go clean -modcache -cache -testcache -i -r
	rm -rf $(DIST_DIR)

build: ## Build the project and output distribution binaries
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) cmd/versionbump/main.go

test: ## Run unit tests
	@echo "Running unit tests..."
	go test $(PKG) --cover

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	go test $(INTEGRATION_TEST_DIR) --tags=integration

test-integration-verbose: ## Run integration tests with verbose output
	@echo "Running integration tests (verbose)..."
	go test $(INTEGRATION_TEST_DIR) --tags=integration -v

test-all: test test-integration ## Run all tests (unit and integration)
	@echo "Running all tests..."

test-cover-serve: ## Run tests with coverage and serve the coverage report
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out $(PKG)
	go tool cover -html=coverage.out

deps: ## Fetch and update project dependencies
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

tidy: ## Tidy up go.mod and remove unused dependencies
	@echo "Tidying up go.mod..."
	go mod tidy

lint: ## Lint the Go source files
	@echo "Linting..."
	golangci-lint run
	go vet ./...


format: ## Format all Go source files
	@echo "Formatting Go source files..."
	gofmt -w $(GO_FILES)

init-project: ## Initialize the project by installing necessary tools
	@echo "Installing necessary Go tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	#go get -u github.com/spf13/cobra@latest
	go mod tidy

dist: clean ## Create binary distributions for common OS and architecture combinations
	@echo "Creating distributions for all supported OS/ARCH combinations..."
	mkdir -p $(DIST_DIR)
	@for os in linux darwin windows; do \
		for arch in amd64 arm64; do \
			EXT=""; \
			if [ "$$os" = "windows" ]; then EXT=".exe"; fi; \
			GOOS=$$os GOARCH=$$arch go build -o $(DIST_DIR)/$(APP_NAME)$$EXT cmd/$(APP_NAME)/*.go; \
			cp README.md $(DIST_DIR)/; \
			tar -czvf $(DIST_DIR)/$(APP_NAME)-$(VERSION)-$$os-$$arch.tgz -C $(DIST_DIR) $(APP_NAME)$$EXT README.md; \
			zip -j $(DIST_DIR)/$(APP_NAME)-$(VERSION)-$$os-$$arch.zip $(DIST_DIR)/$(APP_NAME)$$EXT README.md; \
			rm $(DIST_DIR)/$(APP_NAME)$$EXT; \
			rm $(DIST_DIR)/README.md; \
		done; \
	done
	@echo "All distributions created in the $(DIST_DIR) directory."

sign-dist: dist ## Sign the distribution files
	@echo "Signing distribution files..."
	for file in $(wildcard ./dist/*); do \
		gpg --detach-sign --armor "$$file"; \
	done
	@echo "Done."
