package postgres

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/korjeek/pinger/backend/pkg/apperr"
)

func FromPgError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return apperr.NotFound("resource not found").WithCause(err)
	}

	if pgErr, ok := errors.AsType[*pgconn.PgError](err); ok {
		var appErr *apperr.Error

		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			appErr = apperr.Conflict("resource already exists")
		case pgerrcode.ForeignKeyViolation:
			appErr = apperr.BadRequest("referenced resource does not exist")
		case pgerrcode.NotNullViolation:
			appErr = apperr.Validation("required field is missing", map[string]any{
				"table":  pgErr.TableName,
				"column": pgErr.ColumnName,
			})
		default:
			return apperr.Internal(pgErr)
		}

		return appErr.
			WithDetail("table", pgErr.TableName).
			WithDetail("column", pgErr.ColumnName).
			WithDetail("constraint", pgErr.ConstraintName).
			WithCause(pgErr)
	}

	return apperr.Internal(err)
}
