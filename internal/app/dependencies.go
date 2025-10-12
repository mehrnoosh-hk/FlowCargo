package app

import (
	"fmt"

	db "flowcargo/db/sqlc"
	"flowcargo/internal/app/middleware"
	"flowcargo/internal/shared/logger"
	"flowcargo/internal/tenant"
)

// Dependencies type is based on interfaces, making it easier to replace with new implementations.
type Dependencies struct {
	Logger     logger.Logger
	Repos      Repositories
	Services   Services
	Handlers   Handlers
	Middleware middleware.Middleware // We treat all middlewares the same way, since they are not entity specific.
}

// Repositories groups all repository dependencies.
type Repositories struct {
	TenantRepository tenant.TenantRepository
	// Other repositories here
}

// Services groups all service dependencies.
type Services struct {
	TenantService tenant.TenantService
	// Other services here
}

// Handlers groups all handler dependencies.
type Handlers struct {
	TenantHandler tenant.TenantHandler
	// Other handlers here
}

var wireLogger = func(isDevelopment bool, level logger.LogLevel) (logger.Logger, error) {
	l, err := logger.NewLogger(isDevelopment, level)
	if err != nil {
		return nil, err
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
