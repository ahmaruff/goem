package http

import (
	healthhttp "maka-go/internal/health/http"

	"github.com/labstack/echo/v5"
)

type Dependencies struct {
	Health *healthhttp.Handler
}

func RegisterRoutes(e *echo.Echo, deps Dependencies) {
	api := e.Group("/api/v1")

	healthhttp.RegisterRoutes(
		api.Group("/health"),
		deps.Health,
	)
}
