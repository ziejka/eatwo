package services

import (
	"context"
	"eatwo/db"
	"eatwo/models"
)

type CheckListRepository interface {
	CreateList(ctx context.Context, arg db.CreateListParams) (db.List, error)
	CreateItem(ctx context.Context, arg db.CreateItemParams) (db.Item, error)
	GetCheckListByUser(ctx context.Context, userID string) ([]db.GetCheckListByUserRow, error)
	GetListWithItemsByListId(ctx context.Context, arg db.GetListWithItemsByListIdParams) ([]db.GetListWithItemsByListIdRow, error)
	GetListIDByUser(ctx context.Context, arg db.GetListIDByUserParams) ([]int64, error)
}

type CheckListService struct {
	checkListRepository CheckListRepository
}

func NewCheckListService(c CheckListRepository) *CheckListService {
	return &CheckListService{
		checkListRepository: c,
	}
}

func (c *CheckListService) GetByUser(ctx context.Context, userID string) ([]models.CheckListRecord, error) {
	lists, err := c.checkListRepository.GetCheckListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	var returnList []models.CheckListRecord
	for _, list := range lists {
		returnList = append(returnList, models.CheckListRecord{
			ID: uint(list.ID),
			CheckList: models.CheckList{
				Name:   list.Name,
				UserID: list.UserID,
			},
		})
	}
	return returnList, nil
}

func (c *CheckListService) CreateCheckList(ctx context.Context, list *models.CheckList) error {
	_, err := c.checkListRepository.CreateList(ctx, db.CreateListParams{
		Name:   list.Name,
		UserID: list.UserID,
	})
	return err
}
func (c *CheckListService) CreateCheckListItem(ctx context.Context, userID string, item *models.CheckListItem) error {
	_, err := c.checkListRepository.GetListIDByUser(ctx, db.GetListIDByUserParams{ID: int64(item.ListID), UserID: userID})
	if err != nil {
		return err
	}
	_, err = c.checkListRepository.CreateItem(ctx, db.CreateItemParams{
		Value:  item.Value,
		ListID: int64(item.ListID),
	})
	return err
}

func (c *CheckListService) GetListById(ctx context.Context, userID string, listID uint) (*models.ListWithItems, error) {
	listWithItems, err := c.checkListRepository.GetListWithItemsByListId(ctx, db.GetListWithItemsByListIdParams{
		UserID: userID,
		ID:     int64(listID),
	})
	if err != nil {
		return nil, err
	}
	list := &models.ListWithItems{
		CheckListRecord: models.CheckListRecord{
			ID: uint(listWithItems[0].ID),
		},
	}

	for _, item := range listWithItems {
		list.Items = append(list.Items, models.CheckListItemRecord{
			ID: uint(item.ID_2),
			CheckListItem: models.CheckListItem{
				ListID: uint(item.ID),
				Value:  item.Value,
			},
		})
	}
	return list, nil
}
