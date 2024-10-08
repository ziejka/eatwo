package handlers

import (
	"context"
	"eatwo/models"
	"eatwo/services"
	"eatwo/views/pages"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CheckListService interface {
	CreateCheckListItem(ctx context.Context, userID string, item *models.CheckListItem) error
	CreateCheckList(ctx context.Context, list *models.CheckList) error
	GetByUser(ctx context.Context, userID string) ([]models.CheckListRecord, error)
	GetListById(ctx context.Context, userID string, listID uint) (*models.ListWithItems, error)
}

type CheckListHandler struct {
	checklistService CheckListService
}

func NewCheckListHandler(checklistService CheckListService) *CheckListHandler {
	return &CheckListHandler{
		checklistService: checklistService,
	}
}

type NewCheckList struct {
	Name string `form:"name"`
}

type NewCheckListItem struct {
	Value string `form:"value"`
}

func (l *CheckListHandler) GetCheckListHandler(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	listIDStr := c.Param("id")
	listID, err := strconv.ParseUint(listIDStr, 10, 0)
	if err != nil {
		c.Logger().Errorf("Incorrect list id %", listID)
		return err
	}

	checkList, err := l.checklistService.GetListById(c.Request().Context(), jwtClaims.UserID, uint(listID))
	if err != nil {
		c.Logger().Error(err)
	}
  checkList.ID = uint(listID)
	return renderHTMX(c, http.StatusOK, pages.CheckList(checkList), jwtClaims)
}

func (l *CheckListHandler) GetCheckListsHandler(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	checkLists, err := l.checklistService.GetByUser(c.Request().Context(), jwtClaims.UserID)
	if err != nil {
		c.Logger().Error(err)
		checkLists = []models.CheckListRecord{}
	}
	return renderHTMX(c, http.StatusOK, pages.CheckLists(checkLists), jwtClaims)
}

func (l *CheckListHandler) PostItemHandler(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		// TODO: Render error for API
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		// TODO: Render error for API
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	listIDStr := c.Param("id")
	listID, err := strconv.ParseUint(listIDStr, 10, 0)
	if err != nil {
		c.Logger().Errorf("Incorrect list id %", listID)
		return err
	}

	var listData NewCheckListItem
	if err := c.Bind(&listData); err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	err = l.checklistService.CreateCheckListItem(
		c.Request().Context(),
		jwtClaims.UserID,
		&models.CheckListItem{
			Value:  listData.Value,
			ListID: uint(listID),
		})

	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return l.GetCheckListHandler(c)
}

func (l *CheckListHandler) PostCheckListHandler(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		// TODO: Render error for API
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		// TODO: Render error for API
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	var listData NewCheckList
	if err := c.Bind(&listData); err != nil {
		c.Logger().Error(err.Error())
		return err
	}

	err := l.checklistService.CreateCheckList(
		c.Request().Context(),
		&models.CheckList{Name: listData.Name, UserID: jwtClaims.UserID})
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	return l.GetCheckListsHandler(c)
}
