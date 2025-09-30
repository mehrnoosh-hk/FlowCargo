// Business logic and use cases for tenant management
package tenant

import (
	"context"
	"flowcargo/internal/shared/logger"

	"github.com/google/uuid"
)

type TenantService interface {
	CreateTenant(ctx context.Context, params CreateTenantParams) (*Tenant, error)
	GetTenantByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
	UpdateTenant(ctx context.Context, id uuid.UUID, params UpdateTenantParams) (*Tenant, error)
	SoftDeleteTenant(ctx context.Context, id uuid.UUID) (*Tenant, error)
}

type tenantService struct {
	l          logger.Logger
	tenantRepo TenantRepository
}

func NewTenantService(tenantRepo TenantRepository, l logger.Logger) TenantService {
	return &tenantService{
		l:          l,
		tenantRepo: tenantRepo,
	}
}

func (ts *tenantService) CreateTenant(ctx context.Context, params CreateTenantParams) (*Tenant, error) {
	return ts.tenantRepo.CreateTenant(ctx, params)
}

func (ts *tenantService) GetTenantByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	return ts.tenantRepo.GetTenantByID(ctx, id)
}

func (ts *tenantService) UpdateTenant(ctx context.Context, id uuid.UUID, params UpdateTenantParams) (*Tenant, error) {
	params.ID = id
	return ts.tenantRepo.UpdateTenant(ctx, params)
}

func (ts *tenantService) SoftDeleteTenant(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	isActive := false
	params := UpdateTenantParams{
		ID:       id,
		IsActive: &isActive,
	}
	ts.l.Info("Deleting tenant params: ", "id = ", id)
	return ts.tenantRepo.UpdateTenant(ctx, params)
}
