package middleware

import (
	"net/http"
	"flowcargo/internal/shared/config"
	"flowcargo/internal/shared/logger"
)

type Middleware interface {
	CorsMiddleware() func(next http.Handler) http.Handler
}

type middleware struct {
	CORSConfig config.CORS
	Logger     logger.Logger
}

func NewMiddleware(corsConfig config.CORS, logger logger.Logger) Middleware {
	return &middleware{
		CORSConfig: corsConfig,
		Logger:     logger,
	}
}
