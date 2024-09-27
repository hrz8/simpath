package tokengrant

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hrz8/simpath/config"
	"github.com/hrz8/simpath/internal/authcode"
	"github.com/hrz8/simpath/internal/client"
	"github.com/hrz8/simpath/internal/scope"
	"github.com/hrz8/simpath/internal/token"
	"github.com/hrz8/simpath/internal/user"
	"github.com/hrz8/simpath/jwt"
)

type TokenExchangeBody struct {
	GrantType    string `json:"grant_type"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	Scope        string `json:"scope"` // optional
}

type Service struct {
	db          *sql.DB
	scopeSvc    *scope.Service
	userSvc     *user.Service
	tokenSvc    *token.Service
	authCodeSvc *authcode.Service
}

func NewService(db *sql.DB, sSvc *scope.Service, uSvc *user.Service, tSvc *token.Service, acSvc *authcode.Service) *Service {
	return &Service{db, sSvc, uSvc, tSvc, acSvc}
}

func (s *Service) AuthorizationCodeGrant(body *TokenExchangeBody, client *client.OauthClient) (*AccessTokenResponse, error) {
	authCode, err := s.authCodeSvc.GetValidAuthCode(body.Code, client.ID, body.RedirectURI)
	if err != nil {
		return nil, err
	}
	scope := authCode.Scope

	accessToken, refreshToken, err := s.tokenSvc.Login(
		authCode.ClientID,
		authCode.UserID,
		authCode.Scope,
	)
	if err != nil {
		return nil, err
	}

	var u *user.OauthUser
	u, err = s.userSvc.FindUserByID(accessToken.UserID)
	if err != nil {
		return nil, err
	}

	s.authCodeSvc.DeleteAuthCode(authCode.ID)
	finalResp := &AccessTokenResponse{
		AccessToken:  accessToken.AccessToken,
		ExpiresIn:    config.AccessTokenLifetime,
		TokenType:    "Bearer",
		Scope:        accessToken.Scope,
		UserID:       u.PublicID,
		RefreshToken: refreshToken.RefreshToken,
	}

	if strings.Contains(scope, "openid") {
		jwtClaims := jwt.IDTokenClaims{
			Aud:      client.ClientID, // this token is will be used by client because want to clone the userinfo data
			Sub:      u.PublicID,
			Scope:    authCode.Scope,
			ClientID: client.ClientID,
			AuthTime: time.Now().Add(-3 * time.Minute).Unix(), // TODO: mock for now
			Name:     u.Email,
			Email:    u.Email,
			Picture:  fmt.Sprintf("https://imgur.com/%s.jpg", u.PublicID), // just mock
		}
		idTkn, err := jwt.GenerateIDToken(uuid.NewString(), jwtClaims)
		if err != nil {
			return nil, err
		}

		finalResp.IDToken = idTkn
	}

	return finalResp, nil
}

func (s *Service) RefreshTokenGrant(body *TokenExchangeBody, client *client.OauthClient) (*AccessTokenResponse, error) {
	refreshTkn, err := s.tokenSvc.GetValidRefreshToken(body.RefreshToken, client.ID)
	if err != nil {
		return nil, err
	}

	// get the scope
	scope, err := s.tokenSvc.GetRefreshTokenScope(refreshTkn, body.Scope)
	if err != nil {
		return nil, err
	}

	// log in the user
	accessToken, refreshToken, err := s.tokenSvc.Login(
		refreshTkn.ClientID,
		refreshTkn.UserID,
		scope,
	)
	if err != nil {
		return nil, err
	}

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

func (s *Service) PasswordGrant(body *TokenExchangeBody, client *client.OauthClient) (*AccessTokenResponse, error) {
	scp, err := s.scopeSvc.FindScope(body.Scope)
	if err != nil {
		return nil, err
	}

	u, err := s.userSvc.AuthUser(body.Email, body.Password)
	if err != nil {
		return nil, err
	}

	// log in the user
	accessToken, refreshToken, err := s.tokenSvc.Login(
		client.ID,
		u.ID,
		scp,
	)
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
