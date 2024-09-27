package client

import (
	"database/sql"
	"errors"
)

var (
	ErrClientNotFound = errors.New("Client not found")
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

const findClientByClientUUID = `
SELECT
	id,
	client_id,
	client_secret,
	redirect_uri,
	app_name,
	created_at,
	updated_at
FROM clients
WHERE client_id = $1
`

func (s *Service) FindClientByClientUUID(clientID string) (*OauthClient, error) {
	if clientID == "" {
		clientID = "00000000-0000-0000-0000-000000000000"
	}
	cli := new(OauthClient)
	err := s.db.QueryRow(
		findClientByClientUUID,
		clientID,
	).Scan(
		&cli.ID,
		&cli.ClientID,
		&cli.ClientSecret,
		&cli.RedirectURI,
		&cli.AppName,
		&cli.CreatedAt,
		&cli.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrClientNotFound
		}
		return nil, err
	}

	return cli, nil
}

const findClientByClientID = `
SELECT
	id,
	client_id,
	client_secret,
	redirect_uri,
	app_name,
	created_at,
	updated_at
FROM clients
WHERE id = $1
`

func (s *Service) FindClientByClientID(clientID uint32) (*OauthClient, error) {
	cli := new(OauthClient)
	err := s.db.QueryRow(
		findClientByClientID,
		clientID,
	).Scan(
		&cli.ID,
		&cli.ClientID,
		&cli.ClientSecret,
		&cli.RedirectURI,
		&cli.AppName,
		&cli.CreatedAt,
		&cli.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrClientNotFound
		}
		return nil, err
	}

	return cli, nil
}
