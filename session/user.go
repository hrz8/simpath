package session

import (
	"errors"

	"github.com/hrz8/simpath/config"
)

func (s *Service) StartUserSession() error {
	session, err := s.store.Get(s.r, config.UserSessionName)
	if err != nil {
		return err
	}
	s.userSession = session
	return nil
}

func (s *Service) GetUserData() (*UserData, error) {
	if s.userSession == nil {
		return nil, ErrSessionNotStarted
	}
	userSession, ok := s.userSession.Values[config.UserDataSessionKey].(*UserData)
	if !ok {
		return nil, errors.New("User data type assertion error")
	}
	return userSession, nil
}

func (s *Service) SetUserData(userSession *UserData) error {
	if s.userSession == nil {
		return ErrSessionNotStarted
	}
	s.userSession.Values[config.UserDataSessionKey] = userSession
	return s.userSession.Save(s.r, s.w)
}

func (s *Service) ClearUserData() error {
	if s.userSession == nil {
		return ErrSessionNotStarted
	}

	delete(s.userSession.Values, config.UserDataSessionKey)
	return s.userSession.Save(s.r, s.w)
}
