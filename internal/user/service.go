package user

import (
	"database/sql"
	"errors"

	"github.com/hrz8/simpath/password"
)

var (
	ErrUserNotFound          = errors.New("User not found")
	ErrInvalidUserOrPassword = errors.New("Invalid user or password")
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

const findUserByID = `
SELECT
	id,
	email,
	encrypted_password,
	role_id,
	public_id
FROM users
WHERE id = $1
`

func (s *Service) FindUserByID(userID uint32) (*OauthUser, error) {
	u := new(OauthUser)
	err := s.db.QueryRow(
		findUserByID,
		userID,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
		&u.RoleID,
		&u.PublicID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}

const findUserByEmail = `
SELECT
	id,
	email,
	encrypted_password,
	role_id,
	public_id
FROM users WHERE email = $1`

func (s *Service) FindUserByEmail(email string) (*OauthUser, error) {
	u := new(OauthUser)
	err := s.db.QueryRow(
		findUserByEmail,
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
		&u.RoleID,
		&u.PublicID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}

func (s *Service) AuthUser(email, pwd string) (*OauthUser, error) {
	user, err := s.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if password.VerifyPassword(user.EncryptedPassword, pwd) != nil {
		return nil, ErrInvalidUserOrPassword
	}

	return user, nil
}
