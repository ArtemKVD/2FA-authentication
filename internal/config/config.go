package config

import "os"

type Config struct {
	DBURL            string
	TelegramBotToken string
	ServerPort       string
}

func Load() *Config {
	return &Config{
		DBURL:            os.Getenv("DB_URL"),
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		ServerPort:       os.Getenv("SERVER_PORT"),
	}
}
