package utils

import "github.com/korjeek/pinger/backend/pkg/apperr"

func CodeInt(code apperr.Code) int {
	switch code {
	case apperr.Undefined:
		return 500
	case apperr.NotFound:
		return 404
	case apperr.AlreadyExists:
		return 409
	case apperr.Unauthorized:
		return 401
	case apperr.Forbidden:
		return 403
	default:
		return 500
	}
}
