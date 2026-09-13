package main

import (
	"database/sql"
	"fmt"
	"maka-go/internal/config"
	"maka-go/internal/db"
	"maka-go/internal/health"
	healthhttp "maka-go/internal/health/http"
	apphttp "maka-go/internal/http"
	appmiddleware "maka-go/internal/http/middleware"
	"maka-go/internal/logger"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v5"
	echomiddleware "github.com/labstack/echo/v5/middleware"
)

// version is stamped into every log line. Override at build time with
// -ldflags "-X main.version=...".
var version = "dev"

func main() {
	if err := run(); err != nil {
		logger.L().Fatal().Err(err).Msg("server stopped")
	}
}

// run is the composition root: build every dependency, wire it, then serve.
// Doing it here (instead of inside main) keeps every deferred Close running on
// the way out, no matter which step fails.
func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	if logErr := setupLogger(cfg); logErr != nil {
		logger.L().Warn().Err(logErr).Msg("file logging disabled, stdout only")
	}
	defer logger.Close()

	database, err := openDatabase(cfg)
	if err != nil {
		return err
	}
	defer database.Close()

	e := newServer(buildDependencies(database))

	if err := e.Start(":" + cfg.HTTPPort); err != nil {
		return fmt.Errorf("start http server: %w", err)
	}

	return nil
}

// loadConfig reads .env into the process environment, then into a Config.
func loadConfig() (config.Config, error) {
	if err := godotenv.Load(); err != nil {
		return config.Config{}, fmt.Errorf("load .env file: %w", err)
	}

	return config.Load(), nil
}

// setupLogger installs the process-wide logger. A non-nil error means
// file logging is off; the logger itself still works on stdout.
func setupLogger(cfg config.Config) error {
	return logger.Init(logger.Config{
		Env:     cfg.AppEnv,
		Level:   cfg.LogLevel,
		Service: cfg.AppName,
		Version: version,
		File:    cfg.LogFile,
	})
}

// openDatabase connects and verifies the connection, capturing a stack so the
// fatal log points at the origin of the failure.
func openDatabase(cfg config.Config) (*sql.DB, error) {
	database, err := db.New(db.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Database: cfg.DBDatabase,
		Username: cfg.DBUsername,
		Password: cfg.DBPassword,
	})
	if err != nil {
		return nil, logger.WithStack(fmt.Errorf("connect database: %w", err))
	}

	return database, nil
}

// buildDependencies assembles every module's service + handler. Services are
// built here, never inside the http layer.
func buildDependencies(database *sql.DB) apphttp.Dependencies {
	return apphttp.Dependencies{
		Health: healthhttp.NewHandler(health.NewService(database)),
	}
}

// newServer builds the Echo instance: shared middleware first, then routes.
func newServer(deps apphttp.Dependencies) *echo.Echo {
	e := echo.New()
	// Send Echo's own startup/server logs through zerolog, one JSON shape.
	e.Logger = logger.Slog()

	e.Use(appmiddleware.Logger(appmiddleware.LoggerConfig{
		QuietPaths: []string{"/", "/api/v1/health"},
	}))
	e.Use(echomiddleware.Recover())

	// Module routes first, then app-level ones.
	apphttp.RegisterRoutes(e, deps)

	e.GET("/", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "Hola, Amigos!"})
	})

	return e
}
