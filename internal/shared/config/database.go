package config

import (
	"fmt"
)

type Database struct {
	Host     string `json:"host" mapstructure:"host" validate:"required,hostname_rfc1123"`
	Port     string `json:"port" mapstructure:"port" validate:"required,numeric"`
	User     string `json:"user" mapstructure:"user" validate:"required"`
	Password string `json:"password" mapstructure:"password" validate:"required"`
	Name     string `json:"name" mapstructure:"name" validate:"required"`
	SSLMode  string   `json:"ssl_mode" mapstructure:"ssl-mode" validate:"required"`
}

func (cfg Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name, cfg.Database.SSLMode)
}
