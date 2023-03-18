package util

import (
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func IsInteger(value string) bool {
	if value == "" {
		return true
	}
	_, err := strconv.Atoi(value)
	return err == nil
}

func GetClaimsFromContext(c echo.Context) jwt.MapClaims {
	claims := c.Get("claims")
	return claims.(jwt.MapClaims)
}
