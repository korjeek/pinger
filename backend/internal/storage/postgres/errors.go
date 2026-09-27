package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ErrType int

const (
	Unknown ErrType = iota
	NotFound
	Conflict
	Invalid
)

type ErrDatabase struct {
	errType ErrType
	err     error
}

func NewErrDatabase(err error) ErrDatabase {
	errType := Unknown

	if errors.Is(err, pgx.ErrNoRows) {
		errType = NotFound
	} else if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		switch pgErr.Code {
		case "23505":
			errType = Conflict
		case "23503", "23502":
			errType = Invalid
		}
	}

	return ErrDatabase{
		errType: errType,
		err:     err,
	}
}

func (e ErrDatabase) Error() string {
	return fmt.Sprintf("database error (%v): %v", e.errType, e.err)
}

func (e ErrDatabase) Unwrap() error {
	return e.err
}
