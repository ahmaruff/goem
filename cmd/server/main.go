package main

import (
	"log"
	"maka-go/internal/config"
	"maka-go/internal/db"
	"maka-go/internal/health"
	healthhttp "maka-go/internal/health/http"
	apphttp "maka-go/internal/http"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func main() {
	// load config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.Load()

	// setup DB
	database, err := db.New(db.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Database: cfg.DBDatabase,
		Username: cfg.DBUsername,
		Password: cfg.DBPassword,
	})

	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	defer database.Close()

	// setup http handler
	e := echo.New()

	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	// setup dependencies
	healthService := health.NewService(database)
	healthHandler := healthhttp.NewHandler(healthService)

	// register routes
	apphttp.RegisterRoutes(e, apphttp.Dependencies{
		Health: healthHandler,
	})

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hola, Amigos!"})
	})

	port := ":" + cfg.HTTPPort
	if err := e.Start(port); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}

}
