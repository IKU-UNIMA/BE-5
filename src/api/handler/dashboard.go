package handler

import (
	"be-5/src/api/response"
	"be-5/src/config/database"
	"be-5/src/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type dashboardQueryParam struct {
	Tahun    int `query:"tahun"`
	Fakultas int `query:"fakultas"`
}

func GetDashboardHandler(c echo.Context) error {
	queryParams := &dashboardQueryParam{}
	if err := (&echo.DefaultBinder{}).BindPathParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	fitur := c.Param("fitur")
	if err := checkDashboardFitur(c, fitur); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.Dashboard{}
	condition := createTahunCondition(fitur, queryParams.Tahun)
	query := fmt.Sprintf(
		`SELECT fakultas.id, fakultas.nama, COUNT(*) AS jumlah FROM %s
			JOIN dosen ON dosen.id = %s.id_dosen
			JOIN fakultas ON fakultas.id = dosen.id_fakultas 
			%s GROUP BY fakultas.id ORDER BY jumlah DESC;`,
		fitur, fitur, condition,
	)

	if err := db.WithContext(ctx).Raw(query).Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDetailDashboardHandler(c echo.Context) error {
	fitur := c.Param("fitur")

	if err := checkDashboardFitur(c, fitur); err != nil {
		return err
	}

	queryParams := &dashboardQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.DetailDashboard{}
	condition := createTahunCondition(fitur, queryParams.Tahun)
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
		fitur, fitur, condition,
	)

	if err := db.WithContext(ctx).Raw(query).Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDashboardTotalHandler(c echo.Context) error {
	tahun := c.QueryParam("tahun")
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.DashboardTotal{}
	publikasiQuery := "SELECT COUNT(id) AS total_publikasi FROM publikasi"
	patenQuery := "SELECT COUNT(id) AS total_paten FROM paten"
	pengabdianQuery := "SELECT COUNT(id) AS total_pengabdian FROM pengabdian"

	if tahun != "" {
		publikasiQuery += fmt.Sprintf(" WHERE YEAR(tanggal_terbit) = %s OR YEAR(waktu_pelaksanaan) = %s", tahun, tahun)
		patenQuery += fmt.Sprintf(" WHERE YEAR(tanggal) = %s", tahun)
		pengabdianQuery += fmt.Sprintf(" WHERE tahun_pelaksanaan = %s", tahun)
	}

	// get total publikasi
	if err := db.WithContext(ctx).Raw(publikasiQuery).Find(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	// get total paten
	if err := db.WithContext(ctx).Raw(patenQuery).Find(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	// get total pengabdian
	if err := db.WithContext(ctx).Raw(pengabdianQuery).Find(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDashboardUmum(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.DashboardUmum{}
	fakultasQuery := `SELECT COUNT(id) AS fakultas FROM fakultas`
	prodiQuery := `SELECT COUNT(id) AS fakultas FROM fakultas`
	dosenQuery := `SELECT COUNT(id) AS fakultas FROM fakultas`
	mahasiswaQuery := `SELECT COUNT(id) AS fakultas FROM fakultas`

	if err := db.WithContext(ctx).Raw(fakultasQuery).Find(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Raw(prodiQuery).Find(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Raw(dosenQuery).Find(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Raw(mahasiswaQuery).Find(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
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
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "fitur tidak didukung"})
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
