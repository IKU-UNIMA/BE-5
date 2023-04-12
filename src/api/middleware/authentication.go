package middleware

import (
	"be-5/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Authentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims, err := util.ValidateJWT(c)
		if err != nil {
			return util.FailedResponse(c, http.StatusUnauthorized, map[string]string{"message": err.Error()})
		}

		c.Set("claims", claims)

		return next(c)
	}
}
