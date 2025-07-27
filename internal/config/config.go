package config

import (
	"os"
	"time"
)

type Config struct {
	DBURL                     string
	TelegramBotToken          string
	ServerPort                string
	JWTSecretKey              string
	JWTExpiration             time.Duration
	JWTRefreshSecret          string
	JWTRefreshTokenExpiration time.Duration
}

func Load() *Config {
	return &Config{
		DBURL:                     os.Getenv("DB_URL"),
		TelegramBotToken:          os.Getenv("TELEGRAM_BOT_TOKEN"),
		ServerPort:                os.Getenv("SERVER_PORT"),
		JWTSecretKey:              os.Getenv("JWT_SECRET_KEY"),
		JWTExpiration:             time.Hour,
		JWTRefreshSecret:          os.Getenv("JWTREFRESH_SECRET"),
		JWTRefreshTokenExpiration: time.Hour * 200,
	}
}
