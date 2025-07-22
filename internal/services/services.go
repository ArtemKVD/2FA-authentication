package services

import (
	postgres "2FA/internal/database"
	"2FA/internal/models"
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"log"
	"time"
)

type TelegramBotInterface interface {
	SendCode(chatID int64, code string) error
}

func (s *AuthService) GetByUsername(username string) (*models.User, error) {
	return s.userRepo.GetByUsername(context.Background(), username)
}

type AuthService struct {
	userRepo *postgres.UserRepository
	codeRepo *postgres.CodeRepository
}

func NewAuthService(
	userRepo *postgres.UserRepository,
	codeRepo *postgres.CodeRepository,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		codeRepo: codeRepo,
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
	code, err := generateRandomCode(4)
	if err != nil {
		return "", err
	}

	expiresAt := time.Now().Add(5 * time.Minute)
	err = s.codeRepo.SaveCode(userID, code, expiresAt)
	return code, err
}

func (s *AuthService) VerifyCode(userID int64, code string) (bool, error) {
	savedCode, _, err := s.codeRepo.GetCode(userID)
	if err != nil {
		return false, err
	}

	return code == savedCode, nil
}

func (s *AuthService) DataVerify(username, password string) (bool, error) {

	user, err := s.userRepo.GetByUsername(context.Background(), username)
	if err != nil {
		return false, err
	}
	log.Printf("User login")

	if user.Password != hashPassword(password) {
		return false, errors.New("invalid password")
	}

	return true, nil
}

func hashPassword(password string) string {
	//later
	return password
}
