package session

import (
	"encoding/gob"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/sessions"
	"github.com/hrz8/simpath/config"
)

type Service struct {
	store       sessions.Store
	options     *sessions.Options
	userSession *sessions.Session
	session     *sessions.Session
	r           *http.Request
	w           http.ResponseWriter
}

type UserData struct {
	ClientID        uint32
	ClientUUID      string
	Email           string
	AccessToken     string
	RefreshToken    string
	AuthenticatedAt time.Time
}

var (
	ErrSessionNotStarted = errors.New("Session not started")
)

func init() {
	gob.Register(new(UserData))
}

func NewService() *Service {
	store := sessions.NewCookieStore([]byte(config.SessionSecretKey))
	return &Service{
		store: store, // max age default set to be 30 days
		options: &sessions.Options{
			Path:     config.SessionPath,
			MaxAge:   config.SessionMaxAge,
			HttpOnly: config.SessionHttpOnly,
		},
	}
}

func (s *Service) SetSessionService(w http.ResponseWriter, r *http.Request) {
	s.w = w
	s.r = r
}

func (s *Service) StartSession() error {
	session, err := s.store.Get(s.r, config.SessionName)
	if err != nil {
		return err
	}
	session.Options.MaxAge = s.options.MaxAge
	s.session = session
	return nil
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
