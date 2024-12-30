package handlers

import (
	"eatwo/models"

	"github.com/labstack/echo/v4"
)

type TokenGenerator func(user *models.User) (string, error)
type Services struct {
	AIService       AIService
	DreamService    DreamUpdater
	TokenGenerator  TokenGenerator
	UserAuthService UserAuthService
	SettingsService SettingsService
}

func SetRoutes(e *echo.Echo, services Services) {
	authHandler := NewAuthHandler(services.UserAuthService, services.TokenGenerator)
	dreamHandler := NewDreamHandler(services.AIService, services.DreamService)
	homeHandler := NewHome()
	settings := NewSettingsHandler(services.SettingsService)
	journalHandler := NewJournalHandler(services.DreamService)

	// Auth
	e.POST("/api/v1/sing-up", authHandler.PostSignUp)
	e.POST("/api/v1/login", authHandler.PostLogIn)
	e.DELETE("/api/v1/logout", authHandler.DeleteLogout)

	// home routes
	e.GET("/", homeHandler.GetHome)
	e.GET("/sing-up", homeHandler.GetSignUp)
	e.GET("/login", homeHandler.GetLogIn)

	// journal routes
	e.GET("/journal", journalHandler.GetJournal)

	// account settings
	e.GET("/account-settings", settings.GetAccountSettings)
	e.POST("/api/v1/account-settings/user", settings.PostUserUpdate)
	e.DELETE("/api/v1/account-settings/user", settings.DeleteUser)

	// dream routes
	e.POST("/api/v1/dream", dreamHandler.PostDream)

	// checkList
	// e.GET("/check-list", checkListHandler.GetCheckLists)
	// e.GET("/check-list/:id", checkListHandler.GetCheckList)
	// e.POST("/api/v1/check-list", checkListHandler.PostCheckList)
	// e.POST("/api/v1/check-list/:id", checkListHandler.PostItem)
}
