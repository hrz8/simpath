package user

import (
	"database/sql"
	"errors"
)

var (
	ErrUserNotFound        = errors.New("User not found")
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

const findUserByEmail = "SELECT id, email, encrypted_password FROM users WHERE email = $1"

func (s *Service) FindUserByEmail(email string) (*OauthUser, error) {
	var u OauthUser
	err := s.db.QueryRow(
		findUserByEmail,
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (s *Service) AuthUser(email, password string) (*OauthUser, error) {
	user, err := s.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if VerifyPassword(user.EncryptedPassword, password) != nil {
		return nil, ErrInvalidUserPassword
	}

	return user, nil
}
