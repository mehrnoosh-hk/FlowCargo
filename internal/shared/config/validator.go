package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// ConfigValidator defines the interface for validating configuration.
type ConfigValidator interface {
	Validate(config Config) error
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

// Validate validates the entire Config struct including all nested structs.
// The go-playground/validator package automatically validates nested structs
// recursively when they have validation tags defined.
func (v *configValidator) Validate(config Config) error {
	if err := v.validate.Struct(config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	return nil
}
