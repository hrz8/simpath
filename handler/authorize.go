package handler

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/hrz8/simpath/config"
	"github.com/hrz8/simpath/internal/client"
	"github.com/hrz8/simpath/internal/user"
)

func (h *Handler) AuthorizeFormHandler(w http.ResponseWriter, r *http.Request) {
	// check if client injected to context
	client, err := getClient(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// fetch current session
	cli, usr, eCode, err := h.authorizeCommon(r.Context())
	if err != nil {
		http.Error(w, err.Error(), eCode)
		return
	}

	state := r.Form.Get("state")

	// fetch user existing consent
	uConsent, err := h.consentSvc.IsUserConsent(usr.ID, cli.ID)
	if !uConsent {
		flashMsg, _ := h.sessionSvc.GetFlashMessage()
		csrfToken, _ := h.sessionSvc.GetCSRFToken()

		data := map[string]any{
			"csrf_token":  csrfToken,
			"error":       flashMsg,
			"queryString": getQueryString(r.URL.Query()),
			"clientID":    client.AppName,
		}
		templateRender(w, r, "base.html", "authorize.html", data)
		return
	}

	// user already have a consent before
	redirectURI, err := getRedirectUri(r, cli)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query, err := h.processAuthorize(r, authorizeParams{
		redirectURI: redirectURI,
		clientID:    cli.ID,
		userID:      usr.ID,
	})
	if err != nil {
		redirectError(w, r, redirectURI, err.Error(), state)
		return
	}

	redirectToURL(w, r, redirectURI, query)
}

func (h *Handler) authorizeCommon(ctx context.Context) (*client.OauthClient, *user.OauthUser, int, error) {
	// check if client injected to context
	cli, err := getClient(ctx)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	// check if userData injected to context
	userData, err := getUserDataFromSession(ctx)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	usr, err := h.userSvc.FindUserByEmail(userData.Email)
	if err != nil {
		return nil, nil, http.StatusBadRequest, err
	}

	return cli, usr, 0, nil
}

type authorizeParams struct {
	redirectURI *url.URL
	clientID    uint32
	userID      uint32
}

func (h *Handler) processAuthorize(r *http.Request, params authorizeParams) (url.Values, error) {
	state := r.Form.Get("state")
	reqScope := r.Form.Get("scope")

	scope, err := h.scopeSvc.FindScope(reqScope)
	if err != nil {
		return nil, errors.New("Invalid Scope")
	}

	query := params.redirectURI.Query()
	authCode, err := h.authCodeSvc.GrantAuthorizationCode(
		params.clientID,
		params.userID,
		params.redirectURI.String(),
		scope,
		config.AccessTokenLifetime,
	)
	if err != nil {
		return nil, errors.New("Server Error")
	}

	query.Set("code", authCode.Code)
	if state != "" {
		query.Set("state", state)
	}

	return query, nil
}

func (h *Handler) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	cli, usr, code, err := h.authorizeCommon(r.Context())
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	redirectURI, err := getRedirectUri(r, cli)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	state := r.Form.Get("state")

	authorized := len(r.Form.Get("allow")) > 0
	if !authorized {
		_, err := h.consentSvc.SetUserConsent(usr.ID, cli.ID, false)
		if err != nil {
			redirectError(w, r, redirectURI, "server_error", state)
			return
		}

		redirectError(w, r, redirectURI, "access_denied", state)
		return
	}

	_, err = h.consentSvc.SetUserConsent(usr.ID, cli.ID, true)
	if err != nil {
		redirectError(w, r, redirectURI, "server_error", state)
		return
	}

	query, err := h.processAuthorize(r, authorizeParams{
		redirectURI: redirectURI,
		clientID:    cli.ID,
		userID:      usr.ID,
	})
	if err != nil {
		redirectError(w, r, redirectURI, toSnakeCase(err.Error()), state)
		return
	}

	redirectToURL(w, r, redirectURI, query)
}
