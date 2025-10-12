package config

import (
	"time"
)

// CORS holds the configuration for Cross-Origin Resource Sharing settings.
type CORS struct {
	AllowOrigins     []string      `mapstructure:"allow-origins" validate:"required,min=1,dive,uri"`
	AllowMethods     []string      `mapstructure:"allow-methods" validate:"required,min=1"`
	AllowHeaders     []string      `mapstructure:"allow-headers"`
	AllowCredentials bool          `mapstructure:"allow-credentials"`
	ExposeHeaders    []string      `mapstructure:"expose-headers"`
	MaxAge           time.Duration `mapstructure:"max-age" validate:"gte=0"`
}

// CORS returns the CORS configuration.
func (c Config) CORS() CORS {
	return c.CORSConfig
}
