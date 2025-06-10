package pgutils

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	utilErr "github.com/handmade-jewelry/auth-service/internal/utils/errors"
)

const (
	AlreadyExistsCode = "23505"
)

func MapPostgresError(msg string, err error) *utilErr.HTTPError {
	if errors.Is(err, pgx.ErrNoRows) {
		return utilErr.Error(msg+" not found", http.StatusNotFound)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case AlreadyExistsCode:
			return utilErr.Error(msg+" already exists", http.StatusConflict)
		}
	}
	return utilErr.Error("internal errors", http.StatusInternalServerError)
}
