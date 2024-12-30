package handlers

import (
	"context"
	"eatwo/models"
	"eatwo/services"
	"eatwo/views/pages"

	"net/http"

	"github.com/labstack/echo/v4"
)

type SettingsService interface {
	UpdateUser(ctx context.Context, userUpdateData models.UserUpdate, userID string) (*models.User, error)
	GetUser(ctx context.Context, userID string) (*models.User, error)
	DeleteUserAndData(ctx context.Context, userID string) error
}

type SettingsHandler struct {
	settingsService SettingsService
}

func NewSettingsHandler(settingsService SettingsService) *SettingsHandler {
	return &SettingsHandler{
		settingsService,
	}
}

func (sh *SettingsHandler) GetAccountSettings(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return renderHTMX(c, http.StatusOK, pages.LoginPage())
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		return renderHTMX(c, http.StatusOK, pages.LoginPage())
	}

	user, err := sh.settingsService.GetUser(c.Request().Context(), jwtClaims.UserID)
	if err != nil {
		c.Logger().Error(err)
		return renderError(c, http.StatusInternalServerError, "Failed to get user")
	}

	return renderHTMX(c, http.StatusOK, pages.AccountSettings(*user))
}

func (sh *SettingsHandler) PostUserUpdate(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return renderHTMX(c, http.StatusOK, pages.LoginPage())
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		return renderHTMX(c, http.StatusOK, pages.LoginPage())
	}

	var userUpdate models.UserUpdate
	if err := c.Bind(&userUpdate); err != nil {
		return renderError(c, http.StatusBadRequest, "Invalid input")
	}

	user, err := sh.settingsService.UpdateUser(c.Request().Context(), userUpdate, jwtClaims.UserID)
	if err != nil {
		c.Logger().Error(err)
		return renderError(c, http.StatusInternalServerError, "Failed to update user")
	}

	return renderHTMX(c, http.StatusOK, pages.AccountSettings(*user))
}

func (sh *SettingsHandler) DeleteUser(c echo.Context) error {
	claims := c.Get("claims")
	if claims == nil {
		return renderHTMX(c, http.StatusOK, pages.LoginPage())
	}

	jwtClaims, ok := claims.(*services.CustomClaims)
	if !ok {
		return renderHTMX(c, http.StatusOK, pages.LoginPage())
	}

	err := sh.settingsService.DeleteUserAndData(c.Request().Context(), jwtClaims.UserID)
	if err != nil {
		c.Logger().Error(err)
		return renderError(c, http.StatusInternalServerError, "Failed to delete user")
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
	}

	c.SetCookie(cookie)
	// Probably should redirect to a page that says "Your account has been deleted"
	return redirect(c, http.StatusSeeOther, "/")
}
