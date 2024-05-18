package handlers

import "github.com/labstack/echo/v4"

func SetRoutes(e *echo.Echo) {
	home := Home{}
	e.GET("/", home.GetHome)
}
