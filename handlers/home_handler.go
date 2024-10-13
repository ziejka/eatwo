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
	GetByUserID(ctx context.Context, userID string) ([]*models.DreamRecord, error)
}
type Home struct{
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
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
	}
  dreams, err := h.dreamGetter.GetByUserID(c.Request().Context(), jwtClaims.UserID)
  if err != nil {
    return renderError(c, http.StatusInternalServerError, "Could not get dreams")
  }

	return renderHTMX(c, http.StatusOK, pages.HomePage(jwtClaims.Email, dreams), jwtClaims)
}

func (h *Home) GetHome(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return renderHTMX(c, http.StatusOK, pages.HomePagePublic(), nil)
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		c.Logger().Error("Invalid claims type")
		return renderHTMX(c, http.StatusOK, pages.HomePagePublic(), nil)
	}
  dreams, err := h.dreamGetter.GetByUserID(c.Request().Context(), jwtClaims.UserID)
  if err != nil {
    return renderError(c, http.StatusInternalServerError, "Could not get dreams")
  }

	return renderHTMX(c, http.StatusOK, pages.HomePage(jwtClaims.Email, dreams), jwtClaims)
}

func (h *Home) GetSignUp(c echo.Context) error {
	return renderHTMX(c, http.StatusOK, pages.SignUpPage(), nil)
}

func (h *Home) GetLogIn(c echo.Context) error {
	return renderHTMX(c, http.StatusOK, pages.LoginPage(), nil)
}
