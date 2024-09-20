package scope

import (
	"database/sql"
	"errors"
	"fmt"
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

func toAny(slice []string) []any {
	interfaceSlice := make([]any, len(slice))
	for i, v := range slice {
		interfaceSlice[i] = v
	}
	return interfaceSlice
}

func (s *Service) FindScope(requestedScope string) (string, error) {
	if requestedScope == "" {
		return s.GetDefaultScope(), nil
	}

	scopes := strings.Split(requestedScope, " ")
	placeholders := make([]string, len(scopes))
	for i := range scopes {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	var count int
	sql := fmt.Sprintf("SELECT COUNT(scope) FROM scopes WHERE scope IN (%s)", strings.Join(placeholders, ","))
	s.db.QueryRow(sql, toAny(scopes)...).Scan(&count)

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
