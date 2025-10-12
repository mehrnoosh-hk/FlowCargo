# Development Tools and Utilities
# This file contains targets for installing and managing development tools

.PHONY: install-tools
install-tools: ## Install all development tools
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/gordonklaus/ineffassign@latest
	go install github.com/client9/misspell/cmd/misspell@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "All tools installed successfully!"

.PHONY: install-air
install-air: ## Install air for hot reloading
	go install github.com/cosmtrek/air@latest

.PHONY: install-lint
install-lint: ## Install golangci-lint
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

.PHONY: install-migrate
install-migrate: ## Install migrate tool
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION)

.PHONY: install-goimports
install-goimports: ## Install goimports
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: install-staticcheck
install-staticcheck: ## Install staticcheck
	go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: install-ineffassign
install-ineffassign: ## Install ineffassign
	go install github.com/gordonklaus/ineffassign@latest

.PHONY: install-misspell
install-misspell: ## Install misspell
	go install github.com/client9/misspell/cmd/misspell@latest

.PHONY: install-gosec
install-gosec: ## Install gosec
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest

.PHONY: install-swag
install-swag: ## Install swag CLI tool
	@echo "Installing swag..."
	@go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: swagger-gen
swagger-gen: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
	@echo "Swagger documentation generated successfully!"

.PHONY: swagger-fmt
swagger-fmt: ## Format Swagger annotations
	@echo "Formatting Swagger annotations..."
	@swag fmt
	@echo "Swagger annotations formatted successfully!"

.PHONY: swagger-validate
swagger-validate: swagger-gen ## Validate Swagger documentation
	@echo "Validating Swagger documentation..."
	@test -f docs/swagger.json || (echo "swagger.json not found" && exit 1)
	@echo "Swagger documentation is valid"

.PHONY: sqlc
sqlc: ## Generate Go code from SQL queries
	@echo "Generating SQLC code..."
	sqlc generate -f db/sqlc.yml
	@echo "SQLC code generated successfully!"

.PHONY: tools-version
tools-version: ## Show versions of installed tools
	@echo "=== Tool Versions ==="
	@echo "Go version: $$(go version)"
	@echo "Air version: $$(air -v 2>/dev/null || echo 'Not installed')"
	@echo "Golangci-lint version: $$(golangci-lint version 2>/dev/null || echo 'Not installed')"
	@echo "Migrate version: $$(migrate -version 2>/dev/null || echo 'Not installed')"
	@echo "Goimports version: $$(goimports -version 2>/dev/null || echo 'Not installed')"
	@echo "Staticcheck version: $$(staticcheck -version 2>/dev/null || echo 'Not installed')"
	@echo "Ineffassign version: $$(ineffassign -version 2>/dev/null || echo 'Not installed')"
	@echo "Misspell version: $$(misspell -version 2>/dev/null || echo 'Not installed')"
	@echo "Gosec version: $$(gosec -version 2>/dev/null || echo 'Not installed')"

.PHONY: tools-update
tools-update: ## Update all development tools to latest versions
	@echo "Updating development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/gordonklaus/ineffassign@latest
	go install github.com/client9/misspell/cmd/misspell@latest
	go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@echo "All tools updated successfully!"

.PHONY: tools-check
tools-check: ## Check if all required tools are installed
	@echo "Checking required tools..."
	@command -v air >/dev/null 2>&1 || { echo "air is not installed. Run 'make install-air'"; exit 1; }
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint is not installed. Run 'make install-lint'"; exit 1; }
	@command -v migrate >/dev/null 2>&1 || { echo "migrate is not installed. Run 'make install-migrate'"; exit 1; }
	@command -v goimports >/dev/null 2>&1 || { echo "goimports is not installed. Run 'make install-goimports'"; exit 1; }
	@command -v staticcheck >/dev/null 2>&1 || { echo "staticcheck is not installed. Run 'make install-staticcheck'"; exit 1; }
	@command -v ineffassign >/dev/null 2>&1 || { echo "ineffassign is not installed. Run 'make install-ineffassign'"; exit 1; }
	@command -v misspell >/dev/null 2>&1 || { echo "misspell is not installed. Run 'make install-misspell'"; exit 1; }
	@command -v gosec >/dev/null 2>&1 || { echo "gosec is not installed. Run 'make install-gosec'"; exit 1; }
	@echo "All required tools are installed!"