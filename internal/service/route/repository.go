package route

import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
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

func (r *repository) getRouteByPath(ctx context.Context, path string) (*Route, error) {
	query, args, err := queryBuilder.
		Select(
			"s.host AS host",
			"r.path AS path",
			"r.method AS method",
			"r.scheme AS scheme",
			"r.roles AS roles",
			"r.check_access_token AS check_access_token",
			"r.check_roles AS check_roles",
		).
		From("resource AS r").
		Join("service AS s ON s.id = r.service_id").
		Where(
			squirrel.And{
				squirrel.Eq{"r.path": path},
			},
		).
		ToSql()

	var route Route
	err = pgxscan.Get(ctx, r.dbPool, &route, query, args...)
	if err != nil {
		return nil, err
	}

	return &route, nil
}

func (r *repository) getActiveRoutes(ctx context.Context) ([]*Route, error) {
	query, args, err := queryBuilder.
		Select(
			"s.host AS host",
			"r.path AS path",
			"r.method AS method",
			"r.scheme AS scheme",
			"r.roles AS roles",
			"r.check_access_token AS check_access_token",
			"r.check_roles AS check_roles",
		).
		From("resource AS r").
		Join("service AS s ON s.id = r.service_id").
		Where(
			squirrel.And{
				squirrel.Eq{"s.is_active": true},
				squirrel.Expr("s.deleted_at IS NULL"),
				squirrel.Eq{"r.is_active": true},
				squirrel.Expr("r.deleted_at IS NULL"),
			},
		).
		ToSql()

	var routes []*Route
	err = pgxscan.Select(ctx, r.dbPool, &routes, query, args...)
	if err != nil {
		return nil, err
	}
	return routes, nil
}
