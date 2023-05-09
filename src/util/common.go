package util

import (
	"math"
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

func CountMin(x, y int) int {
	if x < y {
		return x
	}

	return y
}

func RoundFloat(v float64) float64 {
	return math.Round(v*100) / 100
}

func GetClaimsFromContext(c echo.Context) jwt.MapClaims {
	claims := c.Get("claims")
	return claims.(jwt.MapClaims)
}
