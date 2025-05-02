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
	"github.com/jackc/pgx/v5"
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
	dB              *pgx.Conn
}

func NewApp() (*App, error) {
	a := &App{}
	err := a.initDeps()
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
func (a *App) initDeps() error {
	err := a.initConfig()
	if err != nil {
		return err
	}

	a.initCache()

	err = a.initJWTService()
	if err != nil {
		return err
	}

	a.initDb()

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

	a.cfg = &config.Config{
		DBName:              viper.GetString(config.DbName),
		DBUser:              viper.GetString(config.DbUser),
		DBPassword:          viper.GetString(config.DbPassword),
		DbHost:              viper.GetString(config.DbHost),
		DbPort:              viper.GetUint16(config.DbPort),
		SSLMode:             viper.GetString(config.SslMode),
		HTTPServerPort:      viper.GetString(config.HttpServerPort),
		SwaggerURLPath:      viper.GetString(config.SwaggerURLPath),
		SwaggerSpecFilePath: viper.GetString(config.SwaggerSpecFilePath),
		RedisAddress:        viper.GetString(config.RedisAddress),
		RedisPassword:       viper.GetString(config.RedisPassword),
		RedisDB:             viper.GetInt(config.RedisDb),
	}

	return nil
}

func (a *App) initJWTService() error {
	accessTokenExp, err := time.ParseDuration(viper.GetString(config.AccessTokenExpMin))
	if err != nil {
		log.Fatalf("Ошибка при парсинге длительности auth: %v", err)
		return err
	}

	refreshTokenExp, err := time.ParseDuration(viper.GetString(config.RefreshTokenExpMin))
	if err != nil {
		log.Fatalf("Ошибка при парсинге длительности refresh: %v", err)
		return err
	}

	a.jwtService = jwt.NewService(viper.GetString(config.JWTTokenSecret), accessTokenExp, refreshTokenExp)

	return nil
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

	a.resourceService = resourceService.NewService()
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

func (a *App) initDb() {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		a.cfg.DBUser,
		a.cfg.DBPassword,
		a.cfg.DbHost,
		a.cfg.DbPort,
		a.cfg.DBName,
		a.cfg.SSLMode,
	)

	var err error
	a.dB, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	err = a.dB.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping the database: %v\n", err)
	}

	log.Println("Database connection established successfully")
}
