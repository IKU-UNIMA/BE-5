package handler

import (
	"be-5/src/api/request"
	"be-5/src/api/response"
	"be-5/src/config/database"
	"be-5/src/model"
	"be-5/src/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context) error {
	request := &request.Login{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &model.Akun{}

	if err := db.WithContext(ctx).First(data, "email", request.Email).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusUnauthorized, []string{"email atau password salah"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if !util.ValidateHash(request.Password, data.Password) {
		return util.FailedResponse(c, http.StatusUnauthorized, []string{"email atau password salah"})
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

	token := util.GenerateJWT(data.ID, data.Role, bagian)

	return util.SuccessResponse(c, http.StatusOK, response.Login{Nama: nama, Token: token})
}

func ChangePasswordHandler(c echo.Context) error {
	request := &request.ChangePassword{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &model.Akun{}
	claims := util.GetClaimsFromContext(c)
	id := int(claims["id"].(float64))

	if err := db.WithContext(ctx).First(data, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			msg := fmt.Sprintf("user dengan id %d tidak ditemukan", id)
			return util.FailedResponse(c, http.StatusNotFound, []string{msg})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if !util.ValidateHash(request.PasswordLama, data.Password) {
		return util.FailedResponse(c, http.StatusUnauthorized, []string{"password anda berbeda dengan yang lama"})
	}

	if request.PasswordBaru == "" {
		return util.FailedResponse(c, http.StatusBadRequest, []string{"password baru tidak boleh kosong"})
	}

	if err := db.WithContext(ctx).Table("akun").Where("id", id).Update("password", util.HashPassword(request.PasswordBaru)).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
