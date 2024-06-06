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

func (a AuthHandler) Logout(c echo.Context) error {
	// Old token should be stored in some key-value and check if it was logged out
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		MaxAge:   -1,
	}

	c.SetCookie(cookie)
	return render(c, http.StatusOK, pages.HomePageWithNavigation(""))
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

	err = a.setTokenCookie(c, user)
	if err != nil {
		return a.error(c, http.StatusUnauthorized, "something went wrong please try again")
	}
	return render(c, http.StatusOK, pages.HomePageWithNavigation(user.Email))
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

		if errors.Is(err, shared.ErrUserWithEmailExist) {
			return a.error(c, http.StatusUnauthorized, "User with that email already exist")
		}
		return a.error(c, http.StatusUnauthorized, "Could not create user")
	}

	err = a.setTokenCookie(c, user)
	if err != nil {
		return a.error(c, http.StatusUnauthorized, "something went wrong please try again")
	}

	return render(c, http.StatusOK, pages.HomePageWithNavigation(user.Email))
}

func (a AuthHandler) setTokenCookie(c echo.Context, user models.User) error {
	token, err := a.generateToken(user)
	if err != nil {
		return err
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	}
	c.SetCookie(cookie)

	return nil
}

func (a AuthHandler) error(c echo.Context, statusCode int, message string) error {
	c.Response().Header().Set("HX-Retarget", "#auth-error")
	c.Response().Header().Set("HX-Reswap", "outerHTML")
	return render(c, statusCode, pages.AuthError(message))
}
