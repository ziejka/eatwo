package handlers

import (
	"eatwo/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TokenGenerator func(user models.User) (string, error)
type TokenParser func(tokenString string) (*jwt.RegisteredClaims, error)

func SetRoutes(e *echo.Echo, userAuthService UserAuthService, tokenGenerator TokenGenerator) {
	homeHandler := NewHome()
	authHandler := NewAuthHandler(userAuthService, tokenGenerator)

	// home routes
	e.GET("/", homeHandler.GetHome)
	e.GET("/about", homeHandler.GetProtectedAbout)
	e.GET("/signin", homeHandler.GetSignIn)
	e.GET("/login", homeHandler.GetLogIn)

	// api
	e.POST("/api/v1/signin", authHandler.SignInPostHandler)
	e.POST("/api/v1/login", authHandler.LogInPostHandler)
	e.POST("/api/v1/logout", authHandler.Logout)
}
