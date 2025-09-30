package config

import (
	"errors"
)

type Config struct {
	Env      Environment `json:"env"`
	Server   Server      `json:"server"`
	Database Database    `json:"database"`
	Logger   Logger      `json:"logger"`
}

type Environment string

const (
	Dev  Environment = "dev"
	Test Environment = "test"
	Prod Environment = "prod"
)

// NewConfigOrDefault creates tries to create a Config instance from provided environment path.
// If the path is empty, it will use the default environment.
func NewConfigOrDefault(envPath *string) Config {
	// TODO: Implement loading config from file
	// Try to load config from file
	// config, err := NewConfig()
	if envPath == nil || *envPath == "" {
		return DefaultConfig()
	}
	err := errors.New("error loading config") // TODO: Remove after actual implementation
	if err != nil {
		return DefaultConfig()
	}
	//return config
	return Config{}
}

func DefaultConfig() Config {
	return Config{
		Env: Dev,
		Server: Server{
			Host: "localhost",
			Port: "8080",
		},
		Database: Database{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			Name:     "flowcargo_dev",
		},
		Logger: Logger{
			Level:  "INFO",
			Source: true,
		},
	}
}

func (c Config) IsDevelopment() bool {
	return c.Env == Dev
}
