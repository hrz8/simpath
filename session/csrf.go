package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/hrz8/simpath/config"
)

func generateCSRFToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}

func (s *Service) GetCSRFToken() (string, error) {
	if s.session == nil {
		return "", ErrSessionNotStarted
	}
	csrfTkn, ok := s.session.Values[config.CSRFTokenSessionKey].(string)
	if !ok {
		return "", errors.New("CSRF token type assertion error")
	}
	return csrfTkn, nil
}

func (s *Service) SetCSRFToken() error {
	if s.session == nil {
		return ErrSessionNotStarted
	}

	if s.session.Values[config.CSRFTokenSessionKey] == nil {
		csrfToken, err := generateCSRFToken()
		if err != nil {
			return errors.New("CSRF generation error")
		}
		s.session.Values[config.CSRFTokenSessionKey] = csrfToken
	}

	return s.session.Save(s.r, s.w)
}
