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
	Create(ctx context.Context, signUpData models.UserSignUp) (models.User, error)
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

func (a *AuthHandler) DeleteLogout(c echo.Context) error {
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
  return redirect(c, http.StatusSeeOther, "/")
}

func (a *AuthHandler) PostLogIn(c echo.Context) error {
	var logInData models.UserLogIn
	if err := c.Bind(&logInData); err != nil {
		c.Logger().Error(err.Error())
		return renderError(c, http.StatusBadRequest, "Invalid input")
	}

	user, err := a.userAuthService.Validate(c.Request().Context(), logInData)
	if err != nil {
		if errors.Is(err, shared.ErrUserWrongEmailOrPassword) {
			return renderError(c, http.StatusUnauthorized, "Wrong email or password")
		}
		return renderError(c, http.StatusUnauthorized, "something went wrong please try again")
	}

	err = a.setTokenCookie(c, user)
	if err != nil {
		return renderError(c, http.StatusUnauthorized, "something went wrong please try again")
	}

  return redirect(c, http.StatusSeeOther, "/")
}

func (a *AuthHandler) PostSignUp(c echo.Context) error {
	var signUpData models.UserSignUp
	if err := c.Bind(&signUpData); err != nil {
		c.Logger().Error(err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	c.Logger().Infof("signUpData: %+v", signUpData)
	user, err := a.userAuthService.Create(c.Request().Context(), signUpData)
	if err != nil {
		c.Logger().Error(err.Error())

		if errors.Is(err, shared.ErrUserWithEmailExist) {
			// TODO: Don't expose this information to the user
			return renderError(c, http.StatusUnauthorized, "User with that email already exist")
		}
		return renderError(c, http.StatusUnauthorized, "Could not create user")
	}

	err = a.setTokenCookie(c, user)
	if err != nil {
		return renderError(c, http.StatusUnauthorized, "something went wrong please try again")
	}

  return redirect(c, http.StatusSeeOther, "/")
}

func (a *AuthHandler) setTokenCookie(c echo.Context, user models.User) error {
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
