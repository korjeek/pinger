package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash []byte
	ExpiresAt time.Time
}

func NewSession(id, userId uuid.UUID, tokenHash []byte, expiresAt time.Time) *Session {
	return &Session{
		ID:        id,
		UserID:    userId,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}
}
