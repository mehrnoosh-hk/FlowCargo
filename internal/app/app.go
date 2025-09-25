package app

import (
	"context"
	"flowcargo/internal/shared/logger"
	"net/http"
	"os"
	"os/signal"
)

type App struct {
	deps   Dependencies
	srv    Server
	logger logger.Logger // This is just for convenience
}

type wireFunc = func() (Dependencies, error)

func CreateAndRun(envPath string) error {
	app, err := newApp(wireDependencies)
	if err != nil {
		return err
	}
	app.logger.Info("Application created successfully!")
	app.logger.Infof("Env: %s", envPath)
	return app.runApp()
}

func newApp(fn wireFunc) (App, error) {
	deps, err := fn()
	if err != nil {
		return App{}, err
	}
	srv := wireServer()
	return App{
		deps:   deps,
		srv:    srv,
		logger: deps.Logger,
	}, nil
}

func (a App) runApp() error {
	// Start server in a goroutine so it doesn't block
	go func() {
		a.logger.Infof("Starting server on address %s", a.srv.getAddress())
		if err := a.srv.start(a.logger); err != nil && err != http.ErrServerClosed {
			a.logger.Errorf("Server error: %v", err)
		}
	}()

	a.logger.Info("Server started successfully!")
	return a.waitForShutdown(context.Background())
}

func (a App) waitForShutdown(ctx context.Context) error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	select {
	case <-quit:
		a.logger.Info("Received interrupt signal")
	case <-ctx.Done():
		a.logger.Info("Context cancelled")
	}

	return a.shutdown(ctx)
}

func (a App) shutdown(ctx context.Context) error {
	// TODO: Implement dependencies and resources cleanup
	a.logger.Info("Resources released successfully!")
	return a.srv.shutdown(ctx)
}
