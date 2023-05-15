package handler

import (
	"be-5/src/api/request"
	"be-5/src/api/response"
	"be-5/src/config/database"
	"be-5/src/model"
	"be-5/src/util"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const getAdminQuery = "SELECT admin.id, nama, nip, akun.email, bagian FROM admin JOIN akun where admin.id = akun.id"

func GetAllAdminHandler(c echo.Context) error {
	db := database.DB
	ctx := c.Request().Context()
	result := []response.Admin{}

	if err := db.WithContext(ctx).Raw(getAdminQuery).Find(&result).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}

func GetAdminByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()
	result := &response.Admin{}

	condition := getAdminQuery + fmt.Sprintf(" AND admin.id = %d", id)
	if err := db.WithContext(ctx).Raw(condition).First(result).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}

func InsertAdminHandler(c echo.Context) error {
	request := &request.Admin{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.DB
	tx := db.Begin()
	ctx := c.Request().Context()
	akun := &model.Akun{}
	akun.Email = request.Email
	akun.Role = string(util.ADMIN)
	password := util.GeneratePassword()
	akun.Password = util.HashPassword(password)

	if err := tx.WithContext(ctx).Create(akun).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "email sudah digunakan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	admin := request.MapRequest()
	admin.ID = akun.ID

	if err := tx.WithContext(ctx).Create(admin).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "NIP sudah digunakan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.SuccessResponse(c, http.StatusCreated, map[string]string{"password": password})
}

func EditAdminHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	request := &request.Admin{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.DB
	tx := db.Begin()
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).First(new(model.Admin), id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.WithContext(ctx).Table("akun").Where("id", id).Update("email", request.Email).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "email sudah digunakan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	result := request.MapRequest()
	if err := tx.WithContext(ctx).Where("id", id).Omit("password").Updates(result).Error; err != nil {
		if err != nil {
			tx.Rollback()
			if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
				return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "NIP sudah digunakan"})
			}

			return util.FailedResponse(http.StatusInternalServerError, nil)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeleteAdminHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()

	query := db.WithContext(ctx).Delete(new(model.Akun), id)
	if query.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if query.Error == nil && query.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
