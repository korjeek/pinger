package api

import (
	"context"
	"net/http"
	"time"

	"github.com/korjeek/pinger/backend/internal/dto"
	"github.com/korjeek/pinger/backend/internal/usecase"
	"github.com/oapi-codegen/runtime/types"
)

type Handler struct {
	userService usecase.UserService
	authService usecase.AuthService
}

func (h *Handler) CreateUser(ctx context.Context, request CreateUserRequestObject) (CreateUserResponseObject, error) {
	body := request.Body
	input := dto.CreateUserInput{
		Email:    string(body.Email),
		Password: body.Password,
	}

	output, err := h.userService.CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}

	return CreateUser201JSONResponse{
		Email: types.Email(output.Email),
		Id:    output.Id,
	}, nil
}

func (h *Handler) LoginUser(ctx context.Context, request LoginUserRequestObject) (LoginUserResponseObject, error) {
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
		Name:     "refresh_token",
		Value:    refresh.Payload,
		Path:     "/api/auth",
		Expires:  refresh.ExpiresAt,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	return LoginUser200JSONResponse{
		Body: AuthResponse{
			AccessToken: access.Payload,
			TokenType:   "Bearer",
			ExpiresIn:   new(int(time.Until(pair.AccessToken.ExpiresAt).Seconds())),
		},
		Headers: LoginUser200ResponseHeaders{
			SetCookie: new(cookie.String()),
		},
	}, nil
}

func (h *Handler) LogoutUser(ctx context.Context, request LogoutUserRequestObject) (LogoutUserResponseObject, error) {
	cookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/auth",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	return LogoutUser204Response{
		Headers: LogoutUser204ResponseHeaders{
			SetCookie: new(cookie.String()),
		},
	}, nil
}

func (h *Handler) RefreshToken(ctx context.Context, request RefreshTokenRequestObject) (RefreshTokenResponseObject, error) {
	return RefreshToken200JSONResponse{
		Body: AuthResponse{
			AccessToken: "",
			ExpiresIn:   nil,
			TokenType:   "",
		},
		Headers: RefreshToken200ResponseHeaders{
			SetCookie: ,
		}
	}, nil
}
