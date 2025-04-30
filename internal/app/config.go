package app

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	dbName              = "database.name"
	dbUser              = "database.user"
	dbPassword          = "database.password"
	dbHost              = "database.host"
	dbPort              = "database.port"
	sslMode             = "database.sslmode"
	httpServerPort      = "server.http.port"
	swaggerURLPath      = "swagger.url_path"
	swaggerSpecFilePath = "swagger.spec_file_path"
	redisAddress        = "redis.addr"
	redisPassword       = "redis.password"
	redisDb             = "redis.db"
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

func initConfig() *Config {
	return &Config{
		DBName:              viper.GetString(dbName),
		DBUser:              viper.GetString(dbUser),
		DBPassword:          viper.GetString(dbPassword),
		DbHost:              viper.GetString(dbHost),
		DbPort:              viper.GetUint16(dbPort),
		SSLMode:             viper.GetString(sslMode),
		HTTPServerPort:      viper.GetString(httpServerPort),
		SwaggerURLPath:      viper.GetString(swaggerURLPath),
		SwaggerSpecFilePath: viper.GetString(swaggerSpecFilePath),
		RedisAddress:        viper.GetString(redisAddress),
		RedisPassword:       viper.GetString(redisPassword),
		RedisDB:             viper.GetInt(redisDb),
	}
}

func loadConfig() error {
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
