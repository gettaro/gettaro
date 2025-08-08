package errors

import "net/http"

// BadRequestError represents a 400 Bad Request error
type BadRequestError struct {
	Message string
}

// Error implements the error interface
func (e *BadRequestError) Error() string {
	return e.Message
}

// StatusCode returns the HTTP status code for this error
func (e *BadRequestError) StatusCode() int {
	return http.StatusBadRequest
}

// NewBadRequestError creates a new BadRequestError with the given message
func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{
		Message: message,
	}
}
