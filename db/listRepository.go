package db

import (
	"context"
	"database/sql"
	"eatwo/models"
)

type ListRepository struct {
	db *sql.DB
}

func NewListRepository(db *sql.DB) *ListRepository {
	return &ListRepository{
		db,
	}
}

func (l *ListRepository) GetByUser(ctx context.Context, userID string) {
	l.db.ExecContext(ctx,
		`SELECT l.id, l.name, i.id, i."value"
		FROM lists l
		JOIN items i ON i.list_id = l.id
		WHERE l.user_id == ?;
	`, userID)
}

func (l *ListRepository) Create(ctx context.Context, list models.List) error {
	_, err := l.db.ExecContext(ctx, `INSERT INTO lists (name, user_id) VALUES (?, ?)`, list.Name, list.UserID)
	return err
}

func (l *ListRepository) CreateItem(ctx context.Context, item models.Item) error {
	_, err := l.db.ExecContext(ctx, `INSERT INTO item (value, list_id) VALUES (?, ?)`, item.Value, item.ListID)
	return err
}

func (l *ListRepository) Migrate(ctx context.Context) error {
	tx, err := l.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS lists (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		user_id TEXT,
		FOREIGN KEY(user_id) REFERENCES users(id),
		UNIQUE(name, user_id)
	)`)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		value TEXT,
		list_id INTEGER,
        FOREIGN KEY(list_id) REFERENCES lists(id)
	)`)
	if err != nil {
		return err
	}

	return tx.Commit()
}
