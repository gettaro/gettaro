package errors

import "fmt"

// ErrDuplicateConflict is returned when attempting to create a resource with a unique field that already exists
type ErrDuplicateConflict struct {
	Resource string
	Field    string
	Value    string
}

func (e *ErrDuplicateConflict) Error() string {
	return fmt.Sprintf("%s with %s '%s' already exists", e.Resource, e.Field, e.Value)
}

// IsDuplicateConflict checks if an error is an ErrDuplicateConflict
func IsDuplicateConflict(err error) bool {
	_, ok := err.(*ErrDuplicateConflict)
	return ok
}
