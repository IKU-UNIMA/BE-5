package response

type Login struct {
	Nama  string `json:"nama"`
	Token string `json:"token"`
}
