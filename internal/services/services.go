package services

import (
	postgres "2FA/internal/database"
	"errors"
)

type AuthService struct {
	userRepo *postgres.UserRepository
}

func NewAuthService(userRepo *postgres.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) GenerateCode(telegramID int64) (string, error) {
	return "123456", nil
}

func (s *AuthService) VerifyCode(telegramID int64, code string) (bool, error) {
	if code != "123456" {
		return false, errors.New("invalid code")
	}
	return true, nil
}
