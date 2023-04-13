package request

import (
	"be-5/src/model"
	"be-5/src/util"
	"errors"
)

type Paten struct {
	IdKategori        int       `json:"id_kategori" validate:"required"`
	IdJenisPenelitian int       `json:"id_jenis" validate:"required"`
	IdKategoriCapaian int       `json:"id_kategori_capaian"`
	Judul             string    `json:"judul" validate:"required"`
	JumlahHalaman     int       `json:"jumlah_halaman"`
	Tanggal           string    `json:"tanggal" validate:"required"`
	Penyelenggara     string    `json:"penyelenggara"`
	Penerbit          string    `json:"penerbit"`
	Isbn              string    `json:"isbn"`
	TautanEksternal   string    `json:"tautan_eksternal"`
	Keterangan        string    `json:"keterangan"`
	Dokumen           []Dokumen `json:"dokumen"`
	PenulisDosen      []Penulis `json:"penulis_dosen"`
	PenulisMahasiswa  []Penulis `json:"penulis_mahasiswa"`
	PenulisLain       []Penulis `json:"penulis_lain"`
}

func (r *Paten) MapRequest() (*model.Paten, error) {
	tanggal, err := util.ConvertStringToDate(r.Tanggal)
	if err != nil {
		return nil, errors.New("format tanggal salah")
	}

	return &model.Paten{
		IdKategori:        r.IdKategori,
		IdJenisPenelitian: r.IdJenisPenelitian,
		IdKategoriCapaian: r.IdKategoriCapaian,
		Judul:             r.Judul,
		Tanggal:           tanggal,
		JumlahHalaman:     r.JumlahHalaman,
		Penyelenggara:     r.Penyelenggara,
		Penerbit:          r.Penerbit,
		Isbn:              r.Isbn,
		TautanEksternal:   r.TautanEksternal,
		Keterangan:        r.Keterangan,
	}, nil
}
