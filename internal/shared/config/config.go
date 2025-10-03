package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env      Environment `mapstructure:"env"`
	Server   Server      `mapstructure:"server"`
	Database Database    `mapstructure:"database"`
	Logger   Logger      `mapstructure:"logger"`
	CORS     CORS        `mapstructure:"cors"`
}

type Environment string

const (
	Development Environment = "development"
	Test        Environment = "test"
	Production  Environment = "production"
)

// New creates a new Config instance based on the environment.
// It receives the 'env' environment variable to determine which configuration loader to use.
// In Production: reads from env variables and returns error if validation fails.
// In Development/Test: tries to read config file, falls back to defaults if file loading fails.
func New(env Environment, path *string) (Config, error) {
	if env == Production {
		loader := NewEnvLoader()
		return loader.Load()
	} else {
		var cfg Config
		loader := NewFileLoader(path)
		cfg, err := loader.Load()
		if err != nil {
			// Fall back to defaults in development/test environments
			fmt.Println("Failed to load config file, using defaults:", err)
			setDefaults()
			if err := viper.Unmarshal(&cfg); err != nil {
				return Config{}, fmt.Errorf("failed to unmarshal default config: %w", err)
			}
		}
		return cfg, nil
	}
}

func setDefaults() {
	viper.SetDefault("env", "development")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.name", "flowcargo_dev")
	viper.SetDefault("logger.level", "INFO")
	viper.SetDefault("logger.source", true)
	viper.SetDefault("cors.allowOrigins", []string{"http://localhost:3000"})
	viper.SetDefault("cors.allowMethods", []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"})
	viper.SetDefault("cors.allowHeaders", []string{"Content-Type", "Authorization"})
	viper.SetDefault("cors.allowCredentials", true)
	viper.SetDefault("cors.exposeHeaders", []string{})
	viper.SetDefault("cors.maxAge", time.Duration(3600*time.Second))
}

func (c Config) IsDevelopment() bool {
	return (c.Env == Development)
}
