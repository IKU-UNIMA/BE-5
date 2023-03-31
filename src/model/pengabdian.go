package model

import (
	"be-5/src/api/response"
	"fmt"
	"time"
)

type (
	Pengabdian struct {
		ID                      int `gorm:"primaryKey"`
		IdDosen                 int
		IdKategori              int
		Judul                   string `gorm:"type:varchar(255)"`
		Afiliasi                string `gorm:"type:varchar(255)"`
		KelompokBidang          string `gorm:"type:varchar(255)"`
		JenisSkim               string `gorm:"type:varchar(255)"`
		LokasiKegiatan          string `gorm:"type:varchar(255)"`
		TahunUsulan             uint   `gorm:"type:year"`
		TahunKegiatan           uint   `gorm:"type:year"`
		TahunPelaksanaan        uint   `gorm:"type:year"`
		LamaKegiatan            uint
		TahunPelaksanaanKe      uint
		DanaDariDikti           float64
		DanaDariPerguruanTinggi float64
		DanaDariInstitusiLain   float64
		InKind                  string
		NoSkPenugasan           string              `gorm:"type:varchar(255)"`
		TglSkPenugasan          time.Time           `gorm:"type:date"`
		MitraLitabmas           string              `gorm:"type:varchar(255)"`
		Dosen                   Dosen               `gorm:"foreignKey:IdDosen;OnDelete:CASCADES"`
		Penulis                 []PenulisPengabdian `gorm:"foreignKey:IdPengabdian;constraint:OnDelete:CASCADE"`
		Dokumen                 []DokumenPengabdian `gorm:"foreignKey:IdPengabdian;constraint:OnDelete:CASCADE"`
	}

	KategoriPengabdian struct {
		ID                        int `gorm:"primaryKey"`
		IdJenisKategoriPengabdian int
		Nama                      string `gorm:"type:text"`
	}

	JenisKategoriPengabdian struct {
		ID                 int                  `gorm:"primaryKey"`
		Nama               string               `gorm:"type:text"`
		KategoriPengabdian []KategoriPengabdian `gorm:"foreignKey:IdJenisKategoriPengabdian;constraint:OnDelete:CASCADE"`
	}

	PenulisPengabdian struct {
		ID           int `gorm:"primaryKey"`
		IdPengabdian int
		Nama         string `gorm:"type:text"`
		JenisPenulis string `gorm:"type:enum('dosen', 'mahasiswa', 'lain')"`
		Peran        string `gorm:"type:varchar(120)"`
		IsActive     bool
	}

	DokumenPengabdian struct {
		ID             string `gorm:"primaryKey"`
		IdPengabdian   int
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

func (p *Pengabdian) MapToResponse() *response.Pengabdian {
	tahunPelaksanaan := fmt.Sprintf("%d/%d", p.TahunPelaksanaan, p.TahunPelaksanaan+p.LamaKegiatan)
	return &response.Pengabdian{
		ID:               p.ID,
		TahunPelaksanaan: tahunPelaksanaan,
		LamaKegiatan:     p.LamaKegiatan,
	}
}

func MapBatchPengabdianResponse(p []Pengabdian) []*response.Pengabdian {
	res := []*response.Pengabdian{}
	for i := 0; i < len(p)/2; i++ {
		res = append(res, p[i].MapToResponse())
		res = append(res, p[len(p)-1-i].MapToResponse())
	}

	return res
}
