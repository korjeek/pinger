package mapper

import (
	"github.com/korjeek/pinger/backend/internal/domain"
	"github.com/korjeek/pinger/backend/internal/storage/postgres/db"
)

func ToDomainSession(s db.Session) domain.Session {
	return domain.Session{
		ID:         s.ID,
		UserID:     s.UserID,
		TokenHash:  s.TokenHash,
		IssuedAt:   s.IssuedAt,
		ExpiresAt:  s.ExpiresAt,
		RevokedAt:  s.RevokedAt,
		ReplacedBy: s.ReplacedBy,
	}
}

func ToDomainUser(user db.User) domain.User {
	return domain.User{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}
}
