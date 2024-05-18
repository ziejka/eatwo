package main

import (
	"context"
	"database/sql"
	"eatwo/db"
	"eatwo/handlers"
	"eatwo/services"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	gLog "github.com/labstack/gommon/log"
	_ "github.com/mattn/go-sqlite3"
)

const dbFileName = "sqlite.db"

func main() {
	sqlDB, err := sql.Open("sqlite3", dbFileName)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer sqlDB.Close()

	userRepository := db.NewUserRepository(sqlDB)
	err = userRepository.Migrate(context.Background())
	if err != nil {
		log.Fatal(err)
		return
	}
	userAuthService := services.NewAuthService(userRepository)

	e := echo.New()
	defer e.Close()

	e.Use(middleware.Logger())
	e.Logger.SetLevel(gLog.DEBUG)

	e.Static("/", "assets")
	handlers.SetRoutes(e, userAuthService)

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
