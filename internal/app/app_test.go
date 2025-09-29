package app

import (
	"context"
	"errors"
	"flowcargo/internal/shared/config"
	"flowcargo/internal/shared/logger"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppCreateAndRun(t *testing.T) {
	t.Run("Returns error if it can not wire dependencies", func(t *testing.T) {
		testCtx := context.Background()
		original := wireDeps
		defer func() {
			wireDeps = original
		}()
		
		wireDeps = func(ctx context.Context,db *Database, isDev bool, level logger.LogLevel) (Dependencies, error) {
			return Dependencies{}, errors.New("test error")
		}
		
		err := CreateAndRun(testCtx, "path")
		require.Error(t, err)		
	})
	
	t.Run("Creates the application, if all wire functions succeed", func(t *testing.T) {
		testCtx := context.Background()
		originalDeps := wireDeps
		originalDB := wireDB
		originalCfg := wireCfg
		originalSrv := wireSrv
		defer func() {
			wireDeps = originalDeps
			wireDB = originalDB
			wireCfg = originalCfg
			wireSrv = originalSrv
		}()
		
		wireDeps = func(ctx context.Context, db *Database, isDev bool, level logger.LogLevel) (Dependencies, error) {
			return Dependencies{}, nil
		}
		
		wireDB = func(ctx context.Context, dbURL string) (*Database, error) {
			return &Database{}, nil
		}
		
		wireCfg = func() (config.Config) {
			return config.Config{}
		}
				
		_, err := newApp(testCtx, wireCfg, wireDB, wireDeps, wireSrv)
		require.NoError(t, err)	
	})
}