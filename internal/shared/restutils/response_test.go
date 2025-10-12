package restutils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"flowcargo/db/testutils"
)

func TestWriteJSONResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		status         int
		response       any
		expectedStatus int
		expectedHeader string
		shouldContain  string
	}{
		{
			name:           "valid JSON response",
			status:         http.StatusOK,
			response:       map[string]string{"message": "success"},
			expectedStatus: http.StatusOK,
			expectedHeader: "application/json",
			shouldContain:  `"message":"success"`,
		},
		{
			name:           "empty response",
			status:         http.StatusNoContent,
			response:       nil,
			expectedStatus: http.StatusNoContent,
			expectedHeader: "application/json",
			shouldContain:  "null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			writeJSONResponse(rr, tt.status, tt.response)

			require.Equal(t, tt.expectedStatus, rr.Code)
			require.Equal(t, tt.expectedHeader, rr.Header().Get("Content-Type"))
			require.Contains(t, rr.Body.String(), tt.shouldContain)
		})
	}
}

func TestWriteJSONResponseWithMarshalError(t *testing.T) {
	t.Parallel()

	// Create a response that will cause a marshal error
	invalidResponse := make(chan int) // Channels can't be marshaled to JSON

	rr := httptest.NewRecorder()
	writeJSONResponse(rr, http.StatusOK, invalidResponse)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Equal(t, "Failed to marshal response", strings.TrimSpace(rr.Body.String()))
}

func TestWriteSuccessResponse(t *testing.T) {
	t.Parallel()

	type testData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	tests := []struct {
		name           string
		status         int
		data           testData
		message        string
		expectedStatus int
		wantSuccess    bool
		wantData       bool
		wantMessage    string
	}{
		{
			name:           "success response with data",
			status:         http.StatusOK,
			data:           testData{ID: 1, Name: "test"},
			message:        "Operation successful",
			expectedStatus: http.StatusOK,
			wantSuccess:    true,
			wantData:       true,
			wantMessage:    "Operation successful",
		},
		{
			name:           "success response with empty message",
			status:         http.StatusCreated,
			data:           testData{ID: 2, Name: "test2"},
			message:        "",
			expectedStatus: http.StatusCreated,
			wantSuccess:    true,
			wantData:       true,
			wantMessage:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			WriteSuccessResponse(rr, tt.status, tt.data, tt.message)

			require.Equal(t, tt.expectedStatus, rr.Code)

			var response APIResponse[testData]
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))

			require.Equal(t, tt.wantSuccess, response.Success)
			require.Equal(t, tt.wantMessage, response.Message)

			if tt.wantData {
				require.NotNil(t, response.Data)
				require.Equal(t, tt.data.ID, response.Data.ID)
				require.Equal(t, tt.data.Name, response.Data.Name)
			}

			require.False(t, response.Meta.Timestamp.IsZero())
		})
	}
}

func TestWriteErrorResponse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		details        ErrorDetail
		expectedStatus int
		wantSuccess    bool
		wantMessage    string
	}{
		{
			name: "error response",
			details: ErrorDetail{
				Status:   http.StatusBadRequest,
				Code:     ErrCodeValidation,
				Resource: ResourceTenant,
				Errors:   "Invalid input",
			},
			expectedStatus: http.StatusBadRequest,
			wantSuccess:    false,
			wantMessage:    "Error",
		},
		{
			name: "internal server error",
			details: ErrorDetail{
				Status:   http.StatusInternalServerError,
				Code:     ErrCodeInternal,
				Resource: ResourceTenant,
				Errors:   "Database connection failed",
			},
			expectedStatus: http.StatusInternalServerError,
			wantSuccess:    false,
			wantMessage:    "Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			writeErrorResponse(rr, tt.details)

			require.Equal(t, tt.expectedStatus, rr.Code)

			var response APIResponse[any]
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))

			require.Equal(t, tt.wantSuccess, response.Success)
			require.Equal(t, tt.wantMessage, response.Message)
			require.Equal(t, tt.details.Status, response.Error.Status)
			require.Equal(t, tt.details.Code, response.Error.Code)
			require.Equal(t, tt.details.Resource, response.Error.Resource)
			require.Equal(t, tt.details.Errors, response.Error.Errors)
			require.False(t, response.Meta.Timestamp.IsZero())
		})
	}
}

func TestHandleBadRequest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedCode   ErrCode
		expectedLog    string
	}{
		{
			name:           "bad request with validation error",
			err:            errors.New("invalid email format"),
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrCodeValidation,
			expectedLog:    "Bad request error invalid email format",
		},
		{
			name:           "bad request with empty error",
			err:            errors.New(""),
			expectedStatus: http.StatusBadRequest,
			expectedCode:   ErrCodeValidation,
			expectedLog:    "Bad request error ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			mockLogger := testutils.NewTestLogger()

			HandleBadRequest(rr, tt.err, mockLogger)

			require.Equal(t, tt.expectedStatus, rr.Code)

			var response APIResponse[any]
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))

			require.False(t, response.Success)
			require.Equal(t, tt.expectedStatus, response.Error.Status)
			require.Equal(t, tt.expectedCode, response.Error.Code)
			require.Equal(t, ResourceTenant, response.Error.Resource)
			require.Equal(t, tt.err.Error(), response.Error.Errors)

			// Check logger was called
			testLogger := mockLogger.(*testutils.TestLogger)
			require.NotEmpty(t, testLogger.WarnMessages)
			require.Equal(t, tt.expectedLog, testLogger.WarnMessages[0])
		})
	}
}

func TestHandleInternalServerError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		err            error
		resource       Resource
		expectedStatus int
		expectedCode   ErrCode
		expectedLog    string
	}{
		{
			name:           "internal server error with tenant resource",
			err:            errors.New("database connection failed"),
			resource:       ResourceTenant,
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   ErrCodeInternal,
			expectedLog:    "Internal server errorerrordatabase connection failedresourcetenant",
		},
		{
			name:           "internal server error with empty error",
			err:            errors.New(""),
			resource:       ResourceTenant,
			expectedStatus: http.StatusInternalServerError,
			expectedCode:   ErrCodeInternal,
			expectedLog:    "Internal server errorerrorresourcetenant",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			mockLogger := testutils.NewTestLogger()

			HandleInternalServerError(rr, tt.err, tt.resource, mockLogger)

			require.Equal(t, tt.expectedStatus, rr.Code)

			var response APIResponse[any]
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))

			require.False(t, response.Success)
			require.Equal(t, tt.expectedStatus, response.Error.Status)
			require.Equal(t, tt.expectedCode, response.Error.Code)
			require.Equal(t, tt.resource, response.Error.Resource)
			require.Equal(t, tt.err.Error(), response.Error.Errors)

			// Check logger was called
			testLogger := mockLogger.(*testutils.TestLogger)
			require.NotEmpty(t, testLogger.ErrorMessages)
			require.Equal(t, tt.expectedLog, testLogger.ErrorMessages[0])
		})
	}
}

func TestHandleMethodNotAllowed(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		url            *url.URL
		resource       Resource
		expectedStatus int
		expectedCode   ErrCode
		expectedLog    string
	}{
		{
			name:           "method not allowed with POST",
			method:         "POST",
			url:            &url.URL{Path: "/api/tenants"},
			resource:       ResourceTenant,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCode:   ErrCodeMethodNotAllowed,
			expectedLog:    "Method not allowedmethodPOSTurl/api/tenantsresourcetenant",
		},
		{
			name:           "method not allowed with PUT",
			method:         "PUT",
			url:            &url.URL{Path: "/api/tenants/123"},
			resource:       ResourceTenant,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedCode:   ErrCodeMethodNotAllowed,
			expectedLog:    "Method not allowedmethodPUTurl/api/tenants/123resourcetenant",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			mockLogger := testutils.NewTestLogger()

			HandleMethodNotAllowed(rr, tt.method, tt.url, tt.resource, mockLogger)

			require.Equal(t, tt.expectedStatus, rr.Code)

			var response APIResponse[any]
			require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &response))

			require.False(t, response.Success)
			require.Equal(t, tt.expectedStatus, response.Error.Status)
			require.Equal(t, tt.expectedCode, response.Error.Code)
			require.Equal(t, tt.resource, response.Error.Resource)
			require.Equal(t, "Method not allowed", response.Error.Errors)

			// Check logger was called
			testLogger := mockLogger.(*testutils.TestLogger)
			require.NotEmpty(t, testLogger.WarnMessages)
			require.Equal(t, tt.expectedLog, testLogger.WarnMessages[0])
		})
	}
}

func TestAPIResponseSerialization(t *testing.T) {
	t.Parallel()

	type testData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	original := APIResponse[testData]{
		Success: true,
		Message: "Test message",
		Data: &testData{
			ID:   123,
			Name: "test",
		},
		Error: ErrorDetail{
			Status:   0,
			Code:     "",
			Resource: "",
			Errors:   "",
		},
		Meta: Metadata{
			Timestamp: time.Now(),
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal from JSON
	var unmarshaled APIResponse[testData]
	require.NoError(t, json.Unmarshal(data, &unmarshaled))

	// Compare fields
	require.Equal(t, original.Success, unmarshaled.Success)
	require.Equal(t, original.Message, unmarshaled.Message)

	require.NotNil(t, unmarshaled.Data)
	require.Equal(t, original.Data.ID, unmarshaled.Data.ID)
	require.Equal(t, original.Data.Name, unmarshaled.Data.Name)

	// Check that timestamp is preserved (with some tolerance for serialization)
	require.Equal(t, original.Meta.Timestamp.Unix(), unmarshaled.Meta.Timestamp.Unix())
}

// Benchmark tests
func BenchmarkWriteSuccessResponse(b *testing.B) {
	type testData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	data := testData{ID: 1, Name: "test"}
	message := "Success"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		WriteSuccessResponse(rr, http.StatusOK, data, message)
	}
}

func BenchmarkWriteErrorResponse(b *testing.B) {
	details := ErrorDetail{
		Status:   http.StatusBadRequest,
		Code:     ErrCodeValidation,
		Resource: ResourceTenant,
		Errors:   "Invalid input",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		writeErrorResponse(rr, details)
	}
}
