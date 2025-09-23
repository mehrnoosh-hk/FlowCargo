# Build and Development Targets
# This file contains targets for building, running, and development

.PHONY: run
run: ## Run the main application
	go run $(MAIN_PATH)

.PHONY: build
build: ## Build the application
	go build $(BUILD_FLAGS) -o $(BUILD_OUTPUT) $(MAIN_PATH)

.PHONY: build-linux
build-linux: ## Build the application for Linux
	GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_OUTPUT)-linux $(MAIN_PATH)

.PHONY: build-windows
build-windows: ## Build the application for Windows
	GOOS=windows GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_OUTPUT).exe $(MAIN_PATH)

.PHONY: build-darwin
build-darwin: ## Build the application for macOS
	GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_OUTPUT)-darwin $(MAIN_PATH)

.PHONY: build-all
build-all: build-linux build-windows build-darwin ## Build for all platforms

.PHONY: dev
dev: ## Run in development mode with hot reload
	air

.PHONY: dev-config
dev-config: ## Generate air configuration file
	@if [ ! -f $(AIR_CONFIG) ]; then \
		air init; \
		echo "Generated $(AIR_CONFIG)"; \
	else \
		echo "$(AIR_CONFIG) already exists"; \
	fi

.PHONY: clean
clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-linux
	rm -f $(BINARY_NAME).exe
	rm -f $(BINARY_NAME)-darwin
	rm -f $(COVERAGE_OUTPUT) $(COVERAGE_HTML)

.PHONY: mod-tidy
mod-tidy: ## Tidy and verify modules
	go mod tidy
	go mod verify

.PHONY: mod-download
mod-download: ## Download Go modules
	go mod download

.PHONY: mod-graph
mod-graph: ## Show module dependency graph
	go mod graph
