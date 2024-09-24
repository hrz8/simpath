package introspect

import (
	"database/sql"
	"errors"

	"github.com/hrz8/simpath/internal/client"
	"github.com/hrz8/simpath/internal/token"
	"github.com/hrz8/simpath/internal/user"
)

var (
	ErrTokenMissing     = errors.New("Token missing")
	ErrTokenHintInvalid = errors.New("Invalid token hint")
)

const (
	AccessTokenHint  = "access_token"
	RefreshTokenHint = "refresh_token"
)

type TokenIntrospectBody struct {
	Token         string `json:"token"`
	TokenTypeHint string `json:"token_type_hint"`
}

type Service struct {
	db       *sql.DB
	userSvc  *user.Service
	tokenSvc *token.Service
}

func NewService(db *sql.DB, uSvc *user.Service, tSvc *token.Service) *Service {
	return &Service{db, uSvc, tSvc}
}

func (s *Service) IntrospectToken(body *TokenIntrospectBody, client *client.OauthClient) (*IntrospectResponse, error) {
	token := body.Token
	if token == "" {
		return nil, ErrTokenMissing
	}

	// Get token type hint from the query
	tokenTypeHint := body.TokenTypeHint

	// Default to access token hint
	if tokenTypeHint == "" {
		tokenTypeHint = AccessTokenHint
	}

	switch tokenTypeHint {
	case AccessTokenHint:
		accessToken, err := s.tokenSvc.Authenticate(token)
		if err != nil {
			return nil, err
		}
		u, err := s.userSvc.FindUserByID(accessToken.UserID)
		if err != nil {
			return nil, err
		}
		return &IntrospectResponse{
			Active:    true,
			Scope:     accessToken.Scope,
			TokenType: "Bearer",
			ExpiresAt: int(accessToken.ExpiresAt.Unix()),
			ClientID:  client.AppName,
			Username:  u.Email,
		}, nil
	case RefreshTokenHint:
		refreshToken, err := s.tokenSvc.GetValidRefreshToken(token, client.ID)
		if err != nil {
			return nil, err
		}
		u, err := s.userSvc.FindUserByID(refreshToken.UserID)
		if err != nil {
			return nil, err
		}
		return &IntrospectResponse{
			Active:    true,
			Scope:     refreshToken.Scope,
			TokenType: "Bearer",
			ExpiresAt: int(refreshToken.ExpiresAt.Unix()),
			ClientID:  client.AppName,
			Username:  u.Email,
		}, nil
	default:
		return nil, ErrTokenHintInvalid
	}
}
