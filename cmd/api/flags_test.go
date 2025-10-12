package main

import (
	"flag"
	"os"
	"testing"

	"flowcargo/internal/shared/config"
)

func TestParseFlags(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()
	validPath := tmpFile.Name()
	notExist := "./not_exist.yaml"
	tests := []struct {
		name        string
		args        []string
		expected    Flags
		expectError bool
	}{
		{
			name: "valid flags",
			args: []string{"-env=development", "-config=" + validPath},
			expected: Flags{
				Environment: config.Development,
				ConfigPath:  &validPath,
			},
			expectError: false,
		},
		{
			name: "invalid environment",
			args: []string{"-env=invalid", "-config=" + validPath},
			expected: Flags{
				Environment: config.Development,
				ConfigPath:  &validPath,
			},
			expectError: true,
		},
		{
			name: "missing config path",
			args: []string{"-env=test"},
			expected: Flags{
				Environment: config.Test,
				ConfigPath:  nil,
			},
			expectError: false,
		},
		{
			name:        "config path not exist",
			args:        []string{"-env=development", "-config=" + notExist},
			expected:    Flags{},
			expectError: true,
		},
		{
			name:        "config path with stat error",
			args:        []string{"-env=development", "-config=/root/no_permission.yaml"},
			expected:    Flags{},
			expectError: true,
		},
		{
			name:        "missing environment",
			args:        []string{"-config=" + validPath},
			expected:    Flags{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.args)
			fs := flag.NewFlagSet("test", flag.ContinueOnError)

			parsedFlags, err := ParseFlagsFromSet(fs, tt.args)
			if err != nil && !tt.expectError {
				t.Errorf("unexpected error: %v", err)
			}

			if err == nil && tt.expectError {
				t.Errorf("expected error but got none")
			}

			if tt.expectError {
				return
			}

			if parsedFlags.Environment != tt.expected.Environment {
				t.Errorf("unexpected environment: got %v, want %v", parsedFlags.Environment, tt.expected.Environment)
			}

			if parsedFlags.ConfigPath != nil && *parsedFlags.ConfigPath != *tt.expected.ConfigPath {
				t.Errorf("unexpected config path: got %v, want %v", parsedFlags.ConfigPath, tt.expected.ConfigPath)
			}
		})
	}
}
