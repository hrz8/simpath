package handler

import (
	"fmt"
	"net/http"

	"github.com/hrz8/simpath/internal/token"
	"github.com/hrz8/simpath/session"
)

func (h *Handler) LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	flashMsg, _ := h.sessionSvc.GetFlashMessage()

	data := map[string]any{
		"error":       flashMsg,
		"queryString": r.URL.RawQuery,
	}
	templateRender(w, r, "landing.html", "login.html", data)
}

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
	cli, err := h.clientSvc.FindClientByClientUUID(clientID)
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

	at, rt, err := h.login(cli.ID, user.ID, scope)
	if err != nil {
		h.sessionSvc.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	userSession := &session.UserSession{
		ClientID:     cli.ID,
		ClientUUID:   cli.ClientID,
		Email:        user.Email,
		AccessToken:  at.AccessToken,
		RefreshToken: rt.RefreshToken,
	}
	if err := h.sessionSvc.SetUserSession(userSession); err != nil {
		h.sessionSvc.SetFlashMessage(err.Error())
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/v1/authorize%s", getQueryString(r.URL.Query())), http.StatusFound)
}

func (h *Handler) login(clientID uint32, userID uint32, scope string) (*token.OauthAccessToken, *token.OauthRefreshToken, error) {
	at, err := h.tokenSvc.GrantAccessToken(clientID, userID, scope, 3600)
	if err != nil {
		return nil, nil, err
	}

	refreshTokenExp := 1209600 // 14 days
	rt, err := h.tokenSvc.GetOrCreateRefreshToken(clientID, userID, scope, refreshTokenExp)
	if err != nil {
		return nil, nil, err
	}

	return at, rt, err
}
