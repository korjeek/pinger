package apperr

import (
	"errors"
	"fmt"
)

type Code string

const (
	CodeBadRequest       Code = "BAD_REQUEST"
	CodeValidationFailed Code = "VALIDATION_FAILED"
	CodeUnauthorized     Code = "UNAUTHORIZED"
	CodeForbidden        Code = "FORBIDDEN"
	CodeNotFound         Code = "NOT_FOUND"
	CodeConflict         Code = "CONFLICT"
	CodeInternal         Code = "INTERNAL_ERROR"
)

type Error struct {
	Code    Code
	Message string
	Status  int
	Details map[string]any

	cause error
}

func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.cause
}

func (e *Error) Is(target error) bool {
	var t *Error
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

func (e *Error) WithCause(err error) *Error {
	cp := *e
	cp.cause = err
	return &cp
}

func (e *Error) WithMessage(msg string) *Error {
	cp := *e
	cp.Message = msg
	return &cp
}

func (e *Error) WithDetails(d map[string]any) *Error {
	cp := *e
	cp.Details = d
	return &cp
}

func (e *Error) WithDetail(key string, value any) *Error {
	cp := *e
	cp.Details = make(map[string]any, len(e.Details)+1)
	for k, v := range e.Details {
		cp.Details[k] = v
	}
	cp.Details[key] = value
	return &cp
}

func New(code Code, message string) *Error {
	return &Error{Code: code, Message: message}
}

func BadRequest(msg string) *Error {
	return New(CodeBadRequest, msg)
}

func Validation(msg string, details map[string]any) *Error {
	return New(CodeValidationFailed, msg).WithDetails(details)
}

func Unauthorized(msg string) *Error {
	return New(CodeUnauthorized, msg)
}

func Forbidden(msg string) *Error {
	return New(CodeForbidden, msg)
}

func NotFound(msg string) *Error {
	return New(CodeNotFound, msg)
}

func Conflict(msg string) *Error {
	return New(CodeConflict, msg)
}

func Internal(err error) *Error {
	return New(CodeInternal, "internal server error").
		WithCause(err)
}

func As(err error) (*Error, bool) {
	if ae, ok := errors.AsType[*Error](err); ok {
		return ae, true
	}
	return nil, false
}

func From(err error) *Error {
	if err == nil {
		return nil
	}
	if ae, ok := As(err); ok {
		return ae
	}
	return Internal(err)
}

func IsCode(err error, code Code) bool {
	ae, ok := As(err)
	return ok && ae.Code == code
}
