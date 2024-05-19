package handlers

import (
	"eatwo/views/layouts"
	"eatwo/views/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Home struct{}

func (h *Home) GetHome(c echo.Context) error {
	return Render(c, http.StatusOK, layouts.Base(pages.HomePage()))
}
