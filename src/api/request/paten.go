package request

import (
	"be-5/src/model"
	"be-5/src/util"
	"errors"
	"time"
)

type (
	Paten struct {
		IdKategori        int            `json:"id_kategori"`
		IdJenisPenelitian int            `json:"id_jenis"`
		Judul             string         `json:"judul"`
		JumlahHalaman     int            `json:"jumlah_halaman"`
		Tanggal           string         `json:"tanggal"`
		Penyelenggara     string         `json:"penyelenggara"`
		Penerbit          string         `json:"penerbit"`
		Isbn              string         `json:"isbn"`
		TautanEksternal   string         `json:"tautan_eksternal"`
		Keterangan        string         `json:"keterangan"`
		Dokumen           []DokumenPaten `json:"dokumen"`
		Penulis           []PenulisPaten `json:"penulis"`
	}

	DokumenPaten struct {
		IdJenisDokumen int    `json:"id_jenis_dokumen"`
		Nama           string `json:"nama"`
		Keterangan     string `json:"keterangan"`
	}

	DokumenPatenPayload struct {
		IdFile    string
		IdPaten   int
		NamaFile  string
		JenisFile string
		Url       string
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

func (r *DokumenPaten) MapRequest(p *DokumenPatenPayload) *model.DokumenPaten {
	return &model.DokumenPaten{
		ID:             p.IdFile,
		IdPaten:        p.IdPaten,
		IdJenisDokumen: r.IdJenisDokumen,
		Nama:           r.Nama,
		NamaFile:       p.NamaFile,
		JenisFile:      p.JenisFile,
		Url:            p.Url,
		Keterangan:     r.Keterangan,
		TanggalUpload:  time.Now(),
	}
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
