package route

import (
	"be-5/src/api/handler"
	"be-5/src/util/validation"

	customMiddleware "be-5/src/api/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitServer() *echo.Echo {

	app := echo.New()
	app.Use(middleware.CORS())

	app.Validator = &validation.CustomValidator{Validator: validator.New()}

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

	publikasi := v1.Group("/publikasi", customMiddleware.Authentication)
	publikasi.GET("", handler.GetAllPublikasiHandler, customMiddleware.GrantAdminIKU5AndDosen)
	publikasi.GET("/:id", handler.GetPublikasiByIdHandler, customMiddleware.GrantAdminIKU5AndDosen)
	publikasi.POST("", handler.InsertPublikasiHandler, customMiddleware.GrantDosen)
	publikasi.PUT("/:id", handler.EditPublikasiHandler, customMiddleware.GrantAdminIKU5AndDosen)
	publikasi.DELETE("/:id", handler.DeletePublikasiHandler, customMiddleware.GrantAdminIKU5AndDosen)
	publikasi.GET("/kategori", handler.GetAllKategoriPublikasiHandler, customMiddleware.GrantAdminIKU5AndDosen)
	publikasi.GET("/dokumen/:id", handler.GetDokumenPublikasiByIdHandler, customMiddleware.GrantAdminIKU5AndDosen)
	publikasi.PUT("/dokumen/:id", handler.EditDokumenPublikasiHandler, customMiddleware.GrantAdminIKU5AndDosen)
	publikasi.DELETE("/dokumen/:id", handler.DeleteDokumenPublikasiHandler, customMiddleware.GrantAdminIKU5AndDosen)

	paten := v1.Group("/paten", customMiddleware.Authentication)
	paten.GET("", handler.GetAllPatenHandler, customMiddleware.GrantAdminIKU5AndDosen)
	paten.GET("/:id", handler.GetPatenByIdHandler, customMiddleware.GrantAdminIKU5AndDosen)
	paten.POST("", handler.InsertPatenHandler, customMiddleware.GrantDosen)
	paten.PUT("/:id", handler.EditPatenHandler, customMiddleware.GrantAdminIKU5AndDosen)
	paten.DELETE("/:id", handler.DeletePatenHandler, customMiddleware.GrantAdminIKU5AndDosen)
	paten.GET("/kategori", handler.GetAllKategoriPatenHandler, customMiddleware.GrantAdminIKU5AndDosen)
	paten.GET("/dokumen/:id", handler.GetDokumenPatenByIdHandler, customMiddleware.GrantAdminIKU5AndDosen)
	paten.PUT("/dokumen/:id", handler.EditDokumenPatenHandler, customMiddleware.GrantAdminIKU5AndDosen)
	paten.DELETE("/dokumen/:id", handler.DeleteDokumenPatenHandler, customMiddleware.GrantAdminIKU5AndDosen)

	pengabdian := v1.Group("/pengabdian", customMiddleware.Authentication)
	pengabdian.GET("", handler.GetAllPengabdianHandler, customMiddleware.GrantAdminIKU5AndDosen)
	pengabdian.GET("/:id", handler.GetPengabdianByIdHandler, customMiddleware.GrantAdminIKU5AndDosen)
	pengabdian.POST("", handler.InsertPengabdianHandler, customMiddleware.GrantDosen)
	pengabdian.PUT("/:id", handler.EditPengabdianHandler, customMiddleware.GrantAdminIKU5AndDosen)
	pengabdian.DELETE("/:id", handler.DeletePengabdianHandler, customMiddleware.GrantAdminIKU5AndDosen)
	pengabdian.GET("/kategori", handler.GetAllKategoriPengabdianHandler, customMiddleware.GrantAdminIKU5AndDosen)
	pengabdian.GET("/dokumen/:id", handler.GetDokumenPengabdianByIdHandler, customMiddleware.GrantAdminIKU5AndDosen)
	pengabdian.PUT("/dokumen/:id", handler.EditDokumenPengabdianHandler, customMiddleware.GrantAdminIKU5AndDosen)
	pengabdian.DELETE("/dokumen/:id", handler.DeleteDokumenPengabdianHandler, customMiddleware.GrantAdminIKU5AndDosen)

	verif := v1.Group("/verifikasi", customMiddleware.Authentication, customMiddleware.GrantAdminIKU5)
	verif.PATCH("/:fitur/:id", handler.VerifikasiDataHandler)

	dashboard := v1.Group("/dashboard", customMiddleware.Authentication)
	dashboard.GET("", handler.GetDashboardHandler, customMiddleware.GrantAdminIKU5AndRektor)
	dashboard.GET("/fakultas/:id", handler.GetDashboardByFakultasHandler, customMiddleware.GrantAdminIKU5AndRektor)
	dashboard.GET("/total", handler.GetDashboardTotalHandler, customMiddleware.GrantAdminIKU5AndRektor)
	dashboard.GET("/umum", handler.GetDashboardUmum, customMiddleware.GrantAdminUmum)
	dashboard.GET("/dosen", handler.GetDashboardDosen, customMiddleware.GrantDosen)
	dashboard.PATCH("/target", handler.InsertTargetHandler, customMiddleware.GrantAdminIKU5)

	return app
}
