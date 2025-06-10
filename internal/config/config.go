package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

const (
	DBName                       = "database.name"
	DBUser                       = "database.user"
	DBPassword                   = "database.password"
	DBHost                       = "database.host"
	DBPort                       = "database.port"
	SSLMode                      = "database.ssl_mode"
	DBMaxCons                    = "database.max_cons"
	DBMinCons                    = "database.min_cons"
	DBMaxConLifetime             = "database.max_con_lifetime"
	HTTPServerPort               = "server.http.port"
	HTTPProxyPrefix              = "server.http.proxy_prefix"
	HTTPAuthPrefix               = "server.http.auth_prefix"
	HTTPResourcePrefix           = "server.http.resource_prefix"
	SwaggerURL                   = "swagger.url"
	SwaggerAuthURL               = "swagger.auth_url"
	SwaggerResourceURL           = "swagger.resource_url"
	SwaggerAuthSpecURL           = "swagger.auth_spec_url"
	SwaggerAuthSpecPath          = "swagger.auth_spec_path"
	SwaggerResourceSpecURL       = "swagger.resource_spec_url"
	SwaggerResourceSpecPath      = "swagger.resource_spec_path"
	RedisAddress                 = "redis.addr"
	RedisPassword                = "redis.password"
	RedisDb                      = "redis.db"
	GRPCClientMaxRetry           = "grpc_client.max_retry"
	GRPCClientRetryTimeout       = "grpc_client.timeout_per_retry"
	UserServiceHost              = "user_service.host"
	JWTTokenSecret               = "jwt.token_secret"
	AccessTokenTTL               = "jwt.access_token_ttl"
	RefreshTokenTTL              = "jwt.refresh_token_ttl"
	RefreshResourceTTL           = "refresh_resources_ttl"
	RefreshResourcesInterval     = "refresh_resources_interval"
	RefreshCacheRoutesWorkerMode = "worker_mode.refresh_cache_routes"
)

type Config struct {
	DBName                       string
	DBUser                       string
	DBPassword                   string
	DbHost                       string
	DbPort                       uint16
	SSLMode                      string
	DBMaxCons                    int32
	DBMinCons                    int32
	DBMaxConLifetime             time.Duration
	HTTPServerPort               string
	HTTPProxyPrefix              string
	HTTPAuthPrefix               string
	HTTPResourcePrefix           string
	SwaggerURL                   string
	SwaggerAuthURL               string
	SwaggerResourceURL           string
	SwaggerAuthSpecURL           string
	SwaggerAuthSpecPath          string
	SwaggerResourceSpecURL       string
	SwaggerResourceSpecPath      string
	RedisAddress                 string
	RedisPassword                string
	RedisDB                      int
	AccessTokenTTL               time.Duration
	RefreshTokenTTL              time.Duration
	JWTTokenSecret               string
	RefreshResourceTTL           time.Duration
	RefreshResourceInterval      time.Duration
	RefreshCacheRoutesWorkerMode string
}

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error config file: %w", err)
	}

	return nil
}
