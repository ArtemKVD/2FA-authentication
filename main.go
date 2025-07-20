package main

import (
	"log"

	"2FA/internal/config"
	postgres "2FA/internal/database"
	"2FA/internal/handlers"
	"2FA/internal/server"
	"2FA/internal/services"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	db, err := postgres.NewDatabase(cfg.DBURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	authService := services.NewAuthService(postgres.NewUserRepository(db))
	authHandler := handlers.NewAuthHandler(*authService)

	srv := server.NewServer(authHandler)
	if err := srv.Run(cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
