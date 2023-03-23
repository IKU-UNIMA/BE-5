package request

import "be-5/src/model"

type Rektor struct {
	Nama  string `json:"nama"`
	Nip   string `json:"nip"`
	Email string `json:"email"`
}

func (r *Rektor) MapRequest() *model.Rektor {
	return &model.Rektor{
		Nama: r.Nama,
		Nip:  r.Nip,
	}
}
