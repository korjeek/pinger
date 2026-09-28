package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/korjeek/pinger/backend/internal/domain"
	"github.com/korjeek/pinger/backend/internal/dto"
	"github.com/korjeek/pinger/backend/pkg/apperr"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session domain.Session) error
}

type TokenHasher interface {
	Hash(string) ([]byte, error)
}

type AuthService struct {
	user    UserRepository
	session SessionRepository

	tokenGen TokenGenerator

	pwHasher    PasswordHasher
	tokenHasher TokenHasher
}

func (s *AuthService) LoginUser(ctx context.Context, input dto.LoginUserInput) (dto.TokenPair, error) {
	user, err := s.user.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return dto.TokenPair{}, err
	}

	if !s.pwHasher.Equals(user.PasswordHash, input.Password) {
		return dto.TokenPair{}, apperr.Unauthorized("invalid password")
	}

	tokens, err := s.tokenGen.GenerateTokenPair(user.ID)
	if err != nil {
		return dto.TokenPair{}, apperr.Internal(err)
	}

	tokenHash, err := s.tokenHasher.Hash(tokens.RefreshToken.Payload)
	if err != nil {
		return dto.TokenPair{}, apperr.Internal(err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return dto.TokenPair{}, apperr.Internal(err)
	}

	session := domain.NewSession(id, user.ID, tokenHash, tokens.RefreshToken.ExpiresAt)
	if err = s.session.CreateSession(ctx, *session); err != nil {
		return dto.TokenPair{}, err
	}

	return tokens, nil
}

func (s *AuthService) LogoutUser(ctx context.Context) error {

}
