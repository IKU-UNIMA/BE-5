package request

import "be-5/src/model"

type Prodi struct {
	IdFakultas int    `json:"id_fakultas"`
	KodeProdi  int    `json:"kode_prodi"`
	Nama       string `json:"nama"`
	Jenjang    string `json:"jenjang"`
}

func (r *Prodi) MapRequest() *model.Prodi {
	return &model.Prodi{
		IdFakultas: r.IdFakultas,
		KodeProdi:  r.KodeProdi,
		Nama:       r.Nama,
		Jenjang:    r.Jenjang,
	}
}
