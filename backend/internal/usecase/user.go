package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/korjeek/pinger/backend/internal/domain"
	"github.com/korjeek/pinger/backend/internal/dto"
	"github.com/korjeek/pinger/backend/pkg/apperr"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
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
	passwordHash, err := s.pwHasher.Hash(input.Password)
	if err != nil {
		return dto.CreateUserOutput{}, apperr.Internal(err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return dto.CreateUserOutput{}, apperr.Internal(err)
	}

	user := domain.NewUser(id, input.Email, passwordHash)
	if err = s.user.CreateUser(ctx, *user); err != nil {
		return dto.CreateUserOutput{}, err
	}

	return dto.CreateUserOutput{
		Id:    user.ID,
		Email: user.Email,
	}, nil
}
