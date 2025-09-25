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
	Prod Environment = "prod"
)



type Database struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}



func NewConfig() (Config, error) {
	return Config{}, errors.New("Not Implemented")
}

func DefaultConfig() Config {
	return Config{
		Server: Server{
			Host: "localhost",
			Port: "8080",
		},
		Database: Database{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "password",
			Database: "flowcargo",
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
