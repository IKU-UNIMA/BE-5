package request

import "be-5/src/model"

type Target struct {
	Target float32 `json:"target" validate:"required"`
	Tahun  int     `json:"tahun" validate:"required"`
}

func (r *Target) MapRequest() *model.Target {
	return &model.Target{
		Bagian: "IKU 5",
		Target: r.Target,
		Tahun:  r.Tahun,
	}
}
