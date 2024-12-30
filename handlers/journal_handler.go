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
	GetByUserID(ctx context.Context, userID string) (models.DreamsByDate, error)
}
type JournalHandler struct {
	dreamGetter DreamGetter
}

func NewJournalHandler(d DreamGetter) *JournalHandler {
	return &JournalHandler{
		dreamGetter: d,
	}
}

func (j *JournalHandler) GetJournal(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return renderHTMX(c, http.StatusOK, pages.HomePagePublic())
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		c.Logger().Error("Invalid claims type")
		return renderHTMX(c, http.StatusOK, pages.HomePagePublic())
	}
	dreams, err := j.dreamGetter.GetByUserID(c.Request().Context(), jwtClaims.UserID)

	if err != nil {
		return renderError(c, http.StatusInternalServerError, "Could not get dreams")
	}

	return renderHTMX(c, http.StatusOK, pages.GetJournal(dreams))
}
