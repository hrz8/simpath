package tokengrant

import (
	"database/sql"

	"github.com/hrz8/simpath/config"
	"github.com/hrz8/simpath/internal/authcode"
	"github.com/hrz8/simpath/internal/client"
	"github.com/hrz8/simpath/internal/token"
	"github.com/hrz8/simpath/internal/user"
)

type Service struct {
	db          *sql.DB
	userSvc     *user.Service
	tokenSvc    *token.Service
	authCodeSvc *authcode.Service
}

func NewService(db *sql.DB, uSvc *user.Service, tSvc *token.Service, acSvc *authcode.Service) *Service {
	return &Service{db, uSvc, tSvc, acSvc}
}

func (s *Service) AuthorizationCodeGrant(code, redirectURI string, client *client.OauthClient) (*AccessTokenResponse, error) {
	authCode, err := s.authCodeSvc.GetValidAuthCode(client.ID, code, redirectURI)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.tokenSvc.Login(
		authCode.ClientID,
		authCode.UserID,
		authCode.Scope,
	)
	if err != nil {
		return nil, err
	}

	s.authCodeSvc.DeleteAuthCode(authCode.ID)

	var u *user.OauthUser
	u, err = s.userSvc.FindUserByID(accessToken.UserID)
	if err != nil {
		return nil, err
	}

	return &AccessTokenResponse{
		AccessToken:  accessToken.AccessToken,
		ExpiresIn:    config.AccessTokenLifetime,
		TokenType:    "Bearer",
		Scope:        accessToken.Scope,
		UserID:       u.PublicID,
		RefreshToken: refreshToken.RefreshToken,
	}, nil
}

func (s *Service) RefreshTokenGrant(code, redirectURI string, client *client.OauthClient) (*AccessTokenResponse, error) {
	return nil, nil
}
