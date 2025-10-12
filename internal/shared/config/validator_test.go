package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// validConfig returns a fully valid configuration for testing.
// Each test can clone and modify only the parts they need to test.
func validConfig() Config {
	return Config{
		Env: Production,
		Server: Server{
			Host: "localhost",
			Port: "8080",
		},
		Database: Database{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			Name:     "testdb",
			SSLMode:  "disable",
		},
		Logger: Logger{
			Level:  "info",
			Source: true,
		},
		CORSConfig: CORS{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST"},
			AllowHeaders:     []string{"Content-Type"},
			AllowCredentials: true,
			ExposeHeaders:    []string{},
			MaxAge:           3600,
		},
	}
}

func TestConfigValidator_Validate(t *testing.T) {
	validator := NewConfigValidator()

	tests := []struct {
		name        string
		modifyFunc  func(*Config)
		wantErr     bool
		errContains string
	}{
		{
			name:       "valid configuration",
			modifyFunc: func(c *Config) {}, // No modifications
			wantErr:    false,
		},
		{
			name: "invalid server - missing host",
			modifyFunc: func(c *Config) {
				c.Server.Host = ""
			},
			wantErr:     true,
			errContains: "Host",
		},
		{
			name: "invalid server - missing port",
			modifyFunc: func(c *Config) {
				c.Server.Port = ""
			},
			wantErr:     true,
			errContains: "Port",
		},
		{
			name: "invalid server - non-numeric port",
			modifyFunc: func(c *Config) {
				c.Server.Port = "not-a-number"
			},
			wantErr:     true,
			errContains: "Port",
		},
		{
			name: "invalid database - missing host",
			modifyFunc: func(c *Config) {
				c.Database.Host = ""
			},
			wantErr:     true,
			errContains: "Host",
		},
		{
			name: "invalid database - missing user",
			modifyFunc: func(c *Config) {
				c.Database.User = ""
			},
			wantErr:     true,
			errContains: "User",
		},
		{
			name: "invalid database - missing name",
			modifyFunc: func(c *Config) {
				c.Database.Name = ""
			},
			wantErr:     true,
			errContains: "Name",
		},
		{
			name: "invalid database - missing password",
			modifyFunc: func(c *Config) {
				c.Database.Password = ""
			},
			wantErr:     true,
			errContains: "Password",
		},
		{
			name: "invalid database - multiple missing fields",
			modifyFunc: func(c *Config) {
				c.Database.Host = ""
				c.Database.User = ""
				c.Database.Name = ""
			},
			wantErr: true,
		},
		{
			name: "invalid CORS - missing allow origins",
			modifyFunc: func(c *Config) {
				c.CORSConfig.AllowOrigins = []string{}
			},
			wantErr:     true,
			errContains: "AllowOrigins",
		},
		{
			name: "invalid CORS - missing allow methods",
			modifyFunc: func(c *Config) {
				c.CORSConfig.AllowMethods = []string{}
			},
			wantErr:     true,
			errContains: "AllowMethods",
		},
		{
			name: "invalid CORS - invalid origin URI",
			modifyFunc: func(c *Config) {
				c.CORSConfig.AllowOrigins = []string{"not-a-valid-uri"}
			},
			wantErr:     true,
			errContains: "AllowOrigins",
		},
		{
			name: "invalid CORS - negative MaxAge",
			modifyFunc: func(c *Config) {
				c.CORSConfig.MaxAge = -100
			},
			wantErr:     true,
			errContains: "MaxAge",
		},
		{
			name: "valid configuration with wildcard CORS origin",
			modifyFunc: func(c *Config) {
				c.CORSConfig.AllowOrigins = []string{"*"}
			},
			wantErr: false,
		},
		{
			name: "valid configuration with minimal CORS",
			modifyFunc: func(c *Config) {
				c.CORSConfig = CORS{
					AllowOrigins:     []string{"http://example.com"},
					AllowMethods:     []string{"GET"},
					AllowHeaders:     []string{},
					AllowCredentials: true,
					ExposeHeaders:    []string{},
					MaxAge:           0,
				}
			},
			wantErr: false,
		},
		{
			name: "invalid logger - missing level",
			modifyFunc: func(c *Config) {
				c.Logger.Level = ""
			},
			wantErr:     true,
			errContains: "Level",
		},
		{
			name: "valid configuration with different environment",
			modifyFunc: func(c *Config) {
				c.Env = Development
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh valid config for each test
			cfg := validConfig()

			// Apply test-specific modifications
			tt.modifyFunc(&cfg)

			// Validate
			err := validator.Validate(cfg)

			if tt.wantErr {
				require.Error(t, err, "expected an error but got none")
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains, "error message should contain expected string")
				}
			} else {
				assert.NoError(t, err, "expected no error but got: %v", err)
			}
		})
	}
}
