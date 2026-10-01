package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/korjeek/pinger/backend/internal/domain"
	"github.com/korjeek/pinger/backend/internal/storage/postgres/db"
	"github.com/korjeek/pinger/backend/internal/storage/postgres/mapper"
	"github.com/korjeek/pinger/backend/pkg/apperr"
)

type PgUserRepository struct {
	queries *db.Queries
}

func NewPgUserRepository(queries *db.Queries) *PgUserRepository {
	return &PgUserRepository{queries: queries}
}

func (r *PgUserRepository) Create(ctx context.Context, user domain.User) error {
	if err := r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}); err != nil {
		if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok && pgErr.Code == pgerrcode.UniqueViolation {
			return apperr.NewAlreadyExists(err).
				WithMessage("user with this email already exists").
				WithDetail("email", user.Email).
				WithDetail("user_id", user.ID)
		}

		return apperr.NewUndefined(err).
			WithMessage("failed to create user").
			WithDetail("email", user.Email)
	}

	return nil
}

func (r *PgUserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User

	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, apperr.NewNotFound(err).
				WithMessage("user not found by email").
				WithDetail("email", email)
		}

		return user, apperr.NewUndefined(err).
			WithMessage("failed to get user by email").
			WithDetail("email", email)
	}

	return mapper.ToDomainUser(dbUser), nil
}
