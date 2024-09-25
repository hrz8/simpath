package password

import (
	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(passwordHash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	return err
}

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 3)
}
