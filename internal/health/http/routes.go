package http

import "github.com/labstack/echo/v5"

// RegisterRoutes mounts every health route on the given group.
func RegisterRoutes(g *echo.Group, h *Handler) {
	g.GET("", h.Check)
}
