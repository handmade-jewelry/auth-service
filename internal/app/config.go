package app

import (
	"fmt"
	"github.com/spf13/viper"
)

const (
	jwtTokenSecret     = "auth.jwt_token_secret"
	authTokenExpMin    = "auth.auth_token_expiry_minutes"
	refreshTokenExpMin = "auth.refresh_token_expiry_minutes"
	dbName             = "database.name"
	dbUser             = "database.user"
	dbPassword         = "database.password"
	dbHost             = "database.host"
	dbPort             = "database.port"
	sslMode            = "database.sslmode"
)

func initConfig() error {
	viper.SetConfigName("config")   // name of config file (without extension)
	viper.SetConfigType("yaml")     // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("./config") // path to look for the config file in

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	return nil
}
