package response

type Rektor struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Nip   string `json:"nip"`
	Email string `json:"email"`
}
