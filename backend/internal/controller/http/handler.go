package http

import (
	"context"

	"github.com/korjeek/pinger/backend/internal/controller/http/api"
	"github.com/korjeek/pinger/backend/internal/dto"
	"github.com/korjeek/pinger/backend/internal/usecase"
	"github.com/oapi-codegen/runtime/types"
)

type Handler struct {
	userService usecase.UserService
}

func (h *Handler) CreateUser(ctx context.Context, request api.CreateUserRequestObject) (api.CreateUserResponseObject, error) {
	body := request.Body
	input := dto.CreateUserInput{
		Email:    string(body.Email),
		Password: body.Password,
	}

	output, err := h.userService.CreateUser(ctx, input)
	if err != nil {
		return nil, err
	}

	return api.CreateUser201JSONResponse{
		Email: types.Email(output.Email),
		Id:    output.Id,
	}, nil
}
