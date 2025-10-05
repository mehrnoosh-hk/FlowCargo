package app

import (
	"context"
	"flowcargo/internal/app/middleware"
	"net/http"
)

type Server struct {
	srv *http.Server // Define your server fields here
}

var wireSrv = func(address string, middleware middleware.Middleware, handlers Handlers) Server {
	mux := http.NewServeMux()
	mux = wireRoutes(mux, handlers)
	handler := wireMiddleware(mux, middleware)

	return Server{
		srv: &http.Server{
			Addr:    address,
			Handler: handler,
		},
	}
}

func wireRoutes(mux *http.ServeMux, handlers Handlers) *http.ServeMux {
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("Hello, World!"))
	})

	// Route for creating a new tenant
	mux.HandleFunc("POST /tenants", handlers.TenantHandler.CreateTenant)

	// Route for getting a specific tenant by its ID
	mux.HandleFunc("GET /tenants/{id}", handlers.TenantHandler.GetTenant)

	// Route for updating a specific tenant by its ID
	mux.HandleFunc("PUT /tenants/{id}", handlers.TenantHandler.UpdateTenant)

	// Route for deleting (soft delete) a specific tenant by its ID
	mux.HandleFunc("DELETE /tenants/{id}", handlers.TenantHandler.DeleteTenant)

	return mux
}

func wireMiddleware(handler http.Handler, middleware middleware.Middleware) http.Handler {
	// Apply CORS middleware
	handler = middleware.CORS()(handler)
	return handler
}

func (s Server) getAddress() string {
	return s.srv.Addr
}

func (s Server) start() error {
	return s.srv.ListenAndServe()
}

func (s Server) shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
