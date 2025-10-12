package tenant

import (
	"context"
	"testing"

	testutils "flowcargo/db/testutils"
	"flowcargo/internal/shared/logger"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// TenantTestHelper provides everything needed for testing tenant-related functionality.
type TenantTestHelper struct {
	Pool   *pgxpool.Pool
	Tx     pgx.Tx
	Repo   TenantRepository // Interface
	Logger logger.Logger
}

// HelperConfig helps to decouple from specific DBManager implementation.
type HelperConfig struct {
	DBManager testutils.TestDB
	LogLevel  logger.LogLevel
	Debug     bool
}

func DefaultConfig() *HelperConfig {
	return &HelperConfig{
		DBManager: testutils.GetDBManager(),
		LogLevel:  logger.Debug,
		Debug:     true,
	}
}

func NewTenantTestHelper(ctx context.Context, t *testing.T, config *HelperConfig) *TenantTestHelper {
	t.Helper()
	if config == nil {
		config = DefaultConfig()
	}
	manager := config.DBManager
	tx := manager.BeginTx(ctx, t)
	logger, err := logger.NewLogger(false, config.LogLevel)
	require.NoError(t, err)
	repo := NewTenantRepository(tx, logger) // Uses transaction, not pool
	t.Cleanup(func() {
		err := tx.Rollback(ctx)
		require.NoError(t, err)
	})
	return &TenantTestHelper{
		Pool:   manager.GetPool(),
		Tx:     tx,
		Repo:   repo,
		Logger: logger,
	}
}
