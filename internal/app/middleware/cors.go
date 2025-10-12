package middleware

import (
	"net/http"
	"slices"
	"strconv"
	"strings"

	"flowcargo/internal/shared/config"
)

func (m *middleware) CORS() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is allowed
			if origin != "" && !validOrigin(origin, m.CORSConfig) {
				m.Logger.Warnf("Unauthorized origin: %s", origin)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Set CORS headers
			if origin != "" && validOrigin(origin, m.CORSConfig) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			if len(m.CORSConfig.AllowMethods) > 0 {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.CORSConfig.AllowMethods, ", "))
			}

			if len(m.CORSConfig.AllowHeaders) > 0 {
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(m.CORSConfig.AllowHeaders, ", "))
			}

			if m.CORSConfig.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(m.CORSConfig.AllowCredentials))
			}

			if len(m.CORSConfig.ExposeHeaders) > 0 {
				w.Header().Set("Access-Control-Expose-Headers", strings.Join(m.CORSConfig.ExposeHeaders, ", "))
			}

			if m.CORSConfig.MaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", strconv.Itoa(int(m.CORSConfig.MaxAge.Seconds())))
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func validOrigin(origin string, CORScfg config.CORS) bool {
	// Check for wildcard
	if slices.Contains(CORScfg.AllowOrigins, "*") {
		return true
	}

	// Check for exact match
	if slices.Contains(CORScfg.AllowOrigins, origin) {
		return true
	}

	// Check for wildcard subdomain patterns (e.g., "https://*.example.com")
	for _, allowed := range CORScfg.AllowOrigins {
		if strings.Contains(allowed, "*") {
			// Split on * to get prefix and suffix
			// For example: "https://*.example.com" should match "https://api.example.com"
			parts := strings.Split(allowed, "*")
			if len(parts) == 2 {
				prefix := parts[0]
				suffix := parts[1]
				if strings.HasPrefix(origin, prefix) && strings.HasSuffix(origin, suffix) {
					// Make sure there's something between prefix and suffix
					middle := strings.TrimPrefix(strings.TrimSuffix(origin, suffix), prefix)
					if len(middle) > 0 {
						return true
					}
				}
			}
		}
	}

	return false
}
