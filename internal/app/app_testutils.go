package app

import (
	"context"
	"errors"
	"net/http"

	"flowcargo/db/testutils"
	"flowcargo/internal/app/middleware"
	"flowcargo/internal/shared/config"
	"flowcargo/internal/shared/logger"
)

func newFailedWire() Wire {
	return wire{
		wireConfig: func(env config.Environment, envPath *string) (config.Config, error) {
			return config.Config{}, errors.New("test error")
		},
		wireDatabase: func(ctx context.Context, dbURL string) (*Database, error) {
			return &Database{}, nil
		},
		wireDependencies: func(ctx context.Context, db *Database, isDev bool, corsCfg config.CORS, logLevel logger.LogLevel) (Dependencies, error) {
			return Dependencies{}, nil
		},
		wireServer: func(address string, middleware middleware.Middleware, handlers Handlers) Server {
			return Server{srv: nil}
		},
	}
}

func newSucceedWire() Wire {
	return wire{
		wireConfig: func(env config.Environment, envPath *string) (config.Config, error) {
			cfg, err := config.New(config.Test, nil)
			return cfg, err
		},
		wireDatabase: func(ctx context.Context, dbURL string) (*Database, error) {
			dbManager := testutils.GetDBManager()
			pool := dbManager.GetPool()
			return &Database{
				pool: pool,
			}, nil
		},
		wireDependencies: func(ctx context.Context, db *Database, isDev bool, corsCfg config.CORS, logLevel logger.LogLevel) (Dependencies, error) {
			return Dependencies{
				Logger: testutils.NewTestLogger(),
			}, nil
		},
		wireServer: func(address string, middleware middleware.Middleware, handlers Handlers) Server {
			mux := http.NewServeMux()
			return Server{srv: &http.Server{
				Addr:    ":7070",
				Handler: mux,
			}}
		},
	}
}
