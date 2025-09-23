# Database Migration Targets
# This file contains targets for database migrations and schema management

.PHONY: create-migration
create-migration: ## Create a new migration file with sequential numbering (usage: make create-migration NAME=migration_name)
	@if [ -z "$(NAME)" ]; then echo "Usage: make create-migration NAME=migration_name"; exit 1; fi
	migrate create -ext $(MIGRATION_EXT) $(MIGRATION_SEQ) -digits $(MIGRATION_DIGITS) -dir $(MIGRATION_DIR) $(NAME)

.PHONY: migrate-up
migrate-up: ## Apply all pending migrations
	migrate -database $(MIGRATION_DSN) -path $(MIGRATION_DIR) up

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

.PHONY: migrate-version
migrate-version: ## Show migration status
	migrate -database $(MIGRATION_DSN) -path $(MIGRATION_DIR) version

.PHONY: migrate-force
migrate-force: ## Force set migration version (usage: make migrate-force VERSION=1)
	@if [ -z "$(VERSION)" ]; then echo "Usage: make migrate-force VERSION=n"; exit 1; fi
	migrate -database $(MIGRATION_DSN) -path $(MIGRATION_DIR) force $(VERSION)

.PHONY: migrate-drop
migrate-drop: ## Drop all tables (DANGEROUS - use with caution)
	@echo "WARNING: This will drop all tables. Are you sure? (y/N)"
	@read -r confirm && [ "$$confirm" = "y" ] || exit 1
	migrate -database $(MIGRATION_DSN) -path $(MIGRATION_DIR) drop

.PHONY: migrate-reset
migrate-reset: ## Reset database to clean state (DANGEROUS - use with caution)
	@echo "WARNING: This will reset the database. Are you sure? (y/N)"
	@read -r confirm && [ "$$confirm" = "y" ] || exit 1
	migrate -database $(MIGRATION_DSN) -path $(MIGRATION_DIR) drop
	migrate -database $(MIGRATION_DSN) -path $(MIGRATION_DIR) up

.PHONY: migrate-validate
migrate-validate: ## Validate migration files
	migrate -database $(MIGRATION_DSN) -path $(MIGRATION_DIR) validate

.PHONY: migrate-list
migrate-list: ## List all migration files
	@ls -la $(MIGRATION_DIR)/

.PHONY: migrate-goto
migrate-goto: ## Go to specific migration version (usage: make migrate-goto VERSION=1)
	@if [ -z "$(VERSION)" ]; then echo "Usage: make migrate-goto VERSION=n"; exit 1; fi
	migrate -database $(MIGRATION_DSN) -path $(MIGRATION_DIR) goto $(VERSION)
