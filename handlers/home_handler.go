package handlers

import (
	"eatwo/services"
	"eatwo/views/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Home struct{}

func NewHome() *Home {
	return &Home{}
}

func (h *Home) GetHome(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return renderHTMX(c, http.StatusOK, pages.HomePagePublic())
	}

	_, ok := claims.(*services.CustomClaims)
	if !ok {
		c.Logger().Error("Invalid claims type")
		return renderHTMX(c, http.StatusOK, pages.HomePagePublic())
	}
	return renderHTMX(c, http.StatusOK, pages.HomePage())
}

func (h *Home) GetSignUp(c echo.Context) error {
	redirectToHomeWhenLogged(c)
	return renderHTMX(c, http.StatusOK, pages.SignUpPage())
}

func (h *Home) GetLogIn(c echo.Context) error {
	redirectToHomeWhenLogged(c)
	return renderHTMX(c, http.StatusOK, pages.LoginPage())
}

func redirectToHomeWhenLogged(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return nil
	}

	_, ok := claims.(*services.CustomClaims)
	if !ok {
		return nil
	}
	return redirect(c, http.StatusSeeOther, "/")
}
