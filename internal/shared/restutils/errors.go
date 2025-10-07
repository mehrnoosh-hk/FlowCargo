package restutils

import "errors"

type ErrCode string

const (
	ErrCodeValidation       ErrCode = "VALIDATION_ERROR"
	ErrCodeNotFound         ErrCode = "NOT_FOUND"
	ErrCodeInternal         ErrCode = "INTERNAL_SERVER_ERROR"
	ErrCodeUnauthorized     ErrCode = "UNAUTHORIZED"
	ErrCodeForbidden        ErrCode = "FORBIDDEN"
	ErrCodeConflict         ErrCode = "RESOURCE_CONFLICT"
	ErrCodeDatabase         ErrCode = "DATABASE_ERROR"
	ErrCodeExternalService  ErrCode = "EXTERNAL_SERVICE_ERROR"
	ErrCodeMethodNotAllowed ErrCode = "METHOD_NOT_ALLOWED"
)

type Resource string

const (
	ResourceTenant Resource = "tenant"
)

var (
	ErrResourceConflict = errors.New(string(ErrCodeConflict))
	ErrResourceNotFound = errors.New(string(ErrCodeNotFound))
	ErrUnauthorized = errors.New(string(ErrCodeUnauthorized))
)
