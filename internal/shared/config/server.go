package config

import (
	"fmt"
)

type Server struct {
	Host string `mapstructure:"host" validate:"required,hostname_rfc1123"`
	Port string `mapstructure:"port" validate:"required,numeric"`
}

func (c Config) ServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}
