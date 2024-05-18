package handlers

import (
	"eatwo/model"

	"github.com/labstack/echo/v4"
)

type UserGetter interface {
	GetByEmail(email string) (model.User, error)
}

type Auth struct {
	userGetter UserGetter
}

func (a Auth) GetUser(c echo.Context) {
	user, err := a.userGetter.GetByEmail("asdf")
	println(user, err)
}
