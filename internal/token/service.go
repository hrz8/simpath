package token

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrRefreshTokenNotFound = errors.New("Refresh token not found")
	ErrRefreshTokenExpired  = errors.New("Refresh token expired")
	ErrAccessTokenNotFound  = errors.New("Access token not found")
	ErrAccessTokenExpired   = errors.New("Access token expired")
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

const deleteAccessToken = `
	DELETE FROM access_tokens
	WHERE
		client_id = $1 AND
		user_id = $2 AND
		expires_at <= $3
`
const createNewAccessToken = `
	INSERT INTO access_tokens (
		client_id,
		user_id,
		access_token,
		scope,
		expires_at,
		created_at,
		updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7
	) RETURNING *
`

func (s *Service) GrantAccessToken(clientID uint32, userID uint32, scope string, expiresIn int) (*OauthAccessToken, error) {
	tx, _ := s.db.Begin()

	_, err := tx.Query(deleteAccessToken, clientID, userID, time.Now())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var accessToken OauthAccessToken
	tokenVal := generateAccessToken()
	err = tx.QueryRow(
		createNewAccessToken,
		clientID,
		userID,
		tokenVal,
		scope,
		time.Now().UTC().Add(time.Duration(expiresIn)*time.Second),
		time.Now(),
		time.Now(),
	).Scan(
		&accessToken.ID,
		&accessToken.ClientID,
		&accessToken.UserID,
		&accessToken.AccessToken,
		&accessToken.Scope,
		&accessToken.ExpiresAt,
		&accessToken.CreatedAt,
		&accessToken.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return &accessToken, nil
}

const deleteRefreshToken = `
	DELETE FROM refresh_tokens
	WHERE
		client_id = $1 AND
		user_id = $2 AND
		expires_at <= $3
`
const createNewRefreshToken = `
	INSERT INTO refresh_tokens (
		client_id,
		user_id,
		refresh_token,
		scope,
		expires_at,
		created_at,
		updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING *
`

func (s *Service) GetOrCreateRefreshToken(clientID uint32, userID uint32, scope string, expiresIn int) (*OauthRefreshToken, error) {
	tx, _ := s.db.Begin()

	_, err := tx.Query(deleteRefreshToken, clientID, clientID, time.Now())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var tkn OauthRefreshToken
	tokenVal := generateRefreshToken()
	err = tx.QueryRow(
		createNewRefreshToken,
		clientID,
		clientID,
		tokenVal,
		scope,
		time.Now().UTC().Add(time.Duration(expiresIn)*time.Second),
		time.Now(),
		time.Now(),
	).Scan(
		&tkn.ID,
		&tkn.ClientID,
		&tkn.UserID,
		&tkn.RefreshToken,
		&tkn.Scope,
		&tkn.ExpiresAt,
		&tkn.CreatedAt,
		&tkn.UpdatedAt,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return &tkn, nil
}

const getValidRefreshToken = `
	SELECT * FROM refresh_tokens
	WHERE
		client_id = $1 AND
		refresh_token = $2
`

func (s *Service) GetValidRefreshToken(token string, clientID uint32) (*OauthRefreshToken, error) {
	var tkn OauthRefreshToken

	err := s.db.QueryRow(
		getValidRefreshToken,
		clientID,
		token,
	).Scan(
		&tkn.ID,
		&tkn.ClientID,
		&tkn.UserID,
		&tkn.RefreshToken,
		&tkn.Scope,
		&tkn.ExpiresAt,
		&tkn.CreatedAt,
		&tkn.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrRefreshTokenNotFound
		}
		return nil, err
	}

	if time.Now().UTC().After(tkn.ExpiresAt) {
		return nil, ErrRefreshTokenExpired
	}

	return &tkn, nil
}

const getAccessToken = `
	SELECT
		client_id,
		user_id,
		expires_at
	FROM access_tokens
	WHERE access_token = $1
`
const updateRefreshToken = `
	UPDATE refresh_tokens
	SET expires_at = $1
	WHERE
		client_id = $2 AND
		user_id = $3
`

func (s *Service) Authenticate(token string) (*OauthAccessToken, error) {
	var tkn OauthAccessToken

	err := s.db.QueryRow(
		getAccessToken,
		token,
	).Scan(
		&tkn.ClientID,
		&tkn.UserID,
		&tkn.ExpiresAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrAccessTokenNotFound
		}
		return nil, err
	}

	if time.Now().UTC().After(tkn.ExpiresAt) {
		return nil, ErrAccessTokenExpired
	}

	// extends refresh token
	newExpiry := time.Now().Add(1209600 * time.Second) // 14 days
	_, err = s.db.Exec(updateRefreshToken, newExpiry, tkn.ClientID, tkn.UserID)
	if err != nil {
		return nil, err
	}

	return &tkn, nil
}

func generateAccessToken() string {
	return uuid.NewString()
}

func generateRefreshToken() string {
	return uuid.NewString()
}
