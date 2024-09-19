package user

import (
	"database/sql"
	"errors"
)

var (
	ErrInvalidUserPassword = errors.New("Invalid user password")
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

const getUserByEmail = "SELECT email, encrypted_password FROM users WHERE email = $1"

func (s *Service) AuthUser(email, password string) (*OauthUser, error) {
	var user OauthUser
	err := s.db.QueryRow(
		getUserByEmail,
		email,
	).Scan(
		&user.Email,
		&user.EncryptedPassword,
	)
	if err != nil {
		return nil, err
	}

	if VerifyPassword(user.EncryptedPassword, password) != nil {
		return nil, ErrInvalidUserPassword
	}

	return &user, nil
}
