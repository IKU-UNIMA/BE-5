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

	fakultas := v1.Group("/fakultas", customMiddleware.Authentication)
	fakultas.GET("", handler.GetAllFakultasHandler)
	fakultas.GET("/:id", handler.GetFakultasByIdHandler)
	fakultas.POST("", handler.InsertFakultasHandler, customMiddleware.GrantAdminUmum)
	fakultas.PUT("/:id", handler.EditFakultasHandler, customMiddleware.GrantAdminUmum)
	fakultas.DELETE("/:id", handler.DeleteFakultasHandler, customMiddleware.GrantAdminUmum)

	prodi := v1.Group("/prodi", customMiddleware.Authentication)
	prodi.GET("", handler.GetAllProdiHandler)
	prodi.GET("/:id", handler.GetProdiByIdHandler)
	prodi.POST("", handler.InsertProdiHandler, customMiddleware.GrantAdminUmum)
	prodi.PUT("/:id", handler.EditProdiHandler, customMiddleware.GrantAdminUmum)
	prodi.DELETE("/:id", handler.DeleteProdiHandler, customMiddleware.GrantAdminUmum)

	jenisDokumen := v1.Group("/jenis-dokumen", customMiddleware.Authentication)
	jenisDokumen.GET("", handler.GetAllJenisDokumenHandler)

	jenisPenelitian := v1.Group("/jenis-penelitian", customMiddleware.Authentication)
	jenisPenelitian.GET("", handler.GetAllJenisPenelitianHandler)

	kategoriCapaian := v1.Group("/kategori-capaian", customMiddleware.Authentication)
	kategoriCapaian.GET("", handler.GetAllKategoriCapaianHandler)

	akun := v1.Group("/akun")
	akun.POST("/login", handler.LoginHandler)
	akun.PATCH("/password/change", handler.ChangePasswordHandler, customMiddleware.Authentication)
	akun.PATCH("/password/reset/:id", handler.ResetPasswordHandler, customMiddleware.Authentication, customMiddleware.GrantAdminUmum)

	profil := v1.Group("/profil", customMiddleware.Authentication)
	profil.GET("", handler.GetProfilHandler)
	profil.PATCH("", handler.EditProfilHandler)

	admin := v1.Group("/admin", customMiddleware.Authentication, customMiddleware.GrantAdminUmum)
	admin.GET("", handler.GetAllAdminHandler)
	admin.GET("/:id", handler.GetAdminByIdHandler)
	admin.POST("", handler.InsertAdminHandler)
	admin.PUT("/:id", handler.EditAdminHandler)
	admin.DELETE("/:id", handler.DeleteAdminHandler)

	rektor := v1.Group("/rektor", customMiddleware.Authentication, customMiddleware.GrantAdminUmum)
	rektor.GET("", handler.GetAllRektorHandler)
	rektor.GET("/:id", handler.GetRektorByIdHandler)
	rektor.POST("", handler.InsertRektorHandler)
	rektor.PUT("/:id", handler.EditRektorHandler)
	rektor.DELETE("/:id", handler.DeleteRektorHandler)

	dosen := v1.Group("/dosen", customMiddleware.Authentication)
	dosen.GET("", handler.GetAllDosenHandler)
	dosen.GET("/:id", handler.GetDosenByIdHandler)
	dosen.POST("", handler.InsertDosenHandler, customMiddleware.GrantAdminUmum)
	dosen.PUT("/:id", handler.EditDosenHandler, customMiddleware.GrantAdminUmum)
	dosen.DELETE("/:id", handler.DeleteDosenHandler, customMiddleware.GrantAdminUmum)

	paten := v1.Group("/paten", customMiddleware.Authentication)
	paten.GET("", handler.GetAllPatenHandler)
	paten.GET("/:id", handler.GetPatenByIdHandler)
	paten.POST("", handler.InsertPatenHandler, customMiddleware.GrantDosen)
	paten.PUT("/:id", handler.EditPatenHandler, customMiddleware.GrantDosen)
	paten.DELETE("/:id", handler.DeletePatenHandler)
	paten.GET("/kategori", handler.GetAllKategoriPatenHandler)
	paten.GET("/dokumen/:id", handler.GetDokumenPatenByIdHandler)
	paten.PUT("/dokumen/:id", handler.EditDokumenPatenHandler, customMiddleware.GrantDosen)
	paten.DELETE("/dokumen/:id", handler.DeleteDokumenPatenHandler, customMiddleware.GrantDosen)

	return app
}
