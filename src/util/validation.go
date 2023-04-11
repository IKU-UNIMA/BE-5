package util

import (
	"net/http"

	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translation "github.com/go-playground/validator/v10/translations/id"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	id := id.New()
	uni := ut.New(id, id)
	trans, _ := uni.GetTranslator("id")
	id_translation.RegisterDefaultTranslations(cv.Validator, trans)

	if err := cv.Validator.Struct(i); err != nil {
		errs := err.(validator.ValidationErrors)
		httpCode := http.StatusBadRequest
		return echo.NewHTTPError(httpCode, base{
			Status:  httpCode,
			Message: http.StatusText(httpCode),
			Errors:  errs.Translate(trans),
		})
	}

	return nil
}
