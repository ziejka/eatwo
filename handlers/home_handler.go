package handlers

import (
	"eatwo/views"

	"github.com/labstack/echo/v4"
)

type Home struct{}

func (h *Home) GetHome(c echo.Context) error {
	return views.Layout().Render(c.Request().Context(), c.Response().Writer)
}
