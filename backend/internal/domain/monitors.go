package domain

import (
	"time"

	"github.com/google/uuid"
)

type Monitor struct {
	ID           uuid.UUID
	URL          string
	Tracked      bool
	PollInterval time.Duration
}

func NewWebsite(url string, interval time.Duration) *Monitor {
	id, err := uuid.NewV7()
	if err != nil {
		//TODO: error handling
	}

	return &Monitor{
		ID:           id,
		URL:          url,
		Tracked:      true,
		PollInterval: interval,
	}
}
