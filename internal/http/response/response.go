package response

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type APIResp struct {
	Status  string  `json:"status"`
	Code    int     `json:"code"`
	Message *string `json:"message"`
	Data    any     `json:"data"`
}

func Success(c *echo.Context, data any) error {
	return c.JSON(http.StatusOK, APIResp{
		Status: "success",
		Code:   http.StatusOK,
		Data:   data,
	})

}

func Fail(c *echo.Context, code int, message string, data any) error {
	return c.JSON(code, APIResp{
		Status:  "fail",
		Code:    code,
		Message: &message,
		Data:    data,
	})
}
