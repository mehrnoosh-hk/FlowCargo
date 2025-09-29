package tenant

import "github.com/google/uuid"

type CreateTenantRequest struct {
    Name string `json:"name"`
    Email string `json:"email"`
    Domain *string `json:"domain"`
}

type CreateTenantResponse struct {
    ID uuid.UUID `json:"id"`
}

type GetTenantByIDResponse struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Email string `json:"email"`
    Domain *string `json:"domain"`
}

type UpdateTenantRequest struct {
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	Domain *string `json:"domain"`
	IsActive *bool `json:"is_active"`
}
