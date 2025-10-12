package app

import (
	"context"
	"fmt"

	"flowcargo/internal/app/middleware"
	"flowcargo/internal/shared/config"
	"flowcargo/internal/shared/logger"
)

// wire is a specific implementation of the Wire interface, you can create alternative implementations for testing or other purposes.
type wire struct {
	wireConfig       func(env config.Environment, envPath *string) (config.Config, error)
	wireDatabase     func(ctx context.Context, dbURL string) (*Database, error)
	wireDependencies func(
		ctx context.Context,
		db *Database,
		isDev bool,
		corsCfg config.CORS,
		logLevel logger.LogLevel,
	) (Dependencies, error)
	wireServer func(address string, middleware middleware.Middleware, handlers Handlers) Server
}

// NewWire creates a new instance of the Wire interface with interchangeable implementations.
var NewWire = func() Wire {
	return wire{
		wireConfig:       wireConfigFn,
		wireDatabase:     wireDatabaseFn,
		wireDependencies: wireDependenciesFn,
		wireServer:       wireServerFn,
	}
}

// Up orchestrate wiring dependencies
func (w wire) Up(ctx context.Context, env config.Environment, configPath *string) (App, error) {
	// Load config
	cfg, err := w.wireConfig(env, configPath)
	if err != nil {
		return App{}, fmt.Errorf("failed to load config: %w", err)
	}

	// Create database connection
	db, err := w.wireDatabase(ctx, cfg.GetDatabaseURL())
	if err != nil {
		return App{}, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Create dependencies
	deps, err := w.wireDependencies(
		ctx,
		db,
		cfg.IsDevelopment(),
		cfg.CORS(),
		cfg.LogLevel(),
	)
	if err != nil {
		return App{}, fmt.Errorf("failed to create dependencies: %w", err)
	}

	// Create server
	server := w.wireServer(
		cfg.ServerAddress(),
		deps.Middleware,
		deps.Handlers,
	)

	return App{
		cfg:  cfg,
		db:   db,
		deps: deps,
		srv:  server,
	}, nil

}

func wireDependenciesFn(
	ctx context.Context,
	db *Database,
	isDev bool,
	corsCfg config.CORS,
	logLevel logger.LogLevel,
) (Dependencies, error) {
	deps := Dependencies{}
	l, err := wireLogger(isDev, logLevel)
	if err != nil {
		return Dependencies{}, fmt.Errorf("failed to create logger: %w", err)
	}

	repos, err := wireRepos(db.pool, l)
	if err != nil {
		return Dependencies{}, fmt.Errorf("failed to create repositories: %w", err)
	}

	services, err := wireServices(repos, l)
	if err != nil {
		return Dependencies{}, fmt.Errorf("failed to create services: %w", err)
	}

	handlers, err := wireHandlers(services, l)
	if err != nil {
		return Dependencies{}, fmt.Errorf("failed to create handlers: %w", err)
	}

	middleware := middleware.NewMiddleware(corsCfg, l)

	deps.Logger = l
	deps.Repos = repos
	deps.Services = services
	deps.Handlers = handlers
	deps.Middleware = middleware

	return deps, nil
}
