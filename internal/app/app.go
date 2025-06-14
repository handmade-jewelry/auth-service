package app

import (
	"context"
	"fmt"
	"github.com/handmade-jewelry/auth-service/internal/service/auth"
	"github.com/handmade-jewelry/auth-service/internal/service/route"
	pkgAuth "github.com/handmade-jewelry/auth-service/internal/transport/auth"
	"github.com/handmade-jewelry/auth-service/internal/transport/resource"
	"github.com/handmade-jewelry/auth-service/internal/utils/logger"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"

	"github.com/handmade-jewelry/auth-service/internal/cache"
	"github.com/handmade-jewelry/auth-service/internal/config"
	"github.com/handmade-jewelry/auth-service/internal/jwt"
	resourceService "github.com/handmade-jewelry/auth-service/internal/service/resource"
	serviceService "github.com/handmade-jewelry/auth-service/internal/service/service"
	userService "github.com/handmade-jewelry/auth-service/internal/service/user"
	"github.com/handmade-jewelry/auth-service/internal/transport"
	"github.com/handmade-jewelry/auth-service/internal/transport/proxy"
	resourceRefresh "github.com/handmade-jewelry/auth-service/internal/worker/resource-refresh"
)

type App struct {
	cfg                   *config.Config
	redisClient           *cache.RedisClient
	jwtService            *jwt.Service
	userService           *userService.Service
	resourceService       *resourceService.Service
	serviceService        *serviceService.Service
	routeService          *route.Service
	authService           *auth.Service
	authAPIHandler        *pkgAuth.APIHandler
	resourceAPIHandler    *resource.APIHandler
	authMiddleware        *proxy.AuthMiddleware
	server                *transport.Server
	dBPool                *pgxpool.Pool
	refreshResourceTicker *resourceRefresh.Ticker
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	cfg := &transport.SwaggerConfig{
		SwaggerURL:              a.cfg.SwaggerURL,
		SwaggerAuthURL:          a.cfg.SwaggerAuthURL,
		SwaggerResourceURL:      a.cfg.SwaggerResourceURL,
		SwaggerAuthSpecURL:      a.cfg.SwaggerAuthSpecURL,
		SwaggerAuthSpecPath:     a.cfg.SwaggerAuthSpecPath,
		SwaggerResourceSpecURL:  a.cfg.SwaggerResourceSpecURL,
		SwaggerResourceSpecPath: a.cfg.SwaggerResourceSpecPath,
	}

	a.runWorker(ctx)

	return a.server.Run(cfg)
}

func (a *App) initDeps(ctx context.Context) error {
	initFuncs := []func(ctx context.Context) error{
		a.initConfig,
		a.initCache,
		a.initJWTService,
		a.initDb,
		a.initService,
		a.initAPIHandler,
		a.initMiddleware,
		a.initServer,
		a.initWorker,
	}

	for _, initF := range initFuncs {
		err := initF(ctx)
		if err != nil {
			logger.Error("failed to init deps", err)
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.LoadConfig()
	if err != nil {
		return err
	}

	dBMaxConLifetime, err := time.ParseDuration(viper.GetString(config.DBMaxConLifetime))
	if err != nil {
		return fmt.Errorf("failed to parse dBPool max conns lifetime config: %w", err)
	}

	accessTokenTTL, err := time.ParseDuration(viper.GetString(config.AccessTokenTTL))
	if err != nil {
		return fmt.Errorf("failed to parse access token exp config: %w", err)
	}

	refreshTokenTTL, err := time.ParseDuration(viper.GetString(config.RefreshTokenTTL))
	if err != nil {
		return fmt.Errorf("failed to parse refresh token exp config: %w", err)
	}

	refreshResourceTTL, err := time.ParseDuration(viper.GetString(config.RefreshResourceTTL))
	if err != nil {
		return fmt.Errorf("failed to parse refresh resource ttl config: %w", err)
	}

	refreshResourcesInterval, err := time.ParseDuration(viper.GetString(config.RefreshResourcesInterval))
	if err != nil {
		return fmt.Errorf("failed to parse refresh resources interval config: %w", err)
	}

	a.cfg = &config.Config{
		DBName:                       viper.GetString(config.DBName),
		DBUser:                       viper.GetString(config.DBUser),
		DBPassword:                   viper.GetString(config.DBPassword),
		DbHost:                       viper.GetString(config.DBHost),
		DbPort:                       viper.GetUint16(config.DBPort),
		SSLMode:                      viper.GetString(config.SSLMode),
		DBMaxCons:                    viper.GetInt32(config.DBMaxCons),
		DBMinCons:                    viper.GetInt32(config.DBMinCons),
		DBMaxConLifetime:             dBMaxConLifetime,
		HTTPServerPort:               viper.GetString(config.HTTPServerPort),
		HTTPProxyPrefix:              viper.GetString(config.HTTPProxyPrefix),
		HTTPAuthPrefix:               viper.GetString(config.HTTPAuthPrefix),
		HTTPResourcePrefix:           viper.GetString(config.HTTPResourcePrefix),
		SwaggerURL:                   viper.GetString(config.SwaggerURL),
		SwaggerAuthURL:               viper.GetString(config.SwaggerAuthURL),
		SwaggerResourceURL:           viper.GetString(config.SwaggerResourceURL),
		SwaggerAuthSpecURL:           viper.GetString(config.SwaggerAuthSpecURL),
		SwaggerAuthSpecPath:          viper.GetString(config.SwaggerAuthSpecPath),
		SwaggerResourceSpecURL:       viper.GetString(config.SwaggerResourceSpecURL),
		SwaggerResourceSpecPath:      viper.GetString(config.SwaggerResourceSpecPath),
		RedisAddress:                 viper.GetString(config.RedisAddress),
		RedisPassword:                viper.GetString(config.RedisPassword),
		RedisDB:                      viper.GetInt(config.RedisDb),
		AccessTokenTTL:               accessTokenTTL,
		RefreshTokenTTL:              refreshTokenTTL,
		JWTTokenSecret:               viper.GetString(config.JWTTokenSecret),
		RefreshResourceTTL:           refreshResourceTTL,
		RefreshResourceInterval:      refreshResourcesInterval,
		RefreshCacheRoutesWorkerMode: viper.GetString(config.RefreshCacheRoutesWorkerMode),
	}

	return nil
}

func (a *App) initJWTService(_ context.Context) error {
	a.jwtService = jwt.NewService(a.cfg.JWTTokenSecret, a.cfg.AccessTokenTTL, a.cfg.RefreshTokenTTL)
	return nil
}

func (a *App) initService(_ context.Context) error {
	grpcOpts := config.GRPCOptions{
		Host:            viper.GetString(config.UserServiceHost),
		MaxRetry:        viper.GetUint(config.GRPCClientMaxRetry),
		PerRetryTimeout: viper.GetDuration(config.GRPCClientRetryTimeout),
	}

	var err error
	a.userService, err = userService.NewService(&grpcOpts)
	if err != nil {
		return err
	}

	a.routeService = route.NewService(a.dBPool, a.redisClient, a.cfg.RefreshResourceTTL)
	a.serviceService = serviceService.NewService(a.dBPool, a.routeService)
	a.resourceService = resourceService.NewService(a.dBPool, a.serviceService, a.userService, a.routeService)
	a.authService = auth.NewService(a.jwtService, a.redisClient, a.userService, a.cfg.AccessTokenTTL, a.cfg.RefreshTokenTTL)

	return nil
}

func (a *App) initAPIHandler(_ context.Context) error {
	a.authAPIHandler = pkgAuth.NewAPIHandler(a.authService)
	a.resourceAPIHandler = resource.NewAPIHandler(a.serviceService, a.resourceService, a.userService)
	return nil
}

func (a *App) initCache(_ context.Context) error {
	a.redisClient = cache.NewRedisClient(
		a.cfg.RedisAddress,
		a.cfg.RedisPassword,
		a.cfg.RedisDB,
	)

	return nil
}

func (a *App) initMiddleware(_ context.Context) error {
	a.authMiddleware = proxy.NewAuthMiddleware(a.routeService, a.jwtService)
	return nil
}

func (a *App) initServer(_ context.Context) error {
	opts := &transport.Opts{
		HTTPPort:       a.cfg.HTTPServerPort,
		ProxyPrefix:    a.cfg.HTTPProxyPrefix,
		AuthPrefix:     a.cfg.HTTPAuthPrefix,
		ResourcePrefix: a.cfg.HTTPResourcePrefix,
	}
	a.server = transport.NewServer(opts, a.authMiddleware, a.authAPIHandler, a.resourceAPIHandler)
	return nil
}

func (a *App) initDb(ctx context.Context) error {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		a.cfg.DBUser,
		a.cfg.DBPassword,
		a.cfg.DbHost,
		a.cfg.DbPort,
		a.cfg.DBName,
		a.cfg.SSLMode,
	)

	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return fmt.Errorf("failed to parse db config: %w", err)
	}

	cfg.MaxConns = a.cfg.DBMaxCons
	cfg.MinConns = a.cfg.DBMinCons
	cfg.MaxConnLifetime = a.cfg.DBMaxConLifetime

	dbPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return fmt.Errorf("unable to create pool: %w", err)
	}

	if err = dbPool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	a.dBPool = dbPool

	logger.Info(
		"Database connection established",
		a.cfg.DbHost,
		strconv.Itoa(int(a.cfg.DbPort)))

	return nil
}

func (a *App) initWorker(_ context.Context) error {
	a.refreshResourceTicker = resourceRefresh.NewTiker(a.routeService, a.cfg.RefreshCacheRoutesWorkerMode)
	return nil
}

func (a *App) runWorker(ctx context.Context) {
	a.refreshResourceTicker.Run(ctx, a.cfg.RefreshResourceInterval)
}
