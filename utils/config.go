package utils

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
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
