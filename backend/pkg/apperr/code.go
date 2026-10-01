package apperr

import (
	"fmt"
)

type Code int

const (
	Undefined Code = iota
	NotFound
	AlreadyExists
	Unauthorized
	Forbidden
)

func (c Code) String() string {
	switch c {
	case Undefined:
		return "UNDEFINED"
	case NotFound:
		return "NOT_FOUND"
	case AlreadyExists:
		return "ALREADY_EXISTS"
	case Unauthorized:
		return "UNAUTHORIZED"
	case Forbidden:
		return "FORBIDDEN"
	default:
		return fmt.Sprintf("UNKNOWN_CODE(%d)", c)
	}
}
