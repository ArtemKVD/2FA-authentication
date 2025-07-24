package services

import (
	postgres "2FA/internal/database"
	"2FA/internal/models"
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	userRepo    *postgres.UserRepository
	codeRepo    *postgres.CodeRepository
	telegrambot TelegramBotInterface
}

func NewAuthService(userRepo *postgres.UserRepository, codeRepo *postgres.CodeRepository, telegrambot TelegramBotInterface) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		codeRepo:    codeRepo,
		telegrambot: telegrambot,
	}
}
func generateRandomCode(length int) (string, error) {
	randomBytes := make([]byte, length)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	encoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	return encoder.EncodeToString(randomBytes)[:length], nil
}

func (s *AuthService) GenerateCode(userID int64) (string, error) {
	code, err := generateRandomCode(5)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(5 * time.Minute)
	err = s.codeRepo.SaveCode(userID, code, expiresAt)
	return code, err
}

func (s *AuthService) VerifyCode(username, code string) (bool, error) {

	user, err := s.userRepo.GetByUsername(context.Background(), username)
	if err != nil {
		return false, err
	}

	valid, err := s.codeRepo.VerifyCode(user.ID, code)
	log.Printf(code)
	if err != nil {
		log.Printf("error check code: %v", err)
		return false, err
	}

	return valid, nil
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

func (s *AuthService) SendAuthCode(username string) error {
	user, err := s.userRepo.GetByUsername(context.Background(), username)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	code, err := s.GenerateCode(user.ID)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	chatID, err := s.userRepo.GetChatID(user.ID)
	if err != nil {
		return fmt.Errorf("failed to get chat ID: %w", err)
	}

	if s.telegrambot == nil {
		return errors.New("telegram bot not initialized")
	}

	return s.telegrambot.SendAuthCode(chatID, code)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	return string(bytes), err
}
