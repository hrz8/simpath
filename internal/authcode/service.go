package authcode

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAuthorizationCodeInvalid = errors.New("Invalid authorization code")
	ErrAuthorizationCodeExpired = errors.New("Authorization code expired")
	ErrInvalidRedirectURI       = errors.New("Invalid redirect URI")
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

const createNewAuthorizationCode = `
	INSERT INTO authorization_codes (
		client_id,
		user_id,
		code,
		redirect_uri,
		scope,
		expires_at,
		created_at,
		updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8
	) RETURNING *
`

func (s *Service) GrantAuthorizationCode(clientID uint32, userID uint32, redirectURI, scope string, expiresIn int) (*OauthAuthorizationCode, error) {
	var authCode OauthAuthorizationCode
	codeVal := generateAuthorizationCode()
	if err := s.db.QueryRow(
		createNewAuthorizationCode,
		clientID,
		userID,
		codeVal,
		redirectURI,
		scope,
		time.Now().UTC().Add(time.Duration(expiresIn)*time.Second),
		time.Now(),
		time.Now(),
	).Scan(
		&authCode.ID,
		&authCode.ClientID,
		&authCode.UserID,
		&authCode.Code,
		&authCode.RedirectURI,
		&authCode.Scope,
		&authCode.ExpiresAt,
		&authCode.CreatedAt,
		&authCode.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return &authCode, nil
}

const getValidAuthCode = `
	SELECT
		id,
		client_id,
		user_id,
		code,
		redirect_uri,
		scope,
		expires_at
	FROM authorization_codes
	WHERE
		client_id = $1 AND
		code = $2
`

func (s *Service) GetValidAuthCode(clientID uint32, code string, redirectURI string) (*OauthAuthorizationCode, error) {
	authCode := new(OauthAuthorizationCode)
	if err := s.db.QueryRow(getValidAuthCode, clientID, code).Scan(
		&authCode.ID,
		&authCode.ClientID,
		&authCode.UserID,
		&authCode.Code,
		&authCode.RedirectURI,
		&authCode.Scope,
		&authCode.ExpiresAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAuthorizationCodeInvalid
		}
		return nil, err
	}

	// redirect URI must match if it was used to obtain the authorization code
	if redirectURI != authCode.RedirectURI {
		return nil, ErrInvalidRedirectURI
	}

	// check if authorization code expired
	if time.Now().After(authCode.ExpiresAt) {
		return nil, ErrAuthorizationCodeExpired
	}

	return authCode, nil
}

const deleteAuthCode = `
	DELETE FROM authorization_codes
	WHERE id = $1
`

func (s *Service) DeleteAuthCode(id uint32) error {
	_, err := s.db.Exec(deleteAuthCode, id)
	if err != nil {
		return err
	}
	return nil
}

func generateAuthorizationCode() string {
	return uuid.NewString()
}
