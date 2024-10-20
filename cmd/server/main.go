package main

import (
	"context"
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

	userRepository := db.NewUserRepository(sqlDB)
	err = userRepository.Migrate(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	userAuthService := services.NewAuthService(repository)

	checkListRepository := db.NewCheckListRepository(sqlDB)
	err = checkListRepository.Migrate(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	checkListService := services.NewCheckListService(checkListRepository)
	e := echo.New()
	defer e.Close()

	aiService := services.NewMockAIService()

	dreamRepository := db.NewDreamRepository(sqlDB)
	err = dreamRepository.Migrate(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	dreamService := services.NewDreamService(dreamRepository)

	s := handlers.Services{
		UserAuthService:  userAuthService,
		TokenGenerator:   services.GenerateToken,
		CheckListService: checkListService,
		AIService:        aiService,
		DreamService:     dreamService,
	}

	// e.Use(middleware.Logger())
	e.Logger.SetLevel(gLog.DEBUG)
	e.Use(services.JWTMiddleware)

	e.Static("/", "assets")
	handlers.SetRoutes(e, s)

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
