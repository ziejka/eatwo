package handlers

import (
	"context"
	"eatwo/models"
	"eatwo/shared"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserAuthService interface {
	Create(ctx context.Context, signInData models.UserSignIn) error
	LogIn(ctx context.Context, signInData models.UserLogIn) (string, error)
}

type AuthHandler struct {
	userAuthService UserAuthService
}

func NewAuthHandler(userAuthService UserAuthService) *AuthHandler {
	return &AuthHandler{
		userAuthService: userAuthService,
	}
}

func (a AuthHandler) LogInPostHandler(c echo.Context) error {
	var logInData models.UserLogIn
	if err := c.Bind(&logInData); err != nil {
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	token, err := a.userAuthService.LogIn(c.Request().Context(), logInData)
	if err != nil {
		if errors.Is(err, shared.ErrUserWrongEmailOrPassword) {
			return echo.NewHTTPError(http.StatusUnauthorized, "wrong email or password")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "something went wrong please try again")
	}

	c.Response().Header().Set("Authorization", "Bearer "+token)
	return c.String(http.StatusOK, "Logged in successfully")
}

func (a AuthHandler) SignInPostHandler(c echo.Context) error {
	var signInData models.UserSignIn
	if err := c.Bind(&signInData); err != nil {
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	err := a.userAuthService.Create(c.Request().Context(), signInData)
	if err != nil {
		c.Logger().Error(err.Error())
		if errors.Is(err, shared.ErrUserWithEmailExist) {
			return echo.NewHTTPError(http.StatusConflict, "User with such email already exist")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Could not create user")
	}

	return c.String(http.StatusCreated, "user created")
}
