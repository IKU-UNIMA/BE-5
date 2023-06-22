package handler

import (
	"be-5/src/api/request"
	"be-5/src/config/database"
	"be-5/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

type verifPathParam struct {
	ID    int    `param:"id"`
	Fitur string `param:"fitur"`
}

func VerifikasiDataHandler(c echo.Context) error {
	pathParams := &verifPathParam{}
	if err := (&echo.DefaultBinder{}).BindPathParams(c, pathParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	req := &request.VerifikasiData{}
	if err := c.Bind(&req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}
	status := map[string]bool{
		"Belum Diverifikasi":  true,
		"Draft":               true,
		"Tidak Terverifikasi": true,
		"Terverifikasi":       true,
	}

	if !status[req.Status] {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "input anda salah"})
	}

	fitur := map[string]bool{
		"publikasi":  true,
		"paten":      true,
		"pengabdian": true,
	}

	if !fitur[pathParams.Fitur] {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "fitur tidak didukung"})
	}

	db := database.DB
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).Table(pathParams.Fitur).Where(pathParams.ID).Update("status", req.Status); err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
