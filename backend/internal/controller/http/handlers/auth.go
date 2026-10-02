package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/korjeek/pinger/backend/internal/controller/http/api"
	"github.com/korjeek/pinger/backend/internal/dto"
	"github.com/korjeek/pinger/backend/internal/usecase"
	"github.com/korjeek/pinger/backend/pkg/apperr"
	middleware "github.com/oapi-codegen/echo-middleware"
)

const (
	refreshTokenCookieName = "refresh_token"
	refreshTokenCookiePath = "/auth"
	refreshTokenType       = "Bearer"
)

type AuthHandler struct {
	auth usecase.AuthService
}

func (h *AuthHandler) LoginUser(ctx context.Context, request api.LoginUserRequestObject) (api.LoginUserResponseObject, error) {
	body := request.Body
	pair, err := h.auth.LoginUser(ctx, dto.LoginUserInput{
		Email:    string(body.Email),
		Password: body.Password,
	})
	if err != nil {
		return nil, err
	}

	access := pair.AccessToken
	refresh := pair.RefreshToken

	cookie := &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    refresh.Payload,
		Path:     refreshTokenCookiePath,
		Expires:  refresh.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	return api.LoginUser200JSONResponse{
		Body: api.AuthResponse{
			AccessToken: access.Payload,
			TokenType:   refreshTokenType,
			ExpiresIn:   new(int(time.Until(access.ExpiresAt).Seconds())),
		},
		Headers: api.LoginUser200ResponseHeaders{
			SetCookie: new(cookie.String()),
		},
	}, nil
}

func (h *AuthHandler) LogoutUser(ctx context.Context, _ api.LogoutUserRequestObject) (api.LogoutUserResponseObject, error) {
	echoCtx := middleware.GetEchoContext(ctx)
	if echoCtx == nil {
		return nil, apperr.NewUndefined(nil).WithMessage("missing echo context")
	}

	ck, err := echoCtx.Cookie(refreshTokenCookieName)
	if err != nil || ck.Value != "" {
		return nil, apperr.NewUnauthorized(nil).WithMessage("refresh token cookie is missing")
	}

	if err := h.auth.LogoutUser(ctx, ck.Value); err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    "",
		Path:     refreshTokenCookiePath,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	return api.LogoutUser204Response{
		Headers: api.LogoutUser204ResponseHeaders{
			SetCookie: new(cookie.String()),
		},
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, _ api.RefreshTokenRequestObject) (api.RefreshTokenResponseObject, error) {
	echoCtx := middleware.GetEchoContext(ctx)
	if echoCtx == nil {
		return nil, apperr.NewUndefined(nil).WithMessage("missing echo context")
	}

	ck, err := echoCtx.Cookie(refreshTokenCookieName)
	if err != nil || ck.Value != "" {
		return nil, apperr.NewUnauthorized(nil).WithMessage("refresh token cookie is missing")
	}

	pair, err := h.auth.RefreshToken(ctx, ck.Value)
	if err != nil && !apperr.IsCode(err, apperr.Forbidden) && !apperr.IsCode(err, apperr.Unauthorized) {
		return nil, err
	}

	access := pair.AccessToken
	refresh := pair.RefreshToken

	cookie := &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    refresh.Payload,
		Path:     refreshTokenCookiePath,
		Expires:  refresh.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	return api.RefreshToken200JSONResponse{
		Body: api.AuthResponse{
			AccessToken: access.Payload,
			TokenType:   refreshTokenType,
			ExpiresIn:   new(int(time.Until(access.ExpiresAt).Seconds())),
		},
		Headers: api.RefreshToken200ResponseHeaders{
			SetCookie: new(cookie.String()),
		},
	}, nil
}
