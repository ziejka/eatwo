package handlers

import (
	"context"
	"eatwo/models"
	"eatwo/shared"
	"eatwo/views/pages"
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
			c.Response().Header().Set("HX-Retarget", "#auth-error")
			return render(c, http.StatusUnauthorized, pages.AuthError("Wrong email or password"))
		}
		c.Response().Header().Set("HX-Retarget", "#auth-error")
		return render(c, http.StatusUnauthorized, pages.AuthError("something went wrong please try again"))
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
			c.Response().Header().Set("HX-Retarget", "#auth-error")
			return render(c, http.StatusUnauthorized, pages.AuthError("User with that email already exist"))
		}
		c.Response().Header().Set("HX-Retarget", "#auth-error")
		return render(c, http.StatusUnauthorized, pages.AuthError("Could not create user"))
	}

	// TODO Redirect to homepage
	return c.String(http.StatusCreated, "user created")
}
