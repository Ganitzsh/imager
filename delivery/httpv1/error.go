package httpv1

import (
	"fmt"
	"net/http"

	"github.com/ganitzsh/12fact/service"
)

type HandlerError struct {
	service.ServiceError
	Status int
}

func NewHandlerError(code, message string, status int) *HandlerError {
	return &HandlerError{
		ServiceError: service.ServiceError{
			Code:    code,
			Message: message,
		},
		Status: status,
	}
}

func (e HandlerError) Error() string {
	return fmt.Sprintf(
		"%d %s: %s", e.Status, http.StatusText(e.Status), e.ServiceError.Error(),
	)
}

var (
	ErrInternalError = NewHandlerError(
		"internal_error",
		"Internal error occured",
		http.StatusInternalServerError,
	)
	ErrInvalidContentType = NewHandlerError(
		"invalid_content_type",
		"Invalid content type for this request",
		http.StatusBadRequest,
	)
	ErrInvalidInput = NewHandlerError(
		"invalid_input",
		"Invalid input. Check the body of your request",
		http.StatusBadRequest,
	)
)
