package request

import "be-5/src/model"

type Dosen struct {
	Nama              string `json:"nama"`
	Nidn              string `json:"nidn"`
	Nip               string `json:"nip"`
	Email             string `json:"email"`
	IdProdi           int    `json:"id_prodi"`
	StatusKepegawaian string `json:"status_kepegawaian"`
}

func (r *Dosen) MapRequest() *model.Dosen {
	return &model.Dosen{
		Nama:              r.Nama,
		Nidn:              r.Nidn,
		Nip:               r.Nip,
		IdProdi:           r.IdProdi,
		StatusKepegawaian: r.StatusKepegawaian,
	}
}
