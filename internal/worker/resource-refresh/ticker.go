package resource_refresh

import (
	"context"
	"log"
	"time"

	"github.com/handmade-jewelry/auth-service/internal/cache"
	"github.com/handmade-jewelry/auth-service/internal/service/route"
)

type Ticker struct {
	redisClient  *cache.RedisClient
	routeService *route.Service
}

func NewTiker(routeService *route.Service) *Ticker {
	return &Ticker{
		routeService: routeService,
	}
}

func (t *Ticker) Run(ctx context.Context, interval time.Duration, ttl time.Duration) {
	t.run(ctx, ttl)

	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				t.run(ctx, ttl)
			case <-ctx.Done():
				log.Println("Stop resource refresh ticker")
				ticker.Stop()
			}
		}
	}()
}

func (t *Ticker) run(ctx context.Context, ttl time.Duration) {
	err := t.routeService.RefreshCacheRoutes(ctx, ttl)
	if err != nil {
		//todo log error
	}
}
