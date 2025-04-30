package app

import (
	"context"
	"fmt"
	"github.com/handmade-jewellery/auth-service/internal/cache"
	rS "github.com/handmade-jewellery/auth-service/internal/service/resource-service"
	sS "github.com/handmade-jewellery/auth-service/internal/service/service-service"
	uS "github.com/handmade-jewellery/auth-service/internal/service/user-service"
	"github.com/handmade-jewellery/auth-service/internal/transport"
	"github.com/handmade-jewellery/auth-service/internal/transport/proxy"
	"github.com/jackc/pgx/v5"
	"log"
)

type App struct {
	cfg             *Config
	userService     *uS.UserService
	resourceService *rS.ResourceService
	serviceService  *sS.ServiceService
	server          *transport.Server
	dB              *pgx.Conn
	redisClient     *cache.RedisClient
	authMiddleware  *proxy.AuthMiddleware
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

func (a *App) initDeps() error {
	err := a.initConfig()
	if err != nil {
		return err
	}

	a.initCache()
	a.initDb()
	a.initService()
	a.initMiddleware()
	a.initServer()

	return nil
}

func (a *App) initConfig() error {
	err := loadConfig()
	if err != nil {
		return err
	}

	a.cfg = initConfig()

	return nil
}

func (a *App) initService() {
	a.userService = uS.NewService()
	a.resourceService = rS.NewService()
	a.serviceService = sS.NewService()
}

func (a *App) initCache() {
	a.redisClient = cache.NewRedisClient(
		a.cfg.RedisAddress,
		a.cfg.RedisPassword,
		a.cfg.RedisDB,
	)
}

func (a *App) initMiddleware() {
	a.authMiddleware = proxy.NewAuthMiddleware(a.userService, a.resourceService, a.serviceService)
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
