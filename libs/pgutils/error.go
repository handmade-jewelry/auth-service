package pgutils

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AlreadyExistsCode      = "23505"
	FailedPreconditionCode = "23503"
	InvalidArgumentCode    = "23502"
)

func MapPostgresError(msg string, err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return status.Error(codes.NotFound, msg+": not found")
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case AlreadyExistsCode:
			return status.Error(codes.AlreadyExists, msg+": resource already exists")
		case FailedPreconditionCode:
			return status.Error(codes.FailedPrecondition, msg+": foreign key constraint failed")
		case InvalidArgumentCode:
			return status.Error(codes.InvalidArgument, msg+": null value in column that does not allow it")
		}
	}
	return status.Error(codes.Internal, msg+" :"+err.Error())
}
