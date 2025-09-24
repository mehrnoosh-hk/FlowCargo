# FlowCargo - Main Makefile
# This is the root Makefile that includes all domain-specific Makefiles

# Include variables first
include makefiles/variables.mk

# Include all domain-specific Makefiles
include makefiles/build.mk
include makefiles/test.mk
include makefiles/lint.mk
include makefiles/migration.mk
include makefiles/tools.mk

# Default target
.PHONY: help
help: ## Display this help screen
	@echo "FlowCargo - Available Commands:"
	@echo "================================="
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

# Composite targets that combine multiple operations
.PHONY: all
all: fmt vet test build ## Run format, vet, test, and build

.PHONY: ci
ci: quality test-cover build ## Run CI pipeline (quality checks, tests, and build)

.PHONY: dev-setup
dev-setup: install-tools migrate-up ## Setup development environment

.PHONY: pre-commit
pre-commit: quality-fix test ## Run pre-commit checks (format, lint, test)

.PHONY: release
release: clean test-cover build-all ## Prepare release (clean, test, build all platforms)

# Quick development targets
.PHONY: quick-test
quick-test: test-short ## Run quick tests

.PHONY: quick-build
quick-build: build ## Quick build for development

.PHONY: quick-dev
quick-dev: dev ## Start development server with hot reload

# Database management shortcuts
.PHONY: db-reset
db-reset: migrate-reset ## Reset database (DANGEROUS)

.PHONY: db-status
db-status: migrate-status ## Show database migration status

.PHONY: db-up
db-up: migrate-up ## Apply all pending migrations

.PHONY: db-down
db-down: migrate-down ## Rollback last migration

# Docker commands
.PHONY: docker-compose-up
docker-compose-up: ## Start docker containers
	docker-compose -f docker/dev/postgres/docker-compose.yml up -d

.PHONY: docker-compose-down
docker-compose-down: ## Stop docker containers
	docker-compose -f docker/dev/postgres/docker-compose.yml down

# Cleanup targets
.PHONY: clean-all
clean-all: clean test-clean ## Clean all artifacts (build and test)

.PHONY: fresh
fresh: clean-all install-tools ## Fresh start (clean everything and reinstall tools)

# Documentation and information
.PHONY: info
info: ## Show project information
	@echo "FlowCargo Project Information"
	@echo "============================="
	@echo "Binary Name: $(BINARY_NAME)"
	@echo "Main Path: $(MAIN_PATH)"
	@echo "Migration Directory: $(MIGRATION_DIR)"
	@echo "Migration DSN: $(MIGRATION_DSN)"
	@echo "Go Version: $(GO_VERSION)"
	@echo "Build Output: $(BUILD_OUTPUT)"
	@echo ""

.PHONY: targets
targets: ## List all available targets
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:' $(MAKEFILE_LIST) | sed 's/://' | sort