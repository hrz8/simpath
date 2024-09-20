package handler

import (
	"fmt"
	"net/http"
)

const (
	defaultLoginRedirectUri = "/v1/authorize"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// start cookie session | handle guest
	h.sessionSvc.SetSessionService(w, r)
	h.sessionSvc.StartSession()

	// parse form input to perform r.Form.Method()
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clientID := r.Form.Get("client_id")
	cli, err := h.clientSvc.FindClientByClientID(clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	reqScope := r.Form.Get("scope")

	user, err := h.userSvc.AuthUser(email, password)
	if err != nil {
		h.sessionSvc.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	scope, err := h.scopeSvc.FindScope(reqScope)
	if err != nil {
		h.sessionSvc.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	_, err = h.tokenSvc.GrantAccessToken(cli, user, scope, 3600)
	if err != nil {
		h.sessionSvc.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	refreshTokenExp := 1209600 // 14 days
	_, err = h.tokenSvc.GrantRefreshToken(cli, user, scope, refreshTokenExp)
	if err != nil {
		h.sessionSvc.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	loginRedirectURI := r.URL.Query().Get("login_redirect_uri")
	if loginRedirectURI == "" {
		loginRedirectURI = defaultLoginRedirectUri
	}

	http.Redirect(w, r, fmt.Sprintf("%s%s", loginRedirectURI, getQueryString(r.URL.Query())), http.StatusFound)
}
