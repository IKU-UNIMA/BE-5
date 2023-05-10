package handler

import (
	"be-5/src/api/request"
	"be-5/src/api/response"
	"be-5/src/config/database"
	"be-5/src/util"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type dashboardQueryParam struct {
	Tahun int `query:"tahun"`
}

func GetDashboardHandler(c echo.Context) error {
	queryParams := &dashboardQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.Dashboard{}

	dosen := []struct {
		ID       int
		Fakultas string
		Jumlah   int
	}{}
	dosenQuery := `
	SELECT fakultas.id, fakultas.nama AS fakultas, COUNT(dosen.id) AS jumlah FROM fakultas
	left JOIN dosen ON dosen.id_fakultas = fakultas.id
	GROUP BY fakultas.id ORDER BY fakultas.id
	`
	if err := db.WithContext(ctx).Raw(dosenQuery).Find(&dosen).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	var target float64
	targetQuery := `
	SELECT target FROM target
	WHERE bagian = 'IKU 5'
	`
	if err := db.WithContext(ctx).Raw(targetQuery).Find(&target).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data.Target = fmt.Sprintf("%.1f", util.RoundFloat(target))

	query := func(fitur string) string {
		condition := createTahunCondition(fitur, queryParams.Tahun)
		if fitur == "publikasi" && condition != "" {
			splitTahun := strings.Split(condition, " OR ")
			condition = splitTahun[0] + " OR " + "publikasi.id_dosen = dosen.id AND " + splitTahun[1]
		}

		return fmt.Sprintf(
			`SELECT COUNT(%s.id) FROM fakultas
			LEFT JOIN dosen ON dosen.id_fakultas = fakultas.id
			LEFT JOIN %s ON %s.id_dosen = dosen.id %s
			GROUP BY fakultas.id ORDER BY fakultas.id;`,
			fitur, fitur, fitur, condition,
		)
	}

	// get publikasi data
	publikasi := []int{}
	if err := db.WithContext(ctx).Raw(query("publikasi")).Find(&publikasi).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	// get paten data
	paten := []int{}
	if err := db.WithContext(ctx).Raw(query("paten")).Find(&paten).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	// get pengabdian data
	pengabdian := []int{}
	if err := db.WithContext(ctx).Raw(query("pengabdian")).Find(&pengabdian).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	var totalDosen, total int
	for i := 0; i < len(publikasi); i++ {
		jumlah := publikasi[i]
		jumlah += paten[i]
		jumlah += pengabdian[i]

		var persentase float64
		if dosen[i].Jumlah != 0 {
			persentase = util.RoundFloat((float64(jumlah) / float64(dosen[i].Jumlah)) * 100)
		}

		data.Detail = append(data.Detail, response.DashboardDetailPerFakultas{
			ID:          dosen[i].ID,
			Fakultas:    dosen[i].Fakultas,
			JumlahDosen: dosen[i].Jumlah,
			Jumlah:      jumlah,
			Persentase:  fmt.Sprintf("%.2f", persentase) + "%",
		})

		totalDosen += dosen[i].Jumlah
		total += jumlah
	}

	data.Total = total
	data.TotalDosen = totalDosen

	var pencapaian float64
	if totalDosen != 0 {
		pencapaian = util.RoundFloat((float64(total) / float64(totalDosen)) * 100)
	}

	data.Pencapaian = fmt.Sprintf("%.2f", pencapaian) + "%"

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDashboardByFakultasHandler(c echo.Context) error {
	queryParams := &dashboardQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	fakultas, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.DashboardPerProdi{}
	data.Detail = []response.DashboardDetailPerProdi{}

	fakultasConds := ""
	if fakultas != 0 {
		fakultasConds = fmt.Sprintf("WHERE prodi.id_fakultas = %d", fakultas)
	}

	dosen := []struct {
		Jumlah    int
		KodeProdi int
		Prodi     string
		Jenjang   string
	}{}

	dosenQuery := fmt.Sprintf(`
	SELECT COUNT(dosen.id) as jumlah, prodi.kode_prodi, prodi.nama as prodi, prodi.jenjang FROM prodi
	left JOIN dosen ON dosen.id_prodi = prodi.id
	%s GROUP BY prodi.id ORDER BY prodi.id
	`, fakultasConds)
	if err := db.WithContext(ctx).Raw(dosenQuery).Find(&dosen).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	var totalDosen, total int
	if len(dosen) != 0 {
		query := func(fitur string) string {
			condition := createTahunCondition(fitur, queryParams.Tahun)
			if fitur == "publikasi" && condition != "" {
				splitTahun := strings.Split(condition, " OR ")
				condition = splitTahun[0] + " OR " + "publikasi.id_dosen = dosen.id AND " + splitTahun[1]
			}

			return fmt.Sprintf(
				`SELECT COUNT(%s.id) FROM prodi
				LEFT JOIN dosen ON dosen.id_prodi = prodi.id
				LEFT JOIN %s ON %s.id_dosen = dosen.id %s
				%s
				GROUP BY prodi.id ORDER BY prodi.id;`,
				fitur, fitur, fitur, condition, fakultasConds,
			)
		}

		// get publikasi data
		publikasi := []int{}
		if err := db.WithContext(ctx).Debug().Raw(query("publikasi")).Find(&publikasi).Error; err != nil {
			return util.FailedResponse(http.StatusInternalServerError, nil)
		}

		// get paten data
		paten := []int{}
		if err := db.WithContext(ctx).Raw(query("paten")).Find(&paten).Error; err != nil {
			return util.FailedResponse(http.StatusInternalServerError, nil)
		}

		// get pengabdian data
		pengabdian := []int{}
		if err := db.WithContext(ctx).Raw(query("pengabdian")).Find(&pengabdian).Error; err != nil {
			return util.FailedResponse(http.StatusInternalServerError, nil)
		}

		for i := 0; i < len(publikasi); i++ {
			jumlah := publikasi[i]
			jumlah += paten[i]
			jumlah += pengabdian[i]

			var persentase float64
			if dosen[i].Jumlah != 0 {
				persentase = util.RoundFloat((float64(jumlah) / float64(dosen[i].Jumlah)) * 100)
			}

			prodi := fmt.Sprintf("%d - %s (%s)", dosen[i].KodeProdi, dosen[i].Prodi, dosen[i].Jenjang)

			data.Detail = append(data.Detail, response.DashboardDetailPerProdi{
				Prodi:       prodi,
				JumlahDosen: dosen[i].Jumlah,
				Jumlah:      jumlah,
				Persentase:  fmt.Sprintf("%.2f", persentase) + "%",
			})

			total += jumlah
			totalDosen += dosen[i].Jumlah
		}
	}

	data.Total = total
	data.TotalDosen = totalDosen

	var pencapaian float64
	if totalDosen != 0 {
		pencapaian = util.RoundFloat((float64(total) / float64(totalDosen)) * 100)
	}

	data.Pencapaian = fmt.Sprintf("%.2f", pencapaian) + "%"

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDashboardTotalHandler(c echo.Context) error {
	tahun := c.QueryParam("tahun")
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := []response.DashboardTotal{}
	publikasiQuery := "SELECT COUNT(id) AS total FROM publikasi"
	patenQuery := "SELECT COUNT(id) AS total FROM paten"
	pengabdianQuery := "SELECT COUNT(id) AS total FROM pengabdian"

	if tahun != "" {
		publikasiQuery += fmt.Sprintf(" WHERE YEAR(tanggal_terbit) = %s OR YEAR(waktu_pelaksanaan) = %s", tahun, tahun)
		patenQuery += fmt.Sprintf(" WHERE YEAR(tanggal) = %s", tahun)
		pengabdianQuery += fmt.Sprintf(" WHERE tahun_pelaksanaan = %s", tahun)
	}

	total := 0
	// get total publikasi
	if err := db.WithContext(ctx).Raw(publikasiQuery).Find(&total).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data = append(data, response.DashboardTotal{
		Nama:  "Publikasi",
		Total: total,
	})

	// get total paten
	if err := db.WithContext(ctx).Raw(patenQuery).Find(&total).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data = append(data, response.DashboardTotal{
		Nama:  "Paten",
		Total: total,
	})

	// get total pengabdian
	if err := db.WithContext(ctx).Raw(pengabdianQuery).Find(&total).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data = append(data, response.DashboardTotal{
		Nama:  "Pengabdian",
		Total: total,
	})

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDashboardUmum(c echo.Context) error {
	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.DashboardUmum{}
	fakultasQuery := `SELECT COUNT(id) AS fakultas FROM fakultas`
	prodiQuery := `SELECT COUNT(id) AS prodi FROM prodi`
	dosenQuery := `SELECT COUNT(id) AS dosen FROM dosen`
	mahasiswaQuery := `SELECT COUNT(id) AS mahasiswa FROM mahasiswa`

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

func InsertTargetHandler(c echo.Context) error {
	req := &request.Target{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	conds := fmt.Sprintf("bagian='%s' AND tahun=%d", util.IKU5, req.Tahun)

	if err := db.WithContext(ctx).Where(conds).Save(req.MapRequest()).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func createTahunCondition(fitur string, tahun int) string {
	if tahun < 2000 {
		return ""
	}

	conds := "AND "
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
