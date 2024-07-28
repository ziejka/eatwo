package handlers

import (
	"context"
	"eatwo/models"
	"eatwo/services"
	"eatwo/views/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CheckListService interface {
	CreateCheckListItem(ctx context.Context, item *models.CheckListItem) error
	CreateCheckList(ctx context.Context, list *models.CheckList) error
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

func (l *CheckListHandler) PostListHandler(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	var listData NewCheckList
	if err := c.Bind(&listData); err != nil {
		c.Logger().Error(err.Error())
		return err
	}
	l.checklistService.CreateCheckList(c.Request().Context(), &models.CheckList{Name: listData.Name, UserID: jwtClaims.UserID})
	return nil
}
