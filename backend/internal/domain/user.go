package domain

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	ErrFailedGenerateUUID = errors.New("failed to generate UUIDv7")
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
}

func NewUser(email string, passwordHash string) (*User, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("new user: %w", err)
	}

	return &User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}
