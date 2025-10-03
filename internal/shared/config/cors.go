package config

import (
	"time"
)

type CORS struct {
	AllowOrigins     []string      `mapstructure:"allow-origins" validate:"required,dive,uri"`
	AllowMethods     []string      `mapstructure:"allow-methods" validate:"required"`
	AllowHeaders     []string      `mapstructure:"allow-headers"`
	AllowCredentials bool          `mapstructure:"allow-credentials"`
	ExposeHeaders    []string      `mapstructure:"expose-headers"`
	MaxAge           time.Duration `mapstructure:"max-age" validate:"gte=0"`
}