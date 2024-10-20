package util

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER" default:"postgres"`
	DBSource             string        `mapstructure:"DB_SOURCE" default:""`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS" default:":8080"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS" default:":9090"`
	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY" `
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION" `
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION`
}

// LoadConfig load configuration from environment variables
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	// app.env
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("cannot read config file:", err)
		return Config{}, err
	}

	err = viper.Unmarshal(&config)
	fmt.Println("config", config)
	return
}
