package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/korjeek/pinger/backend/internal/domain"
	"github.com/korjeek/pinger/backend/internal/storage/postgres/db"
)

func (r *PgUserRepository) CreateSession(ctx context.Context, session domain.Session) error {
	return r.queries.CreateSession(ctx, db.CreateSessionParams{
		ID:        session.ID,
		UserID:    session.UserID,
		TokenHash: session.TokenHash,
		ExpiresAt: pgtype.Timestamptz{Time: session.ExpiresAt},
	})
}
