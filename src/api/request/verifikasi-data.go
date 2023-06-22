package request

type VerifikasiData struct {
	Status string `json:"status" validate:"required"`
}
