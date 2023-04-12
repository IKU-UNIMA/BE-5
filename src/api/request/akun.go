package request

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ChangePassword struct {
	PasswordBaru string `json:"password_baru" validate:"required"`
}
