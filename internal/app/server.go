package app

import (
	"context"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	"flowcargo/internal/app/middleware"
)

// Server is an interface that defines methods for starting and stopping the server.
type Server struct {
	srv *http.Server // Define your server fields here
}

func wireServerFn(address string, middleware middleware.Middleware, handlers Handlers) Server {
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
	// Swagger UI route - registered first
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

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
	// Apply middleware in reverse order of execution
	// (last applied = first executed)

	handler = middleware.CORS()(handler)
	handler = middleware.ReqLog()(handler)
	handler = middleware.RequestID()(handler)

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
