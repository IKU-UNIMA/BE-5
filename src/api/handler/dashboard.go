package handler

import (
	"be-5/src/api/response"
	"be-5/src/config/database"
	"be-5/src/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type dashboardPathParam struct {
	Fitur string `param:"fitur"`
	Tahun int    `param:"tahun"`
}

type dashboardQueryParam struct {
	Tahun    int `query:"tahun"`
	Fakultas int `query:"fakultas"`
}

func GetDashboardHandler(c echo.Context) error {
	pathParams := &dashboardPathParam{}
	if err := (&echo.DefaultBinder{}).BindPathParams(c, pathParams); err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := checkDashboardFitur(c, pathParams.Fitur); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.Dashboard{}
	condition := createTahunCondition(pathParams.Fitur, pathParams.Tahun)
	query := fmt.Sprintf(
		`SELECT fakultas.id, fakultas.nama, COUNT(*) AS jumlah FROM %s
			JOIN dosen ON dosen.id = %s.id_dosen
			JOIN fakultas ON fakultas.id = dosen.id_fakultas 
			%s GROUP BY fakultas.id ORDER BY jumlah DESC;`,
		pathParams.Fitur, pathParams.Fitur, condition,
	)

	if err := db.WithContext(ctx).Raw(query).Find(&data).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDetailDashboardHandler(c echo.Context) error {
	pathParams := &dashboardPathParam{}
	if err := (&echo.DefaultBinder{}).BindPathParams(c, pathParams); err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := checkDashboardFitur(c, pathParams.Fitur); err != nil {
		return err
	}

	queryParams := &dashboardQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.DetailDashboard{}
	condition := createTahunCondition(pathParams.Fitur, queryParams.Tahun)
	if queryParams.Fakultas > 0 {
		if condition != "" {
			condition += fmt.Sprintf(" AND dosen.id_fakultas = %d", queryParams.Fakultas)
		} else {
			condition = fmt.Sprintf("WHERE dosen.id_fakultas = %d", queryParams.Fakultas)
		}
	}

	query := fmt.Sprintf(
		`SELECT prodi.id AS id_prodi, kode_prodi, prodi.nama AS nama_prodi, prodi.jenjang AS jenjang_prodi,
			fakultas.id AS id_fakultas, fakultas.nama AS nama_fakultas, COUNT(*) AS jumlah FROM %s
			JOIN dosen ON dosen.id = %s.id_dosen
			JOIN prodi ON prodi.id = dosen.id_prodi
			JOIN fakultas ON fakultas.id = dosen.id_fakultas
			%s GROUP BY prodi.id ORDER BY prodi.id;`,
		pathParams.Fitur, pathParams.Fitur, condition,
	)

	if err := db.WithContext(ctx).Raw(query).Find(&data).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDashboardTotalHandler(c echo.Context) error {
	tahun := c.Param("tahun")
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.DashboardTotal{}
	publikasiQuery := fmt.Sprintf(
		`SELECT COUNT(id) AS total_publikasi FROM publikasi
		WHERE YEAR(tanggal_terbit) = %s OR YEAR(waktu_pelaksanaan) = %s`,
		tahun, tahun,
	)
	patenQuery := fmt.Sprintf(
		`SELECT COUNT(id) AS total_paten FROM paten
		WHERE YEAR(tanggal) = %s`,
		tahun,
	)
	pengabdianQuery := fmt.Sprintf(
		`SELECT COUNT(id) AS total_pengabdian FROM pengabdian
		WHERE tahun_pelaksanaan = %s`,
		tahun,
	)

	// get total publikasi
	if err := db.WithContext(ctx).Raw(publikasiQuery).Find(data).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	// get total paten
	if err := db.WithContext(ctx).Raw(patenQuery).Find(data).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	// get total pengabdian
	if err := db.WithContext(ctx).Raw(pengabdianQuery).Find(data).Error; err != nil {
		return util.FailedResponse(c, http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func checkDashboardFitur(c echo.Context, fitur string) error {
	switch fitur {
	case "publikasi":
		break
	case "paten":
		break
	case "pengabdian":
		break
	default:
		httpCode := http.StatusBadRequest
		return echo.NewHTTPError(httpCode, util.Base{
			Status:  httpCode,
			Message: http.StatusText(httpCode),
			Errors:  map[string]string{"message": "fitur tidak didukung"},
		})
	}

	return nil
}

func createTahunCondition(fitur string, tahun int) string {
	if tahun < 2000 {
		return ""
	}

	conds := "WHERE "
	switch fitur {
	case "publikasi":
		conds += fmt.Sprintf("YEAR(tanggal_terbit) = %d OR YEAR(waktu_pelaksanaan) = %d", tahun, tahun)
	case "paten":
		conds += fmt.Sprintf("YEAR(tanggal) = %d", tahun)
	case "pengabdian":
		conds += fmt.Sprintf("tahun_pelaksanaan = %d", tahun)
	}

	return conds
}
