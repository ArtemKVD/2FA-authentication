package postgres

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type CodeRepository struct {
	db *sqlx.DB
}

func NewCodeRepository(db *sqlx.DB) *CodeRepository {
	return &CodeRepository{
		db: db,
	}
}

func (r *CodeRepository) SaveCode(userID int64, code string, expiresAt time.Time) error {
	query := `INSERT INTO auth_codes (user_id, code, expires_at) 
              VALUES ($1, $2, $3) 
              ON CONFLICT (user_id) DO UPDATE 
              SET code = $2, expires_at = $3`
	_, err := r.db.Exec(query, userID, code, expiresAt)
	return err
}

func (r *CodeRepository) GetCode(userID int64) (string, time.Time, error) {
	var code string
	var expiresAt time.Time
	query := `SELECT code, expires_at FROM auth_codes WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&code, &expiresAt)
	return code, expiresAt, err
}
