package request

import (
	"be-5/src/model"
	"be-5/src/util"
	"errors"
	"time"
)

type (
	Pengabdian struct {
		IdKategori              int                 `json:"id_kategori" validate:"number"`
		Judul                   string              `json:"judul"`
		Afiliasi                string              `json:"afiliasi"`
		KelompokBidang          string              `json:"kelompok_bidang"`
		JenisSkim               string              `json:"jenis_skim"`
		LokasiKegiatan          string              `json:"lokasi_kegiatan"`
		TahunUsulan             uint                `json:"tahun_usulan"`
		TahunKegiatan           uint                `json:"tahun_kegiatan"`
		TahunPelaksanaan        uint                `json:"tahun_pelaksanaan"`
		LamaKegiatan            uint                `json:"lama_kegiatan"`
		TahunPelaksanaanKe      uint                `json:"tahun_pelaksanaan_ke"`
		DanaDariDikti           float64             `json:"dana_dari_dikti"`
		DanaDariPerguruanTinggi float64             `json:"dana_dari_perguruan_tinggi"`
		DanaDariInstitusiLain   float64             `json:"dana_dari_institusi_lain"`
		InKind                  string              `json:"in_kind"`
		NoSkPenugasan           string              `json:"no_sk_penugasan"`
		TglSkPenugasan          string              `json:"tgl_sk_penugasan"`
		MitraLitabmas           string              `json:"mitra_litabmas"`
		Dokumen                 []DokumenPengabdian `json:"dokumen"`
		Anggota                 []AnggotaPengabdian `json:"anggota"`
	}

	DokumenPengabdian struct {
		IdJenisDokumen int    `json:"id_jenis_dokumen"`
		Nama           string `json:"nama"`
		Keterangan     string `json:"keterangan"`
	}

	DokumenPengabdianPayload struct {
		IdFile       string
		IdPengabdian int
		NamaFile     string
		JenisFile    string
		Url          string
	}

	AnggotaPengabdian struct {
		Nama         string `json:"nama"`
		JenisAnggota string `json:"jenis_anggota"`
		Peran        string `json:"peran"`
		IsActive     bool   `json:"is_active"`
	}
)

func (r *Pengabdian) MapRequest() (*model.Pengabdian, error) {
	tglSkPenugasan, err := util.ConvertStringToDate(r.TglSkPenugasan)
	if err != nil {
		return nil, errors.New("format tanggal salah")
	}

	return &model.Pengabdian{
		IdKategori:              r.IdKategori,
		Judul:                   r.Judul,
		Afiliasi:                r.Afiliasi,
		KelompokBidang:          r.KelompokBidang,
		JenisSkim:               r.JenisSkim,
		LokasiKegiatan:          r.LokasiKegiatan,
		TahunUsulan:             r.TahunUsulan,
		TahunKegiatan:           r.TahunKegiatan,
		TahunPelaksanaan:        r.TahunPelaksanaan,
		LamaKegiatan:            r.LamaKegiatan,
		TahunPelaksanaanKe:      r.TahunPelaksanaanKe,
		DanaDariDikti:           r.DanaDariDikti,
		DanaDariPerguruanTinggi: r.DanaDariPerguruanTinggi,
		DanaDariInstitusiLain:   r.DanaDariInstitusiLain,
		InKind:                  r.InKind,
		NoSkPenugasan:           r.NoSkPenugasan,
		TglSkPenugasan:          tglSkPenugasan,
		MitraLitabmas:           r.MitraLitabmas,
	}, nil
}

func (r *DokumenPengabdian) MapRequest(p *DokumenPengabdianPayload) *model.DokumenPengabdian {
	return &model.DokumenPengabdian{
		ID:             p.IdFile,
		IdPengabdian:   p.IdPengabdian,
		IdJenisDokumen: r.IdJenisDokumen,
		Nama:           r.Nama,
		NamaFile:       p.NamaFile,
		JenisFile:      p.JenisFile,
		Url:            p.Url,
		Keterangan:     r.Keterangan,
		TanggalUpload:  time.Now(),
	}
}

func (r *AnggotaPengabdian) MapRequest(idPengabdian int) *model.AnggotaPengabdian {
	return &model.AnggotaPengabdian{
		IdPengabdian: idPengabdian,
		Nama:         r.Nama,
		JenisAnggota: r.JenisAnggota,
		Peran:        r.Peran,
		IsActive:     r.IsActive,
	}
}