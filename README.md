# FlowCargo

## Project structure guide

The structure below serves as a guide for the project. It does not reflect the current state of the project.

```
ecommerce/
├── cmd/
│   ├── api/                      # HTTP API server entry point
│   │   └── main.go
│   └── migration/                # Database migration tool
│       └── main.go
├── internal/                     # Private application code
│   ├── tenant/                   # User domain
│   │   ├── tenant.go             # Domain entities & business rules
│   │   ├── repository.go         # Repository interface
│   │   ├── service.go            # Business logic/use cases
│   │   ├── handler.go            # HTTP handlers
│   │   ├── queries/              # SQLC generated code
│   │   │   ├── db.go             # Generated database interface
│   │   │   ├── models.go         # Generated SQL models
│   │   │   ├── queries.sql.go    # Generated query methods
│   │   │   └── tenant.sql        # Raw SQL queries
│   │   └── repository/
│   │       └── postgres.go       # Repository implementation using SQLC
│   ├── shipment_tracking/        
│   │   ├── shipment_tracking.go
│   │   ├── shipment_tracking_item.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   └── repository/
│   │       └── postgres.go
│   ├── notification/          # Notification domain
│   │   ├── notification.go
│   │   ├── sender.go          # Sender interface
│   │   ├── service.go
│   │   └── senders/
│   │       ├── email.go
│   │       └── sms.go
│   └── shared/                # Shared kernel (minimal)
│       ├── events/            # Domain events
│       │   └── event.go
│       ├── errors/            # Common error types
│       │   └── errors.go
│       └── config/            # Configuration
│           └── config.go
├── pkg/                       # Public/reusable packages
│   ├── logger/
│   │   └── logger.go
│   ├── database/
│   │   └── postgres.go
│   └── httputil/
│       └── middleware.go
├── api/                       # API definitions (OpenAPI, proto, etc.)
│   └── openapi.yaml
├── migrations/                # Database migrations
│   ├── 001_create_users.sql
│   ├── 002_create_products.sql
│   └── 003_create_orders.sql
├── sqlc/                      # SQLC configuration
│   └── sqlc.yaml
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
├── scripts/
│   └── setup.sh
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

Dependency direction: 
```
HTTP Handlers → Services → Repositories → Database
     ↓             ↓           ↓
Domain Models ← Domain Models ← Domain Models
```

## Actual implementation

```
FlowCargo/
├── cmd/
│   └── api/
│       └── main.go               # HTTP API server entry point
├── internal/                     # Private application code
│   ├── shared/
│   │   └── config/
│   │       └── config.go         # Configuration management
│   └── tenant/                   # Tenant management domain
│       ├── tenant.go             # Domain entities & business rules
│       ├── repository.go         # Repository interface
│       ├── service.go            # Business logic/use cases
│       ├── handler.go            # HTTP handlers
│       └── db/
│           └── tenant.sql        # Raw SQL queries for SQLC
├── makefiles/                    # Makefile includes
│   ├── build.mk                  # Build targets
│   ├── lint.mk                   # Linting targets
│   ├── migration.mk              # Database migration targets
│   ├── test.mk                   # Testing targets
│   ├── tools.mk                  # Tool installation targets
│   └── variables.mk              # Makefile variables
├── migrations/                   # Database migrations
│   ├── 001_initial_setup.up.sql
│   ├── 001_initial_setup.down.sql
│   ├── 002_create_tenants_table.up.sql
│   └── 002_create_tenants_table.down.sql
├── sqlc/                         # SQLC configuration
│   └── sqlc.yml
└── docker/                       # Docker configuration
│    └── dev/
│        └── postgres/
│            ├── docker-compose.yml
│            ├── Dockerfile
│            └── postgresql.conf
├── .env.dev                      # Development environment variables
├── .gitattributes
├── .gitignore
├── go.mod                        # Go module definition
├── Makefile                      # Build automation
├── README.md                     # This file
```

## Database Migrations

This project uses [golang-migrate](https://github.com/golang-migrate/migrate) for database schema management. Migrations are stored in the `migrations/` directory and follow a sequential naming pattern.

### Migration Structure

Each migration consists of two files:
- `{number}_{name}.up.sql` - Contains the forward migration (applies changes)
- `{number}_{name}.down.sql` - Contains the rollback migration (reverts changes)

### Available Migrations

#### 001_initial_setup
Sets up the foundational database infrastructure:
- Enables PostgreSQL extensions (`uuid-ossp`, `citext`)
- Creates utility functions:
  - `set_updated_at()` - Automatically updates timestamp fields
  - `current_tenant_id()` - Retrieves tenant ID from session context
  - `is_app_admin()` - Checks if current user is an app admin
- Creates the `app_admin` role with appropriate privileges

### Migration Commands

The project provides several Make targets for migration management:

#### Basic Operations
```bash
# Apply all pending migrations
make migrate-up

# Rollback the last migration
make migrate-down

# Show current migration status
make migrate-status

# List all migration files
make migrate-list
```

#### Advanced Operations
```bash
# Create a new migration
make create-migration NAME=add_user_table

# Migrate up by specific number of steps
make migrate-up-steps STEPS=2

# Migrate down by specific number of steps
make migrate-down-steps STEPS=2

# Go to specific migration version
make migrate-goto VERSION=1

# Force set migration version (use with caution)
make migrate-force VERSION=1
```

#### Validation and Maintenance
```bash
# Validate migration files
make migrate-validate

# Reset database to clean state (DANGEROUS)
make migrate-reset

# Drop all tables (DANGEROUS)
make migrate-drop
```

### Configuration

Migration settings are configured in `makefiles/variables.mk`:
- `MIGRATION_DIR` - Directory containing migration files (default: `./migrations`)
- `MIGRATION_DSN` - Database connection string for migrations
- `MIGRATION_EXT` - Migration file extension (default: `sql`)
- `MIGRATION_DIGITS` - Number of digits in migration sequence (default: `3`)

### Environment Setup

Before running migrations, ensure you have:
1. PostgreSQL database running
2. Database connection string configured via `MIGRATION_DSN` environment variable
3. `golang-migrate` tool installed (automatically handled by the Makefile)

### Best Practices

1. **Always create both up and down migrations** - Every change should be reversible
2. **Test migrations thoroughly** - Test both forward and rollback operations
3. **Use transactions** - Wrap complex migrations in transactions when possible
4. **Backup before major changes** - Always backup production data before applying migrations
5. **Sequential numbering** - Use sequential numbers to maintain migration order
6. **Descriptive names** - Use clear, descriptive names for migration files