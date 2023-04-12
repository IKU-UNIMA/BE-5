package handler

import (
	"be-5/src/api/handler/helper"
	"be-5/src/api/request"
	"be-5/src/api/response"
	"be-5/src/config/database"
	"be-5/src/config/storage"
	"be-5/src/model"
	"be-5/src/util"
	"be-5/src/util/validation"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type pengabdianQueryParam struct {
	Tahun int    `query:"tahun"`
	Judul string `query:"Judul"`
	Page  int    `query:"page"`
}

func GetAllPengabdianHandler(c echo.Context) error {
	queryParams := &pengabdianQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	claims := util.GetClaimsFromContext(c)
	role := claims["role"].(string)
	idDosen := int(claims["id"].(float64))

	condition := ""
	if role == string(util.DOSEN) {
		condition = fmt.Sprintf("id_dosen = %d", idDosen)
	} else {
		if queryParams.Tahun != 0 {
			condition = fmt.Sprintf("tahun_kegiatan = %d", queryParams.Tahun)
		}
		if queryParams.Judul != "" {
			if condition != "" {
				condition += " AND UPPER(judul) LIKE '%" + strings.ToUpper(queryParams.Judul) + "%'"
			} else {
				condition = "UPPER(judul) LIKE '%" + strings.ToUpper(queryParams.Judul) + "%'"
			}
		}
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []model.Pengabdian{}

	if err := db.WithContext(ctx).Select("id", "tahun_pelaksanaan", "lama_kegiatan").
		Offset(util.CountOffset(queryParams.Page)).Limit(20).
		Where(condition).Find(&data).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, util.Pagination{
		Page: queryParams.Page,
		Data: model.MapBatchPengabdianResponse(data),
	})
}

func GetPengabdianByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.DetailPengabdian{}

	if !pengabdianAuthorization(c, id, db, ctx) {
		return util.FailedResponse(c, http.StatusUnauthorized, nil)
	}

	if err := db.WithContext(ctx).Table("pengabdian").
		Preload("Dokumen").Preload("Dokumen.JenisDokumen").
		Preload("AnggotaDosen", "jenis_anggota='dosen'").
		Preload("AnggotaMahasiswa", "jenis_anggota='mahasiswa'").
		Preload("AnggotaEksternal", "jenis_anggota='eksternal'").First(&data, id).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	data.TglSkPenugasan = strings.Split(data.TglSkPenugasan, "T")[0]

	return util.SuccessResponse(c, http.StatusOK, data)
}

func InsertPengabdianHandler(c echo.Context) error {
	req := &request.Pengabdian{}
	reqData := c.FormValue("data")
	if err := json.Unmarshal([]byte(reqData), req); err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if len(req.Anggota) < 1 {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": "anggota tidak boleh kosong"})
	}

	claims := util.GetClaimsFromContext(c)
	idDosen := int(claims["id"].(float64))
	pengabdian, err := req.MapRequest()
	if err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	pengabdian.IdDosen = idDosen

	db := database.InitMySQL()
	tx := db.Begin()
	ctx := c.Request().Context()

	// insert pengabdian
	if err := tx.WithContext(ctx).Create(pengabdian).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "id_dosen") {
			return util.FailedResponse(c, http.StatusNotFound, map[string]string{"message": "dosen tidak ditemukan"})
		}

		if strings.Contains(err.Error(), "id_kategori") {
			return util.FailedResponse(c, http.StatusNotFound, map[string]string{"message": "kategori tidak ditemukan"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	// insert dokumen
	idDokumen, err := helper.InsertDokumen(helper.InsertDokumenParam{
		C:       c,
		Ctx:     ctx,
		DB:      db,
		TX:      tx,
		Fitur:   "pengabdian",
		IdFitur: pengabdian.ID,
		Data:    req.Dokumen,
	})

	if err != nil {
		return err
	}

	// insert anggota
	anggota := []model.AnggotaPengabdian{}
	for _, v := range req.Anggota {
		if err := validation.ValidateAnggota(&v); err != nil {
			tx.Rollback()
			helper.DeleteBatchDokumen(idDokumen)
			return err
		}

		anggota = append(anggota, *v.MapRequest(pengabdian.ID))
	}

	if err := tx.WithContext(ctx).Create(&anggota).Error; err != nil {
		tx.Rollback()
		helper.DeleteBatchDokumen(idDokumen)
		if strings.Contains(err.Error(), "jenis_anggota") {
			return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": "jenis anggota tidak valid"})
		}

		if strings.Contains(err.Error(), "peran") {
			return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": "peran tidak valid"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		helper.DeleteBatchDokumen(idDokumen)
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.SuccessResponse(c, http.StatusCreated, nil)
}

func EditPengabdianHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err})
	}

	db := database.InitMySQL()
	tx := db.Begin()
	ctx := c.Request().Context()

	if !pengabdianAuthorization(c, id, db, ctx) {
		return util.FailedResponse(c, http.StatusUnauthorized, nil)
	}

	req := &request.Pengabdian{}
	reqData := c.FormValue("data")
	if err := json.Unmarshal([]byte(reqData), req); err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if len(req.Anggota) < 1 {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": "penulis tidak boleh kosong"})
	}

	pengabdian, errMapping := req.MapRequest()
	if errMapping != nil {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": errMapping.Error()})
	}

	// edit pengabdian
	if err := tx.WithContext(ctx).Omit("id_dosen").Where("id", id).Updates(pengabdian).Error; err != nil {
		tx.Rollback()
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	// insert dokumen
	idDokumen, errDokumen := helper.InsertDokumen(helper.InsertDokumenParam{
		C:       c,
		Ctx:     ctx,
		DB:      db,
		TX:      tx,
		Fitur:   "pengabdian",
		IdFitur: id,
		Data:    req.Dokumen,
	})

	if errDokumen != nil {
		return errDokumen
	}

	anggota := []model.AnggotaPengabdian{}
	for _, v := range req.Anggota {
		if err := validation.ValidateAnggota(&v); err != nil {
			tx.Rollback()
			helper.DeleteBatchDokumen(idDokumen)
			return err
		}

		anggota = append(anggota, *v.MapRequest(id))
	}

	if err := tx.WithContext(ctx).Delete(new(model.AnggotaPengabdian), "id_pengabdian", id).Error; err != nil {
		tx.Rollback()
		helper.DeleteBatchDokumen(idDokumen)
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	// insert anggota
	if err := tx.WithContext(ctx).Create(&anggota).Error; err != nil {
		tx.Rollback()
		helper.DeleteBatchDokumen(idDokumen)
		if strings.Contains(err.Error(), "jenis_anggota") {
			return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": "jenis anggota tidak valid"})
		}

		if strings.Contains(err.Error(), "peran") {
			return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": "peran tidak valid"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		helper.DeleteBatchDokumen(idDokumen)
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeletePengabdianHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if !pengabdianAuthorization(c, id, db, ctx) {
		return util.FailedResponse(c, http.StatusUnauthorized, nil)
	}

	idDokumen := []string{}
	if err := db.WithContext(ctx).Model(&model.DokumenPengabdian{}).Select("id").Where("id_pengabdian", id).Find(&idDokumen).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	query := db.WithContext(ctx).Delete(new(model.Pengabdian), id)
	if query.Error != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if query.RowsAffected < 1 {
		return util.FailedResponse(c, http.StatusNotFound, nil)
	}

	helper.DeleteBatchDokumen(idDokumen)

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func GetAllKategoriPengabdianHandler(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.JenisKategoriPengabdian{}

	if err := db.WithContext(ctx).Preload("KategoriPengabdian").Find(&data).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDokumenPengabdianByIdHandler(c echo.Context) error {
	id := c.Param("id")
	db := database.InitMySQL()
	ctx := c.Request().Context()

	idPengabdian := 0
	if err := db.WithContext(ctx).Model(new(model.DokumenPengabdian)).Select("id_pengabdian").First(&idPengabdian, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, nil)
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if !pengabdianAuthorization(c, idPengabdian, db, ctx) {
		return util.FailedResponse(c, http.StatusUnauthorized, nil)
	}

	data := &response.DokumenPengabdian{}

	if err := db.WithContext(ctx).Preload("JenisDokumen").First(data, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, nil)
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func EditDokumenPengabdianHandler(c echo.Context) error {
	id := c.Param("id")

	db := database.InitMySQL()
	ctx := c.Request().Context()

	idPengabdian := 0
	if err := db.WithContext(ctx).Model(new(model.DokumenPengabdian)).
		Select("id_pengabdian").First(&idPengabdian, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, nil)
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if !pengabdianAuthorization(c, idPengabdian, db, ctx) {
		return util.FailedResponse(c, http.StatusUnauthorized, nil)
	}

	return helper.EditDokumen(helper.EditDokumenParam{
		C:     c,
		Ctx:   ctx,
		DB:    db,
		Fitur: "pengabdian",
		Id:    id,
	})
}

func DeleteDokumenPengabdianHandler(c echo.Context) error {
	id := c.Param("id")
	req := &request.Dokumen{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	idPengabdian := 0
	if err := db.WithContext(ctx).Model(&model.DokumenPengabdian{}).Select("id_pengabdian").First(&idPengabdian, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, nil)
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if !pengabdianAuthorization(c, idPengabdian, db, ctx) {
		return util.FailedResponse(c, http.StatusUnauthorized, nil)
	}

	tx := db.Begin()

	query := tx.WithContext(ctx).Delete(new(model.DokumenPengabdian), "id", id)
	if query.Error != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if query.RowsAffected < 1 {
		return util.FailedResponse(c, http.StatusNotFound, nil)
	}

	if err := storage.DeleteFile(id); err != nil {
		tx.Rollback()
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func pengabdianAuthorization(c echo.Context, id int, db *gorm.DB, ctx context.Context) bool {
	claims := util.GetClaimsFromContext(c)
	role := claims["role"].(string)
	idDosen := int(claims["id"].(float64))
	if role == string(util.ADMIN) {
		return true
	}

	result := 0
	if err := db.WithContext(ctx).Table("pengabdian").Select("id_dosen").
		Where("id", id).Scan(&result).Error; err != nil {
		return false
	}

	return result == idDosen
}
