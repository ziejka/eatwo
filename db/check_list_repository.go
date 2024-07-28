package db

import (
	"context"
	"database/sql"
	"eatwo/models"
)

type CheckListRepository struct {
	db *sql.DB
}

func NewCheckListRepository(db *sql.DB) *CheckListRepository {
	return &CheckListRepository{
		db,
	}
}

type CheckListWithItemsRecord struct {
	listID    uint
	listName  string
	itemID    uint
	itemValue string
}

func (l *CheckListRepository) GetByUser(ctx context.Context, userID string) (*models.ListWithItems, error) {
	rows, err := l.db.QueryContext(ctx,
		`SELECT l.id, l.name, i.id, i."value"
		FROM lists l
		JOIN items i ON i.list_id = l.id
		WHERE l.user_id = ?;
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listWithItems models.ListWithItems
	for rows.Next() {
		var lwi CheckListWithItemsRecord
		if err := rows.Scan(&lwi.listID, &lwi.listName, &lwi.itemID, &lwi.itemValue); err != nil {
			return nil, err
		}
		listWithItems.ID = lwi.listID
		listWithItems.Name = lwi.listName
		listWithItems.Items = append(listWithItems.Items, models.CheckListItemRecord{
			ID: lwi.itemID,
			CheckListItem: models.CheckListItem{
				Value:  lwi.itemValue,
				ListID: lwi.listID,
			},
		})
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &listWithItems, nil

}

func (l *CheckListRepository) Create(ctx context.Context, checkList *models.CheckList) error {
	_, err := l.db.ExecContext(ctx, `INSERT INTO lists (name, user_id) VALUES (?, ?)`, checkList.Name, checkList.UserID)
	return err
}

func (l *CheckListRepository) CreateItem(ctx context.Context, checklistItem *models.CheckListItem) error {
	_, err := l.db.ExecContext(ctx, `INSERT INTO item (value, list_id) VALUES (?, ?)`, checklistItem.Value, checklistItem.ListID)
	return err
}

func (l *CheckListRepository) Migrate(ctx context.Context) error {
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
