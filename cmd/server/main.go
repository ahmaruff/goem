package main

import (
	"log"
	"maka-go/internal/config"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.Load()

	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hola, Amigos!"})
	})

	port := ":" + cfg.HTTPPort
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

}
