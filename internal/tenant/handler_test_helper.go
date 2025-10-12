package tenant

import (
	"context"
	"time"

	"github.com/google/uuid"

	ru "flowcargo/internal/shared/restutils"
)

type mockService struct{}

func newmockService() TenantService {
	return &mockService{}
}

func (ms *mockService) CreateTenant(ctx context.Context, ctp CreateTenantParams) (*Tenant, error) {
	if ctp.Name == "New Tenant" {
		return &Tenant{
			ID:        uuid.New(),
			Name:      ctp.Name,
			Email:     ctp.Email,
			Domain:    ctp.Domain,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil
	}

	if ctp.Name == "Already Exists" {
		return nil, ru.ErrResourceConflict
	}
	return nil, nil
}

func (ms *mockService) GetTenantByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	if id.String() == "11111111-1111-1111-1111-111111111111" {
		return &Tenant{
			ID:        id,
			Name:      "Existing Tenant",
			Email:     "existing@example.com",
			Domain:    nil,
			IsActive:  true,
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-1 * time.Hour),
		}, nil
	}

	if id.String() == "22222222-2222-2222-2222-222222222222" {
		return nil, ru.ErrResourceNotFound
	}
	return nil, nil
}

func (ms *mockService) UpdateTenant(ctx context.Context, id uuid.UUID, params UpdateTenantParams) (*Tenant, error) {
	if id.String() == "11111111-1111-1111-1111-111111111111" {
		return &Tenant{
			ID:        id,
			Name:      *params.Name,
			Email:     *params.Email,
			Domain:    params.Domain,
			IsActive:  true,
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now(),
		}, nil
	}

	if id.String() == "22222222-2222-2222-2222-222222222222" {
		return nil, ru.ErrResourceNotFound
	}
	return nil, nil
}

func (ms *mockService) DeleteTenant(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	if id.String() == "11111111-1111-1111-1111-111111111111" {
		return &Tenant{
			ID:        id,
			Name:      "Deleted Tenant",
			Email:     "deleted@example.com",
			Domain:    nil,
			IsActive:  false,
			CreatedAt: time.Now().Add(-72 * time.Hour),
			UpdatedAt: time.Now(),
		}, nil
	}

	if id.String() == "22222222-2222-2222-2222-222222222222" {
		return nil, ru.ErrResourceNotFound
	}
	return nil, nil
}
