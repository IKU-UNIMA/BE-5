package helper

import (
	"be-5/src/api/request"
	"be-5/src/config/env"
	"be-5/src/config/storage"
	"be-5/src/util"
	"be-5/src/util/validation"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type (
	InsertDokumenParam struct {
		C       echo.Context
		Ctx     context.Context
		DB      *gorm.DB
		TX      *gorm.DB
		Fitur   string
		IdFitur int
		Data    []request.Dokumen
	}

	EditDokumenParam struct {
		C     echo.Context
		Ctx   context.Context
		DB    *gorm.DB
		Fitur string
		Id    string
	}

	DokumenModel struct {
		IdJenisDokumen int
		Nama           string
		NamaFile       string
		JenisFile      string
		Keterangan     string
		Url            string
	}
)

func InsertDokumen(param InsertDokumenParam) ([]string, error) {
	query := ""
	req := param.Data
	form, _ := param.C.MultipartForm()
	files := form.File["files"]
	idDokumen := []string{}
	if files != nil && req != nil {
		query = fmt.Sprintf(`INSERT dokumen_%s
		(id, id_%s, id_jenis_dokumen, nama, nama_file, jenis_file, url, keterangan, tanggal_upload) 
		VALUES `, param.Fitur, param.Fitur)

		minLen := util.CountMin(len(req), len(files))
		for i := 0; i < minLen; i++ {
			if err := validation.ValidateDokumen(&req[i]); err != nil {
				param.TX.Rollback()
				return nil, err
			}

			dFile, err := storage.CreateFile(files[i], env.GetPengabdianFolderId())
			if err != nil {
				param.TX.Rollback()
				if strings.Contains(err.Error(), "unsupported") {
					return nil, util.FailedResponse(param.C, http.StatusBadRequest, map[string]string{"message": err.Error()})
				}

				return nil, util.FailedResponse(param.C, http.StatusInternalServerError, nil)
			}

			idDokumen = append(idDokumen, dFile.Id)

			dokumen := req[i]
			if dokumen.Nama == "" {
				dokumen.Nama = dFile.Name
			}

			year, month, day := time.Now().Date()
			tanggalUpload := fmt.Sprintf("%d-%d-%d", year, month, day)
			query += fmt.Sprintf("('%s',%d,%d,'%s','%s','%s','%s','%s','%s')",
				dFile.Id, param.IdFitur, dokumen.IdJenisDokumen,
				dokumen.Nama, dFile.Name, dFile.MimeType,
				util.CreateFileUrl(dFile.Id), dokumen.Keterangan, tanggalUpload)

			// add , in every end of the value and add ; in the end of query
			if i < minLen-1 {
				query += ","
			} else {
				query += ";"
			}
		}

		if err := param.TX.WithContext(param.Ctx).Exec(query).Error; err != nil {
			param.TX.Rollback()
			DeleteBatchDokumen(idDokumen)
			if strings.Contains(err.Error(), "jenis_dokumen") {
				return nil, util.FailedResponse(param.C, http.StatusBadRequest, map[string]string{"message": "jenis dokumen tidak valid"})
			}

			return nil, util.FailedResponse(param.C, http.StatusInternalServerError, nil)
		}
	}

	return idDokumen, nil
}

func EditDokumen(param EditDokumenParam) error {
	req := &request.Dokumen{}
	reqData := param.C.FormValue("data")
	if err := json.Unmarshal([]byte(reqData), req); err != nil {
		return util.FailedResponse(param.C, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := validation.ValidateDokumen(req); err != nil {
		return err
	}

	var dokumen *DokumenModel
	file, _ := param.C.FormFile("file")
	if file != nil {

		dFile, err := storage.CreateFile(file, getFolderId(param.Fitur))
		if err != nil {
			if strings.Contains(err.Error(), "unsupported") {
				return util.FailedResponse(param.C, http.StatusBadRequest, map[string]string{"message": err.Error()})
			}

			return util.FailedResponse(param.C, http.StatusInternalServerError, nil)
		}
		dokumen = &DokumenModel{
			IdJenisDokumen: req.IdJenisDokumen,
			Nama:           req.Nama,
			NamaFile:       dFile.Name,
			JenisFile:      dFile.MimeType,
			Url:            util.CreateFileUrl(dFile.Id),
			Keterangan:     req.Keterangan,
		}

		if dokumen.Nama == "" {
			dokumen.Nama = dFile.Name
		}

		newId := dFile.Id
		year, month, day := time.Now().Date()
		tanggalUpload := fmt.Sprintf("%d-%d-%d", year, month, day)
		updateDokumenQuery := fmt.Sprintf(`UPDATE dokumen_%s SET
			id='%s',id_jenis_dokumen=%d,nama='%s',nama_file='%s',jenis_file='%s',url='%s',keterangan='%s',tanggal_upload='%s' WHERE id='%s'`,
			param.Fitur, newId, dokumen.IdJenisDokumen, dokumen.Nama, dokumen.NamaFile, dokumen.JenisFile, dokumen.Url, dokumen.Keterangan, tanggalUpload, param.Id)
		if err := param.DB.WithContext(param.Ctx).Exec(updateDokumenQuery).Error; err != nil {
			storage.DeleteFile(newId)
			if strings.Contains(err.Error(), "jenis_dokumen") {
				return util.FailedResponse(param.C, http.StatusBadRequest, map[string]string{"message": "jenis dokumen tidak valid"})
			}

			return util.FailedResponse(param.C, http.StatusInternalServerError, nil)
		}

		storage.DeleteFile(param.Id)
	} else {
		dokumen = &DokumenModel{
			IdJenisDokumen: req.IdJenisDokumen,
			Nama:           req.Nama,
			Keterangan:     req.Keterangan,
		}

		if err := param.DB.WithContext(param.Ctx).Table(fmt.Sprint("dokumen_", param.Fitur)).Where("id", param.Id).Updates(&dokumen).Error; err != nil {
			if strings.Contains(err.Error(), "jenis_dokumen") {
				return util.FailedResponse(param.C, http.StatusBadRequest, map[string]string{"message": "jenis dokumen tidak valid"})
			}

			return util.FailedResponse(param.C, http.StatusInternalServerError, nil)
		}
	}

	return util.SuccessResponse(param.C, http.StatusOK, nil)
}

func getFolderId(fitur string) string {
	switch fitur {
	case "publikasi":
		return env.GetPublikasiFolderId()
	case "paten":
		return env.GetPatenFolderId()
	}

	return env.GetPengabdianFolderId()
}

func DeleteBatchDokumen(id []string) {
	for _, v := range id {
		storage.DeleteFile(v)
	}
}
