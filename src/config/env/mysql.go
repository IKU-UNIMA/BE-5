package env

import "os"

func GetMySQLEnv() string {
	return os.Getenv("MYSQL_DSN")
}
