package middleware

import (
	"be-5/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GrantAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		if claims["role"].(string) != string(util.ADMIN) {
			return util.FailedResponse(c, http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}
