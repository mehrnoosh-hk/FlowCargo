package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"flowcargo/db/testutils"
	"flowcargo/internal/app/middleware"
	"flowcargo/internal/shared/config"
)

// This only tests that middleware are correctly wired to server or not
// The tests for middleware functionality are in the middleware package
func TestWireMiddleware(t *testing.T) {

	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("Hello, World!"))
	})
	corsCfg := config.CORS{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"*"},
	}
	// Create a new CORS middleware
	corsMiddleware := middleware.NewMiddleware(
		corsCfg, testutils.NewTestLogger(),
	)

	// wire the CORS middleware
	handler := wireMiddleware(mux, corsMiddleware)

	// Create a request and response recorder
	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)
	// Set Origin header to trigger CORS middleware
	req.Header.Set("Origin", "http://example.com")
	rr := httptest.NewRecorder()

	// Execute the request
	handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "Hello, World!", rr.Body.String())
	require.Equal(t, "http://example.com", rr.Header().Get("Access-Control-Allow-Origin"))
	require.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
	require.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Headers"))

}
