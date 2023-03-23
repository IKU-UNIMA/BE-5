package response

type Admin struct {
	ID     int    `json:"id"`
	Nama   string `json:"nama"`
	Nip    string `json:"nip"`
	Email  string `json:"email"`
	Bagian string `json:"bagian"`
}
