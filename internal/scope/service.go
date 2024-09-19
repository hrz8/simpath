package scope

import (
	"database/sql"
	"errors"
	"strings"
)

var (
	ErrInvalidScope = errors.New("Invalid scope")
)

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{
		db: db,
	}
}

const findScope = "SELECT COUNT(scope) FROM scopes WHERE scope IN (?)"

func (s *Service) FindScope(requestedScope string) (string, error) {
	if requestedScope == "" {
		return s.GetDefaultScope(), nil
	}
	scopes := strings.Split(requestedScope, " ")
	var count int
	s.db.QueryRow(findScope, scopes).Scan(&count)

	if count == len(scopes) {
		return requestedScope, nil
	}

	return "", ErrInvalidScope
}

const getDefaultScope = "SELECT scope, is_default FROM scopes WHERE is_default = true"

func (s *Service) GetDefaultScope() string {
	var scopes []string
	rows, _ := s.db.Query(getDefaultScope)
	defer rows.Close()
	for rows.Next() {
		var i OauthScope
		rows.Scan(&i.Scope, &i.IsDefault)
		scopes = append(scopes, i.Scope)
	}
	return strings.Join(scopes, " ")
}
