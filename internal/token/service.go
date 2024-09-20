package token

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/hrz8/simpath/internal/client"
	"github.com/hrz8/simpath/internal/user"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

const deleteAccessToken = "DELETE FROM access_tokens WHERE client_id = $1 AND user_id = $2 AND expires_at <= $3"
const createNewAccessToken = "INSERT INTO access_tokens (client_id, user_id, access_token, scope, expires_at, created_at, updated_at) VALUEs ($1, $2, $3, $4, $5, $6, $7) RETURNING *"

func (s *Service) GrantAccessToken(cli *client.OauthClient, usr *user.OauthUser, scope string, expiresIn int) (*OauthAccessToken, error) {
	tx, _ := s.db.Begin()

	_, err := tx.Query(deleteAccessToken, cli.ID, usr.ID, time.Now())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var accessToken OauthAccessToken
	tokenVal := generateAccessToken()
	err = tx.QueryRow(
		createNewAccessToken,
		cli.ID,
		usr.ID,
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

const deleteRefreshToken = "DELETE FROM refresh_tokens WHERE client_id = $1 AND user_id = $2 AND expires_at <= $3"
const createNewRefreshToken = "INSERT INTO refresh_tokens (client_id, user_id, refresh_token, scope, expires_at, created_at, updated_at) VALUEs ($1, $2, $3, $4, $5, $6, $7) RETURNING *"

func (s *Service) GrantRefreshToken(cli *client.OauthClient, usr *user.OauthUser, scope string, expiresIn int) (*OauthRefreshToken, error) {
	tx, _ := s.db.Begin()

	_, err := tx.Query(deleteRefreshToken, cli.ID, usr.ID, time.Now())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var tkn OauthRefreshToken
	tokenVal := generateRefreshToken()
	err = tx.QueryRow(
		createNewRefreshToken,
		cli.ID,
		usr.ID,
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

func generateAccessToken() string {
	return uuid.NewString()
}

func generateRefreshToken() string {
	return uuid.NewString()
}
