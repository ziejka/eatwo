package handlers

import (
	"eatwo/views/layouts"
	"eatwo/views/pages"
	"net/http"

	"github.com/a-h/templ"
	"github.com/golang-jwt/jwt/v5"
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

func (h *Home) GetProtectedAbout(c echo.Context) error {
	cookieToken, err := c.Cookie("token")
	if err != nil {
		c.Logger().Error(err)
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	claims, err := h.parseToken(cookieToken.Value)
	if err != nil {
		c.Logger().Error(err)
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}
	return renderHTMX(c, http.StatusOK, pages.HomePage(claims.Subject), claims)
}

func (h *Home) GetHome(c echo.Context) error {
	claims, err := h.getClaimsFromCookie(c)
	if err != nil {
		c.Logger().Error(err)
		return renderHTMX(c, http.StatusOK, pages.HomePage(""), nil)
	}
	return renderHTMX(c, http.StatusOK, pages.HomePage(claims.Subject), claims)
}

func (h *Home) getClaimsFromCookie(c echo.Context) (*jwt.RegisteredClaims, error) {
	cookieToken, err := c.Cookie("token")
	if err != nil {
		return nil, err
	}
	return h.parseToken(cookieToken.Value)
}

func (h *Home) GetSignIn(c echo.Context) error {
	return renderHTMX(c, http.StatusOK, pages.SignInPage(), nil)
}

func (h *Home) GetLogIn(c echo.Context) error {
	return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
}

func renderHTMX(c echo.Context, statusCode int, t templ.Component, claims *jwt.RegisteredClaims) error {
	if c.Request().Header.Get("HX-Request") != "true" {
		t = layouts.Base(claims, t)
	}
	return render(c, statusCode, t)
}
