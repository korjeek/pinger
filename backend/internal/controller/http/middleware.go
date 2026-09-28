package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/korjeek/pinger/backend/internal/dto"
	"github.com/korjeek/pinger/backend/pkg/apperr"
	"github.com/labstack/echo/v4"
)

const (
	RefreshCookieName = "refresh_token"
	traceIDEchoKey    = "trace_id"
)

func TraceIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.Request().Header.Get("X-Request-Id")
			if id == "" {
				if u, err := uuid.NewV7(); err == nil {
					id = u.String()
				}
			}
			if id != "" {
				c.Response().Header().Set("X-Request-Id", id)
				c.Set(traceIDEchoKey, id)
			}
			return next(c)
		}
	}
}

func TraceIDFromEcho(c echo.Context) string {
	v, _ := c.Get(traceIDEchoKey).(string)
	return v
}

type AccessTokenParser interface {
	ParseAccessToken(token string) (uuid.UUID, error)
}

func AuthMiddleware(parser AccessTokenParser) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			raw := c.Request().Header.Get(echo.HeaderAuthorization)
			if raw == "" {
				return apperr.Unauthorized("authorization header is missing")
			}

			const prefix = "Bearer "
			if len(raw) <= len(prefix) || !strings.EqualFold(raw[:len(prefix)], prefix) {
				return apperr.Unauthorized("authorization header must be 'Bearer <token>'")
			}

			userID, err := parser.ParseAccessToken(strings.TrimSpace(raw[len(prefix):]))
			if err != nil {
				return apperr.Unauthorized("invalid or expired access token").WithCause(err)
			}

			ctx := context.WithValue(c.Request().Context(), ctxKeyUserID, userID)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func RefreshTokenCookieMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cookie, err := c.Cookie(RefreshCookieName); err == nil && cookie.Value != "" {
				ctx := context.WithValue(c.Request().Context(), ctxKeyRefreshToken, cookie.Value)
				c.SetRequest(c.Request().WithContext(ctx))
			}
			return next(c)
		}
	}
}

type cookieConfig struct {
	Path     string
	Domain   string
	Secure   bool
	SameSite http.SameSite
}

type cookieQueue struct {
	set   *dto.Token
	clear bool
}

func RefreshCookieWriterMiddleware(cfg cookieConfig) StrictMiddlewareFunc {
	return func(f StrictHandlerFunc, _ string) StrictHandlerFunc {
		return func(ctx echo.Context, request any) (any, error) {
			q := &cookieQueue{}
			reqCtx := context.WithValue(ctx.Request().Context(), ctxKeyCookieQueue, q)
			ctx.SetRequest(ctx.Request().WithContext(reqCtx))

			resp, err := f(ctx, request)
			if err != nil {
				return resp, err
			}

			switch {
			case q.set != nil:
				http.SetCookie(ctx.Response(), newRefreshCookie(cfg, *q.set))
			case q.clear:
				http.SetCookie(ctx.Response(), expiredRefreshCookie(cfg))
			}
			return resp, nil
		}
	}
}

func QueueSetRefreshCookie(ctx context.Context, token dto.Token) {
	if q, ok := ctx.Value(ctxKeyCookieQueue).(*cookieQueue); ok {
		q.set = &token
	}
}

func QueueClearRefreshCookie(ctx context.Context) {
	if q, ok := ctx.Value(ctxKeyCookieQueue).(*cookieQueue); ok {
		q.clear = true
	}
}

func newRefreshCookie(cfg cookieConfig, token dto.Token) *http.Cookie {
	return &http.Cookie{
		Name:     RefreshCookieName,
		Value:    token.Payload,
		Path:     cfg.Path,
		Domain:   cfg.Domain,
		Expires:  token.ExpiresAt,
		MaxAge:   int(time.Until(token.ExpiresAt).Seconds()),
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
	}
}

func expiredRefreshCookie(cfg cookieConfig) *http.Cookie {
	return &http.Cookie{
		Name:     RefreshCookieName,
		Value:    "",
		Path:     cfg.Path,
		Domain:   cfg.Domain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
	}
}
