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
	GetSessionByHash(ctx context.Context, tokenHash []byte) (domain.Session, error)
	RevokeSession(ctx context.Context, id uuid.UUID) error
	RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error
	RotateSession(ctx context.Context, oldID uuid.UUID, newSession domain.Session) error
}

type TokenGenerator interface {
	GenerateAccessToken(userId uuid.UUID) (dto.Token, error)
	GenerateRefreshToken() (dto.Token, error)
	GenerateTokenPair(userId uuid.UUID) (dto.TokenPair, error)
}

type TokenHasher interface {
	Hash(data string) ([]byte, error)
}

type Clock interface {
	Now() time.Time
}

type AuthService struct {
	user    UserRepository
	session SessionRepository

	tokenGen TokenGenerator

	pwHasher    PasswordHasher
	tokenHasher TokenHasher

	clock Clock
}

func NewAuthService(
	user UserRepository,
	session SessionRepository,
	tokenGen TokenGenerator,
	pwHasher PasswordHasher,
	tokenHasher TokenHasher,
	clock Clock,
) *AuthService {
	return &AuthService{
		user:        user,
		session:     session,
		tokenGen:    tokenGen,
		pwHasher:    pwHasher,
		tokenHasher: tokenHasher,
		clock:       clock,
	}
}

func (s *AuthService) LoginUser(ctx context.Context, input dto.LoginUserInput) (dto.TokenPair, error) {
	user, err := s.user.GetByEmail(ctx, input.Email)
	if err != nil {
		if apperr.IsCode(err, apperr.NotFound) {
			return dto.TokenPair{}, apperr.NewUnauthorized(err).
				WithMessage("invalid email or password")
		}
		return dto.TokenPair{}, err
	}

	if !s.pwHasher.Equals(user.PasswordHash, input.Password) {
		return dto.TokenPair{}, apperr.NewUnauthorized(nil).
			WithMessage("invalid email or password").
			WithDetail("user_id", user.ID)
	}

	pair, err := s.tokenGen.GenerateTokenPair(user.ID)
	if err != nil {
		return dto.TokenPair{}, apperr.NewUndefined(err).
			WithMessage("generate token pair").
			WithDetail("user_id", user.ID)
	}

	if err := s.persistSession(ctx, user.ID, pair.RefreshToken); err != nil {
		return dto.TokenPair{}, err
	}

	return pair, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, raw string) (dto.TokenPair, error) {
	now := s.clock.Now()

	hash, err := s.tokenHasher.Hash(raw)
	if err != nil {
		return dto.TokenPair{}, apperr.NewUndefined(err).
			WithMessage("hash refresh token")
	}

	session, err := s.session.GetSessionByHash(ctx, hash)
	if err != nil {
		if apperr.IsCode(err, apperr.NotFound) {
			return dto.TokenPair{}, apperr.NewUnauthorized(err).
				WithMessage("invalid refresh token")
		}
		return dto.TokenPair{}, err
	}

	switch {
	case session.IsRotated():
		if revokeErr := s.session.RevokeAllUserSessions(ctx, session.UserID); revokeErr != nil {
			return dto.TokenPair{}, apperr.NewUndefined(revokeErr).
				WithMessage("failed to revoke all sessions after reuse breach detected").
				WithDetail("user_id", session.UserID)
		}
		return dto.TokenPair{}, apperr.NewForbidden(nil).
			WithMessage("refresh token reuse detected").
			WithDetail("user_id", session.UserID).
			WithDetail("compromised_session_id", session.ID)

	case session.IsRevoked():
		return dto.TokenPair{}, apperr.NewUnauthorized(nil).
			WithMessage("refresh token has been revoked").
			WithDetail("user_id", session.UserID).
			WithDetail("session_id", session.ID)

	case session.IsExpired(now):
		return dto.TokenPair{}, apperr.NewUnauthorized(nil).
			WithMessage("refresh token has expired").
			WithDetail("user_id", session.UserID).
			WithDetail("session_id", session.ID).
			WithDetail("expired_at", session.ExpiresAt)
	}

	pair, err := s.tokenGen.GenerateTokenPair(session.UserID)
	if err != nil {
		return dto.TokenPair{}, apperr.NewUndefined(err).
			WithMessage("generate token pair").
			WithDetail("user_id", session.UserID)
	}

	newHash, err := s.tokenHasher.Hash(pair.RefreshToken.Payload)
	if err != nil {
		return dto.TokenPair{}, apperr.NewUndefined(err).
			WithMessage("hash refresh token").
			WithDetail("user_id", session.UserID)
	}

	newID, err := uuid.NewV7()
	if err != nil {
		return dto.TokenPair{}, apperr.NewUndefined(err).
			WithMessage("generate session id").
			WithDetail("user_id", session.UserID)
	}

	newSession := domain.NewSession(newID, session.UserID, newHash, pair.RefreshToken.ExpiresAt)
	if err := s.session.RotateSession(ctx, session.ID, newSession); err != nil {
		return dto.TokenPair{}, err
	}

	return pair, nil
}

func (s *AuthService) LogoutUser(ctx context.Context, raw string) error {
	hash, err := s.tokenHasher.Hash(raw)
	if err != nil {
		return apperr.NewUndefined(err).
			WithMessage("hash refresh token")
	}

	session, err := s.session.GetSessionByHash(ctx, hash)
	if err != nil {
		if apperr.IsCode(err, apperr.NotFound) {
			return nil
		}
		return err
	}

	if session.IsRevoked() {
		return nil
	}

	return s.session.RevokeSession(ctx, session.ID)
}

func (s *AuthService) persistSession(ctx context.Context, userID uuid.UUID, refresh dto.Token) error {
	hash, err := s.tokenHasher.Hash(refresh.Payload)
	if err != nil {
		return apperr.NewUndefined(err).
			WithMessage("hash refresh token").
			WithDetail("user_id", userID)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return apperr.NewUndefined(err).
			WithMessage("generate session id").
			WithDetail("user_id", userID)
	}

	session := domain.NewSession(id, userID, hash, refresh.ExpiresAt)
	return s.session.CreateSession(ctx, session)
}
