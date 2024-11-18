package handlers

import (
	"eatwo/models"

	"github.com/labstack/echo/v4"
)

type TokenGenerator func(user models.User) (string, error)
type Services struct {
	AIService        AIService
	CheckListService CheckListService
	DreamService     DreamUpdater
	TokenGenerator   TokenGenerator
	UserAuthService  UserAuthService
  SettingsService  SettingsService
}

func SetRoutes(e *echo.Echo, services Services) {
	authHandler := NewAuthHandler(services.UserAuthService, services.TokenGenerator)
	checkListHandler := NewCheckListHandler(services.CheckListService)
	dreamHander := NewDreamHandler(services.AIService, services.DreamService)
	homeHandler := NewHome(services.DreamService)
	settings := NewSettingsHandler(services.SettingsService)

	// home routes
	e.GET("/", homeHandler.GetHome)
	e.GET("/about", homeHandler.GetProtectedAbout)
	e.GET("/sing-up", homeHandler.GetSignUp)
	e.GET("/login", homeHandler.GetLogIn)

	// account settings
	e.GET("/account-settings", settings.GetAccountSettings)
	e.POST("/api/v1/account-settings/user", settings.PostUserUpdate)
  e.DELETE("/api/v1/account-settings/user", settings.DeleteUser)

	// dream routes
	e.POST("/api/v1/dream", dreamHander.PostDream)

	// checkList
	e.GET("/check-list", checkListHandler.GetCheckLists)
	e.GET("/check-list/:id", checkListHandler.GetCheckList)
	e.POST("/api/v1/check-list", checkListHandler.PostCheckList)
	e.POST("/api/v1/check-list/:id", checkListHandler.PostItem)

	// Auth
	e.POST("/api/v1/sing-up", authHandler.PostSignUp)
	e.POST("/api/v1/user", authHandler.PostLogIn)
	e.DELETE("/api/v1/logout", authHandler.DeleteLogout)
}
