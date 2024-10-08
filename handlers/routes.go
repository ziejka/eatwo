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
	e.GET("/sing-up", homeHandler.GetSignUp)
	e.GET("/login", homeHandler.GetLogIn)

	// checkList
	e.GET("/check-list", checkListHandler.GetCheckListsHandler)
	e.GET("/check-list/:id", checkListHandler.GetCheckListHandler)
	e.POST("/api/v1/check-list", checkListHandler.PostCheckListHandler)
	e.POST("/api/v1/check-list/:id", checkListHandler.PostItemHandler)

	// Auth
	e.POST("/api/v1/sing-up", authHandler.SignUpPostHandler)
	e.POST("/api/v1/login", authHandler.LogInPostHandler)
	e.POST("/api/v1/logout", authHandler.Logout)
}
