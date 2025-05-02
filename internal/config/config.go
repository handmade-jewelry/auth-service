package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	DbName                 = "database.name"
	DbUser                 = "database.user"
	DbPassword             = "database.password"
	DbHost                 = "database.host"
	DbPort                 = "database.port"
	SslMode                = "database.sslmode"
	HttpServerPort         = "server.http.port"
	SwaggerURLPath         = "swagger.url_path"
	SwaggerSpecFilePath    = "swagger.spec_file_path"
	RedisAddress           = "redis.addr"
	RedisPassword          = "redis.password"
	RedisDb                = "redis.db"
	GRPCClientMaxRetry     = "grpc_client.max_retry"
	GRPCClientRetryTimeout = "grpc_client.timeout_per_retry"
	UserServiceHost        = "user_service.host"
	JWTTokenSecret         = "jwt.token_secret"
	AccessTokenExpMin      = "jwt.access_token_expiry_minutes"
	RefreshTokenExpMin     = "jwt.refresh_token_expiry_minutes"
)

type Config struct {
	DBName              string
	DBUser              string
	DBPassword          string
	DbHost              string
	DbPort              uint16
	SSLMode             string
	HTTPServerPort      string
	SwaggerURLPath      string
	SwaggerSpecFilePath string
	RedisAddress        string
	RedisPassword       string
	RedisDB             int
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
