package dto

import "github.com/google/uuid"

type CreateUserInput struct {
	Email    string
	Password string
}

type CreateUserOutput struct {
	Id    uuid.UUID
	Email string
}

type LoginUserInput struct {
	Email    string
	Password string
}
