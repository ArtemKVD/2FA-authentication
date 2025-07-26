package services

import (
	"2FA/internal/auth"
	"2FA/internal/models"
	"fmt"
	"time"
)

func (s *AuthService) GenerateJWT(userID int64) (string, error) {
	return auth.GenerateToken(userID, s.JwtSecret, s.JwtExpiration)
}

func (s *AuthService) GenerateTokenPair(userID int64) (*models.TokenPair, error) {
	accessToken, err := auth.GenerateToken(userID, s.JwtSecret, s.JwtExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, tokenID, err := auth.GenerateRefreshToken(userID, s.JwtRefreshSecret, s.JwtRefreshExpiration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	err = s.RefreshTokenRepo.Create(&models.RefreshToken{
		Token:     tokenID,
		UserID:    userID,
		ExpiresAt: time.Now().Add(s.JwtRefreshExpiration),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &models.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshTokens(refreshToken string) (*models.TokenPair, error) {
	claims, err := auth.ParseRefreshToken(refreshToken, s.JwtRefreshSecret)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	storedToken, err := s.RefreshTokenRepo.GetByToken(claims.TokenID)
	if err != nil || storedToken == nil {
		return nil, fmt.Errorf("refresh token not found")
	}

	if err := s.RefreshTokenRepo.Delete(storedToken.Token); err != nil {
		return nil, fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return s.GenerateTokenPair(claims.UserID)
}
