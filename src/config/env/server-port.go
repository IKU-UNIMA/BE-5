package env

import "os"

func GetServerEnv() string {
	return os.Getenv("SERVER_PORT")
}
