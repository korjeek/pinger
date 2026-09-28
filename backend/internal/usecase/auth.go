package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/korjeek/pinger/backend/internal/domain"
	"github.com/korjeek/pinger/backend/internal/dto"
	"github.com/korjeek/pinger/backend/pkg/apperr"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session domain.Session) error
	GetSessionByHash(ctx context.Context, hash []byte) (domain.Session, error)
	RotateSession(ctx context.Context, oldID, newID uuid.UUID) error
	RevokeSession(ctx context.Context, id uuid.UUID) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
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
		return dto.TokenPair{}, err //TODO: правильно обрабатывать ошибку
	}

	if !s.pwHasher.Equals(user.PasswordHash, input.Password) {
		return dto.TokenPair{}, apperr.Unauthorized("invalid email or password")
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

func (s *AuthService) LogoutUser(ctx context.Context, refreshToken string) error {
	session, err := s.GetSession(ctx, refreshToken)
	if err != nil {
		return apperr.Internal(err) //TODO: правильно обрабатывать ошибку
	}

	return s.session.RevokeSession(ctx, session.ID)
}

func (s *AuthService) GetSession(ctx context.Context, refreshToken string) (domain.Session, error) {
	hash, err := s.tokenHasher.Hash(refreshToken)
	if err != nil {
		return domain.Session{}, apperr.Internal(err)
	}

	return s.session.GetSessionByHash(ctx, hash)
}

func (s *AuthService) RefreshSession(ctx context.Context, refreshToken string) (dto.TokenPair, error) {
	session, err := s.GetSession(ctx, refreshToken)
	if err != nil {
		return dto.TokenPair{}, err //TODO: правильно обрабатывать ошибку
	}

	if session.RevokedAt != nil {
		err = s.session.RevokeAllUserSessions(ctx, session.UserID)
		return dto.TokenPair{}, err //TODO: правильно обрабатывать ошибку
	}

	if !session.IsUsable(time.Now()) {
		return dto.TokenPair{}, err //TODO: правильно обрабатывать ошибку
	}

	tokens, err := s.tokenGen.GenerateTokenPair(session.UserID)
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

	newSession := domain.NewSession(id, session.UserID, tokenHash, tokens.RefreshToken.ExpiresAt)
	if err = s.session.RotateSession(ctx, session.ID, newSession.ID); err != nil {
		return dto.TokenPair{}, err
	}

	return tokens, nil
}
