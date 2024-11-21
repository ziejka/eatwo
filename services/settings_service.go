package services

import (
	"context"
	"database/sql"
	"eatwo/db"
	"eatwo/models"
)

type SettingsService struct {
	queries *db.Queries
	db      *sql.DB
}

func NewSettingsService(queries *db.Queries, db *sql.DB) *SettingsService {
	return &SettingsService{
		queries,
		db,
	}
}

func (s *SettingsService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.queries.GetUserByID(ctx, userID)

	if err != nil {
		return nil, err
	}

	return user.ToModel(), nil
}

func (s *SettingsService) UpdateUser(ctx context.Context, userUpdateData models.UserUpdate, userID string) (*models.User, error) {
	user, err := s.queries.UpdateUser(ctx, db.UpdateUserParams{
		Email: userUpdateData.Email,
		Name:  userUpdateData.Name,
		ID:    userID,
	})

	if err != nil {
		return nil, err
	}

	return user.ToModel(), nil
}

func (s *SettingsService) DeleteUserAndData(ctx context.Context, userID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.queries.WithTx(tx)

	err = qtx.DeleteDreamsForUser(ctx, userID)
	if err != nil {
		return err
	}
	err = qtx.DeleteItemsForUser(ctx, userID)
	if err != nil {
		return err
	}
	err = qtx.DeleteListsForUser(ctx, userID)
	if err != nil {
		return err
	}
	err = qtx.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}
	return tx.Commit()
}
