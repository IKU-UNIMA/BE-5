package request

import "be-5/src/model"

type Prodi struct {
	IdFakultas int    `json:"id_fakultas" validate:"required"`
	KodeProdi  int    `json:"kode_prodi" validate:"required"`
	Nama       string `json:"nama" validate:"required"`
	Jenjang    string `json:"jenjang" validate:"required"`
}

func (r *Prodi) MapRequest() *model.Prodi {
	return &model.Prodi{
		IdFakultas: r.IdFakultas,
		KodeProdi:  r.KodeProdi,
		Nama:       r.Nama,
		Jenjang:    r.Jenjang,
	}
}
