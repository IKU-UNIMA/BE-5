package util

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetId(c echo.Context) (int, string) {
	id, _ := strconv.Atoi(c.Param("id"))
	if id < 1 {
		return 0, "id harus berupa angka lebih dari 1"
	}

	return id, ""
}
