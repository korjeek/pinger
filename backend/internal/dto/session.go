package dto

import "github.com/korjeek/pinger/backend/internal/domain"

type SessionIssue struct {
	Session domain.Session
	Pair    TokenPair
}
