package api

import (
	"context"

	"github.com/google/uuid"
	"github.com/korjeek/pinger/backend/pkg/apperr"
)

type ctxKey int

const (
	ctxKeyUserID ctxKey = iota
	ctxKeyRefreshToken
	ctxKeyCookieQueue
)

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ctxKeyUserID).(uuid.UUID)
	return id, ok
}

func MustUserID(ctx context.Context) (uuid.UUID, error) {
	if id, ok := UserIDFromContext(ctx); ok {
		return id, nil
	}
	return uuid.Nil, apperr.Unauthorized("user is not authenticated")
}

func RefreshTokenFromContext(ctx context.Context) (string, error) {
	v, ok := ctx.Value(ctxKeyRefreshToken).(string)
	if !ok || v == "" {
		return "", apperr.Unauthorized("refresh token is missing")
	}
	return v, nil
}
