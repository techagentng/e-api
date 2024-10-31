package errors

import (
	"fmt"

	validator "github.com/go-playground/validator/v10"
)

type FieldError struct {
	err validator.FieldError
}

func (q *FieldError) String() string {
	// var sb strings.Builder
	err := q.err

	return fmt.Sprintf("%s is invalid: %v", err.Field(), err.Value())

}

// NewFieldError returns a field error
func NewFieldError(err validator.FieldError) *FieldError {
	return &FieldError{err: err}
}

// ValidationError defines error that occur due to validation
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Message)
}
