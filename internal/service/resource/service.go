package resource

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repo *repository
}

func NewService(dbPool *pgxpool.Pool) *Service {
	return &Service{
		repo: newRepository(dbPool),
	}
}

func (s *Service) GetResourceByPath(ctx context.Context, path string) (*Resource, error) {
	return s.repo.getResource(ctx, path)
}

func (s *Service) GetResourceByServiceIDs(ctx context.Context, ids []int) ([]*Resource, error) {
	return s.repo.getResourceByServiceIDs(ctx, ids)
}
