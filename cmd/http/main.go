package main

import (
	"backend-takehome/config"
	"fmt"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	db := config.ConnectDB()
	defer db.Close()

	e := echo.New()
	e.Use(middleware.Logger(), middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": fmt.Sprintf("server running on port %s", os.Getenv("PORT")),
		})
	})

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
