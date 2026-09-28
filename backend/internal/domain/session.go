package domain

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	TokenHash  []byte
	IssuedAt   time.Time
	ExpiresAt  time.Time
	RevokedAt  *time.Time
	ReplacedBy *uuid.UUID
}

func NewSession(id, userID uuid.UUID, tokenHash []byte, expiresAt time.Time) *Session {
	return &Session{
		ID:        id,
		UserID:    userID,
		TokenHash: tokenHash,
		IssuedAt:  time.Now().UTC(),
		ExpiresAt: expiresAt,
	}
}

func (s *Session) IsUsable(now time.Time) bool {
	return s.RevokedAt == nil && now.Before(s.ExpiresAt)
}
