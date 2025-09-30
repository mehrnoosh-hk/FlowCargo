package tenant

import (
	"context"
	"testing"

	testutils "flowcargo/db"
	"flowcargo/internal/shared/logger"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

// TenantTestHelper provides everything needed for testing tenant-related functionality.
type TenantTestHelper struct {
	TX     pgx.Tx
	Repo   TenantRepository // Interface
	Logger logger.Logger
}

func NewTenantTestHelper(ctx context.Context, t *testing.T) *TenantTestHelper {
	t.Helper()
	manager := testutils.GetDBManager()
	tx := manager.BeginTx(ctx, t)
	logger, err := logger.NewLogger(false, logger.Debug)
	require.NoError(t, err)
	repo := NewTenantRepository(tx, logger) // Uses transaction, not pool
	t.Cleanup(func(){
		err := tx.Rollback(ctx)
		require.NoError(t, err)
	})
	return &TenantTestHelper{
		Repo:   repo,
		Logger: logger,
	}
}