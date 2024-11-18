package main

import (
	"database/sql"
	"eatwo/db"
	"eatwo/handlers"
	"eatwo/services"
	"log"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	gLog "github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
)

const dbFileName = "sqlite.db"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sqlDB, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer sqlDB.Close()

	repository := db.New(sqlDB)

	s := handlers.Services{
		AIService:        services.NewMockAIService(),
		CheckListService: services.NewCheckListService(repository),
		DreamService:     services.NewDreamService(repository),
		TokenGenerator:   services.GenerateToken,
		UserAuthService:  services.NewAuthService(repository),
    SettingsService:  services.NewSettingsService(repository, sqlDB),
	}

	e := echo.New()
	defer e.Close()
	// e.Use(middleware.Logger())
	e.Logger.SetLevel(gLog.DEBUG)
	e.Use(services.JWTMiddleware)

	e.Static("/", "assets")
	handlers.SetRoutes(e, s)

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
