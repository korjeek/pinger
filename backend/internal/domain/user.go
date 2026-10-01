package domain

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrSessionNotFound   = errors.New("session not found")
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
}

func NewUser(id uuid.UUID, email string, passwordHash string) User {
	return User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
	}
}
