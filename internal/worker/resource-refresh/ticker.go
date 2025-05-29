package resource_refresh

import (
	"context"
	"github.com/handmade-jewelry/auth-service/logger"
	"log"
	"time"

	"github.com/handmade-jewelry/auth-service/internal/cache"
	"github.com/handmade-jewelry/auth-service/internal/service/route"
)

const (
	enabledMode  = "ENABLED"
	disabledMode = "DISABLED"
)

type Ticker struct {
	redisClient  *cache.RedisClient
	routeService *route.Service
	mode         string
}

func NewTiker(routeService *route.Service, mode string) *Ticker {
	return &Ticker{
		routeService: routeService,
		mode:         mode,
	}
}

func (t *Ticker) Run(ctx context.Context, interval time.Duration, ttl time.Duration) {
	if t.mode == disabledMode {
		return
	}

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
		logger.Error("failed to refresh cached routes", err)
	}
}
