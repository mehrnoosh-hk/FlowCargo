# Variables
BINARY_NAME=flowcargo
MAIN_PATH=./cmd/main.go
MIGRATION_DIR=./db/migrations

# Load environment variables from .env.dev if it exists
ifneq (,$(wildcard .env.dev))
    include .env.dev
    export
endif

# Default migration DSN if not set in environment
MIGRATION_DSN ?= postgres://user:password@localhost/dbname?sslmode=disable

# Default target
.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## Run the main application
	go run $(MAIN_PATH)

.PHONY: build
build: ## Build the application
	go build -o $(BINARY_NAME) $(MAIN_PATH)

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: test-cover
test-cover: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: fmt
fmt: ## Format Go code
	go fmt ./...

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint run

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	golangci-lint run --fix

.PHONY: mod-tidy
mod-tidy: ## Tidy and verify modules
	go mod tidy
	go mod verify

.PHONY: create-migration
create-migration: ## Create a new migration file with sequential numbering (usage: make create-migration NAME=migration_name)
	@if [ -z "$(NAME)" ]; then echo "Usage: make create-migration NAME=migration_name"; exit 1; fi
	migrate create -ext sql -seq -digits 3 -dir $(MIGRATION_DIR) $(NAME)

.PHONY: migrate-down
migrate-down: ## Rollback the last migration
	migrate -database $(MIGRATION_DSN) -path $(MIGRATION_DIR) down

.PHONY: migrate-up-steps
migrate-up-steps: ## Migrate up by specific number of steps (usage: make migrate-up-steps STEPS=2)
	@if [ -z "$(STEPS)" ]; then echo "Usage: make migrate-up-steps STEPS=n"; exit 1; fi
	migrate -database $(MIGRATION_DSN) -path $(MIGRATION_DIR) up $(STEPS)

.PHONY: migrate-down-steps
migrate-down-steps: ## Migrate down by specific number of steps (usage: make migrate-down-steps STEPS=2)
	@if [ -z "$(STEPS)" ]; then echo "Usage: make migrate-down-steps STEPS=n"; exit 1; fi
	migrate -database $(MIGRATION_DSN) -path $(MIGRATION_DIR) down $(STEPS)

.PHONY: migrate-status
migrate-status: ## Show migration status
	migrate -database $(MIGRATION_DSN) -path $(MIGRATION_DIR) version

.PHONY: migrate-force
migrate-force: ## Force set migration version (usage: make migrate-force VERSION=1)
	@if [ -z "$(VERSION)" ]; then echo "Usage: make migrate-force VERSION=n"; exit 1; fi
	migrate -database $(MIGRATION_DSN) -path $(MIGRATION_DIR) force $(VERSION)

.PHONY: clean
clean: ## Clean build artifacts
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

.PHONY: dev
dev: ## Run in development mode with hot reload
	air

.PHONY: install-tools
install-tools: ## Install development tools
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: all
all: fmt vet test build ## Run format, vet, test, and build
