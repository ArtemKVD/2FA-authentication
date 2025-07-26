package database

import (
	"2FA/internal/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type RefreshTokenRepo struct {
	db *sqlx.DB
}

func NewRefreshTokenRepo(db *sqlx.DB) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: db}
}

func (r *RefreshTokenRepo) Create(token *models.RefreshToken) error {
	query := `INSERT INTO refresh_tokens (token, user_id, expires_at) 
              VALUES ($1, $2, $3)`
	_, err := r.db.Exec(query, token.Token, token.UserID, token.ExpiresAt)
	return err
}

func (r *RefreshTokenRepo) GetByToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	err := r.db.Get(&rt,
		`SELECT token, user_id, expires_at 
         FROM refresh_tokens 
         WHERE token = $1`, token)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &rt, err
}

func (r *RefreshTokenRepo) Delete(token string) error {
	_, err := r.db.Exec(
		`DELETE FROM refresh_tokens 
         WHERE token = $1`, token)
	return err
}
