package user

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hrz8/simpath/internal/role"
	"github.com/hrz8/simpath/password"
)

var (
	MinPasswordLength = 8

	// errors
	ErrUserNotFound          = errors.New("User not found")
	ErrInvalidUserOrPassword = errors.New("Invalid user or password")
	ErrEmailTaken            = errors.New("Email taken")
	ErrPasswordTooShort      = fmt.Errorf(
		"Password must be at least %d characters long",
		MinPasswordLength,
	)
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

func (s *Service) IsUserExists(email string) bool {
	_, err := s.FindUserByEmail(email)
	return err == nil
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

const createNewUser = `
	INSERT INTO users (
		email,
		encrypted_password,
		role_id,
		created_at,
		updated_at,
		public_id
	) VALUES (
		$1, $2, $3, $4, $5, $6
	) RETURNING *
`

func (s *Service) CreateUser(roleID role.RoleType, email, pass string) (*OauthUser, error) {
	// TODO: should use transaction to reduce the race condition possibility when check IsUserExist
	u := new(OauthUser)
	u.PublicID = uuid.NewString()

	if pass != "" {
		if len(pass) < MinPasswordLength {
			return nil, ErrPasswordTooShort
		}
		passHash, err := password.HashPassword(pass)
		if err != nil {
			return nil, err
		}
		u.EncryptedPassword = string(passHash)
	}

	if s.IsUserExists(email) {
		return nil, ErrEmailTaken
	}

	err := s.db.QueryRow(
		createNewUser,
		email,
		u.EncryptedPassword,
		roleID,
		time.Now(),
		time.Now(),
		u.PublicID,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
		&u.RoleID,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.PublicID,
	)
	if err != nil {
		return nil, err
	}

	return u, nil
}
