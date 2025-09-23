# Linting and Code Quality Targets
# This file contains targets for code formatting, linting, and quality checks

.PHONY: fmt
fmt: ## Format Go code
	go fmt ./...

.PHONY: fmt-check
fmt-check: ## Check if code is formatted
	@if [ $$(gofmt -l . | wc -l) -ne 0 ]; then \
		echo "Code is not formatted. Run 'make fmt' to fix."; \
		gofmt -l .; \
		exit 1; \
	fi

.PHONY: vet
vet: ## Run go vet
	go vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	golangci-lint $(LINT_FLAGS)

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	golangci-lint run --fix

.PHONY: lint-config
lint-config: ## Generate golangci-lint configuration
	@if [ ! -f $(LINT_CONFIG) ]; then \
		golangci-lint config init; \
		echo "Generated $(LINT_CONFIG)"; \
	else \
		echo "$(LINT_CONFIG) already exists"; \
	fi

.PHONY: lint-check
lint-check: ## Check linting without fixing
	golangci-lint run --new-from-rev=HEAD~1

.PHONY: lint-fast
lint-fast: ## Run fast linting (skip slow linters)
	golangci-lint run --fast

.PHONY: lint-verbose
lint-verbose: ## Run linting with verbose output
	golangci-lint run -v

.PHONY: imports
imports: ## Organize imports
	goimports -w .

.PHONY: imports-check
imports-check: ## Check if imports are organized
	@if [ $$(goimports -l . | wc -l) -ne 0 ]; then \
		echo "Imports are not organized. Run 'make imports' to fix."; \
		goimports -l .; \
		exit 1; \
	fi

.PHONY: staticcheck
staticcheck: ## Run staticcheck
	staticcheck ./...

.PHONY: ineffassign
ineffassign: ## Check for ineffective assignments
	ineffassign ./...

.PHONY: misspell
misspell: ## Check for misspellings
	misspell ./...

.PHONY: gosec
gosec: ## Run security check
	gosec ./...

.PHONY: quality
quality: fmt-check imports-check vet lint ## Run all quality checks

.PHONY: quality-fix
quality-fix: fmt imports lint-fix ## Fix all quality issues
