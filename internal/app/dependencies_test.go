package app

import (
	"errors"
	"flowcargo/internal/shared/logger"
	"testing"
	"context"

	"github.com/stretchr/testify/require"
)

func TestDependenciesWiring(t *testing.T) {
	t.Run("wireDeps fails if wireLogger fails", func(t *testing.T) {
		testCtx := context.Background()
		original := wireLogger
		defer func() {
			wireLogger = original
		}()
		wireLogger = func(isDevelopment bool, level logger.LogLevel) (logger.Logger, error) {
			return logger.Logger{}, errors.New("mock error")
		}
		_, err := wireDeps(testCtx, &Database{}, true, logger.Debug)
		require.Error(t, err)
	})
	
	t.Run("wireDeps succeeds when wireLogger succeeds", func(t *testing.T) {
		testCtx := context.Background()
		original := wireLogger
		defer func() {
			wireLogger = original
		}()
		wireLogger = func(isDevelopment bool, level logger.LogLevel) (logger.Logger, error) {
			return logger.Logger{}, nil
		}
		_, err := wireDeps(testCtx,  &Database{}, true, logger.Debug)
		require.NoError(t, err)
	})
}

func TestWireLogger(t *testing.T) {
	testCases := []struct {
		name     string
		isDev    bool
		level    logger.LogLevel
		expected error
	}{
		{"DevelopmentMode", true, logger.Debug, nil},
		{"ProductionMode", false, logger.Info, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger, err := wireLogger(tc.isDev, tc.level)
			if tc.expected != nil {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, logger)
			}
		})
	}
}
