package models

import "time"

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	UserID       int64  `json:"user_id"`
}

type RefreshToken struct {
	Token     string    `db:"token"`
	UserID    int64     `db:"user_id"`
	ExpiresAt time.Time `db:"expires_at"`
}
