package tenant

import (
	"time"

	"github.com/google/uuid"
)

// Tenant represents a single tenant in the system.
// @Description Tenant account information with contact details and status
type Tenant struct {
	ID        uuid.UUID `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" example:"Acme Corporation"`
	Email     string    `json:"email" example:"contact@acme.com"`
	Domain    *string   `json:"domain" example:"https://acme.com"`
	IsActive  bool      `json:"is_active" example:"true"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-15T10:30:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-15T10:30:00Z"`
}

// CreateTenantParams represents the parameters required to create a new tenant.
// @Description Parameters for creating a new tenant account
type CreateTenantParams struct {
	Name   string  `json:"name" validate:"required,min=2,max=100" example:"Acme Corporation"`
	Email  string  `json:"email" validate:"required,email" example:"contact@acme.com"`
	Domain *string `json:"domain" validate:"omitempty,url" example:"https://acme.com"`
}

// UpdateTenantParams represents the parameters for updating an existing tenant.
// @Description Parameters for updating tenant information (all fields optional except ID)
type UpdateTenantParams struct {
	ID       uuid.UUID `json:"id" validate:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name     *string   `json:"name" validate:"omitempty,min=2,max=100" example:"Acme Corporation Updated"`
	Email    *string   `json:"email" validate:"omitempty,email" example:"newemail@acme.com"`
	Domain   *string   `json:"domain" validate:"omitempty,url" example:"https://newdomain.acme.com"`
	IsActive *bool     `json:"is_active" validate:"omitempty,boolean" example:"false"`
}
