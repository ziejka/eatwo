package handlers

import (
	"context"
	"eatwo/models"
	"eatwo/services"
	"eatwo/views/pages"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DreamGetter interface {
	GetByUserID(ctx context.Context, userID string) ([]models.Dream, error)
}
type Home struct {
	dreamGetter DreamGetter
}

func NewHome(dreamGetter DreamGetter) *Home {
	return &Home{
		dreamGetter: dreamGetter,
	}
}

func (h *Home) GetProtectedAbout(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), false)
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), false)
	}
	dreams, err := h.dreamGetter.GetByUserID(c.Request().Context(), jwtClaims.UserID)
	if err != nil {
		return renderError(c, http.StatusInternalServerError, "Could not get dreams")
	}

	return renderHTMX(c, http.StatusOK, pages.HomePage(jwtClaims.Name, dreams), true)
}

func (h *Home) GetHome(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return renderHTMX(c, http.StatusOK, pages.HomePagePublic(), false)
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		c.Logger().Error("Invalid claims type")
		return renderHTMX(c, http.StatusOK, pages.HomePagePublic(), false)
	}
	dreams, err := h.dreamGetter.GetByUserID(c.Request().Context(), jwtClaims.UserID)

	if err != nil {
		return renderError(c, http.StatusInternalServerError, "Could not get dreams")
	}

	return renderHTMX(c, http.StatusOK, pages.HomePage(jwtClaims.Name, dreams), true)
}

func (h *Home) GetSignUp(c echo.Context) error {
	redirectToHomeWhenLogged(c)
	return renderHTMX(c, http.StatusOK, pages.SignUpPage(), false)
}

func (h *Home) GetLogIn(c echo.Context) error {
	redirectToHomeWhenLogged(c)
	return renderHTMX(c, http.StatusOK, pages.LoginPage(), false)
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
	return c.Redirect(http.StatusSeeOther, "/")
}
