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

type dosenQueryParam struct {
	Fakultas int    `query:"fakultas"`
	Prodi    int    `query:"prodi"`
	Nidn     string `query:"nidn"`
	Nama     string `query:"nama"`
	Page     int    `query:"page"`
}

func GetAllDosenHandler(c echo.Context) error {
	queryParams := &dosenQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	result := []response.Dosen{}
	limit := 20
	condition := ""

	if queryParams.Fakultas != 0 && queryParams.Prodi == 0 {
		condition = fmt.Sprintf("id_fakultas = %d", queryParams.Fakultas)
	}

	if queryParams.Prodi != 0 {
		condition = fmt.Sprintf("id_prodi = %d", queryParams.Prodi)
	}

	if queryParams.Nidn != "" {
		queryParams.Nama = ""
		if queryParams.Fakultas != 0 || queryParams.Prodi != 0 {
			condition += " AND nidn = " + queryParams.Nidn
		} else {
			condition = "nidn = " + queryParams.Nidn
		}
	}

	if queryParams.Nama != "" {
		if queryParams.Fakultas != 0 || queryParams.Prodi != 0 {
			condition += " AND UPPER(nama) LIKE '%" + strings.ToUpper(queryParams.Nama) + "%'"
		} else {
			condition = "UPPER(nama) LIKE '%" + strings.ToUpper(queryParams.Nama) + "%'"
		}
	}

	if err := db.WithContext(ctx).Preload("Fakultas").Preload("Prodi").Where(condition).
		Offset(util.CountOffset(queryParams.Page, limit)).Limit(limit).
		Find(&result).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	var totalResult int64
	if err := db.WithContext(ctx).Table("dosen").Count(&totalResult).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, util.Pagination{
		Limit:       limit,
		Page:        queryParams.Page,
		TotalPage:   util.CountTotalPage(int(totalResult), limit),
		TotalResult: int(totalResult),
		Data:        result,
	})
}

func GetDosenByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	result := &response.DetailDosen{}

	email := ""
	if err := db.WithContext(ctx).Table("akun").Select("email").Where("id", id).Scan(&email).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	result.Email = email

	if err := db.WithContext(ctx).
		Preload("Fakultas").Preload("Prodi").
		Table("dosen").First(result, id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}

func InsertDosenHandler(c echo.Context) error {
	request := &request.Dosen{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.InitMySQL()
	tx := db.Begin()
	ctx := c.Request().Context()
	akun := &model.Akun{}
	akun.Email = request.Email
	akun.Role = string(util.DOSEN)
	password := util.GeneratePassword()
	akun.Password = util.HashPassword(password)

	idFakultas := 0
	prodiQuery := db.WithContext(ctx).
		Table("prodi").Select("id_fakultas").Where("id", request.IdProdi).Scan(&idFakultas)
	if prodiQuery.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if prodiQuery.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "prodi tidak ditemukan"})
	}

	if err := tx.WithContext(ctx).Create(akun).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "email sudah digunakan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	dosen := request.MapRequest()
	dosen.ID = akun.ID
	dosen.IdFakultas = idFakultas

	if err := tx.WithContext(ctx).Create(dosen).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			if strings.Contains(err.Error(), "nidn") {
				return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "NIDN sudah digunakan"})
			}

			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "NIP sudah digunakan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.SuccessResponse(c, http.StatusCreated, map[string]string{"password": password})
}

func EditDosenHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	request := &request.Dosen{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.InitMySQL()
	tx := db.Begin()
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).First(new(model.Dosen), id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	idFakultas := 0
	prodiQuery := db.WithContext(ctx).
		Table("prodi").Select("id_fakultas").Where("id", request.IdProdi).Scan(&idFakultas)
	if prodiQuery.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if prodiQuery.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "prodi tidak ditemukan"})
	}

	result := request.MapRequest()
	result.IdFakultas = idFakultas

	if err := tx.WithContext(ctx).Table("akun").Where("id", id).Update("email", request.Email).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "email sudah digunakan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.WithContext(ctx).Where("id", id).Omit("password").Updates(result).Error; err != nil {
		if err != nil {
			tx.Rollback()
			if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
				if strings.Contains(err.Error(), "nidn") {
					return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "NIDN sudah digunakan"})
				}

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

func DeleteDosenHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.InitMySQL()
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
