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
	Tahun  int    `query:"tahun"`
	Status string `query:"status"`
	Judul  string `query:"judul"`
	Page   int    `query:"page"`
}

func GetAllPublikasiHandler(c echo.Context) error {
	queryParams := &publikasiQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	claims := util.GetClaimsFromContext(c)
	role := claims["role"].(string)
	idDosen := int(claims["id"].(float64))

	var order, condition string
	if role == string(util.DOSEN) {
		order = "tanggal_terbit DESC"
		condition = fmt.Sprintf("id_dosen = %d", idDosen)
	} else {
		order = "created_at DESC"
		if queryParams.Tahun != 0 {
			condition = fmt.Sprintf(`YEAR(tanggal_terbit) = %d`, queryParams.Tahun)
		}

		if queryParams.Status != "" {
			if condition != "" {
				condition += fmt.Sprintf(" AND status = '%s'", queryParams.Status)
			} else {
				condition = fmt.Sprintf("status = '%s'", queryParams.Status)
			}
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
		Where(condition).Order(order).Find(&data).Error; err != nil {
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
	ctx := c.Request().Context()
	publikasi, err := req.MapRequest()
	if err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	// mapping penulis
	penulis := []model.PenulisPublikasi{}
	for _, v := range req.PenulisDosen {
		if err := validation.ValidatePenulis(&v); err != nil {
			return err
		}

		penulis = append(penulis, *v.MapRequestToPublikasi("dosen"))
	}

	for _, v := range req.PenulisMahasiswa {
		if len(req.PenulisMahasiswa) == 1 && req.PenulisMahasiswa[0].Nama == "" {
			break
		}

		if err := validation.ValidatePenulis(&v); err != nil {
			return err
		}

		penulis = append(penulis, *v.MapRequestToPublikasi("mahasiswa"))
	}

	for _, v := range req.PenulisLain {
		if len(req.PenulisLain) == 1 && req.PenulisLain[0].Nama == "" {
			break
		}

		if err := validation.ValidatePenulis(&v); err != nil {
			return err
		}

		penulis = append(penulis, *v.MapRequestToPublikasi("lain"))
	}

	publikasi.IdDosen = idDosen
	publikasi.Penulis = penulis

	tx := db.Begin()
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
		tx.Rollback()
		helper.DeleteBatchDokumen(idDokumen)
		return err
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

	publikasi, errMapping := req.MapRequest()
	if errMapping != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": errMapping.Error()})
	}

	// mapping penulis
	penulis := []model.PenulisPublikasi{}
	for _, v := range req.PenulisDosen {
		if err := validation.ValidatePenulis(&v); err != nil {
			return err
		}

		penulis = append(penulis, *v.MapRequestToPublikasi("dosen"))
	}

	for _, v := range req.PenulisMahasiswa {
		if len(req.PenulisMahasiswa) == 1 && req.PenulisMahasiswa[0].Nama == "" {
			break
		}

		if err := validation.ValidatePenulis(&v); err != nil {
			return err
		}

		penulis = append(penulis, *v.MapRequestToPublikasi("mahasiswa"))
	}

	for _, v := range req.PenulisLain {
		if len(req.PenulisLain) == 1 && req.PenulisLain[0].Nama == "" {
			break
		}

		if err := validation.ValidatePenulis(&v); err != nil {
			return err
		}

		penulis = append(penulis, *v.MapRequestToPublikasi("lain"))
	}

	kategoriCapaian := ""
	if publikasi.IdKategoriCapaian == 0 {
		kategoriCapaian = "null"
	} else {
		kategoriCapaian = fmt.Sprintf("%d", publikasi.IdKategoriCapaian)
	}

	// edit publikasi
	query := fmt.Sprintf(`
	UPDATE publikasi SET
		id_kategori=%d, id_jenis_penelitian=%d, id_kategori_capaian=%s,
		judul='%s', judul_asli='%s', judul_chapter='%s',
		nama_jurnal='%s', nama_koran_majalah='%s', nama_seminar='%s', tautan_laman_jurnal='%s',
		tanggal_terbit='%s', volume='%s', edisi='%s', nomor='%s', halaman='%s', jumlah_halaman=%d,
		penerbit='%s', penyelenggara='%s', kota_penyelenggaraan='%s', is_seminar=%t, is_prosiding=%t,
		bahasa='%s', doi='%s', isbn='%s', issn='%s', e_issn='%s', tautan='%s', keterangan='%s'
	WHERE id=%d
	`, publikasi.IdKategori, publikasi.IdJenisPenelitian,
		kategoriCapaian,
		publikasi.Judul, publikasi.JudulAsli, publikasi.JudulChapter,
		publikasi.NamaJurnal, publikasi.NamaKoranMajalah, publikasi.NamaSeminar, publikasi.TautanLamanJurnal,
		req.TanggalTerbit,
		publikasi.Volume, publikasi.Edisi, publikasi.Nomor, publikasi.Halaman, publikasi.JumlahHalaman,
		publikasi.Penerbit, publikasi.Penyelenggara, publikasi.KotaPenyelenggaraan, publikasi.IsSeminar, publikasi.IsProsiding,
		publikasi.Bahasa, publikasi.Doi, publikasi.Isbn, publikasi.Issn, publikasi.EIssn, publikasi.Tautan, publikasi.Keterangan,
		id,
	)

	tx := db.Begin()
	if err := tx.WithContext(ctx).Exec(query).Error; err != nil {
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
		IdFitur: id,
		Data:    req.Dokumen,
	})

	if err != nil {
		tx.Rollback()
		helper.DeleteBatchDokumen(idDokumen)
		return err
	}

	// delete old penulis
	if err := tx.WithContext(ctx).Delete(new(model.PenulisPublikasi), "id_publikasi", id).Error; err != nil {
		tx.Rollback()
		helper.DeleteBatchDokumen(idDokumen)
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	// insert penulis
	if err := tx.WithContext(ctx).Model(&model.Publikasi{ID: id}).Association("Penulis").Replace(&penulis); err != nil {
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

	if strings.Contains(err.Error(), "jenis_penulis") {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "jenis penulis tidak valid"})
	}

	if strings.Contains(err.Error(), "peran") {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "peran tidak valid"})
	}

	return util.FailedResponse(http.StatusInternalServerError, nil)
}
