package app

import (
	"fmt"
	"flowcargo/internal/shared/config"
	"flowcargo/internal/shared/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Dependencies struct {
	config config.Config
	dbPool *pgxpool.Pool
	Logger logger.Logger
}

var wireDependencies = func() (Dependencies, error) {
	deps := Dependencies{}

	cfg := wireConfig()
	l, err := wireLogger(cfg.IsDevelopment(), cfg.Logger.Level)
	if err != nil {
		return Dependencies{}, err
	}

	deps.config = cfg
	deps.Logger = l

	return deps, nil
}

var wireConfig = func () config.Config {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		fmt.Println("Fallback to default config")
		cfg = config.DefaultConfig() // TODO: implement actual logic
	}
	return cfg
}

var wireLogger = func(isDevelopment bool, level logger.LogLevel) (logger.Logger, error) {
	l, err := logger.NewLogger(isDevelopment, level)
	if err != nil {
		return logger.Logger{}, err
	}
	return l, nil
}
