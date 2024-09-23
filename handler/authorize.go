package handler

import (
	"context"
	"net/http"
	"net/url"

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

	flashMsg, _ := h.sessionSvc.GetFlashMessage()
	data := map[string]any{
		"error":       flashMsg,
		"queryString": getQueryString(r.URL.Query()),
		"clientID":    client.AppName,
	}
	templateRender(w, r, "base.html", "authorize.html", data)
}

func getRedirectUri(r *http.Request, cli *client.OauthClient) (*url.URL, error) {
	redirectURI := r.Form.Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = cli.RedirectURI
	}

	// parse the redirect URL
	parsedRedirectURI, err := url.ParseRequestURI(redirectURI)
	if err != nil {
		return nil, err
	}

	return parsedRedirectURI, nil
}

func (h *Handler) authorizeCommon(ctx context.Context) (*client.OauthClient, *user.OauthUser, int, error) {
	// check if client injected to context
	cli, err := getClient(ctx)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	// check if userSession injected to context
	userSession, err := getUserSession(ctx)
	if err != nil {
		return nil, nil, http.StatusInternalServerError, err
	}

	usr, err := h.userSvc.FindUserByEmail(userSession.Email)
	if err != nil {
		return nil, nil, http.StatusBadRequest, err
	}

	return cli, usr, 0, nil
}

func (h *Handler) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	cli, user, code, err := h.authorizeCommon(r.Context())
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
	reqScope := r.Form.Get("scope")

	scope, err := h.scopeSvc.FindScope(reqScope)
	if err != nil {
		redirectError(w, r, redirectURI, "invalid_scope", state)
		return
	}

	query := redirectURI.Query()
	authCode, err := h.authCodeSvc.GrantAuthorizationCode(
		cli.ID,
		user.ID,
		redirectURI.String(),
		scope,
		3600, // 1 hour
	)
	if err != nil {
		redirectError(w, r, redirectURI, "server_error", state)
		return
	}

	query.Set("code", authCode.Code)
	if state != "" {
		query.Set("state", state)
	}

	redirectToURL(w, r, redirectURI, query)
}
