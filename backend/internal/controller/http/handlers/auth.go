package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/korjeek/pinger/backend/internal/controller/http/api"
	"github.com/korjeek/pinger/backend/internal/dto"
	"github.com/korjeek/pinger/backend/internal/usecase"
)

const (
	refreshTokenCookieName = "refresh_token"
	refreshTokenCookiePath = "/auth"
	refreshTokenType = "Bearer"
)

type AuthHandler struct {
	authService usecase.AuthService
}

func (h *AuthHandler) LoginUser(ctx context.Context, request api.LoginUserRequestObject) (api.LoginUserResponseObject, error) {
	body := request.Body
	input := dto.LoginUserInput{
		Email:    string(body.Email),
		Password: body.Password,
	}

	pair, err := h.authService.LoginUser(ctx, input)
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

func (h *AuthHandler) LogoutUser(ctx context.Context, request api.LogoutUserRequestObject) (api.LogoutUserResponseObject, error) {
	token, _ := ctx.Value("refreshToken").(string)
	if err :=  h.authService.LogoutUser(ctx, token); err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    "",
		Path:     refreshTokenCookiePath,
		MaxAge:   -1,
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

func (h *AuthHandler) RefreshToken(ctx context.Context, request api.RefreshTokenRequestObject) (api.RefreshTokenResponseObject, error) {
	return api.RefreshToken200JSONResponse{
		Body: api.AuthResponse{
			AccessToken: "",
			ExpiresIn:   nil,
			TokenType:   "",
		},
		Headers: api.RefreshToken200ResponseHeaders{
			SetCookie: ,
		}
	}, nil
}
