package middleware

import (
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/korjeek/pinger/backend/internal/usecase/token"
	"github.com/labstack/echo/v4"
)

var (
	ErrNoAuthHeader      = errors.New("authorization header is missing")
	ErrInvalidAuthHeader = errors.New("authorization header is malformed")
	ErrClaimsInvalid     = errors.New("provided claims do not match expected scopes")
)

const (
	authHeaderKey    = "Authorization"
	bearerPrefix     = "Bearer "
	ContextKeyUserID = "userID"
)

func JwtAuthMiddleware(manger token.Generator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			authHeader := ctx.Request().Header.Get(authHeaderKey)
			if authHeader == "" {
				return ErrNoAuthHeader
			}

			if !strings.HasPrefix(authHeader, bearerPrefix) {
				return ErrInvalidAuthHeader
			}

			jws := strings.TrimPrefix(authHeader, bearerPrefix)

			id, err := manger.Validate(jws)

			ctx.Set(ContextKeyUserID, id)

			return next(ctx)
		}
	}
}
