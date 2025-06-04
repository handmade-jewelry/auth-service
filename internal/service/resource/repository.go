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

func (r *repository) createResource(ctx context.Context, dto ResourceDTO) (*Resource, error) {
	query, args, err := queryBuilder.
		Insert(resourceTable).
		Columns(
			"service_id",
			"public_path",
			"service_path",
			"method",
			"scheme",
			"roles",
			"is_active",
			"check_access_token",
			"check_roles",
		).
		Values(
			dto.ServiceID,
			dto.PublicPath,
			dto.ServicePath,
			dto.Method,
			dto.Scheme,
			dto.Roles,
			dto.IsActive,
			dto.CheckAccessToken,
			dto.CheckRoles,
		).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, err
	}

	var resource Resource
	if err = pgxscan.Get(ctx, r.dbPool, &resource, query, args...); err != nil {
		return nil, err
	}

	return &resource, nil
}

func (r *repository) updateResource(ctx context.Context, dto ResourceDTO, id int64) (*Resource, error) {
	query, args, err := queryBuilder.
		Update(resourceTable).
		Set("service_id", dto.ServiceID).
		Set("public_path", dto.PublicPath).
		Set("service_path", dto.ServicePath).
		Set("method", dto.Method).
		Set("scheme", dto.Scheme).
		Set("roles", dto.Roles).
		Set("is_active", dto.IsActive).
		Set("check_access_token", dto.CheckAccessToken).
		Set("check_roles", dto.CheckRoles).
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, err
	}

	var resource Resource
	if err = pgxscan.Get(ctx, r.dbPool, &resource, query, args...); err != nil {
		return nil, err
	}

	return &resource, nil
}

func (r *repository) resource(ctx context.Context, id int64) (*Resource, error) {
	query, args, err := queryBuilder.
		Select("*").
		From(resourceTable).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var resource Resource
	err = pgxscan.Get(ctx, r.dbPool, &resource, query, args...)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}

func (r *repository) deleteResource(ctx context.Context, id int64) (*Resource, error) {
	query, args, err := queryBuilder.
		Update(resourceTable).
		Set("deleted_at", squirrel.Expr("CURRENT_TIMESTAMP")).
		Set("is_active", false).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, err
	}

	var resource Resource
	err = pgxscan.Get(ctx, r.dbPool, &resource, query, args...)
	if err != nil {
		return nil, err
	}

	return &resource, nil
}

func (r *repository) resourceByPublicPath(ctx context.Context, publicPath string) (*Resource, error) {
	query, args, err := queryBuilder.
		Select("*").
		From(resourceTable).
		Where(squirrel.Eq{"public_path": publicPath}).
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

func (r *repository) resourceByServiceIDs(ctx context.Context, ids []int) ([]*Resource, error) {
	query, args, err := queryBuilder.
		Select("*").
		From(resourceTable).
		Where(squirrel.Eq{"service_id": ids}).
		Where(squirrel.Eq{"is_active": true}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var resources []*Resource
	err = pgxscan.Get(ctx, r.dbPool, &resources, query, args...)
	if err != nil {
		return nil, err
	}
	return resources, nil
}
