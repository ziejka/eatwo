package services

import (
	"context"
	"eatwo/models"
)

type CheckListRepository interface {
	Create(ctx context.Context, checkList *models.CheckList) error
	CreateItem(ctx context.Context, checklistItem *models.CheckListItem) error
	GetByUser(ctx context.Context, email string) (*models.ListWithItems, error)
}

type CheckListService struct {
	checkListRepository CheckListRepository
}

func NewCheckListService(c CheckListRepository) *CheckListService {
	return &CheckListService{
		checkListRepository: c,
	}
}

func (c *CheckListService) GetByUser(ctx context.Context, userID string) (*models.ListWithItems, error) {
	return c.checkListRepository.GetByUser(ctx, userID)
}

func (c *CheckListService) CreateCheckList(ctx context.Context, list *models.CheckList) error {
	return c.checkListRepository.Create(ctx, list)
}

func (c *CheckListService) CreateCheckListItem(ctx context.Context, item *models.CheckListItem) error {
	return c.checkListRepository.CreateItem(ctx, item)
}
