package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/hrz8/simpath/internal/authcode"
	"github.com/hrz8/simpath/internal/client"
	"github.com/hrz8/simpath/internal/introspect"
	"github.com/hrz8/simpath/internal/tokengrant"
	"github.com/hrz8/simpath/internal/user"
	"github.com/hrz8/simpath/password"
	"github.com/hrz8/simpath/response"
)

var (
	ErrInvalidGrantType        = errors.New("Invalid grant type")
	ErrInvalidClientIDOrSecret = errors.New("Invalid client id or client secret")

	errStatusCodeMap = map[error]int{
		authcode.ErrAuthorizationCodeInvalid: http.StatusBadRequest,
		authcode.ErrAuthorizationCodeExpired: http.StatusUnprocessableEntity,
		authcode.ErrInvalidRedirectURI:       http.StatusUnprocessableEntity,
		user.ErrUserNotFound:                 http.StatusNotFound,
		user.ErrInvalidUserOrPassword:        http.StatusUnprocessableEntity,
		introspect.ErrTokenMissing:           http.StatusUnprocessableEntity,
		introspect.ErrTokenHintInvalid:       http.StatusUnprocessableEntity,
	}
)

func getErrStatusCode(err error) int {
	code, ok := errStatusCodeMap[err]
	if ok {
		return code
	}

	return http.StatusInternalServerError
}

func (h *Handler) clientBasicAuth(r *http.Request) (*client.OauthClient, error) {
	// get client credentials from basic auth
	clientID, secret, ok := r.BasicAuth()
	if !ok {
		return nil, ErrInvalidClientIDOrSecret
	}

	// authenticate the client
	client, err := h.AuthClient(clientID, secret)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (h *Handler) AuthClient(clientID, secret string) (*client.OauthClient, error) {
	// fetch client
	cli, err := h.clientSvc.FindClientByClientUUID(clientID)
	if err != nil {
		return nil, ErrInvalidClientIDOrSecret
	}

	// verify the secret
	if password.VerifyPassword(cli.ClientSecret, secret) != nil {
		return nil, ErrInvalidClientIDOrSecret
	}

	return cli, nil
}

func (h *Handler) TokenHandler(w http.ResponseWriter, r *http.Request) {
	body := new(tokengrant.TokenExchangeBody)
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// map of grant types against each functions
	grantTypes := map[string]func(body *tokengrant.TokenExchangeBody, client *client.OauthClient) (*tokengrant.AccessTokenResponse, error){
		"authorization_code": h.tokenGrantSvc.AuthorizationCodeGrant,
		"refresh_token":      h.tokenGrantSvc.RefreshTokenGrant,
		"password":           h.tokenGrantSvc.PasswordGrant,
		// "client_credentials": h.tokenGrantSvc.ClientCredentialsGrant,
	}

	// check the grant type
	grantFn, ok := grantTypes[body.GrantType]
	if !ok {
		response.Error(w, ErrInvalidGrantType.Error(), http.StatusBadRequest)
		return
	}

	cli, err := h.clientBasicAuth(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	// processing access token grant
	resp, err := grantFn(body, cli)
	if err != nil {
		response.Error(w, err.Error(), getErrStatusCode(err))
		return
	}

	// write response to json
	response.WriteJSON(w, resp, 200)
}

func (h *Handler) IntrospectHandler(w http.ResponseWriter, r *http.Request) {
	body := new(introspect.TokenIntrospectBody)
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	cli, err := h.clientBasicAuth(r)
	if err != nil {
		response.UnauthorizedError(w, err.Error())
		return
	}

	resp, err := h.introspectSvc.IntrospectToken(body, cli)
	if err != nil {
		response.Error(w, err.Error(), getErrStatusCode(err))
		return
	}

	// write response to json
	response.WriteJSON(w, resp, 200)
}
