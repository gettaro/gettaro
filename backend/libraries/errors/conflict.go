package errors

import "net/http"

// ConflictError represents a 409 Conflict error
type ConflictError struct {
	Message string
}

// Error implements the error interface
func (e *ConflictError) Error() string {
	return e.Message
}

// StatusCode returns the HTTP status code for this error
func (e *ConflictError) StatusCode() int {
	return http.StatusConflict
}

// NewConflictError creates a new ConflictError with the given message
func NewConflictError(message string) *ConflictError {
	return &ConflictError{
		Message: message,
	}
}
