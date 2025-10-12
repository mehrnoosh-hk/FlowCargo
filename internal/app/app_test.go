package app

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"flowcargo/internal/shared/config"
)

func TestNewApp(t *testing.T) {
	testCases := []struct {
		name        string
		wireFunc    func() Wire
		expectError bool
	}{
		{
			name:        "Creating new app fails if it can not wire dependencies",
			wireFunc:    newFailedWire,
			expectError: true,
		},
		{
			name:        "Creates the application, if all wire functions succeed",
			wireFunc:    newSucceedWire,
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testCtx := context.Background()
			configFile := "path"
			env := config.Test

			wire := tc.wireFunc()

			app, err := wire.Up(testCtx, env, &configFile)
			if tc.expectError {
				require.Equal(t, App{}, app)
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, app)
			}
		})
	}
}

func TestCreateAndRun(t *testing.T) {

	tests := []struct {
		name     string // description of this test case
		wireFunc func() Wire
		// Named input parameters for target function.
		environment config.Environment
		configPath  *string
		wantErr     bool
	}{
		{
			name:        "CreateAndRun fails if it can not wire dependencies",
			wireFunc:    newFailedWire,
			environment: config.Development,
			configPath:  nil,
			wantErr:     true,
		},
		{
			name:        "CreateAndRun succeeds if it can wire dependencies",
			wireFunc:    newSucceedWire,
			environment: config.Development,
			configPath:  nil,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Start goroutine to cancel context BEFORE calling CreateAndRun
			go func() {
				time.Sleep(100 * time.Millisecond)
				cancel()
			}()

			gotErr := CreateAndRun(
				ctx,
				tt.environment,
				tt.configPath,
				tt.wireFunc(),
			)

			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CreateAndRun() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CreateAndRun() succeeded unexpectedly")
			}
			t.Log(tt.name + " Passed")
		})
	}
}
