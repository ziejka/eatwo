package db

import (
	"context"
	"database/sql"
	"eatwo/models"
	"eatwo/shared"
	"errors"
	"fmt"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r UserRepository) GetByEmail(ctx context.Context, email string) (*models.UserRecord, error) {
	rows := r.db.QueryRow("SELECT * FROM users WHERE email = ?", email)
	var user models.UserRecord

	if err := rows.Scan(&user.Email, &user.Name, &user.HashPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, shared.ErrNotExists
		}
		return nil, err
	}

	return &user, nil
}

func (r UserRepository) Create(ctx context.Context, user *models.UserRecord) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (email, name, hash_password) VALUES (?, ?, ?)",
		user.Email, user.Name, user.HashPassword)

	if err != nil {
		return fmt.Errorf("CreateUser: %v", err)
	}

	return nil
}

func (r UserRepository) Migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS users(
		email TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		hash_password TEXT NOT NULL
	)`)

	return err
}
