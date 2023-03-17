package util

import "golang.org/x/crypto/bcrypt"

func HashPassword(pass string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hashed)
}

func ValidateHash(pass, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
	return err == nil
}
