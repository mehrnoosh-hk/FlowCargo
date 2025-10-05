package app

import (
	"context"
	db "flowcargo/db/sqlc"
	"flowcargo/internal/app/middleware"
	"flowcargo/internal/shared/logger"
	"flowcargo/internal/tenant"
	"fmt"
)

// Dependencies type is based on interfaces, making it easier to replace with new implementations.
type Dependencies struct {
	Logger   logger.Logger
	Repos    Repositories
	Services Services
	Handlers Handlers
	Middleware middleware.Middleware // We treat all middlewares the same way, since they are not entity specific.
}

type Repositories struct {
	TenantRepository tenant.TenantRepository
	// Other repositories here
}

type Services struct {
	TenantService tenant.TenantService
	// Other services here
}

type Handlers struct {
	TenantHandler tenant.TenantHandler
	// Other handlers here
}

var wireDeps = func(
	ctx context.Context,
	db *Database,
	isDev bool,
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

	deps.Logger = l
	deps.Repos = repos
	deps.Services = services
	deps.Handlers = handlers
	return deps, nil
}

var wireLogger = func(isDevelopment bool, level logger.LogLevel) (logger.Logger, error) {
	l, err := logger.NewLogger(isDevelopment, level)
	if err != nil {
		return logger.Logger{}, err
	}
	return l, nil
}

var wireRepos = func(conn db.DBTX, l logger.Logger) (Repositories, error) {
	tr, err := wireTenantRepository(conn, l)
	if err != nil {
		return Repositories{}, err
	}
	return Repositories{
		TenantRepository: tr,
	}, nil
}

var wireServices = func(repos Repositories, l logger.Logger) (Services, error) {
	ts, err := wireTenantService(repos.TenantRepository, l)
	if err != nil {
		return Services{}, err
	}
	return Services{
		TenantService: ts,
	}, nil
}

var wireHandlers = func(services Services, l logger.Logger) (Handlers, error) {
	th, err := wireTenantHandler(services.TenantService, l)
	if err != nil {
		return Handlers{}, err
	}
	return Handlers{
		TenantHandler: th,
	}, nil
}

var wireTenantRepository = func(conn db.DBTX, l logger.Logger) (tenant.TenantRepository, error) {
	if conn == nil {
		return nil, fmt.Errorf("conn is nil")
	}
	tenantRepo := tenant.NewTenantRepository(conn, l)
	return tenantRepo, nil
}

func (d Dependencies) getLogger() logger.Logger {
	return d.Logger
}

var wireTenantService = func(repo tenant.TenantRepository, l logger.Logger) (tenant.TenantService, error) {
	ts := tenant.NewTenantService(repo, l)
	return ts, nil
}

var wireTenantHandler = func(service tenant.TenantService, l logger.Logger) (tenant.TenantHandler, error) {
	th := tenant.NewTenantHandler(service, l)
	return th, nil
}
