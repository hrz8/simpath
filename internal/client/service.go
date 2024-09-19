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

const findClientByClientID = "SELECT id, redirect_uri, app_name FROM clients WHERE client_id = $1"

func (s *Service) FindClientByClientID(clientID string) (*OauthClient, error) {
	if clientID == "" {
		clientID = "00000000-0000-0000-0000-000000000000"
	}
	var cli OauthClient
	err := s.db.QueryRow(
		findClientByClientID,
		clientID,
	).Scan(
		&cli.ID,
		&cli.RedirectURI,
		&cli.AppName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrClientNotFound
		}
		return nil, err
	}

	return &cli, nil
}
