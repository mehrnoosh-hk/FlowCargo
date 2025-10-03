package app

import (
	"context"
	"flowcargo/internal/shared/config"
)

var wireCfg = func(ctx context.Context, env config.Environment, path *string) (config.Config, error) {
	return config.New(env, path)
}
