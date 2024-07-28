package handlers

import (
	"eatwo/models"

	"github.com/labstack/echo/v4"
)

type TokenGenerator func(user models.User) (string, error)

func SetRoutes(e *echo.Echo, userAuthService UserAuthService, tokenGenerator TokenGenerator, checklistService CheckListService) {
	homeHandler := NewHome()
	authHandler := NewAuthHandler(userAuthService, tokenGenerator)
	checkListHandler := NewCheckListHandler(checklistService)

	// home routes
	e.GET("/", homeHandler.GetHome)
	e.GET("/about", homeHandler.GetProtectedAbout)
	e.GET("/signin", homeHandler.GetSignIn)
	e.GET("/login", homeHandler.GetLogIn)

	// API v1
	// Auth
	e.POST("/api/v1/signin", authHandler.SignInPostHandler)
	e.POST("/api/v1/login", authHandler.LogInPostHandler)
	e.POST("/api/v1/logout", authHandler.Logout)

	// checkList
	e.POST("/api/v1/check-list", checkListHandler.PostListHandler)
}
