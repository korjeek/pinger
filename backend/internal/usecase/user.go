package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/korjeek/pinger/backend/internal/domain"
	"github.com/korjeek/pinger/backend/internal/dto"
	"github.com/korjeek/pinger/backend/pkg/apperr"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Equals(hashedPassword, password string) bool
}

type UserService struct {
	user     UserRepository
	pwHasher PasswordHasher
}

func (s *UserService) CreateUser(ctx context.Context, input dto.CreateUserInput) (dto.CreateUserOutput, error) {
	var output dto.CreateUserOutput

	passwordHash, err := s.pwHasher.Hash(input.Password)
	if err != nil {
		return output, apperr.NewUndefined(err).
			WithMessage("failed to process user security credentials").
			WithDetail("email", input.Email)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return output, apperr.NewUndefined(err).
			WithMessage("failed to generate unique user identifier")
	}

	user := domain.NewUser(id, input.Email, passwordHash)
	if err = s.user.Create(ctx, user); err != nil {
		return output, err
	}

	return dto.CreateUserOutput{
		Id:    user.ID,
		Email: user.Email,
	}, nil
}
