package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/korjeek/pinger/backend/internal/domain"
	"github.com/korjeek/pinger/backend/internal/storage/postgres/db"
	"github.com/korjeek/pinger/backend/internal/storage/postgres/mapper"
	"github.com/korjeek/pinger/backend/pkg/apperr"
)

type PgSessionRepository struct {
	queries *db.Queries
}

func NewPgSessionRepository(queries *db.Queries) *PgSessionRepository {
	return &PgSessionRepository{queries: queries}
}

func (r *PgSessionRepository) Create(ctx context.Context, session domain.Session) error {
	if err := r.queries.CreateSession(ctx, db.CreateSessionParams{
		ID:        session.ID,
		UserID:    session.UserID,
		TokenHash: session.TokenHash,
		ExpiresAt: session.ExpiresAt,
	}); err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return apperr.NewAlreadyExists(err).
				WithMessage("session already exists").
				WithDetail("session_id", session.ID).
				WithDetail("user_id", session.UserID)
		}
		return apperr.NewUndefined(err).
			WithMessage("failed to create session").
			WithDetail("user_id", session.UserID)
	}

	return nil
}

func (r *PgSessionRepository) GetByHash(ctx context.Context, hash []byte) (domain.Session, error) {
	var session domain.Session

	dbSess, err := r.queries.GetSessionByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return session, apperr.NewNotFound(err).
				WithMessage("session not found by provided token hash")
		}
		return session, apperr.NewUndefined(err).
			WithMessage("failed to get session by hash")
	}

	return mapper.ToDomainSession(dbSess), nil
}

func (r *PgSessionRepository) Rotate(ctx context.Context, oldID, newID uuid.UUID) error {
	if err := r.queries.RotateSession(ctx, db.RotateSessionParams{
		ID:         oldID,
		ReplacedBy: &newID,
	}); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperr.NewNotFound(err).
				WithMessage("target session for rotation not found").
				WithDetail("old_session_id", oldID)
		}
		return apperr.NewUndefined(err).
			WithMessage("failed to rotate session").
			WithDetail("old_session_id", oldID).
			WithDetail("new_session_id", newID)
	}

	return nil
}

func (r *PgSessionRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	if err := r.queries.RevokeSession(ctx, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperr.NewNotFound(err).
				WithMessage("session to revoke not found").
				WithDetail("session_id", id)
		}
		return apperr.NewUndefined(err).
			WithMessage("failed to revoke session").
			WithDetail("session_id", id)
	}

	return nil
}

func (r *PgSessionRepository) RevokeAll(ctx context.Context, userID uuid.UUID) error {
	if err := r.queries.RevokeAllUserSessions(ctx, userID); err != nil {
		return apperr.NewUndefined(err).
			WithMessage("failed to revoke all user sessions").
			WithDetail("user_id", userID)
	}

	return nil
}
