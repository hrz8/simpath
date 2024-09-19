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
	err = tx.QueryRow(
		createNewAccessToken,
		cli.ID,
		usr.ID,
		uuid.NewString(),
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
