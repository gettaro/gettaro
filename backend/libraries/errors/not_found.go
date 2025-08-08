package errors

import "net/http"

// NotFoundError represents a 404 Not Found error
type NotFoundError struct {
	Message string
}

// Error implements the error interface
func (e *NotFoundError) Error() string {
	return e.Message
}

// StatusCode returns the HTTP status code for this error
func (e *NotFoundError) StatusCode() int {
	return http.StatusNotFound
}

// NewNotFoundError creates a new NotFoundError with the given message
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		Message: message,
	}
}
