package middleware

import (
	"flowcargo/internal/shared/config"
	"flowcargo/internal/shared/logger"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCorsMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		corsConfig     config.CORS
		requestMethod  string
		requestOrigin  string
		requestHeaders map[string]string

		// Expected results
		expectStatus            int
		expectAllowOrigin       string
		expectAllowMethods      string
		expectAllowHeaders      string
		expectAllowCredentials  string
		expectExposeHeaders     string
		expectMaxAge            string
		expectNextHandlerCalled bool
	}{
		// Test Case 1: Valid origin with all CORS features enabled
		{
			name: "valid origin with all features",
			corsConfig: config.CORS{
				AllowOrigins:     []string{"https://example.com"},
				AllowMethods:     []string{"GET", "POST", "PUT"},
				AllowHeaders:     []string{"Content-Type", "Authorization"},
				AllowCredentials: true,
				ExposeHeaders:    []string{"X-Custom-Header"},
				MaxAge:           1 * time.Hour,
			},
			requestMethod:           "GET",
			requestOrigin:           "https://example.com",
			expectStatus:            http.StatusOK,
			expectAllowOrigin:       "https://example.com",
			expectAllowMethods:      "GET, POST, PUT",
			expectAllowHeaders:      "Content-Type, Authorization",
			expectAllowCredentials:  "true",
			expectExposeHeaders:     "X-Custom-Header",
			expectMaxAge:            "3600",
			expectNextHandlerCalled: true,
		},

		// Test Case 2: Invalid origin - should be rejected
		{
			name: "invalid origin rejected",
			corsConfig: config.CORS{
				AllowOrigins: []string{"https://example.com"},
				AllowMethods: []string{"GET"},
			},
			requestMethod:           "GET",
			requestOrigin:           "https://malicious.com",
			expectStatus:            http.StatusForbidden,
			expectNextHandlerCalled: false,
		},

		// Test Case 3: Empty origin (same-origin request)
		{
			name: "empty origin allowed (same-origin)",
			corsConfig: config.CORS{
				AllowOrigins: []string{"https://example.com"},
				AllowMethods: []string{"GET"},
			},
			requestMethod:           "GET",
			requestOrigin:           "",
			expectStatus:            http.StatusOK,
			expectAllowOrigin:       "",
			expectNextHandlerCalled: true,
		},

		// Test Case 4: Wildcard origin
		{
			name: "wildcard origin allows any origin",
			corsConfig: config.CORS{
				AllowOrigins: []string{"*"},
				AllowMethods: []string{"GET"},
			},
			requestMethod:           "GET",
			requestOrigin:           "https://any-domain.com",
			expectStatus:            http.StatusOK,
			expectAllowOrigin:       "https://any-domain.com",
			expectNextHandlerCalled: true,
		},

		// Test Case 5: Wildcard subdomain pattern
		{
			name: "wildcard subdomain pattern",
			corsConfig: config.CORS{
				AllowOrigins: []string{"https://*.example.com"},
				AllowMethods: []string{"GET"},
			},
			requestMethod:           "GET",
			requestOrigin:           "https://api.example.com",
			expectStatus:            http.StatusOK,
			expectAllowOrigin:       "https://api.example.com",
			expectNextHandlerCalled: true,
		},

		// Test Case 6: Wildcard subdomain pattern - non-matching
		{
			name: "wildcard subdomain pattern non-matching",
			corsConfig: config.CORS{
				AllowOrigins: []string{"https://*.example.com"},
				AllowMethods: []string{"GET"},
			},
			requestMethod:           "GET",
			requestOrigin:           "https://other.com",
			expectStatus:            http.StatusForbidden,
			expectNextHandlerCalled: false,
		},

		// Test Case 7: Preflight request (OPTIONS) - valid origin
		{
			name: "preflight request with valid origin",
			corsConfig: config.CORS{
				AllowOrigins: []string{"https://example.com"},
				AllowMethods: []string{"GET", "POST"},
				AllowHeaders: []string{"Content-Type"},
			},
			requestMethod:           "OPTIONS",
			requestOrigin:           "https://example.com",
			expectStatus:            http.StatusNoContent,
			expectAllowOrigin:       "https://example.com",
			expectAllowMethods:      "GET, POST",
			expectAllowHeaders:      "Content-Type",
			expectNextHandlerCalled: false, // Preflight stops here
		},

		// Test Case 8: Preflight request (OPTIONS) - invalid origin
		{
			name: "preflight request with invalid origin",
			corsConfig: config.CORS{
				AllowOrigins: []string{"https://example.com"},
				AllowMethods: []string{"GET"},
			},
			requestMethod:           "OPTIONS",
			requestOrigin:           "https://malicious.com",
			expectStatus:            http.StatusForbidden,
			expectNextHandlerCalled: false,
		},

		// Test Case 9: AllowCredentials false
		{
			name: "credentials not allowed",
			corsConfig: config.CORS{
				AllowOrigins:     []string{"https://example.com"},
				AllowMethods:     []string{"GET"},
				AllowCredentials: false,
			},
			requestMethod:           "GET",
			requestOrigin:           "https://example.com",
			expectStatus:            http.StatusOK,
			expectAllowOrigin:       "https://example.com",
			expectAllowCredentials:  "", // Should not be set
			expectNextHandlerCalled: true,
		},

		// Test Case 10: Multiple allowed origins
		{
			name: "multiple allowed origins - second match",
			corsConfig: config.CORS{
				AllowOrigins: []string{"https://example.com", "https://another.com"},
				AllowMethods: []string{"GET"},
			},
			requestMethod:           "GET",
			requestOrigin:           "https://another.com",
			expectStatus:            http.StatusOK,
			expectAllowOrigin:       "https://another.com",
			expectNextHandlerCalled: true,
		},

		// Test Case 11: No ExposeHeaders configured
		{
			name: "no expose headers",
			corsConfig: config.CORS{
				AllowOrigins:  []string{"https://example.com"},
				AllowMethods:  []string{"GET"},
				ExposeHeaders: []string{},
			},
			requestMethod:           "GET",
			requestOrigin:           "https://example.com",
			expectStatus:            http.StatusOK,
			expectAllowOrigin:       "https://example.com",
			expectExposeHeaders:     "",
			expectNextHandlerCalled: true,
		},

		// Test Case 12: Zero MaxAge
		{
			name: "zero max age not set",
			corsConfig: config.CORS{
				AllowOrigins: []string{"https://example.com"},
				AllowMethods: []string{"GET"},
				MaxAge:       0,
			},
			requestMethod:           "GET",
			requestOrigin:           "https://example.com",
			expectStatus:            http.StatusOK,
			expectAllowOrigin:       "https://example.com",
			expectMaxAge:            "",
			expectNextHandlerCalled: true,
		},

		// Test Case 13: POST request with valid origin
		{
			name: "POST request with valid origin",
			corsConfig: config.CORS{
				AllowOrigins: []string{"https://example.com"},
				AllowMethods: []string{"GET", "POST"},
			},
			requestMethod:           "POST",
			requestOrigin:           "https://example.com",
			expectStatus:            http.StatusOK,
			expectAllowOrigin:       "https://example.com",
			expectNextHandlerCalled: true,
		},

		// Test Case 14: Empty AllowMethods
		{
			name: "empty allow methods",
			corsConfig: config.CORS{
				AllowOrigins: []string{"https://example.com"},
				AllowMethods: []string{},
			},
			requestMethod:           "GET",
			requestOrigin:           "https://example.com",
			expectStatus:            http.StatusOK,
			expectAllowOrigin:       "https://example.com",
			expectAllowMethods:      "",
			expectNextHandlerCalled: true,
		},

		// Test Case 15: Empty AllowHeaders
		{
			name: "empty allow headers",
			corsConfig: config.CORS{
				AllowOrigins: []string{"https://example.com"},
				AllowMethods: []string{"GET"},
				AllowHeaders: []string{},
			},
			requestMethod:           "GET",
			requestOrigin:           "https://example.com",
			expectStatus:            http.StatusOK,
			expectAllowOrigin:       "https://example.com",
			expectAllowHeaders:      "",
			expectNextHandlerCalled: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test logger
			testLogger, err := logger.NewLogger(true, logger.Info)
			if err != nil {
				t.Fatalf("Failed to create logger: %v", err)
			}

			// Create middleware instance
			m := NewMiddleware(tt.corsConfig, testLogger)

			// Track if next handler was called
			nextHandlerCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextHandlerCalled = true
				w.WriteHeader(http.StatusOK)
			})

			// Create the CORS middleware handler
			handler := m.CorsMiddleware()(nextHandler)

			// Create request
			req := httptest.NewRequest(tt.requestMethod, "http://example.com/test", nil)
			if tt.requestOrigin != "" {
				req.Header.Set("Origin", tt.requestOrigin)
			}

			// Add any additional headers
			for key, value := range tt.requestHeaders {
				req.Header.Set(key, value)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute the handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if rr.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, rr.Code)
			}

			// Check if next handler was called
			if nextHandlerCalled != tt.expectNextHandlerCalled {
				t.Errorf("Expected nextHandlerCalled=%v, got %v", tt.expectNextHandlerCalled, nextHandlerCalled)
			}

			// Check CORS headers
			if tt.expectAllowOrigin != "" {
				if got := rr.Header().Get("Access-Control-Allow-Origin"); got != tt.expectAllowOrigin {
					t.Errorf("Expected Access-Control-Allow-Origin='%s', got '%s'", tt.expectAllowOrigin, got)
				}
			} else {
				// If we expect empty, make sure it's not set (or is empty)
				if got := rr.Header().Get("Access-Control-Allow-Origin"); got != "" && tt.expectStatus == http.StatusOK {
					t.Errorf("Expected Access-Control-Allow-Origin to be empty, got '%s'", got)
				}
			}

			if tt.expectAllowMethods != "" {
				if got := rr.Header().Get("Access-Control-Allow-Methods"); got != tt.expectAllowMethods {
					t.Errorf("Expected Access-Control-Allow-Methods='%s', got '%s'", tt.expectAllowMethods, got)
				}
			}

			if tt.expectAllowHeaders != "" {
				if got := rr.Header().Get("Access-Control-Allow-Headers"); got != tt.expectAllowHeaders {
					t.Errorf("Expected Access-Control-Allow-Headers='%s', got '%s'", tt.expectAllowHeaders, got)
				}
			}

			if tt.expectAllowCredentials != "" {
				if got := rr.Header().Get("Access-Control-Allow-Credentials"); got != tt.expectAllowCredentials {
					t.Errorf("Expected Access-Control-Allow-Credentials='%s', got '%s'", tt.expectAllowCredentials, got)
				}
			} else {
				// Verify it's not set when credentials are false
				if got := rr.Header().Get("Access-Control-Allow-Credentials"); got == "true" && !tt.corsConfig.AllowCredentials {
					t.Errorf("Expected Access-Control-Allow-Credentials to not be 'true', got '%s'", got)
				}
			}

			if tt.expectExposeHeaders != "" {
				if got := rr.Header().Get("Access-Control-Expose-Headers"); got != tt.expectExposeHeaders {
					t.Errorf("Expected Access-Control-Expose-Headers='%s', got '%s'", tt.expectExposeHeaders, got)
				}
			}

			if tt.expectMaxAge != "" {
				if got := rr.Header().Get("Access-Control-Max-Age"); got != tt.expectMaxAge {
					t.Errorf("Expected Access-Control-Max-Age='%s', got '%s'", tt.expectMaxAge, got)
				}
			}
		})
	}
}

func TestValidOrigin(t *testing.T) {
	tests := []struct {
		name          string
		origin        string
		allowOrigins  []string
		expectedValid bool
	}{
		{
			name:          "exact match",
			origin:        "https://example.com",
			allowOrigins:  []string{"https://example.com"},
			expectedValid: true,
		},
		{
			name:          "no match",
			origin:        "https://malicious.com",
			allowOrigins:  []string{"https://example.com"},
			expectedValid: false,
		},
		{
			name:          "wildcard allows all",
			origin:        "https://anything.com",
			allowOrigins:  []string{"*"},
			expectedValid: true,
		},
		{
			name:          "wildcard subdomain match",
			origin:        "https://api.example.com",
			allowOrigins:  []string{"https://*.example.com"},
			expectedValid: true,
		},
		{
			name:          "wildcard subdomain no match",
			origin:        "https://other.com",
			allowOrigins:  []string{"https://*.example.com"},
			expectedValid: false,
		},
		{
			name:          "multiple origins first match",
			origin:        "https://example.com",
			allowOrigins:  []string{"https://example.com", "https://another.com"},
			expectedValid: true,
		},
		{
			name:          "multiple origins second match",
			origin:        "https://another.com",
			allowOrigins:  []string{"https://example.com", "https://another.com"},
			expectedValid: true,
		},
		{
			name:          "empty origin",
			origin:        "",
			allowOrigins:  []string{"https://example.com"},
			expectedValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			corsConfig := config.CORS{
				AllowOrigins: tt.allowOrigins,
			}

			result := validOrigin(tt.origin, corsConfig)

			if result != tt.expectedValid {
				t.Errorf("Expected validOrigin()=%v, got %v", tt.expectedValid, result)
			}
		})
	}
}
