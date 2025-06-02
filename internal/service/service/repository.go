package service

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

func (r *repository) createService(ctx context.Context, dto *ServiceDTO) (*ServiceEntity, error) {
	query, args, err := queryBuilder.
		Insert(serviceTable).
		Columns("name", "is_active", "host").
		Values(dto.Name, dto.IsActive, dto.Host).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, err
	}

	var svc ServiceEntity
	if err = pgxscan.Get(ctx, r.dbPool, &svc, query, args...); err != nil {
		return nil, err
	}

	return &svc, nil
}

func (r *repository) updateService(ctx context.Context, dto *ServiceDTO, id int64) (*ServiceEntity, error) {
	query, args, err := queryBuilder.
		Update(serviceTable).
		Set("name", dto.Name).
		Set("is_active", dto.IsActive).
		Set("host", dto.Host).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING *").
		ToSql()
	if err != nil {
		return nil, err
	}

	var svc ServiceEntity
	if err = pgxscan.Get(ctx, r.dbPool, &svc, query, args...); err != nil {
		return nil, err
	}

	return &svc, nil
}

func (r *repository) deleteService(ctx context.Context, id int64) error {
	query, args, err := queryBuilder.
		Delete("service").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.dbPool.Exec(ctx, query, args...)
	return err
}

func (r *repository) serviceByID(ctx context.Context, id int64) (*ServiceEntity, error) {
	query, args, err := queryBuilder.
		Select("*").
		From(serviceTable).
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var service ServiceEntity
	err = pgxscan.Get(ctx, r.dbPool, &service, query, args...)
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (r *repository) serviceByName(ctx context.Context, name string) (*ServiceEntity, error) {
	query, args, err := queryBuilder.
		Select("*").
		From(serviceTable).
		Where(squirrel.Eq{"name": name}).
		ToSql()
	if err != nil {
		return nil, err
	}

	var service ServiceEntity
	err = pgxscan.Get(ctx, r.dbPool, &service, query, args...)
	if err != nil {
		return nil, err
	}

	return &service, nil
}
