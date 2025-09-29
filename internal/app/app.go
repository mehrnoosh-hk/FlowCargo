package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"flowcargo/internal/shared/config"
	"flowcargo/internal/shared/logger"
)

type App struct {
	cfg  config.Config
	db   *Database
	deps Dependencies
	srv  Server
}

type wireConfig = func() config.Config

// wireDB is a function that wires up the database dependency for the application.
// It helps for creating alternative implementations of the database dependency. (especially for testing)
type wireDataBase = func(ctx context.Context, URL string) (*Database, error)

// wireDeps is a function that wires up the dependencies for the application.
// It helps for creating alternative implementations of the dependencies. (especially for testing)
type wireDependencies = func(ctx context.Context, db *Database, isDev bool, logLevel logger.LogLevel) (Dependencies, error)

// wireServer is a function that wires up the server dependency for the application.
// It helps for creating alternative implementations of the server dependency. (especially for testing)
type wireServer = func(address string, handlers Handlers) Server

// CreateAndRun is lifecycle management for the application.
// It depends on function type to be replacable
func CreateAndRun(ctx context.Context, envPath string) error {
	app, err := newApp(
		ctx,
		wireCfg,
		wireDB,
		wireDeps,
		wireSrv,
	)
	if err != nil {
		return err
	}
	app.Logger().Info("Application created successfully!")
	app.Logger().Infof("Env: %s", envPath)
	return app.runApp()
}

// newApp is a factory function that creates a new App instance with all dependencies.
// By replacing the passing functions, you can change the implementation of the App
func newApp(
	ctx context.Context,
	fConfig wireConfig,
	fDB wireDataBase,
	fDeps wireDependencies,
	fServer wireServer,
) (App, error) {
	cfg := fConfig()
	db, err := fDB(ctx, cfg.GetDatabaseURL())
	if err != nil {
		return App{}, err
	}
	deps, err := fDeps(ctx, db, cfg.IsDevelopment(), cfg.LogLevel())
	if err != nil {
		return App{}, err
	}

	srv := fServer(cfg.ServerAddress(), deps.Handlers)
	return App{
		cfg:  cfg,
		db:   db,
		deps: deps,
		srv:  srv,
	}, nil
}

// runApp starts all the resources and runs the application.
// It waits for the application to finish.
func (a App) runApp() error {
	// Start server in a goroutine so it doesn't block
	go func() {
		if err := a.srv.start(); err != nil && err != http.ErrServerClosed {
			a.Logger().Errorf("Server error: %v", err)
		}
	}()

	a.Logger().Infof("Server started successfully! Listening on address %s", a.srv.getAddress())
	return a.waitForShutdown(context.Background())
}

func (a App) waitForShutdown(ctx context.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		a.Logger().Info("Received interrupt signal")
	case <-ctx.Done():
		a.Logger().Info("Context cancelled")
	}

	return a.shutdown(ctx)
}

// shutdown stops all the resources.
func (a App) shutdown(ctx context.Context) error {
	// TODO: Implement dependencies and resources cleanup
	a.db.Close()
	err := a.srv.shutdown(ctx)
	if err != nil {
		return err
	}
	a.Logger().Info("Resources released successfully!")
	return nil
}

func (a App) Logger() logger.Logger {
	return a.deps.getLogger()
}
