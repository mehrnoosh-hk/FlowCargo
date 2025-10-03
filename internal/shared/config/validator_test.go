package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigValidator_Validate(t *testing.T) {
	validator := NewConfigValidator()

	t.Run("Valid configuration", func(t *testing.T) {
		cfg := Config{
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
			CORS: CORS{
				AllowOrigins:     []string{"http://localhost:3000"},
				AllowMethods:     []string{"GET", "POST"},
				AllowHeaders:     []string{"Content-Type"},
				AllowCredentials: true,
				ExposeHeaders:    []string{},
				MaxAge:           3600,
			},
		}

		err := validator.Validate(cfg)
		assert.NoError(t, err)
	})

	t.Run("Invalid server - missing host", func(t *testing.T) {
		cfg := Config{
			Env: Production,
			Server: Server{
				Host: "", // Missing required field
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
		}

		err := validator.Validate(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config validation failed")
	})

	t.Run("Invalid database - missing required fields", func(t *testing.T) {
		cfg := Config{
			Env: Production,
			Server: Server{
				Host: "localhost",
				Port: "8080",
			},
			Database: Database{
				Host:     "", // Missing
				Port:     "5432",
				User:     "", // Missing
				Password: "password",
				Name:     "", // Missing
				SSLMode:  "disable",
			},
		}

		err := validator.Validate(cfg)
		assert.Error(t, err)
	})

	t.Run("Invalid port - not numeric", func(t *testing.T) {
		cfg := Config{
			Env: Production,
			Server: Server{
				Host: "localhost",
				Port: "abc", // Should be numeric
			},
			Database: Database{
				Host:     "localhost",
				Port:     "5432",
				User:     "postgres",
				Password: "password",
				Name:     "testdb",
				SSLMode:  "disable",
			},
		}

		err := validator.Validate(cfg)
		assert.Error(t, err)
	})
}

func TestConfigValidator_ValidateServer(t *testing.T) {
	validator := NewConfigValidator()

	t.Run("Valid server config", func(t *testing.T) {
		server := Server{
			Host: "localhost",
			Port: "8080",
		}

		err := validator.ValidateServer(server)
		assert.NoError(t, err)
	})

	t.Run("Invalid server - missing host", func(t *testing.T) {
		server := Server{
			Host: "",
			Port: "8080",
		}

		err := validator.ValidateServer(server)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "server config validation failed")
	})

	t.Run("Invalid server - invalid port", func(t *testing.T) {
		server := Server{
			Host: "localhost",
			Port: "not-a-number",
		}

		err := validator.ValidateServer(server)
		assert.Error(t, err)
	})
}

func TestConfigValidator_ValidateDatabase(t *testing.T) {
	validator := NewConfigValidator()

	t.Run("Valid database config", func(t *testing.T) {
		db := Database{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			Name:     "testdb",
			SSLMode:  "disable",
		}

		err := validator.ValidateDatabase(db)
		assert.NoError(t, err)
	})

	t.Run("Invalid database - missing fields", func(t *testing.T) {
		db := Database{
			Host:     "",
			Port:     "5432",
			User:     "",
			Password: "",
			Name:     "",
			SSLMode:  "disable",
		}

		err := validator.ValidateDatabase(db)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database config validation failed")
	})
}

func TestConfigValidator_ValidateCORS(t *testing.T) {
	validator := NewConfigValidator()

	t.Run("Valid CORS config", func(t *testing.T) {
		cors := CORS{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST"},
			AllowHeaders:     []string{"Content-Type"},
			AllowCredentials: true,
			ExposeHeaders:    []string{},
			MaxAge:           3600,
		}

		err := validator.ValidateCORS(cors)
		assert.NoError(t, err)
	})

	t.Run("Invalid CORS - missing required fields", func(t *testing.T) {
		cors := CORS{
			AllowOrigins:     []string{},
			AllowMethods:     []string{},
			AllowHeaders:     []string{},
			AllowCredentials: false,
			ExposeHeaders:    []string{},
			MaxAge:           0,
		}

		err := validator.ValidateCORS(cors)
		assert.Error(t, err)
	})

	t.Run("Invalid CORS - invalid origin URI", func(t *testing.T) {
		cors := CORS{
			AllowOrigins:     []string{"not-a-valid-uri"},
			AllowMethods:     []string{"GET"},
			AllowHeaders:     []string{},
			AllowCredentials: false,
			ExposeHeaders:    []string{},
			MaxAge:           3600,
		}

		err := validator.ValidateCORS(cors)
		assert.Error(t, err)
	})

	t.Run("Invalid CORS - negative MaxAge", func(t *testing.T) {
		cors := CORS{
			AllowOrigins:     []string{"http://localhost:3000"},
			AllowMethods:     []string{"GET"},
			AllowHeaders:     []string{},
			AllowCredentials: false,
			ExposeHeaders:    []string{},
			MaxAge:           -100, // Should be >= 0
		}

		err := validator.ValidateCORS(cors)
		assert.Error(t, err)
	})
}
