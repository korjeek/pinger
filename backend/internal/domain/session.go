package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash []byte
	ExpiresAt time.Time
}

func NewSession(userId uuid.UUID, tokenHash []byte, expiresAt time.Time) (*Session, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("new session: %w", err)
	}

	return &Session{
		ID:        id,
		UserID:    userId,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	}, nil
}
