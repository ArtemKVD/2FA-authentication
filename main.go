package main

import (
	"context"
	"log"

	config "2FA/internal/config"
	telegram "2FA/internal/telegram"
)

func main() {
	cfg := config.Load()
	bot, err := telegram.BotCreate(cfg.TelegramBotToken, nil)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go bot.Start(ctx)
}
