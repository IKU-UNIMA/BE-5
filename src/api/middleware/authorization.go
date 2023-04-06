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

func GrantAdminUmum(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		if claims["role"].(string) != string(util.ADMIN) ||
			claims["bagian"].(string) != "UMUM" {
			return util.FailedResponse(c, http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}

func GrantDosen(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		if claims["role"].(string) != string(util.DOSEN) {
			return util.FailedResponse(c, http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}

func GrantAdminAndDosen(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		role := claims["role"].(string)
		if role != string(util.ADMIN) && role != string(util.DOSEN) {
			return util.FailedResponse(c, http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}

func GrantAdminIKU5AndDosen(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		role := claims["role"].(string)
		bagian := claims["bagian"].(string)
		if role != string(util.ADMIN) && role != string(util.DOSEN) {
			return util.FailedResponse(c, http.StatusUnauthorized, nil)
		}

		if role == string(util.ADMIN) && bagian != "IKU 5" {
			return util.FailedResponse(c, http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}
