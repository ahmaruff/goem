package http

import (
	nethttp "net/http"

	"maka-go/internal/health"
	httpresponse "maka-go/internal/http/response"
	"maka-go/internal/logger"

	"github.com/labstack/echo/v5"
)

type Handler struct {
	service *health.Service
}

func NewHandler(service *health.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Check(c *echo.Context) error {
	ctx := c.Request().Context()

	if err := h.service.Check(ctx); err != nil {
		// Boundary: log the cause once (with stack), then answer the client.
		logger.Err(ctx, err).
			Str("dependency", "mysql").
			Msg("health check failed")

		return httpresponse.Fail(c, nethttp.StatusServiceUnavailable, "database unreachable", nil)
	}

	return httpresponse.Success(c, map[string]string{
		"status": "ok",
	})
}
