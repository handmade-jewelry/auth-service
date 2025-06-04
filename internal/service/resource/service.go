package resource

import (
	"context"
	"fmt"
	"github.com/handmade-jewelry/auth-service/internal/service/route"
	"github.com/handmade-jewelry/auth-service/internal/service/service"
	"github.com/handmade-jewelry/auth-service/internal/service/user"
	"github.com/handmade-jewelry/auth-service/libs/pgutils"
	"github.com/handmade-jewelry/auth-service/logger"
	"github.com/jackc/pgx/v5/pgxpool"
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

func (s *Service) ResourceByPublicPath(ctx context.Context, path string) (*Resource, error) {
	return s.repo.resourceByPublicPath(ctx, path)
}

func (s *Service) ResourceByServiceIDs(ctx context.Context, ids []int) ([]*Resource, error) {
	return s.repo.resourceByServiceIDs(ctx, ids)
}

func (s *Service) CreateResource(ctx context.Context, dto ResourceDTO) (*Resource, error) {
	srv, err := s.serviceService.ServiceByID(ctx, dto.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get service by id: %d: %w", dto.ServiceID, err)
	}

	err = s.validateResourceDTO(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("invalid resource data: %w", err)
	}

	resource, err := s.repo.createResource(ctx, dto)
	if err != nil {
		logger.ErrorWithFields("failed to create resource", err, "data", dto)
		return nil, pgutils.MapPostgresError("failed to create resource", err)
	}

	s.setCacheResource(ctx, resource, srv)

	return resource, nil
}

func (s *Service) validateResourceDTO(ctx context.Context, dto ResourceDTO) error {
	err := s.validateRoles(ctx, dto)
	if err != nil {
		return fmt.Errorf("failed to validate roles: %w", err)
	}

	err = dto.Validate()
	if err != nil {
		return fmt.Errorf("invalid resource data: %w", err)
	}

	return nil
}

func (s *Service) validateRoles(ctx context.Context, dto ResourceDTO) error {
	if !dto.CheckRoles {
		return nil
	}

	if len(dto.Roles) == 0 {
		return fmt.Errorf("roles list cannot be empty")
	}

	roleMap, err := s.userService.RoleMap(ctx)
	if err != nil {
		return fmt.Errorf("failed to get role list: %w", err)
	}

	for _, role := range dto.Roles {
		if _, ok := roleMap[role]; !ok {
			return fmt.Errorf("invalid role: %s", role)
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
		logger.Error("failed to set route in cache", err)
	}
}

func (s *Service) UpdateResource(ctx context.Context, dto ResourceDTO, id int64) (*Resource, error) {
	srv, err := s.serviceService.ServiceByID(ctx, dto.ServiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get service by id: %d: %w", dto.ServiceID, err)
	}

	err = s.validateResourceDTO(ctx, dto)
	if err != nil {
		return nil, fmt.Errorf("invalid resource data: %w", err)
	}

	resource, err := s.repo.updateResource(ctx, dto, id)
	if err != nil {
		logger.ErrorWithFields("failed to update resource", err, "data", dto)
		return nil, pgutils.MapPostgresError("failed to update resource", err)
	}

	s.setCacheResource(ctx, resource, srv)

	return resource, nil
}

func (s *Service) Resource(ctx context.Context, id int64) (*Resource, error) {
	resource, err := s.repo.resource(ctx, id)
	if err != nil {
		logger.ErrorWithFields("failed to get resource", err, "id", id)
		return nil, pgutils.MapPostgresError("failed to get resource", err)
	}

	return resource, nil
}

func (s *Service) DeleteResource(ctx context.Context, id int64) error {
	resource, err := s.repo.deleteResource(ctx, id)
	if err != nil {
		logger.ErrorWithFields("failed to delete resource", err, "id", id)
		return pgutils.MapPostgresError("failed to delete resource", err)
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
