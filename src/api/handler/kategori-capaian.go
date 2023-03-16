package handler

import (
	"be-5/src/api/response"
	"be-5/src/config/database"
	"be-5/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllKategoriCapaianHandler(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	result := []response.KategoriCapaian{}

	if err := db.WithContext(ctx).Find(&result).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}
