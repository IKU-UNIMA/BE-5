package handler

import (
	"be-5/src/api/request"
	"be-5/src/api/response"
	"be-5/src/config/database"
	"be-5/src/model"
	"be-5/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetProfilHandler(c echo.Context) error {
	db := database.DB
	ctx := c.Request().Context()
	data := &response.Profil{}
	claims := util.GetClaimsFromContext(c)
	id := int(claims["id"].(float64))
	role := claims["role"].(string)

	if err := db.WithContext(ctx).Table(role).First(data, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusUnauthorized, map[string]string{"message": "user tidak ditemukan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func EditProfilHandler(c echo.Context) error {
	request := &request.Profil{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()
	data := &model.Akun{}
	claims := util.GetClaimsFromContext(c)
	id := int(claims["id"].(float64))

	if err := db.WithContext(ctx).First(data, id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusUnauthorized, map[string]string{"message": "user tidak ditemukan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Table(data.Role).Where("id", id).Update("nama", request.Nama).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
