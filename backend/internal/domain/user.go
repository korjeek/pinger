package domain

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
}

func NewUser(id uuid.UUID, email string, passwordHash string) *User {
	return &User{
		ID:           id,
		Email:        email,
		PasswordHash: passwordHash,
	}
}
