package route

import (
	"context"
	"encoding/json"
	"github.com/handmade-jewelry/auth-service/logger"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/handmade-jewelry/auth-service/internal/cache"
	"github.com/handmade-jewelry/auth-service/libs/pgutils"
)

const serviceRoutePrefix = "service_route:"

type Service struct {
	repo        *repository
	redisClient *cache.RedisClient
}

func NewService(dbPool *pgxpool.Pool, redisClient *cache.RedisClient) *Service {
	return &Service{
		repo:        newRepository(dbPool),
		redisClient: redisClient,
	}
}

func (s *Service) GetRouteByPath(ctx context.Context, path string) (*Route, error) {
	var route *Route
	val, err := s.redisClient.GetBytes(ctx, serviceRoutePrefix+path)
	if err == nil {
		err = json.Unmarshal(val, &route)
		if err == nil {
			return route, nil
		}
	}

	logger.ErrorWithFields("failed to unmarshal route from cache", err, "path", path)

	route, err = s.repo.getRouteByPath(ctx, path)
	if err != nil {
		return nil, pgutils.MapPostgresError("failed to get route", err)
	}

	return route, nil
}

func (s *Service) RefreshCacheRoutes(ctx context.Context, ttl time.Duration) error {
	routes, err := s.repo.getActiveRoutes(ctx)
	if err != nil {
		return pgutils.MapPostgresError("failed to get active routes", err)
	}

	logger.Info("fetched active routes", "count", strconv.Itoa(len(routes)))

	for _, route := range routes {
		value, err := json.Marshal(route)
		if err != nil {
			logger.ErrorWithFields("failed to marshal route", err, "route.public_path", route.PublicPath)
			continue
		}

		err = s.redisClient.Set(ctx, serviceRoutePrefix+route.PublicPath, string(value), ttl)
		if err != nil {
			logger.Error("failed to set route in Redis", err)
		}
	}

	logger.Info("route cache refresh complete", "count", strconv.Itoa(len(routes)))

	return nil
}
