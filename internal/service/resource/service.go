package resource

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/handmade-jewelry/auth-service/internal/service/route"
	"github.com/handmade-jewelry/auth-service/internal/service/service"
	"github.com/handmade-jewelry/auth-service/internal/service/user"
	"github.com/handmade-jewelry/auth-service/internal/utils/errors"
	"github.com/handmade-jewelry/auth-service/internal/utils/logger"
	"github.com/handmade-jewelry/auth-service/internal/utils/pgutils"
)

type Service struct {
	repo           *repository
	serviceService *service.Service
	userService    *user.Service
	routeService   *route.Service
}

func NewService(
	dbPool *pgxpool.Pool,
	serviceService *service.Service,
	userService *user.Service,
	routeService *route.Service,
) *Service {
	return &Service{
		repo:           newRepository(dbPool),
		serviceService: serviceService,
		userService:    userService,
		routeService:   routeService,
	}
}

func (s *Service) CreateResource(ctx context.Context, dto ResourceDTO) (*Resource, *errors.HTTPError) {
	srv, httpErr := s.serviceService.ServiceByID(ctx, dto.ServiceID)
	if httpErr != nil {
		return nil, httpErr
	}

	httpErr = s.validateResourceDTO(ctx, dto)
	if httpErr != nil {
		return nil, httpErr
	}

	resource, err := s.repo.createResource(ctx, dto)
	if err != nil {
		logger.ErrorWithFields("failed to create resource", err, "resource_dto", dto)
		return nil, pgutils.MapPostgresError("resource", err)
	}

	s.setCacheResource(ctx, resource, srv)

	return resource, nil
}

func (s *Service) validateResourceDTO(ctx context.Context, dto ResourceDTO) *errors.HTTPError {
	err := s.validateRoles(ctx, dto)
	if err != nil {
		return err
	}

	err = dto.Validate()
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) validateRoles(ctx context.Context, dto ResourceDTO) *errors.HTTPError {
	if !dto.CheckRoles {
		return nil
	}

	if len(dto.Roles) == 0 {
		return errors.Error("roles list cannot be empty", http.StatusBadRequest)
	}

	roleMap, err := s.userService.RoleMap(ctx)
	if err != nil {
		logger.Error("failed to get role list: %w", err)
		return errors.InternalError()
	}

	for _, role := range dto.Roles {
		if _, ok := roleMap[role]; !ok {
			return errors.Error("invalid role", http.StatusBadRequest)
		}
	}

	return nil
}
func (s *Service) setCacheResource(ctx context.Context, resource *Resource, service *service.ServiceEntity) {
	route := &route.Route{
		Host:             service.Host,
		PublicPath:       resource.PublicPath,
		ServicePath:      resource.PublicPath,
		Method:           string(resource.Method),
		Scheme:           string(resource.Scheme),
		Roles:            resource.Roles,
		CheckAccessToken: resource.CheckAccessToken,
		CheckRoles:       resource.CheckRoles,
	}

	err := s.routeService.SetCacheRoute(ctx, route)
	if err != nil {
		logger.ErrorWithFields("failed to set route in cache", err, "route", route)
	}
}

func (s *Service) UpdateResource(ctx context.Context, dto ResourceDTO, id int64) (*Resource, *errors.HTTPError) {
	srv, httpErr := s.serviceService.ServiceByID(ctx, dto.ServiceID)
	if httpErr != nil {
		return nil, httpErr
	}

	httpErr = s.validateResourceDTO(ctx, dto)
	if httpErr != nil {
		return nil, httpErr
	}

	resource, err := s.repo.updateResource(ctx, dto, id)
	if err != nil {
		logger.ErrorWithFields("failed to update resource", err, "data", dto)
		return nil, pgutils.MapPostgresError("resource", err)
	}

	s.setCacheResource(ctx, resource, srv)

	return resource, nil
}

func (s *Service) Resource(ctx context.Context, id int64) (*Resource, *errors.HTTPError) {
	resource, err := s.repo.resource(ctx, id)
	if err != nil {
		logger.ErrorWithFields("failed to get resource", err, "id", id)
		return nil, pgutils.MapPostgresError("resource", err)
	}

	return resource, nil
}

func (s *Service) DeleteResource(ctx context.Context, id int64) *errors.HTTPError {
	resource, err := s.repo.deleteResource(ctx, id)
	if err != nil {
		logger.ErrorWithFields("failed to delete resource", err, "id", id)
		return pgutils.MapPostgresError("resource", err)
	}

	s.deleteCacheResource(ctx, resource)
	return nil
}

func (s *Service) deleteCacheResource(ctx context.Context, resource *Resource) {
	err := s.routeService.DeleteCacheRoute(ctx, resource.PublicPath)
	if err != nil {
		logger.Error("failed to delete cache route", err)
	}
}
