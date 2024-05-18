package main

import (
	"database/sql"
	"eatwo/handlers"
	"log"

	"github.com/labstack/echo/v4"
)

const dbFileName = "sqlite.db"

func main() {
	db, err := sql.Open("sqlite", dbFileName)
	if err != nil {
		log.Fatal(err)
		return
	}


	e := echo.New()
	e.Static("/", "assets")
	handlers.SetRoutes(e)

	e.Logger.Fatal(e.Start("127.0.0.1:8080"))
}
