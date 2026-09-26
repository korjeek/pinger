package dto

import (
	"time"
)

type Token struct {
	Payload   string
	ExpiresAt time.Time
}

type TokenPair struct {
	AccessToken  Token
	RefreshToken Token
}
