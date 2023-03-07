package main

import (
	"be-5/src/api/route"
	"be-5/src/config/database"
	"be-5/src/config/env"

	"github.com/joho/godotenv"
)

func main() {
	// load env file
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	// migrate gorm
	database.MigrateMySQL()

	app := route.InitServer()
	app.Logger.Fatal(app.Start(":" + env.GetServerEnv()))
}
