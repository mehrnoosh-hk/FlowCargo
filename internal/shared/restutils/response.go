package restutils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"flowcargo/internal/shared/logger"
)

type APIResponse[T any] struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    *T           `json:"data"`
	Error   ErrorDetail `json:"error"`
	Meta    Metadata    `json:"meta"`
}

type Metadata struct {
	Timestamp time.Time `json:"timestamp"`
}

type ErrorDetail struct {
	Status   int      `json:"status"`
	Code     ErrCode  `json:"code"`
	Resource Resource `json:"resource"`
	Errors   string   `json:"errors"`
}

func WriteJSONResponse(w http.ResponseWriter, status int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

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
	w.WriteHeader(status)
	WriteJSONResponse(w, status, response)
}

func WriteErrorResponse(w http.ResponseWriter, details ErrorDetail) {
	w.WriteHeader(details.Status)
	response := APIResponse[any]{
		Success: false,
		Message: "Error",
		Data:    nil,
		Error:   details,
		Meta:    Metadata{Timestamp: time.Now()},
	}
	WriteJSONResponse(w, details.Status, response)
}

func HandleBadRequest(w http.ResponseWriter, err error, logger logger.Logger) {
	logger.Warn("Bad request ", "error ", err.Error())
	details := ErrorDetail{
		Status:   http.StatusBadRequest,
		Code:     ErrCodeValidation,
		Resource: ResourceTenant,
		Errors:   err.Error(),
	}
	WriteErrorResponse(w, details)
}

func HandleInternalServerError(w http.ResponseWriter, err error, resource Resource, logger logger.Logger) {
	logger.Error("Internal server error", "error", err.Error(), "resource", resource)
	details := ErrorDetail{
		Status:   http.StatusInternalServerError,
		Code:     ErrCodeInternal,
		Resource: resource,
		Errors:   err.Error(),
	}
	WriteErrorResponse(w, details)
}

func HandleMethodNotAllowed(w http.ResponseWriter, method string, url *url.URL, resource Resource, logger logger.Logger) {
	logger.Warn("Method not allowed", "method", method, "url", url, "resource", resource)
	details := ErrorDetail{
		Status:   http.StatusMethodNotAllowed,
		Code:     ErrCodeMethodNotAllowed,
		Resource: resource,
		Errors:   "Method not allowed",
	}
	WriteErrorResponse(w, details)
}
