package handlers

import "github.com/labstack/echo/v4"

func SetRoutes(e *echo.Echo, userAuthService UserAuthService) {
	homeHandler := Home{}
	authHandler := NewAuthHandler(userAuthService)

	// home routes
	e.GET("/", homeHandler.GetHome)

	// auth routes
	e.POST("/signIn", authHandler.SignInPostHandler)
}
