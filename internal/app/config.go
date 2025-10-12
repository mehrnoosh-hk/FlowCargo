package app

import (
	"flowcargo/internal/shared/config"
)

func wireConfigFn(env config.Environment, envPath *string) (config.Config, error) {
	return config.New(env, envPath)
}
