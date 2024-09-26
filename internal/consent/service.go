package consent

import (
	"database/sql"
	"time"
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

const findUserConsent = `
SELECT
	id,
	client_id,
	user_id,
	consent,
	created_at,
	updated_at
FROM consents
WHERE
	user_id = $1 AND
	client_id = $2 AND
	deleted_at IS NULL
`

func (s *Service) IsUserConsent(userID, clientID uint32) (bool, error) {
	con := new(OauthConsent)
	if err := s.db.QueryRow(
		findUserConsent,
		userID,
		clientID,
	).Scan(
		&con.ID,
		&con.ClientID,
		&con.UserID,
		&con.Consent,
		&con.CreatedAt,
		&con.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}

	return con.Consent, nil
}

const softDeleteConsent = `
UPDATE consents
SET deleted_at = $1
WHERE
	user_id = $2 AND
	client_id = $3 AND
	deleted_at IS NULL
`

const createNewConsent = `
INSERT INTO consents (
	client_id,
	user_id,
	consent,
	created_at,
	updated_at
) VALUES (
	$1, $2, $3, $4, $5
) RETURNING *
`

func (s *Service) SetUserConsent(userID, clientID uint32, val bool) (*OauthConsent, error) {
	tx, _ := s.db.Begin()

	_, err := tx.Exec(softDeleteConsent, time.Now(), userID, clientID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	con := new(OauthConsent)
	if err := tx.QueryRow(
		createNewConsent,
		clientID,
		userID,
		val,
		time.Now(),
		time.Now(),
	).Scan(
		&con.ID,
		&con.ClientID,
		&con.UserID,
		&con.Consent,
		&con.CreatedAt,
		&con.UpdatedAt,
		&con.DeletedAt,
	); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	return con, nil
}
