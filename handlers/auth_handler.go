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
	Create(ctx context.Context, signInData models.UserSignIn) (models.User, error)
	Validate(ctx context.Context, logInData models.UserLogIn) (models.User, error)
}

type AuthHandler struct {
	userAuthService UserAuthService
	generateToken   TokenGenerator
}

func NewAuthHandler(userAuthService UserAuthService, tokenGenerator TokenGenerator) *AuthHandler {
	return &AuthHandler{
		userAuthService,
		tokenGenerator,
	}
}

func (a AuthHandler) LogInPostHandler(c echo.Context) error {
	var logInData models.UserLogIn
	if err := c.Bind(&logInData); err != nil {
		c.Logger().Error(err.Error())
		return a.error(c, http.StatusBadRequest, "Invalid input")
	}

	user, err := a.userAuthService.Validate(c.Request().Context(), logInData)
	if err != nil {
		if errors.Is(err, shared.ErrUserWrongEmailOrPassword) {
			return a.error(c, http.StatusUnauthorized, "Wrong email or password")
		}
		return a.error(c, http.StatusUnauthorized, "something went wrong please try again")
	}

	token, err := a.generateToken(user)
	if err != nil {
		return a.error(c, http.StatusUnauthorized, "something went wrong please try again")
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	c.SetCookie(cookie)
	return render(c, http.StatusOK, pages.HomePage(user.Email))
}

func (a AuthHandler) SignInPostHandler(c echo.Context) error {
	var signInData models.UserSignIn
	if err := c.Bind(&signInData); err != nil {
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	user, err := a.userAuthService.Create(c.Request().Context(), signInData)
	if err != nil {
		c.Logger().Error(err.Error())
		c.Response().Header().Set("HX-Swap", "outerHTML")
		c.Response().Header().Set("HX-Retarget", "#auth-error")

		if errors.Is(err, shared.ErrUserWithEmailExist) {
			return render(c, http.StatusUnauthorized, pages.AuthError("User with that email already exist"))
		}
		return render(c, http.StatusUnauthorized, pages.AuthError("Could not create user"))
	}

	// TODO Redirect to homepage
	token, err := a.generateToken(user)
	if err != nil {
		return a.error(c, http.StatusUnauthorized, "something went wrong please try again")
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	c.SetCookie(cookie)
	return render(c, http.StatusOK, pages.HomePage(user.Email))
}

func (a AuthHandler) error(c echo.Context, statusCode int, message string) error {
	c.Response().Header().Set("HX-Retarget", "#auth-error")
	c.Response().Header().Set("HX-Reswap", "outerHTML")
	return render(c, statusCode, pages.AuthError(message))
}
