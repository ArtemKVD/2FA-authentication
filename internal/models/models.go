package models

import "time"

type User struct {
	ID         int64  `db:"id"`
	Username   string `db:"username"`
	Password   string `db:"password_hash"`
	TelegramID int64  `db:"telegram_id"`
	ChatID     int64  `db:"chat_id"`
}

type TelegramAuth struct {
	ID       int    `db:"id"`
	UserID   int    `db:"user_id"`
	Username string `db:"username"`
	ChatID   int64  `db:"chat_id"`
}

type AuthSession struct {
	ID        int       `db:"id"`
	UserID    int       `db:"user_id"`
	Token     string    `db:"token"`
	ExpiresAt time.Time `db:"expires_at"`
	Verified  bool      `db:"verified"`
}
