package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"flowcargo/db/testutils"
	"flowcargo/internal/shared/config"
)

func TestRecoveryMiddleware(t *testing.T) {
	testLogger := testutils.NewTestLogger()
	m := NewMiddleware(config.CORS{}, testLogger)

	req := httptest.NewRequest(http.MethodGet, "/recover", nil)
	rr := httptest.NewRecorder()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Something went wrong")
	})

	handler := m.Recovery()(nextHandler)
	handler.ServeHTTP(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	testLoggerTyped := testLogger.(*testutils.TestLogger)
	
	require.Len(t, testLoggerTyped.ErrorMessages, 1)
	t.Log(testLoggerTyped.ErrorMessages[0])
}
