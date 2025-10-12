package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"flowcargo/db/testutils"
	"flowcargo/internal/shared/config"
)


func Test_middleware_RequestID(t *testing.T) {
	t.Run("Existing ID preserved", func(t *testing.T) {
		m := NewMiddleware(config.CORS{}, testutils.NewTestLogger())
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		existingID := uuid.New().String()
		req.Header.Set("X-Request-ID", existingID)
		rr := httptest.NewRecorder()
		handler := m.RequestID()
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, existingID, r.Context().Value(requestIDKey))
			require.Equal(t,rr.Code, http.StatusOK)
		})
		handler(nextHandler).ServeHTTP(rr, req)
	})

	t.Run("New ID generated", func(t *testing.T) {
		m := NewMiddleware(config.CORS{}, testutils.NewTestLogger())
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		handler := m.RequestID()
		nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.NotEqual(t, "", r.Context().Value(requestIDKey))
			require.Equal(t,rr.Code, http.StatusOK)
		})
		handler(nextHandler).ServeHTTP(rr, req)
	})
}

func Test_middleware_RequestID_Concurrency(t *testing.T) {
	const numberOfConcurrentReqests = 100
	results := make(chan string, numberOfConcurrentReqests)
	m := NewMiddleware(config.CORS{}, testutils.NewTestLogger())
	handler := m.RequestID()
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(requestIDKey).(string)
		results <- id
	})
	
	var wg sync.WaitGroup

	for range numberOfConcurrentReqests {
		wg.Go(func(){
			req := httptest.NewRequest(http.MethodGet, "/", nil)
        	rr := httptest.NewRecorder()
			handler(nextHandler).ServeHTTP(rr, req)
		})
	}
	wg.Wait()
	close(results)

	collectedIDs := make([]string, 0, numberOfConcurrentReqests)
	for id := range results {
		collectedIDs = append(collectedIDs, id)
	}
	require.Len(t, collectedIDs, numberOfConcurrentReqests)

	uniqueIDs := make(map[string]bool)
	for _, id := range collectedIDs {
		require.False(t, uniqueIDs[id], "duplicate ID found: %s", id)
		uniqueIDs[id] = true
	}
}
