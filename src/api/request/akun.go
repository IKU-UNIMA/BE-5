package request

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePassword struct {
	PasswordBaru string `json:"password_baru"`
}
