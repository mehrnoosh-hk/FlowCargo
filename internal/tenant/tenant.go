package tenant

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Tenant represents a single tenant in the system.
type Tenant struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Schema    string    `json:"schema"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository defines the interface for tenant persistence.
type Repository interface {
	CreateTenant(ctx context.Context, tenant *Tenant) error
	GetTenantByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
}

// Service defines the interface for the tenant service.
type Service interface {
	CreateTenant(ctx context.Context, name string) (*Tenant, error)
	GetTenantByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
}