package request

import "be-5/src/model"

type Rektor struct {
	Nama  string `json:"nama" validate:"required"`
	Nip   string `json:"nip" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

func (r *Rektor) MapRequest() *model.Rektor {
	return &model.Rektor{
		Nama: r.Nama,
		Nip:  r.Nip,
	}
}
