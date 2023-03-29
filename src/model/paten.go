package model

import "time"

type (
	Paten struct {
		ID                int `gorm:"primaryKey"`
		IdDosen           int
		IdKategori        int
		IdJenisPenelitian int
		IdKategoriCapaian int             `gorm:"default:null"`
		Judul             string          `gorm:"type:varchar(255)"`
		Tanggal           time.Time       `gorm:"type:date"`
		JumlahHalaman     int             `gorm:"type:smallint"`
		Penyelenggara     string          `gorm:"type:varchar(255)"`
		Penerbit          string          `gorm:"type:varchar(255)"`
		Isbn              string          `gorm:"type:varchar(255)"`
		TautanEksternal   string          `gorm:"type:text"`
		Keterangan        string          `gorm:"type:text"`
		Dosen             Dosen           `gorm:"foreignKey:IdDosen;constraint:OnDelete:CASCADE"`
		Kategori          KategoriPaten   `gorm:"foreignKey:IdKategori"`
		JenisPenelitian   JenisPenelitian `gorm:"foreignKey:IdJenisPenelitian"`
		KategoriCapaian   KategoriCapaian `gorm:"foreignKey:IdKategoriCapaian"`
		Penulis           []PenulisPaten  `gorm:"foreignKey:IdPaten;constraint:OnDelete:CASCADE"`
		Dokumen           []DokumenPaten  `gorm:"foreignKey:IdPaten;constraint:OnDelete:CASCADE"`
	}

	KategoriPaten struct {
		ID                   int `gorm:"primaryKey"`
		IdJenisKategoriPaten int
		Nama                 string `gorm:"type:text"`
	}

	JenisKategoriPaten struct {
		ID            int             `gorm:"primaryKey"`
		Nama          string          `gorm:"type:text"`
		KategoriPaten []KategoriPaten `gorm:"foreignKey:IdJenisKategoriPaten;constraint:OnDelete:CASCADE"`
	}

	PenulisPaten struct {
		ID           int `gorm:"primaryKey"`
		IdPaten      int
		Nama         string `gorm:"type:text"`
		JenisPenulis string `gorm:"type:enum('dosen', 'mahasiswa', 'lain')"`
		Urutan       int    `gorm:"type:tinyint unsigned"`
		Afiliasi     string `gorm:"type:varchar(255)"`
		Peran        string `gorm:"type:varchar(120)"`
		IsAuthor     bool
	}

	DokumenPaten struct {
		ID             string `gorm:"primaryKey"`
		IdPaten        int
		IdJenisDokumen int
		Nama           string       `gorm:"type:varchar(255)"`
		NamaFile       string       `gorm:"type:varchar(255)"`
		JenisFile      string       `gorm:"type:varchar(255)"`
		Keterangan     string       `gorm:"type:text"`
		Url            string       `gorm:"type:text"`
		TanggalUpload  time.Time    `gorm:"type:date"`
		JenisDokumen   JenisDokumen `gorm:"foreignKey:IdJenisDokumen"`
	}
)
