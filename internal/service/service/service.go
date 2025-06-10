package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/handmade-jewelry/auth-service/internal/service/route"
	"github.com/handmade-jewelry/auth-service/internal/utils/errors"
	"github.com/handmade-jewelry/auth-service/internal/utils/logger"
	"github.com/handmade-jewelry/auth-service/internal/utils/pgutils"
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

func (s *Service) CreateService(ctx context.Context, dto *ServiceDTO) (*ServiceEntity, *errors.HTTPError) {
	srv, err := s.repo.createService(ctx, dto)
	if err != nil {
		logger.Error("failed to create service", err)
		return nil, pgutils.MapPostgresError("service", err)
	}

	err = s.routeService.RefreshCacheRoutes(ctx)
	if err != nil {
		logger.Error("failed to refresh cache routes", err)
	}

	return srv, nil
}

func (s *Service) UpdateService(ctx context.Context, dto *ServiceDTO, id int64) (*ServiceEntity, *errors.HTTPError) {
	srv, err := s.repo.updateService(ctx, dto, id)
	if err != nil {
		return nil, pgutils.MapPostgresError("service", err)
	}

	err = s.routeService.RefreshCacheRoutes(ctx)
	if err != nil {
		logger.Error("failed to refresh cache routes", err)
	}

	return srv, nil
}

func (s *Service) DeleteService(ctx context.Context, id int64) *errors.HTTPError {
	err := s.repo.deleteService(ctx, id)
	if err != nil {
		logger.ErrorWithFields("failed to delete service", err, "id", id)
		return pgutils.MapPostgresError("service", err)
	}

	err = s.routeService.RefreshCacheRoutes(ctx)
	if err != nil {
		logger.Error("failed to refresh cache routes", err)
	}

	return nil
}

func (s *Service) ServiceByID(ctx context.Context, id int64) (*ServiceEntity, *errors.HTTPError) {
	srv, err := s.repo.serviceByID(ctx, id)
	if err != nil {
		logger.ErrorWithFields("failed to get service by id", err, "service_id", id)
		return nil, pgutils.MapPostgresError("service", err)
	}

	return srv, nil
}

func (s *Service) ServiceByName(ctx context.Context, name string) (*ServiceEntity, *errors.HTTPError) {
	srv, err := s.repo.serviceByName(ctx, name)
	if err != nil {
		logger.ErrorWithFields("failed to get service by name", err, "service_name", name)
		return nil, pgutils.MapPostgresError("service", err)
	}

	return srv, nil
}
