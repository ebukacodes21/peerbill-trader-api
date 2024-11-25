package utils

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPServerAddr  string        `mapstructure:"HTTP_SERVER_ADDR"`
	GRPCServerAddr  string        `mapstructure:"GRPC_SERVER_ADDR"`
	REDISServerAddr string        `mapstructure:"REDIS_SERVER_ADDR"`
	WebsocketAddr   string        `mapstructure:"WEBSOCKET_SERVER_ADDR"`
	AllowedOrigins  []string      `mapstructure:"ALLOWED_ORIGINS"`
	MigrationURL    string        `mapstructure:"MIGRATION_URL"`
	DBDriver        string        `mapstructure:"DB_DRIVER"`
	DBSource        string        `mapstructure:"DB_SOURCE"`
	EmailSender     string        `mapstructure:"EMAIL_SENDER"`
	EmailAddress    string        `mapstructure:"EMAIL_SENDER_ADDRESS"`
	EmailPassword   string        `mapstructure:"EMAIL_SENDER_PASSWORD"`
	TokenKey        string        `mapstructure:"TOKEN_KEY"`
	TokenAccess     time.Duration `mapstructure:"TOKEN_ACCESS"`
	RefreshAccess   time.Duration `mapstructure:"REFRESH_ACCESS"`
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
