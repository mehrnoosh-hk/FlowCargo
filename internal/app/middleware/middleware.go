package middleware

import (
	"net/http"

	"flowcargo/internal/shared/config"
	"flowcargo/internal/shared/logger"
)

// Middleware is an interface for app middlewares
type Middleware interface {
	Recovery() func(next http.Handler) http.Handler
	CORS() func(next http.Handler) http.Handler
	RequestID() func(next http.Handler) http.Handler
	ReqLog() func(next http.Handler) http.Handler
}

type middleware struct {
	CORSConfig config.CORS
	Logger     logger.Logger
}

// NewMiddleware creates an implementation instance of Middleware interface
func NewMiddleware(corsConfig config.CORS, logger logger.Logger) Middleware {
	return &middleware{
		CORSConfig: corsConfig,
		Logger:     logger,
	}
}

type contextKey string

const requestIDKey contextKey = "requestID"
