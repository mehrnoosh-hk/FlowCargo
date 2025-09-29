package app

import (
	"flowcargo/internal/shared/config"
)

var wireCfg = func () config.Config {
	return config.NewConfigOrDefault()
}