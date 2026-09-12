package http

import (
	nethttp "net/http"

	"maka-go/internal/health"
	httpresponse "maka-go/internal/http/response"

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
	if err := h.service.Check(c.Request().Context()); err != nil {
		return httpresponse.Fail(c, nethttp.StatusServiceUnavailable, "database unreachable", nil)
	}

	return httpresponse.Success(c, map[string]string{
		"status": "ok",
	})
}
