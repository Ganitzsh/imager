package service

import (
	"fmt"
)

type ServiceError struct {
	Message string
	Code    string
}

func NewServiceError(code, message string) *ServiceError {
	return &ServiceError{
		Message: message,
		Code:    code,
	}
}

func (e ServiceError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

var (
	ErrInvalidInput = NewServiceError(
		"invalid_input",
		"Invalid input",
	)
	ErrInvalidConfig = NewServiceError(
		"invalid_configuration",
		"Invalid configuration, check again",
	)
	ErrUnknownStoreType = NewServiceError(
		"unknown_store_type",
		"The given store type is not recognized. You can use 'redis'",
	)
	ErrResourceNotFound = NewServiceError(
		"resource_not_found",
		"The given resource could not be found",
	)
	ErrResourceAlreadyExists = NewServiceError(
		"resource_already_exists",
		"A resource with the same values already exists",
	)
	ErrTokenInvalid = NewServiceError(
		"token_invalid",
		"The given token is invalid",
	)
	ErrTokenExpired = NewServiceError(
		"token_expired",
		"The given token is expired",
	)
	ErrInvalidTLSconfiguration = NewServiceError(
		"config_invalid_tls",
		"Invalid TLS configuration",
	)
	ErrFileSizeExceeded = NewServiceError(
		"file_size_exceeded",
		"File size exceeds the maximum size difined by the server",
	)
	ErrUnsupportedFormat = NewServiceError(
		"unsupported_format",
		"Unsupported image format",
	)
	ErrInternalError = NewServiceError(
		"internal_error",
		"An internal error occured",
	)
	ErrDivideByZeroError = NewServiceError(
		"internal_error",
		"An internal error occured",
	)
)
