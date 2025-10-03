package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Loader defines the interface for loading configuration.
type Loader interface {
	Load() (Config, error)
}

// --- File Loader ---

// FileLoader loads configuration from a file.
type FileLoader struct {
	path *string
}

// NewFileLoader creates a new file loader.
func NewFileLoader(path *string) *FileLoader {
	return &FileLoader{path: path}
}

// Load implements the Loader interface for FileLoader.
// It attempts to load from the specified file path, falling back to defaults if no path is provided.
func (l *FileLoader) Load() (Config, error) {
	var cfg Config

	// If no path is provided, unmarshal the defaults and return.
	if l.path == nil || *l.path == "" {
		setDefaults()
		fmt.Println("No config file path provided, using defaults.")
		if err := viper.Unmarshal(&cfg); err != nil {
			return Config{}, fmt.Errorf("failed to unmarshal default config: %w", err)
		}
		return cfg, nil
	}

	viper.SetConfigName(*l.path)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("could not read config file '%s': %w", *l.path, err)
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config from file: %w", err)
	}

	fmt.Println("Config loaded successfully from file:", *l.path)
	return cfg, nil
}

// --- Environment Loader ---

// EnvLoader loads configuration from environment variables.
type EnvLoader struct{}

// NewEnvLoader creates a new environment variable loader.
func NewEnvLoader() *EnvLoader {
	return &EnvLoader{}
}

// Load implements the Loader interface for EnvLoader.
// It configures viper to read from environment variables and validates them.
func (l *EnvLoader) Load() (Config, error) {
	var cfg Config

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config from environment variables: %w", err)
	}

	// Validate the configuration for production environment
	validator := NewConfigValidator()
	if err := validator.Validate(cfg); err != nil {
		return Config{}, fmt.Errorf("environment variables validation failed: %w", err)
	}

	fmt.Println("Config loaded successfully from environment variables.")
	return cfg, nil
}
