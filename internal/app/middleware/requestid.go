package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (m *middleware) RequestID() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// check wether request id exist or not
			id := r.Header.Get("X-Request-ID")
			var ctx context.Context
			if id == "" {
				id = uuid.New().String()
			}
			ctx = context.WithValue(r.Context(), requestIDKey, id)
			w.Header().Set("X-Request-ID", id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
