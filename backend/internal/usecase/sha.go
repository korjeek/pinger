package usecase

import (
	"crypto/sha256"
)

type SHA256Hasher struct{}

func (h *SHA256Hasher) Hash(data string) ([]byte, error) {
	sum := sha256.Sum256([]byte(data))
	return sum[:], nil
}
