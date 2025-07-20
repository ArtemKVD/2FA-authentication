package postgres

import (
	"context"
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

func (repo *UserRepository) GetByID(ctx context.Context, telegramID int64) (*models.User, error) {
	var user models.User

	err := repo.database.conn.GetContext(ctx, &user, `
        SELECT id, username 
        FROM users 
        WHERE telegram_id = $1
    `, telegramID)

	if err != nil {
		log.Printf("error db request")
		return nil, err
	}

	return &user, nil
}
