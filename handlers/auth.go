package handlers

import (
	"context"
	"eatwo/model"
	"eatwo/shared"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserAuthService interface {
	Create(ctx context.Context, signInData *model.UserSignIn) error
	Validate(ctx context.Context, signInData *model.UserSignIn) error
}

type AuthHandler struct {
	userAuthService UserAuthService
}

func NewAuthHandler(userAuthService UserAuthService) *AuthHandler {
	return &AuthHandler{
		userAuthService: userAuthService,
	}
}

func (a AuthHandler) PostUserHandler(c echo.Context) error {
	var signInData model.UserSignIn
	if err := c.Bind(&signInData); err != nil {
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	c.Logger().Debug(fmt.Sprintf("%+v", signInData))

	err := a.userAuthService.Create(c.Request().Context(), &signInData)
	if err != nil {
		c.Logger().Error(err.Error())
		if errors.Is(err, shared.ErrUserWithEmailExist) {
			return echo.NewHTTPError(http.StatusConflict, "User with such email already exist")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create user")
	}

	return c.String(http.StatusCreated, "user created")
}
