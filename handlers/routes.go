package handlers

import (
	"eatwo/models"

	"github.com/labstack/echo/v4"
)

type TokenGenerator func(user models.User) (string, error)
type Services struct {
	UserAuthService  UserAuthService
	TokenGenerator   TokenGenerator
	CheckListService CheckListService
	AIService        AIService
	DreamService     DreamService
}

func SetRoutes(e *echo.Echo, services Services) {
	homeHandler := NewHome()
	authHandler := NewAuthHandler(services.UserAuthService, services.TokenGenerator)
	checkListHandler := NewCheckListHandler(services.CheckListService)
	dreamHander := NewDreamHandler(services.AIService, services.DreamService)

	// home routes
	e.GET("/", homeHandler.GetHome)
	e.GET("/about", homeHandler.GetProtectedAbout)
	e.GET("/sing-up", homeHandler.GetSignUp)
	e.GET("/login", homeHandler.GetLogIn)

	// dream routes
	e.POST("/api/v1/dream", dreamHander.PostDream)

	// checkList
	e.GET("/check-list", checkListHandler.GetCheckLists)
	e.GET("/check-list/:id", checkListHandler.GetCheckList)
	e.POST("/api/v1/check-list", checkListHandler.PostCheckList)
	e.POST("/api/v1/check-list/:id", checkListHandler.PostItem)

	// Auth
	e.POST("/api/v1/sing-up", authHandler.PostSignUp)
	e.POST("/api/v1/login", authHandler.PostLogIn)
	e.DELETE("/api/v1/logout", authHandler.DeleteLogout)
}
