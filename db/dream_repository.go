package db

import (
	"context"
	"database/sql"
	"eatwo/models"
	"fmt"
	"time"
)

type DreamRepository struct {
	db *sql.DB
}

func NewDreamRepository(db *sql.DB) *DreamRepository {
	return &DreamRepository{
		db: db,
	}
}

func (r *DreamRepository) Create(ctx context.Context, dream *models.DreamRecord) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO dreams (id, user_id, description, explanation, date) VALUES (?, ?, ?, ?, ?)",
		dream.ID, dream.UserID, dream.Description, dream.Explanation, dream.Date.Format(time.RFC822))

	return err
}

func (r *DreamRepository) UpdateExplanation(ctx context.Context, dreamID, explanation string, userID string) error {
	_, err := r.db.ExecContext(ctx, "UPDATE dreams SET explanation = ? WHERE id = ? AND user_id = ?", explanation, dreamID, userID)
	return err
}

func (r *DreamRepository) GetByUserID(ctx context.Context, userID string) ([]*models.DreamRecord, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT * FROM dreams WHERE user_id = ? ORDER BY date", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dreams := make([]*models.DreamRecord, 0)
	var dateStr string

	for rows.Next() {
		dream := new(models.DreamRecord)
		if err := rows.Scan(&dream.ID, &dream.UserID, &dream.Description, &dream.Explanation, &dateStr); err != nil {
			return nil, err
		}
		dateTime, err := time.Parse(time.RFC822, dateStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing date: %v", err)
		}

		dream.Date = dateTime
		dreams = append(dreams, dream)
	}

	return dreams, nil
}

func (r *DreamRepository) Migrate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS dreams (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    description TEXT NOT NULL,
    explanation TEXT NOT NULL,
    date TEXT NOT NULL
  )`)

	return err
}
