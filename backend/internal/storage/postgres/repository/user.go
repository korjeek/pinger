package repository

import (
	"context"

	"github.com/korjeek/pinger/backend/internal/domain"
	"github.com/korjeek/pinger/backend/internal/storage/postgres"
	"github.com/korjeek/pinger/backend/internal/storage/postgres/db"
)

type PgUserRepository struct {
	queries *db.Queries
}

func NewPgUserRepository(queries *db.Queries) *PgUserRepository {
	return &PgUserRepository{queries: queries}
}

func (r *PgUserRepository) CreateUser(ctx context.Context, user domain.User) error {
	err := r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	})

	return postgres.FromPgError(err)
}

func (r *PgUserRepository) GetUserByEmail(ctx context.Context, email string) (user domain.User, err error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return user, postgres.FromPgError(err)
	}

	return domain.User(dbUser), nil
}
