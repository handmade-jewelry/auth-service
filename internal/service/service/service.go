package service

import (
	"context"
	"github.com/handmade-jewelry/auth-service/internal/service/route"
	"github.com/handmade-jewelry/auth-service/libs/pgutils"
	"github.com/handmade-jewelry/auth-service/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	repo         *repository
	routeService *route.Service
}

func NewService(dbPool *pgxpool.Pool, routeService *route.Service) *Service {
	return &Service{
		repo:         newRepository(dbPool),
		routeService: routeService,
	}
}

func (s *Service) CreateService(ctx context.Context, dto *ServiceDTO) (*ServiceEntity, error) {
	srv, err := s.repo.createService(ctx, dto)
	if err != nil {
		return nil, pgutils.MapPostgresError("failed to create service", err)
	}

	err = s.routeService.RefreshCacheRoutes(ctx)
	if err != nil {
		logger.Error("failed to refresh cache routes", err)
	}

	return srv, nil
}

func (s *Service) UpdateService(ctx context.Context, dto *ServiceDTO, id int64) (*ServiceEntity, error) {
	srv, err := s.repo.updateService(ctx, dto, id)
	if err != nil {
		return nil, pgutils.MapPostgresError("failed to update service", err)
	}

	err = s.routeService.RefreshCacheRoutes(ctx)
	if err != nil {
		logger.Error("failed to refresh cache routes", err)
	}

	return srv, nil
}

func (s *Service) DeleteService(ctx context.Context, id int64) error {
	err := s.repo.deleteService(ctx, id)
	if err != nil {
		return pgutils.MapPostgresError("failed to delete service", err)
	}

	err = s.routeService.RefreshCacheRoutes(ctx)
	if err != nil {
		logger.Error("failed to refresh cache routes", err)
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
