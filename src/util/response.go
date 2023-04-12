package util

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Base struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
	Data    interface{}       `json:"data"`
}

func SuccessResponse(c echo.Context, httpCode int, data interface{}) error {
	return c.JSON(
		httpCode,
		Base{
			Status:  httpCode,
			Message: http.StatusText(httpCode),
			Data:    data,
		},
	)
}

func FailedResponse(c echo.Context, httpCode int, errors map[string]string) error {
	return c.JSON(
		httpCode,
		Base{
			Status:  httpCode,
			Message: http.StatusText(httpCode),
			Errors:  errors,
		},
	)
}
