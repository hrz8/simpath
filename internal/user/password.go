package user

import (
	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(passwordHash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err
}
