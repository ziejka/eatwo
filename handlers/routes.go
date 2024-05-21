package handlers

import "github.com/labstack/echo/v4"

func SetRoutes(e *echo.Echo, userAuthService UserAuthService) {
	homeHandler := Home{}
	authHandler := NewAuthHandler(userAuthService)

	// home routes
	e.GET("/", homeHandler.GetHome)
	e.GET("/signin", homeHandler.GetSignIn)
	e.GET("/login", homeHandler.GetLogIn)

	// api
	e.POST("/api/v1/signin", authHandler.SignInPostHandler)
	e.POST("/api/v1/login", authHandler.LogInPostHandler)
}
