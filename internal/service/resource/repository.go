package resource

import (
	"context"
	"errors"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var queryBuilder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type repository struct {
	dbPool *pgxpool.Pool
}

func newRepository(dbPool *pgxpool.Pool) *repository {
	return &repository{
		dbPool: dbPool,
	}
}

func (r *repository) getResource(ctx context.Context, path string) (*Resource, error) {
	query, args, err := queryBuilder.
		Select("*").
		From("resource").
		Where(squirrel.Eq{"path": path}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var resource Resource
	err = pgxscan.Get(ctx, r.dbPool, &resource, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &resource, nil
}
