package handlers

import (
	"eatwo/views/layouts"
	"eatwo/views/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Home struct {
	parseToken TokenParser
}

func NewHome(parseToken TokenParser) *Home {
	return &Home{
		parseToken,
	}
}

func (h *Home) GetHome(c echo.Context) error {
	cookieToken, err := c.Cookie("token")
	if err != nil {
		c.Logger().Error(err)
		return render(c, http.StatusOK, layouts.Base(pages.LoginPage()))
	}

	claims, err := h.parseToken(cookieToken.Value)
	if err != nil {
		c.Logger().Error(err)
		return render(c, http.StatusOK, layouts.Base(pages.LoginPage()))
	}
	return render(c, http.StatusOK, layouts.Base(pages.HomePage(claims.Subject)))
}

func (h *Home) GetSignIn(c echo.Context) error {
	return render(c, http.StatusOK, layouts.Base(pages.SignInPage()))
}

func (h *Home) GetLogIn(c echo.Context) error {
	return render(c, http.StatusOK, layouts.Base(pages.LoginPage()))
}
