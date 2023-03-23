package response

type Admin struct {
	ID     int    `json:"id"`
	Nama   string `json:"nama"`
	Nidn   string `json:"nidn"`
	Nip    string `json:"nip"`
	Email  string `json:"email"`
	Bagian string `json:"bagian"`
}
