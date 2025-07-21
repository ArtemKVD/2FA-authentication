package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	models "2FA/internal/models"
)

type Database struct {
	conn *sqlx.DB
}

func NewDatabase(connectionString string) (*Database, error) {
	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		log.Printf("error connect to db")
		return nil, err
	}

	return &Database{
		conn: db,
	}, nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}

type UserRepository struct {
	database *Database
}

func NewUserRepository(db *Database) *UserRepository {
	return &UserRepository{
		database: db,
	}
}

func (r *UserRepository) GetChatID(userID int64) (int64, error) {
	const query = `SELECT chat_id FROM users WHERE id = $1`

	var chatID int64
	err := r.database.conn.Get(&chatID, query, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("user not found")
		}
		return 0, fmt.Errorf("failed to get chat_id: %w", err)
	}

	return chatID, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, password_hash, telegram_id, chat_id FROM users WHERE username = $1`
	var user models.User
	err := r.database.conn.GetContext(ctx, &user, query, username)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	return &user, err
}
