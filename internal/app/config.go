package app

import (
	"context"
	"flowcargo/internal/shared/config"
)

var wireCfg = func(ctx context.Context, path *string) config.Config {
	return config.NewConfigOrDefault(path)
}
