package util

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type base struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Errors  []string    `json:"errors"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(c echo.Context, httpCode int, data interface{}) error {
	return c.JSON(
		httpCode,
		base{
			Status:  httpCode,
			Message: http.StatusText(httpCode),
			Data:    data,
		},
	)
}

func FailedResponse(c echo.Context, httpCode int, errors []string) error {
	return c.JSON(
		httpCode,
		base{
			Status:  httpCode,
			Message: http.StatusText(httpCode),
			Errors:  errors,
		},
	)
}
