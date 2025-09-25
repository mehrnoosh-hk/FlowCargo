package config

import "flowcargo/internal/shared/logger"

type Logger struct {
	Level logger.LogLevel `json:"level"`
	Source bool `json:"source"`
}

func (c Config) LogLevel() logger.LogLevel {
	return c.Logger.Level
}
