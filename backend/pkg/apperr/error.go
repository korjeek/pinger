// Package apperr provides a type-safe, fluent error wrapping mechanism
// designed for backend applications to communicate domain-specific errors.
package apperr

import (
	"errors"
	"fmt"
)

// Error represents a structured domain error.
// It contains a machine-readable Code, an optional human-readable Message,
// key-value Details for additional context, and the underlying Cause.
type Error struct {
	Code    Code
	Message string
	Details map[string]any
	Cause   error
}

// NewError creates a new base Error with the given Code and underlying cause.
func NewError(code Code, cause error) Error {
	return Error{
		Code:  code,
		Cause: cause,
	}
}

// NewUndefined helper function initialized with an Undefined error code.
func NewUndefined(cause error) Error {
	return NewError(Undefined, cause)
}

// NewNotFound helper function initialized with a NotFound error code.
func NewNotFound(cause error) Error {
	return NewError(NotFound, cause)
}

// NewAlreadyExists helper function initialized with an AlreadyExists error code.
func NewAlreadyExists(cause error) Error {
	return NewError(AlreadyExists, cause)
}

// NewUnauthorized helper function initialized with an Unauthorized error code.
func NewUnauthorized(cause error) Error {
	return NewError(Unauthorized, cause)
}

// NewForbidden helper function initialized with an Forbidden error code.
func NewForbidden(cause error) Error {
	return NewError(Forbidden, cause)
}

// Error implements the standard error interface.
// It formats the output string by dynamically omitting empty Message field.
func (e Error) Error() string {
	res := fmt.Sprintf("[%s]", e.Code)
	if e.Message != "" {
		res += fmt.Sprintf(": %s", e.Message)
	}
	return res
}

// Unwrap returns the underlying cause of the error.
// It enables compatibility with the standard library's errors.Unwrap, errors.Is, and errors.As.
func (e Error) Unwrap() error {
	return e.Cause
}

// WithMessage attaches a custom human-readable message to the error.
// It returns a modified copy of the Error, supporting fluent chaining.
func (e Error) WithMessage(msg string) Error {
	e.Message = msg
	return e
}

// WithDetail adds a single key-value pair to the error's contextual details.
// It isolates the map mutation by copying existing details, supporting fluent chaining.
func (e Error) WithDetail(key string, value any) Error {
	newDetails := make(map[string]any, len(e.Details)+1)
	for k, v := range e.Details {
		newDetails[k] = v
	}

	newDetails[key] = value
	e.Details = newDetails
	return e
}

// As extracts the first *Error from the error chain.
func As(err error) (*Error, bool) {
	return errors.AsType[*Error](err)
}

// IsCode reports whether the error chain contains an *Error with the given Code.
func IsCode(err error, code Code) bool {
	e, ok := As(err)
	return ok && e.Code == code
}
