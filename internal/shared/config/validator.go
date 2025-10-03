package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ConfigValidator defines the interface for validating configuration.
type ConfigValidator interface {
	Validate(config Config) error
	ValidateServer(server Server) error
	ValidateDatabase(database Database) error
	ValidateCORS(cors CORS) error
}

// configValidator is the concrete implementation of ConfigValidator.
type configValidator struct {
	validate *validator.Validate
}

// NewConfigValidator creates a new instance of ConfigValidator.
func NewConfigValidator() ConfigValidator {
	return &configValidator{
		validate: validator.New(),
	}
}

// Validate validates the entire Config struct.
func (v *configValidator) Validate(config Config) error {
	if err := v.validate.Struct(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	return nil
}

// ValidateServer validates the Server configuration.
func (v *configValidator) ValidateServer(server Server) error {
	if err := v.validate.Struct(server); err != nil {
		return fmt.Errorf("server config validation failed: %w", err)
	}
	return nil
}

// ValidateDatabase validates the Database configuration.
func (v *configValidator) ValidateDatabase(database Database) error {
	if err := v.validate.Struct(database); err != nil {
		return fmt.Errorf("database config validation failed: %w", err)
	}
	return nil
}

// ValidateCORS validates the CORS configuration.
func (v *configValidator) ValidateCORS(cors CORS) error {
	if err := v.validate.Struct(cors); err != nil {
		return fmt.Errorf("CORS config validation failed: %w", err)
	}
	return nil
}
