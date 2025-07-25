package config

import (
	"os"
	"time"
)

type Config struct {
	DBURL            string
	TelegramBotToken string
	ServerPort       string
	JWTSecretKey     string        `mapstructure:"JWT_SECRET_KEY"`
	JWTExpiration    time.Duration `mapstructure:"JWT_EXPIRATION"`
}

func Load() *Config {
	return &Config{
		DBURL:            os.Getenv("DB_URL"),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		ServerPort:       os.Getenv("SERVER_PORT"),
		JWTSecretKey:     os.Getenv("JWT_SECRET_KEY"),
		JWTExpiration:    time.Hour,
	}
}
