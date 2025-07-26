package main

import (
	"2FA/internal/config"
	database "2FA/internal/database"
	"2FA/internal/handlers"
	"2FA/internal/server"
	"2FA/internal/services"
	"2FA/internal/telegram"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	cfg := config.Load()

	db, err := sqlx.Connect("postgres", cfg.DBURL)
	if err != nil {
		log.Println("fail connect db", err)
	}
	defer db.Close()

	log.Println("postgres connected")

	codeRepo := database.NewCodeRepository(db)
	userRepo := database.NewUserRepository(db)
	RefreshTokenRepo := database.NewRefreshTokenRepo(db)
	bot, err := telegram.BotCreate(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("failed create bot: %v", err)
	}
	log.Println("Authorized Telegram bot", bot.Api.Self.UserName)

	authService := services.NewAuthService(
		userRepo,
		codeRepo,
		bot,
		cfg.JWTSecretKey,
		cfg.JWTExpiration,
		cfg.JWTRefreshSecret,
		cfg.JWTRefreshTokenExpiration,
		RefreshTokenRepo,
	)

	authHandler := handlers.NewAuthHandler(*authService)
	router := server.NewServer(authHandler, authService, cfg.JWTSecretKey, cfg.JWTRefreshSecret)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()
	log.Println("server started on port", cfg.ServerPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

}
