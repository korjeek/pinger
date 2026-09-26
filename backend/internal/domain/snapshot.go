package domain

import (
	"time"
	"uuid"
)

type Snapshot struct {
	ID           uuid.UUID
	WebsiteID    uuid.UUID
	StatusCode   int
	ResponseTime time.Duration
	ResponseSize int64
	ServerName   string
	SSLExpiresAt *time.Time
}
