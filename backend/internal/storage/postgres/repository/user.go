package repository

import (
	"context"
	"fmt"

	"github.com/korjeek/pinger/backend/internal/domain"
	"github.com/korjeek/pinger/backend/internal/storage/postgres/db"
)

type PgUserRepository struct {
	queries db.Queries
}

func (r *PgUserRepository) CreateUser(ctx context.Context, user domain.User) error {
	return r.queries.CreateUser(ctx, db.CreateUserParams{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	})
}

func (r *PgUserRepository) GetUserByEmail(ctx context.Context, email string) (user domain.User, err error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return user, fmt.Errorf("get user by id: %w", err)
	}

	return domain.User(dbUser), nil
}
