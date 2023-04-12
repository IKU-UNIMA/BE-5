package handler

import (
	"be-5/src/api/request"
	"be-5/src/api/response"
	"be-5/src/config/database"
	"be-5/src/model"
	"be-5/src/util"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetAllProdiHandler(c echo.Context) error {
	idFakultas := c.QueryParam("fakultas")
	if !util.IsInteger(idFakultas) {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": "id fakultas harus berupa angka"})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	result := []response.Prodi{}
	condition := ""

	if idFakultas != "" && idFakultas != "0" {
		condition = "id_fakultas = " + idFakultas
	}

	if err := db.WithContext(ctx).Preload("Fakultas").Where(condition).Order("id_fakultas").Find(&result).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}

func GetProdiByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	result := &response.Prodi{}

	if err := db.WithContext(ctx).Preload("Fakultas").First(result, id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, nil)
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}

func InsertProdiHandler(c echo.Context) error {
	request := &request.Prodi{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	data := request.MapRequest()

	if err := db.WithContext(ctx).Create(data).Error; err != nil {
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": "fakultas sudah ada"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusCreated, data.ID)
}

func EditProdiHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err})
	}

	request := &request.Prodi{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).First(new(model.Prodi), id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, nil)
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Where("id", id).Updates(request.MapRequest()).Error; err != nil {
		if err != nil {
			if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
				return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": "prodi sudah ada"})
			}

			return util.FailedResponse(c, http.StatusInternalServerError, nil)
		}
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeleteProdiHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	query := db.WithContext(ctx).Delete(new(model.Prodi), id)
	if query.Error == nil && query.RowsAffected < 1 {
		return util.FailedResponse(c, http.StatusNotFound, nil)
	}

	if query.Error != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
