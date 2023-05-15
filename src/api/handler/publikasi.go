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

type publikasiQueryParam struct {
	Tahun int    `query:"tahun"`
	Judul string `query:"judul"`
	Page  int    `query:"page"`
}

func GetAllPublikasiHandler(c echo.Context) error {
	queryParams := &publikasiQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	claims := util.GetClaimsFromContext(c)
	role := claims["role"].(string)
	idDosen := int(claims["id"].(float64))

	condition := ""
	if role == string(util.DOSEN) {
		condition = fmt.Sprintf("id_dosen = %d", idDosen)
	} else {
		if queryParams.Tahun != 0 {
			condition = fmt.Sprintf(`YEAR(tanggal_terbit) = %d OR YEAR(waktu_pelaksanaan) = %d`,
				queryParams.Tahun, queryParams.Tahun)
		}

		if queryParams.Judul != "" {
			if condition != "" {
				condition += " AND UPPER(judul) LIKE '%" + strings.ToUpper(queryParams.Judul) + "%'"
			} else {
				condition = "UPPER(judul) LIKE '%" + strings.ToUpper(queryParams.Judul) + "%'"
			}
		}
	}

	db := database.DB
	ctx := c.Request().Context()
	limit := 20
	data := []response.Publikasi{}

	if err := db.WithContext(ctx).Preload("Dosen").
		Preload("JenisPenelitian").Preload("Kategori").
		Offset(util.CountOffset(queryParams.Page, limit)).Limit(limit).
		Where(condition).Order("tanggal_terbit DESC, waktu_pelaksanaan DESC").Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	var totalResult int64
	if err := db.WithContext(ctx).Table("publikasi").Where(condition).Count(&totalResult).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, util.Pagination{
		Limit:       limit,
		Page:        queryParams.Page,
		TotalPage:   util.CountTotalPage(int(totalResult), limit),
		TotalResult: int(totalResult),
		Data:        data,
	})
}

func GetPublikasiByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()
	data := &response.DetailPublikasi{}

	if err := publikasiAuthorization(c, id, db, ctx); err != nil {
		return err
	}

	if err := db.WithContext(ctx).Table("publikasi").
		Preload("Dosen.Fakultas").Preload("Dosen.Prodi").
		Preload("JenisPenelitian").Preload("KategoriCapaian").
		Preload("Dokumen.JenisDokumen").
		Preload("PenulisDosen", "jenis_penulis = 'dosen'").
		Preload("PenulisMahasiswa", "jenis_penulis = 'mahasiswa'").
		Preload("PenulisLain", "jenis_penulis = 'lain'").
		Where("id", id).First(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func InsertPublikasiHandler(c echo.Context) error {
	req := &request.Publikasi{}
	reqData := c.FormValue("data")
	if err := json.Unmarshal([]byte(reqData), req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if len(req.PenulisDosen) < 1 {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "penulis dosen tidak boleh kosong"})
	}

	claims := util.GetClaimsFromContext(c)
	idDosen := int(claims["id"].(float64))

	db := database.DB
	tx := db.Begin()
	ctx := c.Request().Context()
	publikasi, err := req.MapRequest()
	if err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	fmt.Println(req)
	fmt.Println(publikasi)

	publikasi.IdDosen = idDosen

	// insert publikasi
	if err := tx.WithContext(ctx).Create(publikasi).Error; err != nil {
		tx.Rollback()
		return checkPublikasiError(c, err)
	}

	// insert dokumen
	idDokumen, err := helper.InsertDokumen(helper.InsertDokumenParam{
		C:       c,
		Ctx:     ctx,
		DB:      db,
		TX:      tx,
		Fitur:   "publikasi",
		IdFitur: publikasi.ID,
		Data:    req.Dokumen,
	})

	if err != nil {
		return err
	}

	// mapping penulis
	penulis := []model.PenulisPublikasi{}
	for _, v := range req.PenulisDosen {
		if err := validation.ValidatePenulis(&v); err != nil {
			tx.Rollback()
			helper.DeleteBatchDokumen(idDokumen)
			return err
		}

		penulis = append(penulis, *v.MapRequestToPublikasi(publikasi.ID, "dosen"))
	}

	for _, v := range req.PenulisMahasiswa {
		if len(req.PenulisMahasiswa) == 1 && req.PenulisMahasiswa[0].Nama == "" {
			break
		}

		if err := validation.ValidatePenulis(&v); err != nil {
			tx.Rollback()
			helper.DeleteBatchDokumen(idDokumen)
			return err
		}

		penulis = append(penulis, *v.MapRequestToPublikasi(publikasi.ID, "mahasiswa"))
	}

	for _, v := range req.PenulisLain {
		if len(req.PenulisLain) == 1 && req.PenulisLain[0].Nama == "" {
			break
		}

		if err := validation.ValidatePenulis(&v); err != nil {
			tx.Rollback()
			helper.DeleteBatchDokumen(idDokumen)
			return err
		}

		penulis = append(penulis, *v.MapRequestToPublikasi(publikasi.ID, "lain"))
	}

	// insert penulis
	if err := tx.WithContext(ctx).Create(&penulis).Error; err != nil {
		tx.Rollback()
		helper.DeleteBatchDokumen(idDokumen)
		if strings.Contains(err.Error(), "jenis_penulis") {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "jenis penulis tidak valid"})
		}

		if strings.Contains(err.Error(), "peran") {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "peran tidak valid"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		helper.DeleteBatchDokumen(idDokumen)
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.SuccessResponse(c, http.StatusCreated, nil)
}

func EditPublikasiHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()

	if err := publikasiAuthorization(c, id, db, ctx); err != nil {
		return err
	}

	req := &request.Publikasi{}
	reqData := c.FormValue("data")
	if err := json.Unmarshal([]byte(reqData), req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	if len(req.PenulisDosen) < 1 {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "penulis dosen tidak boleh kosong"})
	}

	tx := db.Begin()
	publikasi, errMapping := req.MapRequest()
	if errMapping != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": errMapping.Error()})
	}

	// edit publikasi
	if err := tx.WithContext(ctx).Where("id", id).Updates(publikasi).Error; err != nil {
		tx.Rollback()
		return checkPublikasiError(c, err)
	}

	// insert dokumen
	idDokumen, errDokumen := helper.InsertDokumen(helper.InsertDokumenParam{
		C:       c,
		Ctx:     ctx,
		DB:      db,
		TX:      tx,
		Fitur:   "publikasi",
		IdFitur: id,
		Data:    req.Dokumen,
	})

	if errDokumen != nil {
		return errDokumen
	}

	// mapping penulis
	penulis := []model.PenulisPublikasi{}
	for _, v := range req.PenulisDosen {
		if err := validation.ValidatePenulis(&v); err != nil {
			tx.Rollback()
			helper.DeleteBatchDokumen(idDokumen)
			return err
		}

		penulis = append(penulis, *v.MapRequestToPublikasi(id, "dosen"))
	}

	for _, v := range req.PenulisMahasiswa {
		if len(req.PenulisMahasiswa) == 1 && req.PenulisMahasiswa[0].Nama == "" {
			break
		}

		if err := validation.ValidatePenulis(&v); err != nil {
			tx.Rollback()
			helper.DeleteBatchDokumen(idDokumen)
			return err
		}

		penulis = append(penulis, *v.MapRequestToPublikasi(id, "mahasiswa"))
	}

	for _, v := range req.PenulisLain {
		if len(req.PenulisLain) == 1 && req.PenulisLain[0].Nama == "" {
			break
		}

		if err := validation.ValidatePenulis(&v); err != nil {
			tx.Rollback()
			helper.DeleteBatchDokumen(idDokumen)
			return err
		}

		penulis = append(penulis, *v.MapRequestToPublikasi(id, "lain"))
	}

	// delete old penulis
	if err := tx.WithContext(ctx).Delete(new(model.PenulisPublikasi), "id_publikasi", id).Error; err != nil {
		tx.Rollback()
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	// insert penulis
	if err := tx.WithContext(ctx).Create(&penulis).Error; err != nil {
		tx.Rollback()
		helper.DeleteBatchDokumen(idDokumen)
		if strings.Contains(err.Error(), "jenis_penulis") {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "jenis penulis tidak valid"})
		}

		if strings.Contains(err.Error(), "peran") {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "peran tidak valid"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		helper.DeleteBatchDokumen(idDokumen)
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeletePublikasiHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()

	if err := publikasiAuthorization(c, id, db, ctx); err != nil {
		return err
	}

	idDokumen := []string{}
	if err := db.WithContext(ctx).Model(&model.DokumenPublikasi{}).Select("id").Where("id_publikasi", id).Find(&idDokumen).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	query := db.WithContext(ctx).Delete(new(model.Publikasi), id)
	if query.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if query.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, nil)
	}

	helper.DeleteBatchDokumen(idDokumen)

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func GetAllKategoriPublikasiHandler(c echo.Context) error {
	db := database.DB
	ctx := c.Request().Context()
	data := []response.JenisKategoriPublikasi{}

	if err := db.WithContext(ctx).Preload("KategoriPublikasi").Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDokumenPublikasiByIdHandler(c echo.Context) error {
	id := c.Param("id")
	db := database.DB
	ctx := c.Request().Context()

	idPublikasi := 0
	if err := db.WithContext(ctx).Model(new(model.DokumenPublikasi)).
		Select("id_publikasi").First(&idPublikasi, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := publikasiAuthorization(c, idPublikasi, db, ctx); err != nil {
		return err
	}

	data := &response.DokumenPublikasi{}

	if err := db.WithContext(ctx).Preload("JenisDokumen").First(data, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func EditDokumenPublikasiHandler(c echo.Context) error {
	id := c.Param("id")

	db := database.DB
	ctx := c.Request().Context()

	idPublikasi := 0
	if err := db.WithContext(ctx).Model(new(model.DokumenPublikasi)).
		Select("id_publikasi").First(&idPublikasi, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := publikasiAuthorization(c, idPublikasi, db, ctx); err != nil {
		return err
	}

	return helper.EditDokumen(helper.EditDokumenParam{
		C:     c,
		Ctx:   ctx,
		DB:    db,
		Fitur: "publikasi",
		Id:    id,
	})
}

func DeleteDokumenPublikasiHandler(c echo.Context) error {
	id := c.Param("id")
	db := database.DB
	ctx := c.Request().Context()

	idPublikasi := 0
	if err := db.WithContext(ctx).Model(new(model.DokumenPublikasi)).
		Select("id_publikasi").First(&idPublikasi, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := publikasiAuthorization(c, idPublikasi, db, ctx); err != nil {
		return err
	}

	tx := db.Begin()

	query := tx.WithContext(ctx).Delete(new(model.DokumenPublikasi), "id", id)
	if query.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if query.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, nil)
	}

	if err := storage.DeleteFile(id); err != nil {
		tx.Rollback()
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func publikasiAuthorization(c echo.Context, id int, db *gorm.DB, ctx context.Context) error {
	claims := util.GetClaimsFromContext(c)
	role := claims["role"].(string)
	idDosen := int(claims["id"].(float64))
	if role == string(util.ADMIN) {
		return nil
	}

	result := 0
	query := db.WithContext(ctx).Table("publikasi").Select("id_dosen").
		Where("id", id).Scan(&result)
	if query.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if query.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "data publikasi tidak ditemukan"})
	}

	if result == idDosen {
		return nil
	}

	return util.FailedResponse(http.StatusUnauthorized, nil)
}

// checkPublikasiError used to check the error while inserting or updating publikasi
func checkPublikasiError(c echo.Context, err error) error {
	if strings.Contains(err.Error(), "id_dosen") {
		return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "dosen tidak ditemukan"})
	}

	if strings.Contains(err.Error(), "id_kategori_capaian") {
		return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "kategori capaian tidak ditemukan"})
	}

	if strings.Contains(err.Error(), "id_kategori") {
		return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "kategori tidak ditemukan"})
	}

	if strings.Contains(err.Error(), "id_jenis_penelitian") {
		return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "jenis penelitian tidak ditemukan"})
	}

	return util.FailedResponse(http.StatusInternalServerError, nil)
}
