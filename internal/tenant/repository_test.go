package tenant

import (
	"context"
	"os"
	"testing"

	testutils "flowcargo/db/testutils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// TestMain sets up and tears down the test database
func TestMain(m *testing.M) {

	manager := testutils.GetDBManager()

	// Run tests
	code := m.Run()

	// Clean shutdown
	if manager != nil {
		manager.Close()
	}
	os.Exit(code)
}

func TestTenantRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("create tenant with correct data without domain", func(t *testing.T) {
		helper := NewTenantTestHelper(ctx, t, nil)
		fixture := NewTenantFixture(helper.Repo)
		tenant, err := fixture.Tenant().WithName("John Doe").WithEmail("john@test.com").Create(ctx)
		require.NoError(t, err)
		require.NotEmpty(t, tenant.ID)
	})

	t.Run("Create a new tenanr repository instance with pool", func(t *testing.T) {
		helper := NewTenantTestHelper(ctx, t, nil)
		require.NotNil(t, helper.Pool)
		require.NotNil(t, helper.Logger)
		repo := NewTenantRepository(helper.Pool, helper.Logger)
		require.NotNil(t, repo)
	})

	t.Run("create tenant with correct data with domain", func(t *testing.T) {
		helper := NewTenantTestHelper(ctx, t, nil)
		domain := "test.com"
		tenant, err := helper.Repo.CreateTenant(ctx, CreateTenantParams{
			Name:   "John",
			Email:  "John_Doe_test_2@test.com",
			Domain: &domain,
		})
		require.NoError(t, err)
		require.NotEmpty(t, tenant.ID)
	})

	t.Run("Can not create two tenants with same email", func(t *testing.T) {
		helper := NewTenantTestHelper(ctx, t, nil)
		email := "John_Doe_test_3@test.com"
		_, err := helper.Repo.CreateTenant(ctx, CreateTenantParams{
			Name:   "John",
			Email:  email,
			Domain: nil,
		})
		require.NoError(t, err)

		_, err = helper.Repo.CreateTenant(ctx, CreateTenantParams{
			Name:   "John",
			Email:  email,
			Domain: nil,
		})
		require.Error(t, err)
		t.Log(err)
	})

	t.Run("Update an existing tenant", func(t *testing.T) {
		helper := NewTenantTestHelper(ctx, t, nil)
		email := "John_Doe_test_4@test.com"
		tenant, err := helper.Repo.CreateTenant(ctx, CreateTenantParams{
			Name:   "John",
			Email:  email,
			Domain: nil,
		})
		require.NoError(t, err)
		require.NotEmpty(t, tenant.ID)

		name := "John Doe"
		domain := "test.com"
		updatedTenant, err := helper.Repo.UpdateTenant(ctx, UpdateTenantParams{
			ID:     tenant.ID,
			Name:   &name,
			Email:  &email,
			Domain: &domain,
		})
		if err != nil {
			t.Log(err)
		}
		require.NoError(t, err)
		require.Equal(t, tenant.ID, updatedTenant.ID)
		require.Equal(t, "John Doe", updatedTenant.Name)
		require.Equal(t, email, updatedTenant.Email)
		require.Equal(t, updatedTenant.Domain, &domain)
	})

	t.Run("Update only one column of an existing tenant", func(t *testing.T) {
		helper := NewTenantTestHelper(ctx, t, nil)
		email := "John_Doe_test_4@test.com"
		tenant, err := helper.Repo.CreateTenant(ctx, CreateTenantParams{
			Name:   "John",
			Email:  email,
			Domain: nil,
		})
		require.NoError(t, err)
		require.NotEmpty(t, tenant.ID)

		name := "John Doe"
		updatedTenant, err := helper.Repo.UpdateTenant(ctx, UpdateTenantParams{
			ID:     tenant.ID,
			Name:   &name,
			Email:  nil,
			Domain: nil,
		})
		require.NoError(t, err)
		require.Equal(t, tenant.ID, updatedTenant.ID)
		require.Equal(t, "John Doe", updatedTenant.Name)
		require.Nil(t, updatedTenant.Domain)
	})

	t.Run("Can not update a tenant which is not exist", func(t *testing.T) {
		helper := NewTenantTestHelper(ctx, t, nil)
		name := "John Doe"
		updatedTenant, err := helper.Repo.UpdateTenant(ctx, UpdateTenantParams{
			ID:   uuid.New(),
			Name: &name,
		})
		require.Error(t, err)
		require.Empty(t, updatedTenant, "The tenant object should be an empty struct")
	})

	t.Run("Get a tenant by ID", func(t *testing.T) {
		helper := NewTenantTestHelper(ctx, t, nil)
		email := "John_Doe_test_4@test.com"
		tenant, err := helper.Repo.CreateTenant(ctx, CreateTenantParams{
			Name:   "John",
			Email:  email,
			Domain: nil,
		})
		require.NoError(t, err)
		require.NotEmpty(t, tenant.ID)

		tenantByID, err := helper.Repo.GetTenantByID(ctx, tenant.ID)
		require.NoError(t, err)
		require.Equal(t, tenant.ID, tenantByID.ID)
		require.Equal(t, tenant.Name, tenantByID.Name)
		require.Equal(t, tenant.Email, tenantByID.Email)
		require.Equal(t, tenant.Domain, tenantByID.Domain)
	})

	t.Run("Get tenant by ID for non existing tenant", func(t *testing.T) {
		helper := NewTenantTestHelper(ctx, t, nil)
		ID := uuid.New()
		tenant, err := helper.Repo.GetTenantByID(ctx, ID)
		require.Error(t, err)
		require.Nil(t, tenant)
	})
}
