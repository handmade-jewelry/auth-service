package app

import (
	"context"
	"fmt"
	"github.com/handmade-jewelry/auth-service/internal/cache"
	"github.com/handmade-jewelry/auth-service/internal/config"
	"github.com/handmade-jewelry/auth-service/internal/jwt"
	resourceService "github.com/handmade-jewelry/auth-service/internal/service/resource"
	serviceService "github.com/handmade-jewelry/auth-service/internal/service/service"
	userService "github.com/handmade-jewelry/auth-service/internal/service/user"
	"github.com/handmade-jewelry/auth-service/internal/transport"
	"github.com/handmade-jewelry/auth-service/internal/transport/proxy"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"log"
	"time"
)

type App struct {
	cfg             *config.Config
	userService     *userService.Service
	resourceService *resourceService.Service
	serviceService  *serviceService.Service
	redisClient     *cache.RedisClient
	authMiddleware  *proxy.AuthMiddleware
	jwtService      *jwt.Service
	server          *transport.Server
	dBPool          *pgxpool.Pool
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	cfg := &transport.Config{
		HTTPPort:            a.cfg.HTTPServerPort,
		SwaggerURLPath:      a.cfg.SwaggerURLPath,
		SwaggerSpecFilePath: a.cfg.SwaggerSpecFilePath,
	}

	return a.server.Run(cfg)
}

// todo
func (a *App) initDeps(ctx context.Context) error {
	err := a.initConfig()
	if err != nil {
		return err
	}

	a.initCache()
	a.initJWTService()

	err = a.initDb(ctx)
	if err != nil {
		return err
	}

	err = a.initService()
	if err != nil {
		return err
	}

	a.initMiddleware()
	a.initServer()

	return nil
}

func (a *App) initConfig() error {
	err := config.LoadConfig()
	if err != nil {
		return err
	}

	dBMaxConLifetime, err := time.ParseDuration(viper.GetString(config.DBMaxConLifetime))
	if err != nil {
		log.Fatalf("Failed to parse dBPool max conns lifetime duration config: %v", err)
		return err
	}

	accessTokenExp, err := time.ParseDuration(viper.GetString(config.AccessTokenExpMin))
	if err != nil {
		log.Fatalf("Failed to parse access token exp config: %v", err)
		return err
	}

	refreshTokenExp, err := time.ParseDuration(viper.GetString(config.RefreshTokenExpMin))
	if err != nil {
		log.Fatalf("Failed to parse refresh token exp config: %v", err)
		return err
	}

	a.cfg = &config.Config{
		DBName:              viper.GetString(config.DBName),
		DBUser:              viper.GetString(config.DBUser),
		DBPassword:          viper.GetString(config.DBPassword),
		DbHost:              viper.GetString(config.DBHost),
		DbPort:              viper.GetUint16(config.DBPort),
		SSLMode:             viper.GetString(config.SSLMode),
		DBMaxCons:           viper.GetInt32(config.DBMaxCons),
		DBMinCons:           viper.GetInt32(config.DBMinCons),
		DBMaxConLifetime:    dBMaxConLifetime,
		HTTPServerPort:      viper.GetString(config.HTTPServerPort),
		SwaggerURLPath:      viper.GetString(config.SwaggerURLPath),
		SwaggerSpecFilePath: viper.GetString(config.SwaggerSpecFilePath),
		RedisAddress:        viper.GetString(config.RedisAddress),
		RedisPassword:       viper.GetString(config.RedisPassword),
		RedisDB:             viper.GetInt(config.RedisDb),
		AccessTokenExp:      accessTokenExp,
		RefreshTokenExp:     refreshTokenExp,
		JWTTokenSecret:      viper.GetString(config.JWTTokenSecret),
	}

	return nil
}

func (a *App) initJWTService() {
	a.jwtService = jwt.NewService(a.cfg.JWTTokenSecret, a.cfg.AccessTokenExp, a.cfg.RefreshTokenExp)
}

func (a *App) initService() error {
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

	a.resourceService = resourceService.NewService(a.dBPool)
	a.serviceService = serviceService.NewService()

	return nil
}

func (a *App) initCache() {
	a.redisClient = cache.NewRedisClient(
		a.cfg.RedisAddress,
		a.cfg.RedisPassword,
		a.cfg.RedisDB,
	)
}

func (a *App) initMiddleware() {
	a.authMiddleware = proxy.NewAuthMiddleware(a.userService, a.resourceService, a.serviceService, a.jwtService, a.redisClient)
}

func (a *App) initServer() {
	a.server = transport.NewServer(a.authMiddleware)
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
		//todo log error
		return fmt.Errorf("failed to parse db config: %w", err)
	}

	cfg.MaxConns = a.cfg.DBMaxCons
	cfg.MinConns = a.cfg.DBMinCons
	cfg.MaxConnLifetime = a.cfg.DBMaxConLifetime

	dbPool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		//todo log error
		return fmt.Errorf("unable to create pool: %w", err)
	}

	if err = dbPool.Ping(ctx); err != nil {
		//todo log error
		return fmt.Errorf("failed to ping db: %w", err)
	}

	a.dBPool = dbPool

	log.Println("Database connection established successfully")

	return nil
}
