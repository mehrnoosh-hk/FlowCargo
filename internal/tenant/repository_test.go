package tenant

import (
	"context"
	"os"
	"testing"

	testutils "flowcargo/db"
	"flowcargo/internal/shared/logger"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

type TestHelper struct {
	TX     pgx.Tx
	Repo   TenantRepository
	Logger logger.Logger
}

var testDBManager *testutils.TestDBManager

// TestMain sets up and tears down the test database
func TestMain(m *testing.M) {
	// Create test database manager and assign it to global variable
	DBManager := testutils.NewTestDBManager(testutils.GetTestDSN())
	testDBManager = DBManager

	// Initialize database
	if err := testDBManager.Initialize(); err != nil {
		panic("Failed to initialize test database: " + err.Error())
	}

	// Run tests
	code := m.Run()

	// Clean shutdown
	if testDBManager != nil {
		testDBManager.Close()
	}
	os.Exit(code)
}

// repository_test.go - TESTING
func NewTenantTestHelper(ctx context.Context, t *testing.T) *TestHelper {
	t.Helper()
	db := testDBManager.GetDB(t) // Get the pool
	tx, err := db.Begin(ctx)
	require.NoError(t, err)
	logger, err := logger.NewLogger(false, logger.Debug)
	require.NoError(t, err)
	repo := NewTenantRepository(tx, logger) // Use transaction, not pool
	t.Cleanup(func(){
		err := tx.Rollback(ctx)
		require.NoError(t, err)
	})
	return &TestHelper{
		Repo:   repo,
		Logger: logger,
	}
}

func TestTenantRepository(t *testing.T) {
	ctx := context.Background()

	t.Run("create tenant with correct data without domain", func(t *testing.T) {
		helper := NewTenantTestHelper(ctx, t)
		tenant, err := helper.Repo.CreateTenant(ctx, CreateTenantParams{
			Name:   "John",
			Email:  "John_Doe_test_1@test.com",
			Domain: nil,
		})
		require.NoError(t, err)
		require.NotEmpty(t, tenant.ID)
	})

	t.Run("create tenant with correct data with domain", func(t *testing.T) {
		helper := NewTenantTestHelper(ctx, t)
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
		helper := NewTenantTestHelper(ctx, t)
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
		// require.ErrorIs(t, err, ErrDuplicateEmail)
	})
	
	t.Run("Update an existing tenant", func(t *testing.T) {
		helper := NewTenantTestHelper(ctx, t)
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
	
	t.Run("Update only one column of an existing tenanat", func(t *testing.T) {
		helper := NewTenantTestHelper(ctx, t)
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

}
