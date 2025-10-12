package config

import "flowcargo/internal/shared/logger"

type Logger struct {
	Level  logger.LogLevel `json:"level" mapstructure:"level" validate:"required,oneof=debug info warn error"`
	Source bool            `json:"source" mapstructure:"source"`
}

func (c Config) LogLevel() logger.LogLevel {
	return c.Logger.Level
}
