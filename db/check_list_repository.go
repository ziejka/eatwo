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

type CheckListIDRecor struct {
	listID uint
}

func (l *CheckListRepository) GetByUser(ctx context.Context, userID string) ([]models.CheckListRecord, error) {
	rows, err := l.db.QueryContext(ctx,
		`SELECT l.id, l.user_id, l.name
		FROM lists l
		WHERE l.user_id = ?;
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var checkLists []models.CheckListRecord
	for rows.Next() {
		var checkList models.CheckListRecord

		if err := rows.Scan(&checkList.ID, &checkList.UserID, &checkList.Name); err != nil {
			return nil, err
		}
		checkLists = append(checkLists, checkList)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return checkLists, nil

}

func (l *CheckListRepository) GetListWithItemsById(ctx context.Context, userID string, listID uint) (*models.ListWithItems, error) {
	rows, err := l.db.QueryContext(ctx,
		`SELECT l.id, l.name, i.id, i."value"
		FROM lists l
		JOIN items i ON i.list_id = l.id
		WHERE l.id = ? AND l.user_id = ?;
	`, listID, userID)

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

func (l *CheckListRepository) GetListIDByUser(ctx context.Context, userID string, listID uint) (uint, error) {
	row := l.db.QueryRowContext(ctx, "SELECT id FROM lists WHERE id = ? AND user_id = ?", listID, userID)
	var listIDRecors CheckListIDRecor
	if err := row.Scan(&listIDRecors.listID); err != nil {
		return 0, err
	}

	return listIDRecors.listID, nil
}

func (l *CheckListRepository) CreateItem(ctx context.Context, checklistItem *models.CheckListItem) error {
	_, err := l.db.ExecContext(ctx, `INSERT INTO items (value, list_id) VALUES (?, ?)`, checklistItem.Value, checklistItem.ListID)
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
