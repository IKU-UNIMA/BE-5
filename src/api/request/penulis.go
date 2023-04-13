package request

import "be-5/src/model"

type Penulis struct {
	Nama     string `json:"nama"`
	Urutan   int    `json:"urutan"`
	Afiliasi string `json:"afiliasi"`
	Peran    string `json:"peran"`
	IsAuthor bool   `json:"is_author"`
}

func (r *Penulis) MapRequestToPaten(idPaten int, jenisPenulis string) *model.PenulisPaten {
	return &model.PenulisPaten{
		IdPaten:      idPaten,
		Nama:         r.Nama,
		JenisPenulis: jenisPenulis,
		Urutan:       r.Urutan,
		Afiliasi:     r.Afiliasi,
		Peran:        r.Peran,
		IsAuthor:     r.IsAuthor,
	}
}
