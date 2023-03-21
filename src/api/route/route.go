package route

import (
	"be-5/src/api/handler"

	customMiddleware "be-5/src/api/middleware"

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

	prodi := v1.Group("/prodi")
	prodi.GET("", handler.GetAllProdiHandler)
	prodi.GET("/:id", handler.GetProdiByIdHandler)
	prodi.POST("", handler.InsertProdiHandler)
	prodi.PUT("/:id", handler.EditProdiHandler)
	prodi.DELETE("/:id", handler.DeleteProdiHandler)

	jenisDokumen := v1.Group("/jenis-dokumen")
	jenisDokumen.GET("", handler.GetAllJenisDokumenHandler)

	jenisPenelitian := v1.Group("/jenis-penelitian")
	jenisPenelitian.GET("", handler.GetAllJenisPenelitianHandler)

	kategoriCapaian := v1.Group("/kategori-capaian")
	kategoriCapaian.GET("", handler.GetAllKategoriCapaianHandler)

	akun := v1.Group("/akun")
	akun.POST("/login", handler.LoginHandler)
	akun.PATCH("/password/change", handler.ChangePasswordHandler, customMiddleware.Authentication)
	akun.PATCH("/password/reset/:id", handler.ResetPasswordHandler, customMiddleware.Authentication, customMiddleware.GrantAdmin)

	profil := v1.Group("/profil")
	profil.PATCH("", handler.EditProfilHandler)

	dosen := v1.Group("/dosen", customMiddleware.Authentication)
	dosen.GET("", handler.GetAllDosenHandler)
	dosen.GET("/:id", handler.GetDosenByIdHandler)
	dosen.POST("", handler.InsertDosenHandler, customMiddleware.GrantAdmin)
	dosen.PUT("/:id", handler.EditDosenHandler, customMiddleware.GrantAdmin)
	dosen.DELETE("/:id", handler.DeleteDosenHandler, customMiddleware.GrantAdmin)

	return app
}
