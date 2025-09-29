package config

import (
	"errors"
)


type Config struct {
	Env Environment `json:"env"`
	Server Server `json:"server"`
	Database Database `json:"database"`
	Logger Logger `json:"logger"`
}

type Environment string

const (
	Dev Environment = "dev"
	Test Environment = "test"
	Prod Environment = "prod"
)

func NewConfigOrDefault() (Config) {
	// Try to load config from file
	// config, err := NewConfig()
	err := errors.New("error loading config")
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
			Name: "flowcargo_dev",
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
