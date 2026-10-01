package usecase

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCost = errors.New("bcrypt: invalid cost")

type BcryptHasher struct {
	cost int
}

func NewBcryptHasher(cost int) (*BcryptHasher, error) {
	if cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return nil, ErrInvalidCost
	}
	return &BcryptHasher{cost: cost}, nil
}

func (h *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	return string(hash), err
}

func (h *BcryptHasher) Equals(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
