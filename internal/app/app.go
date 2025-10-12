package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"flowcargo/internal/shared/config"
	"flowcargo/internal/shared/logger"
)

// App is the main application struct that holds all dependencies and configurations.
type App struct {
	cfg  config.Config
	db   *Database
	deps Dependencies
	srv  Server
}

// Wire is an interface that defines methods for wiring up App struct fields.
// It helps for creating alternative implementations of the dependencies. (especially for testing)
type Wire interface {
	Up(ctx context.Context, env config.Environment, configPath *string) (App, error)
}

// CreateAndRun is lifecycle management for the application.
// It uses Wire interface to create the App instance and then runs it.
func CreateAndRun(
	ctx context.Context,
	environment config.Environment,
	configPath *string,
	wire Wire,
) error {
	app, err := wire.Up(ctx, environment, configPath)
	if err != nil {
		return err
	}
	app.Logger().Info("Application created successfully!")
	return app.runApp(ctx)
}

// runApp starts all the resources and runs the application.
// It waits for the application to finish.
func (a App) runApp(ctx context.Context) error {
	// Start server in a goroutine so it doesn't block
	go func() {
		if err := a.srv.start(); err != nil && err != http.ErrServerClosed {
			a.Logger().Errorf("Server error: %v", err)
		}
	}()

	a.Logger().Infof("Server started successfully! Listening on address %s", a.srv.getAddress())
	return a.waitForShutdown(ctx)
}

func (a App) waitForShutdown(ctx context.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

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
	a.db.Close() // Does not return an error, and waits for acquired connections to be released
	err := a.srv.shutdown(ctx)
	if err != nil {
		return err
	}
	err = a.Logger().Close()
	if err != nil {
		return err
	}
	fmt.Println("Resources cleaned successfully")
	return nil
}

// Logger returns the application's logger instance.
func (a App) Logger() logger.Logger {
	return a.deps.getLogger()
}
