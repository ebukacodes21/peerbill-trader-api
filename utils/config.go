package utils

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	ServerAddr    string        `mapstructure:"SERVER_ADDR"`
	DBDriver      string        `mapstructure:"DB_DRIVER"`
	DBSource      string        `mapstructure:"DB_SOURCE"`
	EmailSender   string        `mapstructure:"EMAIL_SENDER"`
	EmailAddress  string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailPassword string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	TokenAccess   time.Duration `mapstructure:"TOKEN_ACCESS"`
	TokenKey      string        `mapstructure:"TOKEN_KEY"`
}

func LoadConfig(pathname string) (config Config, err error) {
	viper.AddConfigPath(pathname)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()

	if err != nil {
		log.Fatal(err)
	}

	err = viper.Unmarshal(&config)
	return
}
