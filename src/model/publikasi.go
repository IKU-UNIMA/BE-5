package model

import "time"

type (
	Publikasi struct {
		ID                  int `gorm:"primaryKey"`
		IdDosen             int
		IdKategori          int
		IdJenisPenelitian   int
		IdKategoriCapaian   int        `gorm:"default:null"`
		Judul               string     `gorm:"type:varchar(255)"`
		JudulAsli           string     `gorm:"type:varchar(255)"`
		JudulChapter        string     `gorm:"type:varchar(255)"`
		NamaJurnal          string     `gorm:"type:varchar(255)"`
		NamaKoranMajalah    string     `gorm:"type:varchar(255)"`
		NamaSeminar         string     `gorm:"type:varchar(255)"`
		TautanLamanJurnal   string     `gorm:"type:varchar(255)"`
		TanggalTerbit       *time.Time `gorm:"type:date"`
		WaktuPelaksanaan    *time.Time `gorm:"type:date"`
		Volume              string     `gorm:"type:varchar(255)"`
		Edisi               string     `gorm:"type:varchar(255)"`
		Nomor               string     `gorm:"type:varchar(255)"`
		Halaman             string     `gorm:"type:varchar(255)"`
		JumlahHalaman       int        `gorm:"type:smallint"`
		Penerbit            string     `gorm:"type:varchar(255)"`
		Penyelenggara       string     `gorm:"type:varchar(255)"`
		KotaPenyelenggaraan string     `gorm:"type:varchar(255)"`
		IsSeminar           bool
		IsProsiding         bool
		Bahasa              string `gorm:"type:varchar(255)"`
		Doi                 string `gorm:"type:varchar(255)"`
		Isbn                string `gorm:"type:varchar(255)"`
		Issn                string `gorm:"type:varchar(255)"`
		EIssn               string `gorm:"type:varchar(255)"`
		Tautan              string `gorm:"type:varchar(255)"`
		Keterangan          string `gorm:"type:text"`
		Status              string `gorm:"type:enum('Belum Diverifikasi','Draft','Tidak Terverifikasi','Terverifikasi')"`
		CreatedAt           time.Time
		Kategori            KategoriPublikasi  `gorm:"foreignKey:IdKategori"`
		Dosen               Dosen              `gorm:"foreignKey:IdDosen;OnDelete:CASCADES"`
		JenisPenelitian     JenisPenelitian    `gorm:"foreignKey:IdJenisPenelitian"`
		KategoriCapaian     KategoriCapaian    `gorm:"foreignKey:IdKategoriCapaian"`
		Dokumen             []DokumenPublikasi `gorm:"foreignKey:IdPublikasi;constraint:OnDelete:CASCADE"`
		Penulis             []PenulisPublikasi `gorm:"foreignKey:IdPublikasi;constraint:OnDelete:CASCADE"`
	}

	KategoriPublikasi struct {
		ID                       int `gorm:"primaryKey"`
		IdJenisKategoriPublikasi int
		Nama                     string `gorm:"type:text"`
	}

	JenisKategoriPublikasi struct {
		ID                int                 `gorm:"primaryKey"`
		Nama              string              `gorm:"type:text"`
		KategoriPublikasi []KategoriPublikasi `gorm:"foreignKey:IdJenisKategoriPublikasi;constraint:OnDelete:CASCADE"`
	}

	PenulisPublikasi struct {
		ID           int `gorm:"primaryKey"`
		IdPublikasi  int
		Nama         string `gorm:"type:text"`
		JenisPenulis string `gorm:"type:enum('dosen', 'mahasiswa', 'lain')"`
		Urutan       int    `gorm:"type:tinyint unsigned"`
		Afiliasi     string `gorm:"type:varchar(255)"`
		Peran        string `gorm:"type:enum('Penulis', 'Editor', 'Penerjemah', 'Penemu/Inventor')"`
		IsAuthor     bool
	}

	DokumenPublikasi struct {
		ID             string `gorm:"primaryKey"`
		IdPublikasi    int
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
