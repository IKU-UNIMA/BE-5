package request

import (
	"be-5/src/model"
	"be-5/src/util"
	"errors"
)

type (
	Paten struct {
		IdKategori        int            `json:"id_kategori" validate:"required"`
		IdJenisPenelitian int            `json:"id_jenis" validate:"required"`
		Judul             string         `json:"judul" validate:"required"`
		JumlahHalaman     int            `json:"jumlah_halaman"`
		Tanggal           string         `json:"tanggal" validate:"required"`
		Penyelenggara     string         `json:"penyelenggara"`
		Penerbit          string         `json:"penerbit"`
		Isbn              string         `json:"isbn"`
		TautanEksternal   string         `json:"tautan_eksternal"`
		Keterangan        string         `json:"keterangan"`
		Dokumen           []Dokumen      `json:"dokumen"`
		Penulis           []PenulisPaten `json:"penulis"`
	}

	PenulisPaten struct {
		Nama         string `json:"nama"`
		JenisPenulis string `json:"jenis_penulis"`
		Urutan       int    `json:"urutan"`
		Afiliasi     string `json:"afiliasi"`
		Peran        string `json:"peran"`
		IsAuthor     bool   `json:"is_author"`
	}
)

func (r *Paten) MapRequest() (*model.Paten, error) {
	tanggal, err := util.ConvertStringToDate(r.Tanggal)
	if err != nil {
		return nil, errors.New("format tanggal salah")
	}

	return &model.Paten{
		IdKategori:        r.IdKategori,
		IdJenisPenelitian: r.IdJenisPenelitian,
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

func (r *PenulisPaten) MapRequest(idPaten int) *model.PenulisPaten {
	return &model.PenulisPaten{
		IdPaten:      idPaten,
		Nama:         r.Nama,
		JenisPenulis: r.JenisPenulis,
		Urutan:       r.Urutan,
		Afiliasi:     r.Afiliasi,
		Peran:        r.Peran,
		IsAuthor:     r.IsAuthor,
	}
}
