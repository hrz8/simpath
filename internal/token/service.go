package token

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/hrz8/simpath/config"
	"github.com/hrz8/simpath/internal/scope"
)

var (
	ErrRefreshTokenNotFound          = errors.New("Refresh token not found")
	ErrRefreshTokenExpired           = errors.New("Refresh token expired")
	ErrAccessTokenNotFound           = errors.New("Access token not found")
	ErrAccessTokenExpired            = errors.New("Access token expired")
	ErrRequestedScopeCannotBeGreater = errors.New("Requested scope cannot be greater")
)

type Service struct {
	db       *sql.DB
	scopeSvc *scope.Service
}

func NewService(db *sql.DB, sSvc *scope.Service) *Service {
	return &Service{db, sSvc}
}

func (s *Service) Login(clientID uint32, userID uint32, scope string) (*OauthAccessToken, *OauthRefreshToken, error) {
	at, err := s.GrantAccessToken(clientID, userID, scope, config.AccessTokenLifetime)
	if err != nil {
		return nil, nil, err
	}

	refreshTokenExp := config.RefreshTokenLifetime
	rt, err := s.GetOrCreateRefreshToken(clientID, userID, scope, refreshTokenExp)
	if err != nil {
		return nil, nil, err
	}

	return at, rt, err
}

const softDeleteExpiredAccessToken = `
	UPDATE access_tokens
	SET deleted_at = $1
	WHERE
		client_id = $2 AND
		user_id = $3 AND
		expires_at <= $4 AND
		deleted_at IS NULL
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

	_, err := tx.Exec(softDeleteExpiredAccessToken, time.Now(), clientID, userID, time.Now())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var tkn OauthAccessToken
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
		&tkn.ID,
		&tkn.ClientID,
		&tkn.UserID,
		&tkn.AccessToken,
		&tkn.Scope,
		&tkn.ExpiresAt,
		&tkn.CreatedAt,
		&tkn.UpdatedAt,
		&tkn.DeletedAt,
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

const softDeleteExpiredRefreshToken = `
	UPDATE refresh_tokens
	SET deleted_at = $1
	WHERE
		client_id = $2 AND
		user_id = $3 AND
		expires_at <= $4 AND
		deleted_at IS NULL
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
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7
	) RETURNING *
`

func (s *Service) GetOrCreateRefreshToken(clientID uint32, userID uint32, scope string, expiresIn int) (*OauthRefreshToken, error) {
	tx, _ := s.db.Begin()

	_, err := tx.Exec(softDeleteExpiredRefreshToken, time.Now(), clientID, clientID, time.Now())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tkn := new(OauthRefreshToken)
	tokenVal := generateRefreshToken()
	err = tx.QueryRow(
		createNewRefreshToken,
		clientID,
		userID,
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
		&tkn.DeletedAt,
	)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return tkn, nil
}

const getValidRefreshToken = `
	SELECT * FROM refresh_tokens
	WHERE
		client_id = $1 AND
		refresh_token = $2 AND
		deleted_at IS NULL
`

func (s *Service) GetValidRefreshToken(token string, clientID uint32) (*OauthRefreshToken, error) {
	tkn := new(OauthRefreshToken)
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
		&tkn.DeletedAt,
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

	return tkn, nil
}

func (s *Service) GetRefreshTokenScope(refreshTkn *OauthRefreshToken, reqScope string) (string, error) {
	var err error
	scope := refreshTkn.Scope

	if reqScope != "" {
		scope, err = s.scopeSvc.FindScope(reqScope)
		if err != nil {
			return "", err
		}
	}

	// Requested scope CANNOT include any scope not originally granted
	if !spaceDelimitedStringNotGreater(scope, refreshTkn.Scope) {
		return "", ErrRequestedScopeCannotBeGreater
	}

	return scope, nil
}

const getAccessToken = `
	SELECT
		client_id,
		user_id,
		expires_at
	FROM access_tokens
	WHERE
		access_token = $1 AND
		deleted_at IS NULL
`
const updateRefreshToken = `
	UPDATE refresh_tokens
	SET expires_at = $1
	WHERE
		client_id = $2 AND
		user_id = $3 AND
		deleted_at IS NULL
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
	newExpiry := time.Now().Add(config.RefreshTokenLifetime * time.Second) // 14 days
	_, err = s.db.Exec(updateRefreshToken, newExpiry, tkn.ClientID, tkn.UserID)
	if err != nil {
		return nil, err
	}

	return &tkn, nil
}

const getRefreshTokenByToken = `
	SELECT id, refresh_token, client_id, user_id FROM refresh_tokens
	WHERE
		refresh_token = $1 AND
		deleted_at IS NULL
`
const softDeleteRefreshToken = `
	UPDATE refresh_tokens
	SET deleted_at = $1
	WHERE
		client_id = $2 AND
		user_id = $3 AND
		deleted_at IS NULL
`
const getAccessTokenByToken = `
	SELECT id, access_token, client_id, user_id FROM access_tokens
	WHERE
		access_token = $1 AND
		deleted_at IS NULL
`
const softDeleteAccessToken = `
	UPDATE access_tokens
	SET deleted_at = $1
	WHERE
		client_id = $2 AND
		user_id = $3 AND
		deleted_at IS NULL
`

func (s *Service) ClearUserTokens(refreshToken string, accessToken string) {
	refreshTkn := new(OauthRefreshToken)
	err := s.db.QueryRow(
		getRefreshTokenByToken,
		refreshToken,
	).Scan(
		&refreshTkn.ID,
		&refreshTkn.RefreshToken,
		&refreshTkn.ClientID,
		&refreshTkn.UserID,
	)
	if err == nil {
		s.db.Exec(softDeleteRefreshToken, time.Now(), refreshTkn.ClientID, refreshTkn.UserID)
	}

	accessTkn := new(OauthAccessToken)
	err = s.db.QueryRow(
		getAccessTokenByToken,
		accessToken,
	).Scan(
		&accessTkn.ID,
		&accessTkn.AccessToken,
		&accessTkn.ClientID,
		&accessTkn.UserID,
	)
	if err == nil {
		s.db.Exec(softDeleteAccessToken, time.Now(), accessTkn.ClientID, accessTkn.UserID)
	}
}

func generateAccessToken() string {
	return uuid.NewString()
}

func generateRefreshToken() string {
	return uuid.NewString()
}
