package config

import (
	"fmt"
)

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

func (c Config) ServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}

func (c Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}
