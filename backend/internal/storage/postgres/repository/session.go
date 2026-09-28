package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/korjeek/pinger/backend/internal/domain"
	"github.com/korjeek/pinger/backend/internal/storage/postgres"
	"github.com/korjeek/pinger/backend/internal/storage/postgres/db"
)

type PgSessionRepository struct {
	queries *db.Queries
}

func NewPgSessionRepository(queries *db.Queries) *PgSessionRepository {
	return &PgSessionRepository{queries: queries}
}

func (r *PgSessionRepository) CreateSession(ctx context.Context, session domain.Session) error {
	err := r.queries.CreateSession(ctx, db.CreateSessionParams{
		ID:        session.ID,
		UserID:    session.UserID,
		TokenHash: session.TokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: session.ExpiresAt},
	})

	return postgres.FromPgError(err)
}

func (r *PgSessionRepository) GetSessionByHash(ctx context.Context, hash []byte) (domain.Session, error) {
	session, err := r.queries.GetSessionByHash(ctx, hash)
	if err != nil {
		return domain.Session{}, postgres.FromPgError(err)
	}
	return fromSession(session), nil
}

func (r *PgSessionRepository) RotateSession(ctx context.Context, oldID, newID uuid.UUID) error {
	return r.queries.RotateSession(ctx, db.RotateSessionParams{
		ID:         oldID,
		ReplacedBy: pgtype.UUID{Bytes: newID, Valid: true},
	})
}

func (r *PgSessionRepository) RevokeSession(ctx context.Context, id uuid.UUID) error {
	return postgres.FromPgError(r.queries.RevokeSession(ctx, id))
}

func (r *PgSessionRepository) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	return postgres.FromPgError(r.queries.RevokeAllUserSessions(ctx, userID))
}

func fromSession(s db.Session) domain.Session {
	out := domain.Session{
		ID:        s.ID,
		UserID:    s.UserID,
		TokenHash: s.TokenHash,
	}

	if s.IssuedAt.Valid {
		out.IssuedAt = s.IssuedAt.Time
	}
	if s.ExpiresAt.Valid {
		out.ExpiresAt = s.ExpiresAt.Time
	}
	if s.RevokedAt.Valid {
		out.RevokedAt = &s.RevokedAt.Time
	}
	if s.ReplacedBy.Valid {
		out.ReplacedBy = new(uuid.UUID(s.ReplacedBy.Bytes))
	}

	return out
}
