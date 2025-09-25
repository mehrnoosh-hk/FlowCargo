package app

import (
	"errors"
	"flowcargo/internal/shared/logger"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDependenciesWiering(t *testing.T) {
	t.Run("WireDependenciesFailsWhenWireLoggerFails", func(t *testing.T) {
		original := wireLogger
		defer func() {
			wireLogger = original
		}()
		wireLogger = func(isDevelopment bool, level logger.LogLevel) (logger.Logger, error) {
			return logger.Logger{}, errors.New("mock error")
		}
		_, err := wireDependencies()
		require.Error(t, err)
	})
	
	t.Run("WireDependenciesSucceedWhenWireLoggerSucceeds", func(t *testing.T) {
		original := wireLogger
		defer func() {
			wireLogger = original
		}()
		wireLogger = func(isDevelopment bool, level logger.LogLevel) (logger.Logger, error) {
			return logger.Logger{}, nil
		}
		_, err := wireDependencies()
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
