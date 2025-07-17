package services

import (
	models "2FA/internal/models"
	"context"
)

type AuthService interface {
	GenerateCode(telegramID int) (string, error)
	VerifCode(telegramID int, code string) (bool, error)
}

type UserService interface {
	GetAuthByID(ctx context.Context, telegramID int) (*models.TelegramAuth, error)
}
