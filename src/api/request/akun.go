package request

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePassword struct {
	PasswordLama string `json:"password_lama"`
	PasswordBaru string `json:"password_baru"`
}
