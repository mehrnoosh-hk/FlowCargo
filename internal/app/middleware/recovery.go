package middleware

import (
	"net/http"
	"runtime/debug"
)

func (m *middleware) Recovery() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					stack := debug.Stack()
					m.Logger.Errorf(
						"PANIC recovered:\n"+
						"Error: %v\n"+
						"Method: %s\n"+
						"URL: %s\n"+
						"Remote: %s\n"+
						"Stack Trace:\n%s",
						err,
						r.Method,
						r.URL.String(),
						r.RemoteAddr,
						stack,
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
