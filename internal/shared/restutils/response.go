package restutils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"flowcargo/internal/shared/logger"
)

// APIResponse defines the standard structure for API responses.
// @Description Standard API response wrapper containing success status, message, data, error details, and metadata
type APIResponse[T any] struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty" example:"Operation completed successfully"`
	Data    *T          `json:"data"`
	Error   ErrorDetail `json:"error"`
	Meta    Metadata    `json:"meta"`
}

// SuccessResponse represents a successful API response wrapper.
// @Description Standard success response wrapper with data payload
type SuccessResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data"`
	Error   ErrorDetail `json:"error"`
	Meta    Metadata    `json:"meta"`
}

// ErrorResponse represents an error API response wrapper.
// @Description Standard error response wrapper with error details
type ErrorResponse struct {
	Success bool        `json:"success" example:"false"`
	Message string      `json:"message" example:"Error"`
	Data    interface{} `json:"data"`
	Error   ErrorDetail `json:"error"`
	Meta    Metadata    `json:"meta"`
}

// Metadata holds additional information about the response.
// @Description Response metadata containing timestamp and other contextual information
type Metadata struct {
	Timestamp time.Time `json:"timestamp" example:"2024-01-01T00:00:00Z"`
}

// ErrorDetail provides detailed information about an error.
// @Description Detailed error information including status code, error code, resource, and error messages
type ErrorDetail struct {
	Status   int      `json:"status" example:"400"`
	Code     ErrCode  `json:"code" example:"VALIDATION_ERROR"`
	Resource Resource `json:"resource" example:"tenant"`
	Errors   string   `json:"errors" example:"validation failed: name is required"`
}

func writeJSONResponse(w http.ResponseWriter, status int, response any) {
	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}

// WriteSuccessResponse writes a standardized success response to the http.ResponseWriter.
func WriteSuccessResponse[T any](w http.ResponseWriter, status int, data T, message string) {
	response := APIResponse[T]{
		Success: true,
		Message: message,
		Data:    &data,
		Error:   ErrorDetail{},
		Meta: Metadata{
			Timestamp: time.Now(),
		},
	}
	writeJSONResponse(w, status, response)
}

func writeErrorResponse(w http.ResponseWriter, details ErrorDetail) {
	response := APIResponse[any]{
		Success: false,
		Message: "Error",
		Data:    nil,
		Error:   details,
		Meta:    Metadata{Timestamp: time.Now()},
	}
	writeJSONResponse(w, details.Status, response)
}

// HandleBadRequest writes a standardized bad request error response and logs the error.
func HandleBadRequest(w http.ResponseWriter, err error, logger logger.Logger) {
	logger.Warn("Bad request ", "error ", err.Error())
	details := ErrorDetail{
		Status:   http.StatusBadRequest,
		Code:     ErrCodeValidation,
		Resource: ResourceTenant,
		Errors:   err.Error(),
	}
	writeErrorResponse(w, details)
}

// HandleInternalServerError writes a standardized internal server error response and logs the error.
func HandleInternalServerError(w http.ResponseWriter, err error, resource Resource, logger logger.Logger) {
	logger.Error("Internal server error", "error", err.Error(), "resource", resource)
	details := ErrorDetail{
		Status:   http.StatusInternalServerError,
		Code:     ErrCodeInternal,
		Resource: resource,
		Errors:   err.Error(),
	}
	writeErrorResponse(w, details)
}

// HandleMethodNotAllowed writes a standardized method not allowed error response and logs the attempt.
func HandleMethodNotAllowed(w http.ResponseWriter, method string, url *url.URL, resource Resource, logger logger.Logger) {
	logger.Warn("Method not allowed", "method", method, "url", url, "resource", resource)
	details := ErrorDetail{
		Status:   http.StatusMethodNotAllowed,
		Code:     ErrCodeMethodNotAllowed,
		Resource: resource,
		Errors:   "Method not allowed",
	}
	writeErrorResponse(w, details)
}
