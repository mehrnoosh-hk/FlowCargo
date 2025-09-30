package tenant

import (
	"time"

	"github.com/google/uuid"
)

// Tenant represents a single tenant in the system.
type Tenant struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Domain    *string   `json:"domain"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateTenantParams struct {
	Name   string  `json:"name" validate:"required,min=2,max=100"`
	Email  string  `json:"email" validate:"required,email"`
	Domain *string `json:"domain" validate:"omitempty,url"`
}

type UpdateTenantParams struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Name     *string   `json:"name" validate:"omitempty,min=2,max=100"`
	Email    *string   `json:"email" validate:"omitempty,email"`
	Domain   *string   `json:"domain" validate:"omitempty,url"`
	IsActive *bool     `json:"is_active" validate:"omitempty,boolean"`
}
