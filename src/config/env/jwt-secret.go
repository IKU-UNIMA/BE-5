package env

import "os"

func GetSecretJWTEnv() string {
	return os.Getenv("SECRET_JWT")
}
