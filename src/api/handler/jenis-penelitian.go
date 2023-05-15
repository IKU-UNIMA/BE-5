package handler

import (
	"be-5/src/api/response"
	"be-5/src/config/database"
	"be-5/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllJenisPenelitianHandler(c echo.Context) error {
	db := database.DB
	ctx := c.Request().Context()
	result := []response.JenisPenelitian{}

	if err := db.WithContext(ctx).Find(&result).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}
