package route

import (
	"be-5/src/api/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitServer() *echo.Echo {
	app := echo.New()
	app.Use(middleware.CORS())

	app.GET("", func(c echo.Context) error {
		return c.JSON(200, "Welcome to IKU 5 API")
	})

	v1 := app.Group("/api/v1")

	fakultas := v1.Group("/fakultas")
	fakultas.GET("", handler.GetAllFakultasHandler)
	fakultas.GET("/:id", handler.GetFakultasByIdHandler)
	fakultas.POST("", handler.InsertFakultasHandler)
	fakultas.PUT("/:id", handler.EditFakultasHandler)
	fakultas.DELETE("/:id", handler.DeleteFakultasHandler)

	return app
}
