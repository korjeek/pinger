package repository

import (
	"context"

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

	r.queries.RotateSession(ctx, db.RotateSessionParams{
		ID:         session.ID,
		ReplacedBy: pgtype.UUID{},
	})

	return postgres.FromPgError(err)
}
