package request

import "be-5/src/model"

type Dosen struct {
	Nama              string `json:"nama" validate:"required"`
	Nidn              string `json:"nidn" validate:"required"`
	Nip               string `json:"nip" validate:"required"`
	Email             string `json:"email" validate:"required,email"`
	IdProdi           int    `json:"id_prodi" validate:"required"`
	StatusKepegawaian string `json:"status_kepegawaian" validate:"required"`
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
