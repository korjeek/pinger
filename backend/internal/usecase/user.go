package usecase

import (
	"context"
	"fmt"

	"github.com/korjeek/pinger/backend/internal/domain"
	"github.com/korjeek/pinger/backend/internal/dto"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session domain.Session) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Equals(hashedPassword, password string) bool
}

type TokenHasher interface {
	Hash(string) ([]byte, error)
}

type UserService struct {
	user    UserRepository
	session SessionRepository

	generator TokenGenerator

	pwHasher    PasswordHasher
	tokenHasher TokenHasher
}

func (s *UserService) CreateUser(ctx context.Context, input dto.CreateUserInput) (res dto.CreateUserOutput, err error) {
	passwordHash, err := s.pwHasher.Hash(input.Password)
	if err != nil {
		return res, fmt.Errorf("hash password: %w", err)
	}

	user, err := domain.NewUser(input.Email, passwordHash)
	if err != nil {
		return res, fmt.Errorf("create domain user: %w", err)
	}

	if err = s.user.CreateUser(ctx, *user); err != nil {
		return res, fmt.Errorf("db create user: %w", err)
	}

	return dto.CreateUserOutput{
		Id:    user.ID,
		Email: user.Email,
	}, nil
}

func (s *UserService) LoginUser(ctx context.Context, input dto.LoginUserInput) (pair dto.TokenPair, err error) {
	user, err := s.user.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return pair, err
	}

	if !s.pwHasher.Equals(user.PasswordHash, input.Password) {
		return pair, fmt.Errorf("invalid password")
	}

	tokens, err := s.generator.GenerateTokenPair(user.ID)
	if err != nil {
		return pair, err
	}

	tokenHash, err := s.tokenHasher.Hash(tokens.RefreshToken.Payload)
	if err != nil {
		return pair, err
	}

	session, err := domain.NewSession(user.ID, tokenHash, tokens.RefreshToken.ExpiresAt)
	if err != nil {
		return pair, err
	}

	if err = s.session.CreateSession(ctx, *session); err != nil {
		return pair, err
	}

	return tokens, nil
}
