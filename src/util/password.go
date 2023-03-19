package util

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword() string {
	pass := make([]byte, 8)
	rand.Seed(time.Now().UnixNano())
	rand.Read(pass)

	return fmt.Sprintf("%x", pass)
}

func HashPassword(pass string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hashed)
}

func ValidateHash(pass, hashed string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pass))
	return err == nil
}
