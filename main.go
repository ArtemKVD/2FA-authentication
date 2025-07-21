package main

import (
	"log"
	"net/http"

	"2FA/internal/config"
	postgres "2FA/internal/database"
	"2FA/internal/handlers"
	"2FA/internal/server"
	"2FA/internal/services"
	"2FA/internal/telegram"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db, err := sqlx.Connect("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("connect error: %v", err)
	}
	defer db.Close()

	log.Println("postgres connect succes")

	userRepo := postgres.NewUserRepository(db)
	codeRepo := postgres.NewCodeRepository(db)

	authService := services.NewAuthService(userRepo, codeRepo)
	authHandler := handlers.NewAuthHandler(*authService)
	bot, err := telegram.BotCreate(cfg.TelegramBotToken, *authService)
	if err != nil {
		log.Fatalf("Failed to create Telegram bot: %v", err)
	}

	log.Printf("Authorized Telegram bot")

	router := server.NewServer(authHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	log.Println("Server exited properly")
}
