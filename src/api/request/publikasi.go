package request

import (
	"be-5/src/model"
	"be-5/src/util"
	"errors"
	"time"
)

type Publikasi struct {
	IdKategori          int       `json:"id_kategori" validate:"required"`
	IdJenisPenelitian   int       `json:"id_jenis" validate:"required"`
	IdKategoriCapaian   int       `json:"id_kategori_capaian"`
	Judul               string    `json:"judul" validate:"required"`
	JudulAsli           string    `json:"judul_asli"`
	JudulChapter        string    `json:"judul_chapter"`
	NamaJurnal          string    `json:"nama_jurnal"`
	NamaKoranMajalah    string    `json:"nama_koran_majalah"`
	NamaSeminar         string    `json:"nama_seminar"`
	TautanLamanJurnal   string    `json:"tautan_laman_jurnal"`
	TanggalTerbit       string    `json:"tanggal_terbit"`
	WaktuPelaksanaan    string    `json:"waktu_pelaksanaan"`
	Volume              string    `json:"volume"`
	Edisi               string    `json:"edisi"`
	Nomor               string    `json:"nomor"`
	Halaman             string    `json:"halaman"`
	JumlahHalaman       int       `json:"jumlah_halaman"`
	Penerbit            string    `json:"penerbit"`
	Penyelenggara       string    `json:"penyelenggara"`
	KotaPenyelenggaraan string    `json:"kota_penyelenggaraan"`
	IsSeminar           bool      `json:"is_seminar"`
	IsProsiding         bool      `json:"is_prosiding"`
	Bahasa              string    `json:"bahasa"`
	Doi                 string    `json:"doi"`
	Isbn                string    `json:"isbn"`
	Issn                string    `json:"issn"`
	EIssn               string    `json:"e_issn"`
	Tautan              string    `json:"tautan"`
	Keterangan          string    `json:"keterangan"`
	Dokumen             []Dokumen `json:"dokumen"`
	PenulisDosen        []Penulis `json:"penulis_dosen"`
	PenulisMahasiswa    []Penulis `json:"penulis_mahasiswa"`
	PenulisLain         []Penulis `json:"penulis_lain"`
}

func (r *Publikasi) MapRequest() (*model.Publikasi, error) {
	var (
		tanggalTerbit, waktuPelaksanaan *time.Time
	)

	if r.TanggalTerbit != "" {
		tanggal, err := util.ConvertStringToDate(r.TanggalTerbit)
		if err != nil {
			return nil, errors.New("format tanggal terbit salah")
		}

		if !tanggal.IsZero() {
			tanggalTerbit = &tanggal
		}
	}

	if r.WaktuPelaksanaan != "" {
		tanggal, err := util.ConvertStringToDate(r.WaktuPelaksanaan)
		if err != nil {
			return nil, errors.New("format waktu pelaksanaan salah")
		}

		if !tanggal.IsZero() {
			tanggalTerbit = &tanggal
		}
	}

	return &model.Publikasi{
		IdKategori:          r.IdKategori,
		IdJenisPenelitian:   r.IdJenisPenelitian,
		IdKategoriCapaian:   r.IdKategoriCapaian,
		Judul:               r.Judul,
		JudulAsli:           r.JudulAsli,
		JudulChapter:        r.JudulChapter,
		NamaJurnal:          r.NamaJurnal,
		NamaKoranMajalah:    r.NamaKoranMajalah,
		NamaSeminar:         r.NamaSeminar,
		TautanLamanJurnal:   r.TautanLamanJurnal,
		TanggalTerbit:       tanggalTerbit,
		WaktuPelaksanaan:    waktuPelaksanaan,
		Volume:              r.Volume,
		Edisi:               r.Edisi,
		Nomor:               r.Nomor,
		Halaman:             r.Halaman,
		JumlahHalaman:       r.JumlahHalaman,
		Penerbit:            r.Penerbit,
		Penyelenggara:       r.Penyelenggara,
		KotaPenyelenggaraan: r.KotaPenyelenggaraan,
		IsSeminar:           r.IsSeminar,
		IsProsiding:         r.IsProsiding,
		Bahasa:              r.Bahasa,
		Doi:                 r.Doi,
		Isbn:                r.Isbn,
		Issn:                r.Issn,
		EIssn:               r.EIssn,
		Tautan:              r.Tautan,
		Keterangan:          r.Keterangan,
		Status:              util.BELUM_DIVERIFIKASI,
	}, nil
}
