package handlers

import (
	"eatwo/services"
	"eatwo/views/layouts"
	"eatwo/views/pages"
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

type Home struct{}

func NewHome() *Home {
	return &Home{}
}

func (h *Home) GetProtectedAbout(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	return renderHTMX(c, http.StatusOK, pages.HomePage(jwtClaims.Subject), jwtClaims)
}

func (h *Home) GetHome(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		c.Logger().Error("Unauthorize")
		return renderHTMX(c, http.StatusOK, pages.HomePage(""), nil)
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		c.Logger().Error("Invalid claims type")
		return renderHTMX(c, http.StatusOK, pages.HomePage(""), nil)
	}
	return renderHTMX(c, http.StatusOK, pages.HomePage(jwtClaims.Subject), jwtClaims)
}

func (h *Home) GetSignIn(c echo.Context) error {
	return renderHTMX(c, http.StatusOK, pages.SignInPage(), nil)
}

func (h *Home) GetLogIn(c echo.Context) error {
	return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
}

func renderHTMX(c echo.Context, statusCode int, t templ.Component, claims *services.CustomClaims) error {
	if c.Request().Header.Get("HX-Request") != "true" {
		t = layouts.Base(claims, t)
	}
	return render(c, statusCode, t)
}
