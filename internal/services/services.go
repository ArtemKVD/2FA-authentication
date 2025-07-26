package services

import (
	postgres "2FA/internal/database"
	"2FA/internal/models"
	"context"
	"errors"
	"log"
	"time"
)

const (
	hashCost = 16
)

type TelegramBotInterface interface {
	SendAuthCode(chatID int64, code string) error
}

func (s *AuthService) GetByUsername(username string) (*models.User, error) {
	return s.userRepo.GetByUsername(context.Background(), username)
}

type AuthService struct {
	userRepo             *postgres.UserRepository
	codeRepo             *postgres.CodeRepository
	RefreshTokenRepo     *postgres.RefreshTokenRepo
	telegrambot          TelegramBotInterface
	JwtSecret            string
	JwtExpiration        time.Duration
	JwtRefreshSecret     string
	JwtRefreshExpiration time.Duration
}

func NewAuthService(userRepo *postgres.UserRepository, codeRepo *postgres.CodeRepository,
	telegrambot TelegramBotInterface, JwtSecret string, JwtExpiration time.Duration, JwtRefreshSecret string, JwtRefreshExpiration time.Duration, RefreshTokenRepo *postgres.RefreshTokenRepo) *AuthService {
	return &AuthService{
		userRepo:             userRepo,
		codeRepo:             codeRepo,
		RefreshTokenRepo:     RefreshTokenRepo,
		telegrambot:          telegrambot,
		JwtSecret:            JwtSecret,
		JwtExpiration:        JwtExpiration,
		JwtRefreshSecret:     JwtRefreshSecret,
		JwtRefreshExpiration: JwtRefreshExpiration,
	}
}

func (s *AuthService) DataVerify(username, password string) (bool, error) {

	user, err := s.userRepo.GetByUsername(context.Background(), username)
	if err != nil {
		return false, err
	}
	log.Printf("User login")
	passwordcheck, err := hashPassword(password)
	if err != nil {
		log.Printf("hash password error")
	}
	if user.Password != passwordcheck {
		return false, errors.New("invalid password")
	}

	return true, nil
}
