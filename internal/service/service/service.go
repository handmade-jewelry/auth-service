package service

import (
	"context"
	"github.com/handmade-jewelry/auth-service/libs/pgutils"
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

func (s *Service) CreateService(ctx context.Context, dto *ServiceDTO) (*ServiceEntity, error) {
	srv, err := s.repo.createService(ctx, dto)
	if err != nil {
		return nil, pgutils.MapPostgresError("failed to create service", err)
	}

	return srv, nil
}

func (s *Service) UpdateService(ctx context.Context, dto *ServiceDTO, id int64) (*ServiceEntity, error) {
	srv, err := s.repo.updateService(ctx, dto, id)
	if err != nil {
		return nil, pgutils.MapPostgresError("failed to update service", err)
	}

	return srv, nil
}

func (s *Service) DeleteService(ctx context.Context, id int64) error {
	err := s.repo.deleteService(ctx, id)
	if err != nil {
		return pgutils.MapPostgresError("failed to delete service", err)
	}

	return nil
}

func (s *Service) ServiceByID(ctx context.Context, id int64) (*ServiceEntity, error) {
	srv, err := s.repo.serviceByID(ctx, id)
	if err != nil {
		return nil, pgutils.MapPostgresError("failed to get service by id", err)
	}

	return srv, nil
}

func (s *Service) ServiceByName(ctx context.Context, name string) (*ServiceEntity, error) {
	srv, err := s.repo.serviceByName(ctx, name)
	if err != nil {
		return nil, pgutils.MapPostgresError("failed to get service by name", err)
	}

	return srv, nil
}
