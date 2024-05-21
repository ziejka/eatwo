package handlers

import (
	"eatwo/views/layouts"
	"eatwo/views/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Home struct{}

func (h *Home) GetHome(c echo.Context) error {
	return render(c, http.StatusOK, layouts.Base(pages.HomePage()))
}

func (h *Home) GetSignIn(c echo.Context) error {
	return render(c, http.StatusOK, layouts.Base(pages.SignInPage()))
}

func (h *Home) GetLogIn(c echo.Context) error {
	return render(c, http.StatusOK, layouts.Base(pages.LoginPage()))
}
