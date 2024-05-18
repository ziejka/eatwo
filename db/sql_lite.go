package db

import (
	"database/sql"
	"eatwo/model"
	"errors"
)

var (
	ErrNotExists = errors.New("user does not exist")
)

type SQLiteRepository struct {
	db *sql.DB
}

func (r SQLiteRepository) GetUserByEmail(email string) (*model.User, error) {
	rows := r.db.QueryRow("SELECT * FROM users WHERE email = ?", email)
	var user model.User

	if err := rows.Scan(&user.Email, &user.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}

	return &user, nil
}
