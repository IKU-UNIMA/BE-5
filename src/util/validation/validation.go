package validation

import (
	"be-5/src/api/request"
	"be-5/src/util"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		errs := err.(validator.ValidationErrors)
		httpCode := http.StatusBadRequest
		return echo.NewHTTPError(httpCode, util.Base{
			Status:  httpCode,
			Message: http.StatusText(httpCode),
			Errors:  translate(errs),
		})
	}

	return nil
}

func translate(errs validator.ValidationErrors) map[string]string {
	errors := map[string]string{}
	for _, e := range errs {
		key := createErrorKey(e.Field())
		errors[key] = getTagMessage(e)
	}

	return errors
}

func createErrorKey(key string) string {
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchAllCap.ReplaceAllString(key, "${1}_${2}")
	return strings.ToLower(snake)
}

func getTagMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "field ini wajib diisi"
	case "email":
		return "email harus berupa alamat email yang valid"
	}

	return ""
}

func ValidateDokumen(req *request.Dokumen) error {
	httpCode := http.StatusBadRequest
	if req.IdJenisDokumen < 1 {
		return echo.NewHTTPError(httpCode, util.Base{
			Status:  httpCode,
			Message: http.StatusText(httpCode),
			Errors:  map[string]string{"message": "jenis dokumen wajib diisi"},
		})
	}

	return nil
}

func ValidatePenulis(req *request.Penulis) error {
	errs := map[string]string{}
	if req.JenisPenulis == "" {
		errs["message"] = "jenis penulis wajib diisi"
	} else if req.Nama == "" {
		errs["message"] = "nama penulis wajib diisi"
	} else if req.Peran == "" {
		errs["message"] = "peran wajib diisi"
	}

	if len(errs) < 1 {
		return nil
	}

	httpCode := http.StatusBadRequest
	return echo.NewHTTPError(httpCode, util.Base{
		Status:  httpCode,
		Message: http.StatusText(httpCode),
		Errors:  errs,
	})
}

func ValidateAnggota(req *request.AnggotaPengabdian) error {
	errs := map[string]string{}
	if req.JenisAnggota == "" {
		errs["message"] = "jenis penulis wajib diisi"
	} else if req.Nama == "" {
		errs["message"] = "nama penulis wajib diisi"
	} else if req.Peran == "" {
		errs["message"] = "peran wajib diisi"
	}

	if len(errs) < 1 {
		return nil
	}

	httpCode := http.StatusBadRequest
	return echo.NewHTTPError(httpCode, util.Base{
		Status:  httpCode,
		Message: http.StatusText(httpCode),
		Errors:  errs,
	})
}
