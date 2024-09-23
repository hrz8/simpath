package authcode

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
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

func generateAuthorizationCode() string {
	return uuid.NewString()
}
