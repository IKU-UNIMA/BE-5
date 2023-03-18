package util

import (
	"be-5/src/config/env"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const JWT_ERROR = "login to process"

var secret = env.GetSecretJWTEnv()

func GenerateJWT(id int, role, bagian string) string {
	claims := jwt.MapClaims{
		"id":         id,
		"role":       role,
		"bagian":     bagian,
		"expires_at": time.Now().Add(12 * time.Hour).Unix(),
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, _ := rawToken.SignedString([]byte(secret))
	return token
}

func ValidateJWT(c echo.Context) (interface{}, error) {
	header := c.Request().Header["Authorization"][0]
	bearer := strings.HasPrefix(strings.ToLower(header), "bearer")
	if !bearer {
		return nil, errors.New(JWT_ERROR)
	}

	strToken := strings.Split(header, " ")
	if len(strToken) != 2 {
		return nil, errors.New(JWT_ERROR)
	}

	token, _ := jwt.Parse(strToken[1], func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, errors.New(JWT_ERROR)
	}

	var mapClaims jwt.MapClaims
	if v, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, errors.New(JWT_ERROR)
	} else {
		mapClaims = v
	}

	if exp, ok := mapClaims["expires_at"].(float64); !ok {
		return nil, errors.New(JWT_ERROR)
	} else {
		if int64(exp)-time.Now().Unix() <= 0 {
			return nil, errors.New(JWT_ERROR)
		}
	}

	if _, ok := mapClaims["id"].(float64); !ok {
		return nil, errors.New(JWT_ERROR)
	}

	if _, ok := mapClaims["role"].(string); !ok {
		return nil, errors.New(JWT_ERROR)
	}

	if _, ok := mapClaims["bagian"].(string); !ok {
		return nil, errors.New(JWT_ERROR)
	}

	return mapClaims, nil
}
