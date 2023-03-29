package handler

import (
	"be-5/src/api/request"
	"be-5/src/api/response"
	"be-5/src/config/database"
	"be-5/src/config/env"
	"be-5/src/config/storage"
	"be-5/src/model"
	"be-5/src/util"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type patenQueryParam struct {
	Tahun int    `query:"tahun"`
	Nama  string `query:"nama"`
	Page  int    `query:"page"`
}

var patenFolderId = env.GetPatenFolderId()

func GetAllPatenHandler(c echo.Context) error {
	queryParams := &patenQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err.Error()})
	}

	claims := util.GetClaimsFromContext(c)
	role := claims["role"].(string)
	idDosen := int(claims["id"].(float64))

	condition := ""
	if role == string(util.ADMIN) {
		if queryParams.Tahun != 0 {
			condition = fmt.Sprintf(`YEAR("tanggal") = %d`, queryParams.Tahun)
		}

		if queryParams.Nama != "" {
			if queryParams.Tahun != 0 {
				condition += " AND UPPER(nama) LIKE '%" + strings.ToUpper(queryParams.Nama) + "%'"
			} else {
				condition = "UPPER(nama) LIKE '%" + strings.ToUpper(queryParams.Nama) + "%'"
			}
		}
	} else {
		condition = fmt.Sprintf("id_dosen = %d", idDosen)
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.Paten{}

	if err := db.WithContext(ctx).
		Preload("JenisPenelitian").Preload("Kategori").
		Offset(util.CountOffset(queryParams.Page)).Limit(20).
		Where(condition).Find(&data).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, util.Pagination{Page: queryParams.Page, Data: data})
}

func GetPatenByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.DetailPaten{}

	if !patenAuthorization(c, id, db, ctx) {
		return util.FailedResponse(c, http.StatusUnauthorized, nil)
	}

	if err := db.WithContext(ctx).Table("paten").
		Preload("JenisPenelitian").Preload("KategoriCapaian").Preload("Penulis").
		Preload("Dokumen").Preload("Dokumen.JenisDokumen").
		Where("id", id).First(data).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func InsertPatenHandler(c echo.Context) error {
	req := &request.Paten{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err.Error()})
	}

	claims := util.GetClaimsFromContext(c)
	idDosen := int(claims["id"].(float64))

	db := database.InitMySQL()
	tx := db.Begin()
	ctx := c.Request().Context()
	paten, err := req.MapRequest()
	if err != nil {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err.Error()})
	}

	paten.IdDosen = idDosen

	if err := tx.WithContext(ctx).Create(paten).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "id_dosen") {
			return util.FailedResponse(c, http.StatusNotFound, []string{"dosen tidak ditemukan"})
		}
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	dokumenPaten := []model.DokumenPaten{}
	form, _ := c.MultipartForm()
	if form != nil && paten.Dokumen != nil {
		files := form.File["files"]
		minLen := util.CountMin(len(paten.Dokumen), len(files))
		for i := 0; i < minLen; i++ {
			file := files[i]
			dFile, err := storage.CreateFile(file, patenFolderId)
			if err != nil {
				tx.Rollback()
				return c.JSON(http.StatusBadRequest, err.Error())
			}

			dokumenPaten = append(dokumenPaten, *req.Dokumen[i].MapRequest(&request.DokumenPatenPayload{
				IdFile:    dFile.Id,
				IdPaten:   paten.ID,
				NamaFile:  dFile.Name,
				JenisFile: dFile.MimeType,
				Url:       util.CreateFileUrl(dFile.Id),
			}))

			if paten.Dokumen[i].Nama == "" {
				dokumenPaten[i].Nama = dFile.Name
			}
		}

		if err := tx.WithContext(ctx).Create(&dokumenPaten).Error; err != nil {
			tx.Rollback()
			deleteBatchDokumenPaten(dokumenPaten)
			if strings.Contains(err.Error(), "jenis_dokumen") {
				return util.FailedResponse(c, http.StatusBadRequest, []string{"jenis dokumen tidak valid"})
			}

			return util.FailedResponse(c, http.StatusInternalServerError, nil)
		}
	}

	penulis := []model.PenulisPaten{}
	for _, v := range req.Penulis {
		penulis = append(penulis, *v.MapRequest(paten.ID))
	}

	if err := tx.WithContext(ctx).Create(&penulis).Error; err != nil {
		tx.Rollback()
		deleteBatchDokumenPaten(dokumenPaten)
		if strings.Contains(err.Error(), "jenis_penulis") {
			return util.FailedResponse(c, http.StatusBadRequest, []string{"jenis penulis tidak valid"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		deleteBatchDokumenPaten(dokumenPaten)
		return util.FailedResponse(c, http.StatusBadRequest, []string{err.Error()})
	}

	return util.SuccessResponse(c, http.StatusCreated, nil)
}

func EditPatenHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	claims := util.GetClaimsFromContext(c)

	idDosen := 0
	if err := db.WithContext(ctx).Table("paten").Select("id_dosen").
		Where("id", id).Scan(&idDosen).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if idDosen != int(claims["id"].(float64)) {
		return util.FailedResponse(c, http.StatusUnauthorized, nil)
	}

	req := &request.Paten{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err.Error()})
	}

	tx := db.Begin()
	paten, errMapping := req.MapRequest()
	if errMapping != nil {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{errMapping.Error()})
	}

	paten.IdDosen = idDosen

	if err := tx.WithContext(ctx).Create(paten).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "id_dosen") {
			return util.FailedResponse(c, http.StatusNotFound, []string{"dosen tidak ditemukan"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	dokumenPaten := []model.DokumenPaten{}
	form, _ := c.MultipartForm()
	if form != nil && paten.Dokumen != nil {
		files := form.File["files"]
		minLen := util.CountMin(len(paten.Dokumen), len(files))
		for i := 0; i < minLen; i++ {
			file := files[i]
			dFile, err := storage.CreateFile(file, patenFolderId)
			if err != nil {
				tx.Rollback()
				return c.JSON(http.StatusBadRequest, err.Error())
			}

			dokumenPaten = append(dokumenPaten, *req.Dokumen[i].MapRequest(&request.DokumenPatenPayload{
				IdFile:    dFile.Id,
				IdPaten:   paten.ID,
				NamaFile:  dFile.Name,
				JenisFile: dFile.MimeType,
				Url:       util.CreateFileUrl(dFile.Id),
			}))

			if paten.Dokumen[i].Nama == "" {
				dokumenPaten[i].Nama = dFile.Name
			}
		}

		if err := tx.WithContext(ctx).Create(&dokumenPaten).Error; err != nil {
			tx.Rollback()
			deleteBatchDokumenPaten(dokumenPaten)
			if strings.Contains(err.Error(), "jenis_dokumen") {
				return util.FailedResponse(c, http.StatusBadRequest, []string{"jenis dokumen tidak valid"})
			}

			return util.FailedResponse(c, http.StatusInternalServerError, nil)
		}
	}

	penulis := []model.PenulisPaten{}
	for _, v := range req.Penulis {
		penulis = append(penulis, *v.MapRequest(paten.ID))
	}

	if err := tx.WithContext(ctx).Delete(new(model.PenulisPaten), "id_paten", paten.ID).Error; err != nil {
		tx.Rollback()
		deleteBatchDokumenPaten(dokumenPaten)
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if err := tx.WithContext(ctx).Create(&penulis).Error; err != nil {
		tx.Rollback()
		deleteBatchDokumenPaten(dokumenPaten)
		if strings.Contains(err.Error(), "jenis_penulis") {
			return util.FailedResponse(c, http.StatusBadRequest, []string{"jenis penulis tidak valid"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		deleteBatchDokumenPaten(dokumenPaten)
		return util.FailedResponse(c, http.StatusBadRequest, []string{err.Error()})
	}

	return util.SuccessResponse(c, http.StatusCreated, nil)
}

func DeletePatenHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != "" {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if !patenAuthorization(c, id, db, ctx) {
		return util.FailedResponse(c, http.StatusUnauthorized, nil)
	}

	dokumen := []string{}
	if err := db.WithContext(ctx).Select("id").Where("id_paten", id).Find(&dokumen).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	for _, id := range dokumen {
		if err := storage.DeleteFile(id); err != nil {
			return util.FailedResponse(c, http.StatusInternalServerError, nil)
		}
	}

	query := db.WithContext(ctx).Delete(new(model.Paten), id)
	if query.Error != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if query.RowsAffected < 1 {
		return util.FailedResponse(c, http.StatusNotFound, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func GetAllKategoriPatenHandler(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.JenisKategoriPaten{}

	if err := db.WithContext(ctx).Preload("KategoriPaten").Find(&data).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDokumenPatenByIdHandler(c echo.Context) error {
	id := c.Param("id")
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.DokumenPaten{}

	if err := db.WithContext(ctx).First(data, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, nil)
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func EditDokumenPatenHandler(c echo.Context) error {
	id := c.Param("id")
	req := &request.DokumenPaten{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	idPaten := 0
	if err := db.WithContext(ctx).Table("dokumen_paten").First(&idPaten, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, nil)
		}
	}

	if !patenAuthorization(c, idPaten, db, ctx) {
		return util.FailedResponse(c, http.StatusUnauthorized, nil)
	}

	file, _ := c.FormFile("file")
	if file == nil {
		return util.FailedResponse(c, http.StatusBadRequest, []string{"file tidak boleh kosong"})
	}

	if err := storage.DeleteFile(id); err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, []string{err.Error()})
	}

	dFile, err := storage.CreateFile(file, patenFolderId)
	if err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, []string{err.Error()})
	}

	dokumenPaten := req.MapRequest(&request.DokumenPatenPayload{
		IdFile:    dFile.Id,
		IdPaten:   idPaten,
		NamaFile:  dFile.Name,
		JenisFile: dFile.MimeType,
		Url:       util.CreateFileUrl(dFile.Id),
	})

	if dokumenPaten.Nama == "" {
		dokumenPaten.Nama = dFile.Name
	}

	if err := db.WithContext(ctx).Create(&dokumenPaten).Error; err != nil {
		if strings.Contains(err.Error(), "jenis_dokumen") {
			return util.FailedResponse(c, http.StatusBadRequest, []string{"jenis dokumen tidak valid"})
		}

		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeleteDokumenPatenHandler(c echo.Context) error {
	id := c.Param("id")
	req := &request.DokumenPaten{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(c, http.StatusUnprocessableEntity, []string{err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	idPaten := 0
	if err := db.WithContext(ctx).Table("dokumen_paten").First(&idPaten, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(c, http.StatusNotFound, nil)
		}
	}

	if !patenAuthorization(c, idPaten, db, ctx) {
		return util.FailedResponse(c, http.StatusUnauthorized, nil)
	}

	idDokumen := ""
	if err := db.WithContext(ctx).Select("id").Where("id", id).Find(&idDokumen).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	if idDokumen == "" {
		return util.FailedResponse(c, http.StatusNotFound, nil)
	}

	tx := db.Begin()

	query := tx.WithContext(ctx).Delete(new(model.DokumenPaten), id)
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
		return util.FailedResponse(c, http.StatusBadRequest, []string{err.Error()})
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func patenAuthorization(c echo.Context, id int, db *gorm.DB, ctx context.Context) bool {
	claims := util.GetClaimsFromContext(c)
	role := claims["role"].(string)
	idDosen := int(claims["id"].(float64))
	if role == string(util.ADMIN) {
		return claims["bagian"].(string) == util.IKU5
	}

	result := 0
	if err := db.WithContext(ctx).Table("paten").Select("id_dosen").
		Where("id", id).Scan(&result).Error; err != nil {
		return false
	}

	return result == idDosen
}

func deleteBatchDokumenPaten(dokumen []model.DokumenPaten) {
	for _, v := range dokumen {
		storage.DeleteFile(v.ID)
	}
}