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

func ValidateDokumen(c echo.Context, req *request.Dokumen) error {
	if req.IdJenisDokumen < 1 {
		return util.FailedResponse(c, http.StatusBadRequest, map[string]string{"message": "jenis dokumen wajib diisi"})
	}

	return nil
}

func ValidatePenulis(c echo.Context, req *request.Penulis) error {
	errors := map[string]string{}
	if req.JenisPenulis == "" {
		errors["message"] = "jenis penulis wajib diisi"
	} else if req.Nama == "" {
		errors["message"] = "nama penulis wajib diisi"
	} else if req.Peran == "" {
		errors["message"] = "peran wajib diisi"
	}

	if len(errors) < 1 {
		return nil
	}

	return util.FailedResponse(c, http.StatusBadRequest, errors)
}
