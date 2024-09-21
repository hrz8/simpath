package session

import (
	"encoding/gob"
	"errors"
	"net/http"

	"github.com/gorilla/sessions"
)

type Service struct {
	sessionStore   sessions.Store
	sessionOptions *sessions.Options
	session        *sessions.Session
	r              *http.Request
	w              http.ResponseWriter
}

type UserSession struct {
	ClientID     string
	Email        string
	AccessToken  string
	RefreshToken string
}

var (
	StorageSessionName   = "simpath_session"
	UserSessionKey       = "simpath_user"
	ErrSessionNotStarted = errors.New("Session not started")
)

func init() {
	gob.Register(new(UserSession))
}

func NewService() *Service {
	cookieStore := sessions.NewCookieStore([]byte("some_secret"))

	return &Service{
		sessionStore: cookieStore,
		sessionOptions: &sessions.Options{
			Path:     "/",
			MaxAge:   604800,
			HttpOnly: true,
		},
	}
}

func (s *Service) SetSessionService(w http.ResponseWriter, r *http.Request) {
	s.w = w
	s.r = r
}

func (s *Service) StartSession() error {
	session, err := s.sessionStore.Get(s.r, StorageSessionName)
	if err != nil {
		return err
	}
	s.session = session
	return nil
}

func (s *Service) GetUserSession() (*UserSession, error) {
	if s.session == nil {
		return nil, ErrSessionNotStarted
	}
	userSession, ok := s.session.Values[UserSessionKey].(*UserSession)
	if !ok {
		return nil, errors.New("User session type assertion error")
	}

	return userSession, nil
}

func (s *Service) SetUserSession(userSession *UserSession) error {
	if s.session == nil {
		return ErrSessionNotStarted
	}
	s.session.Values[UserSessionKey] = userSession
	return s.session.Save(s.r, s.w)
}

func (s *Service) ClearUserSession() error {
	if s.session == nil {
		return ErrSessionNotStarted
	}

	delete(s.session.Values, UserSessionKey)
	return s.session.Save(s.r, s.w)
}

func (s *Service) SetFlashMessage(msg string) error {
	if s.session == nil {
		return ErrSessionNotStarted
	}
	s.session.AddFlash(msg)
	return s.session.Save(s.r, s.w)
}

func (s *Service) GetFlashMessage() (any, error) {
	if s.session == nil {
		return nil, ErrSessionNotStarted
	}
	if flashes := s.session.Flashes(); len(flashes) > 0 {
		s.session.Save(s.r, s.w)
		return flashes[0], nil
	}
	return nil, nil
}
