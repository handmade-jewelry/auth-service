package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

const (
	DBName                   = "database.name"
	DBUser                   = "database.user"
	DBPassword               = "database.password"
	DBHost                   = "database.host"
	DBPort                   = "database.port"
	SSLMode                  = "database.ssl_mode"
	DBMaxCons                = "database.max_cons"
	DBMinCons                = "database.min_cons"
	DBMaxConLifetime         = "database.max_con_lifetime"
	HTTPServerPort           = "server.http.port"
	SwaggerURLPath           = "swagger.url_path"
	SwaggerSpecFilePath      = "swagger.spec_file_path"
	RedisAddress             = "redis.addr"
	RedisPassword            = "redis.password"
	RedisDb                  = "redis.db"
	GRPCClientMaxRetry       = "grpc_client.max_retry"
	GRPCClientRetryTimeout   = "grpc_client.timeout_per_retry"
	UserServiceHost          = "user_service.host"
	JWTTokenSecret           = "jwt.token_secret"
	AccessTokenExpMin        = "jwt.access_token_expiry_minutes"
	RefreshTokenExpMin       = "jwt.refresh_token_expiry_minutes"
	RefreshResourceTTL       = "refresh_resources_ttl"
	RefreshResourcesInterval = "refresh_resources_interval"
)

type Config struct {
	DBName                  string
	DBUser                  string
	DBPassword              string
	DbHost                  string
	DbPort                  uint16
	SSLMode                 string
	DBMaxCons               int32
	DBMinCons               int32
	DBMaxConLifetime        time.Duration
	HTTPServerPort          string
	SwaggerURLPath          string
	SwaggerSpecFilePath     string
	RedisAddress            string
	RedisPassword           string
	RedisDB                 int
	AccessTokenExp          time.Duration
	RefreshTokenExp         time.Duration
	JWTTokenSecret          string
	RefreshResourceTTL      time.Duration
	RefreshResourceInterval time.Duration
}

func LoadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		//todo panic?..
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return nil
}
