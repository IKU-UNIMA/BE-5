package request

import (
	"be-5/src/model"
	"be-5/src/util"
	"errors"
)

type (
	Pengabdian struct {
		IdKategori              int                 `json:"id_kategori" validate:"required"`
		Judul                   string              `json:"judul" validate:"required"`
		Afiliasi                string              `json:"afiliasi" validate:"required"`
		KelompokBidang          string              `json:"kelompok_bidang"`
		JenisSkim               string              `json:"jenis_skim"`
		LokasiKegiatan          string              `json:"lokasi_kegiatan"`
		TahunUsulan             uint                `json:"tahun_usulan" validate:"required"`
		TahunKegiatan           uint                `json:"tahun_kegiatan" validate:"required"`
		TahunPelaksanaan        uint                `json:"tahun_pelaksanaan" validate:"required"`
		LamaKegiatan            uint                `json:"lama_kegiatan" validate:"required"`
		TahunPelaksanaanKe      uint                `json:"tahun_pelaksanaan_ke" validate:"required"`
		DanaDariDikti           float64             `json:"dana_dari_dikti"`
		DanaDariPerguruanTinggi float64             `json:"dana_dari_perguruan_tinggi"`
		DanaDariInstitusiLain   float64             `json:"dana_dari_institusi_lain"`
		InKind                  string              `json:"in_kind"`
		NoSkPenugasan           string              `json:"no_sk_penugasan" validate:"required"`
		TglSkPenugasan          string              `json:"tgl_sk_penugasan" validate:"required"`
		MitraLitabmas           string              `json:"mitra_litabmas"`
		Dokumen                 []Dokumen           `json:"dokumen"`
		AnggotaDosen            []AnggotaPengabdian `json:"anggota_dosen"`
		AnggotaMahasiswa        []AnggotaPengabdian `json:"anggota_mahasiswa"`
		AnggotaEksternal        []AnggotaPengabdian `json:"anggota_eksternal"`
	}

	AnggotaPengabdian struct {
		Nama     string `json:"nama"`
		Peran    string `json:"peran"`
		IsActive bool   `json:"is_active"`
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
		Status:                  util.BELUM_DIVERIFIKASI,
	}, nil
}

func (r *AnggotaPengabdian) MapRequest(jenisAnggota string) *model.AnggotaPengabdian {
	return &model.AnggotaPengabdian{
		Nama:         r.Nama,
		JenisAnggota: jenisAnggota,
		Peran:        r.Peran,
		IsActive:     r.IsActive,
	}
}
