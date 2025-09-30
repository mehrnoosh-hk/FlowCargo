package tenant

import (
	"context"
)

// type TenantFixture interface {}

type tenantFixture struct {
	repo   TenantRepository
	params CreateTenantParams
}

func NewTenantFixture(repo TenantRepository) tenantFixture {
	return tenantFixture{
		repo:   repo,
		params: CreateTenantParams{},
	}
}

func (f *tenantFixture) Tenant() *tenantFixture {
	f.params = CreateTenantParams{
		Name:   "john",
		Email:  "john@example.com",
		Domain: nil,
	}
	return f
}

func (f *tenantFixture) WithName(name string) *tenantFixture {
	f.params.Name = name
	return f
}

func (f *tenantFixture) WithDomain(domain string) *tenantFixture {
	f.params.Domain = &domain
	return f
}

func (f *tenantFixture) WithEmail(email string) *tenantFixture {
	f.params.Email = email
	return f
}

func (f *tenantFixture) Create(ctx context.Context) (*Tenant, error) {
	return f.repo.CreateTenant(ctx, f.params)
}
