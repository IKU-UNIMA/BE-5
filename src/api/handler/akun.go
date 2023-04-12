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

func LoginHandler(c echo.Context) error {
	request := &request.Login{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &model.Akun{}

	if err := db.WithContext(ctx).First(data, "email", request.Email).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusUnauthorized, map[string]string{"message": "email atau password salah"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if !util.ValidateHash(request.Password, data.Password) {
		return util.FailedResponse(c, http.StatusUnauthorized, map[string]string{"message": "email atau password salah"})
	}

	var bagian string
	if data.Role == string(util.ADMIN) {
		if err := db.WithContext(ctx).Table("admin").Select("bagian").Where("id", data.ID).Scan(&bagian).Error; err != nil {
			return util.FailedResponse(c, http.StatusInternalServerError, nil)
		}
	}

	var nama string
	if err := db.WithContext(ctx).Table(data.Role).Select("nama").Where("id", data.ID).Scan(&nama).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	token := util.GenerateJWT(data.ID, nama, data.Role, bagian)

	return util.SuccessResponse(c, http.StatusOK, response.Login{Token: token})
}

func ChangePasswordHandler(c echo.Context) error {
	request := &request.ChangePassword{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	claims := util.GetClaimsFromContext(c)
	id := int(claims["id"].(float64))

	if err := db.WithContext(ctx).First(new(model.Akun), "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, map[string]string{"message": "user tidak ditemukan"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Table("akun").Where("id", id).Update("password", util.HashPassword(request.PasswordBaru)).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func ResetPasswordHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).First(new(model.Akun), "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, map[string]string{"message": "user tidak ditemukan"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	password := util.GeneratePassword()

	if err := db.WithContext(ctx).Table("akun").Where("id", id).Update("password", util.HashPassword(password)).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, map[string]string{"password": password})
}
