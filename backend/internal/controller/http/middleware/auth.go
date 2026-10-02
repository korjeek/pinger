package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const ContextKeyUserID = "userID"

type AccessTokenParser interface {
	Validate(tokenString string) (uuid.UUID, error)
}

var (
	ErrNoAuthHeader      = errors.New("authorization header is missing")
	ErrInvalidAuthHeader = errors.New("authorization header is malformed")
)

func JWTAuth(parser AccessTokenParser) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			raw := c.Request().Header.Get(echo.HeaderAuthorization)
			if raw == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing Authorization header")
			}

			const prefix = "Bearer "
			if len(raw) <= len(prefix) || !strings.EqualFold(raw[:len(prefix)], prefix) {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid Authorization header format")
			}

			userID, err := parser.Validate(raw[len(prefix):])
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired access token")
			}

			c.Set(ContextKeyUserID, userID)
			return next(c)
		}
	}
}
