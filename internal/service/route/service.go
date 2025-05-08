package route

import (
	"context"
	"encoding/json"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/handmade-jewelry/auth-service/internal/cache"
)

const refreshTokenPrefix = "refresh_token:"

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
	val, err := s.redisClient.GetBytes(ctx, refreshTokenPrefix+path)
	if err == nil {
		err = json.Unmarshal(val, &route)
		if err == nil {
			return route, nil
		}
	}

	//todo log cache error

	route, err = s.repo.getRouteByPath(ctx, path)
	if err != nil {
		return nil, err
	}

	return route, nil
}

func (s *Service) RefreshCacheRoutes(ctx context.Context, ttl time.Duration) error {
	routes, err := s.repo.getActiveRoutes(ctx)
	if err != nil {
		return err
	}

	for _, route := range routes {
		value, err := json.Marshal(route)
		if err != nil {
			//todo log err
			continue
		}

		err = s.redisClient.Set(ctx, refreshTokenPrefix+route.Path, string(value), ttl)
		if err != nil {
			//todo log err
		}
	}

	return nil
}
