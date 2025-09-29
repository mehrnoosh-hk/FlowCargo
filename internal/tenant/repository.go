// Repository interface for tenant operations
package tenant

import (
	"context"

	db "flowcargo/db/sqlc"
	"flowcargo/internal/shared/logger"

	"github.com/google/uuid"
)

type TenantRepository interface {
	CreateTenant(ctx context.Context, p CreateTenantParams) (*Tenant, error)
	GetTenantByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
	UpdateTenant(ctx context.Context, p UpdateTenantParams) (*Tenant, error)
	DeleteTenant(ctx context.Context, id uuid.UUID) error
}

type tenantRepository struct {
	l  logger.Logger
	queries *db.Queries // This can accept both pool and tx
}

// NewTenantRepository accepts both pool and tx
func NewTenantRepository(conn db.DBTX, l logger.Logger) TenantRepository {
	return &tenantRepository{
		l:  l,
		queries: db.New(conn),
	}
}

func (t *tenantRepository) GetTenantByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	dbTenant, err := t.queries.GetTenantByID(ctx, db.UUIDToPgUUID(id))
	if err != nil {
		return nil, err
	}
	tenant := DatabaseTenantToDomain(dbTenant)
	return &tenant, nil
}

func (t *tenantRepository) CreateTenant(ctx context.Context, p CreateTenantParams) (*Tenant, error) {
	dbTenant, err := t.queries.CreateTenant(ctx, db.CreateTenantParams{
		Name:   p.Name,
		Email:  p.Email,
		Domain: db.StringToPgText(p.Domain),
	})
	if err != nil {
		return nil, err
	}
	tenant := DatabaseTenantToDomain(dbTenant)
	return &tenant, nil
}

func (t *tenantRepository) UpdateTenant(ctx context.Context, p UpdateTenantParams) (*Tenant, error) {
	dbParams := db.UpdateTenantParams{
		ID:       db.UUIDToPgUUID(p.ID),
		Name:     db.StringToPgText(p.Name),
		Email:    db.StringToPgText(p.Email),
		Domain:   db.StringToPgText(p.Domain),
		IsActive: db.BoolToPgBool(p.IsActive),
	}
	dbTenant, err := t.queries.UpdateTenant(ctx, dbParams)
	if err != nil {
		return nil, err
	}
	tenant := DatabaseTenantToDomain(dbTenant)
	return &tenant, nil
}

func (t *tenantRepository) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	return t.queries.DeleteTenant(ctx, db.UUIDToPgUUID(id))
}

func DatabaseTenantToDomain(t db.Tenant) Tenant {
	return Tenant{
		ID:        t.ID.Bytes,
		Name:      t.Name,
		Email:     t.Email,
		Domain:    db.PgTextToString(t.Domain),
		IsActive:  t.IsActive,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func DomainCreateParamToDBCreateParam(p CreateTenantParams) db.CreateTenantParams {
	return db.CreateTenantParams{}
}
