package database

import (
	"be-5/src/config/env"
	"be-5/src/model"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitMySQL() {
	var err error
	DB, err = gorm.Open(mysql.Open(env.GetMySQLEnv()), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
}

func MigrateMySQL() {
	DB.AutoMigrate(
		&model.Fakultas{},
		&model.Prodi{},
		&model.JenisDokumen{},
		&model.JenisPenelitian{},
		&model.KategoriCapaian{},
		&model.Akun{},
		&model.Admin{},
		&model.Dosen{},
		&model.Rektor{},
		&model.Paten{},
		&model.JenisKategoriPaten{},
		&model.KategoriPaten{},
		&model.DokumenPaten{},
		&model.PenulisPaten{},
		&model.Pengabdian{},
		&model.JenisKategoriPengabdian{},
		&model.KategoriPengabdian{},
		&model.DokumenPengabdian{},
		&model.AnggotaPengabdian{},
		&model.Publikasi{},
		&model.JenisKategoriPublikasi{},
		&model.KategoriPublikasi{},
		&model.DokumenPublikasi{},
		&model.PenulisPublikasi{},
		&model.Target{},
	)
}
