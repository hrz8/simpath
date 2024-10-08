package handler

import (
	"net/http"
	"time"

	"github.com/hrz8/simpath/session"
)

var (
	defaultLoginRedirectURI = "/v1/dashboard"
)

func (h *Handler) LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	flashMsg, _ := h.sessionSvc.GetFlashMessage()
	csrfToken, _ := h.sessionSvc.GetCSRFToken()

	data := map[string]any{
		"csrf_token":  csrfToken,
		"error":       flashMsg,
		"queryString": getQueryString(r.URL.Query()),
	}
	templateRender(w, r, "landing.html", "login.html", data)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	reqScope := r.Form.Get("scope")

	// check if client injected to context
	cli, err := getClient(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := h.userSvc.AuthUser(email, password)
	if err != nil {
		h.sessionSvc.SetFlashMessage(err.Error())
		redirectSelf(w, r)
		return
	}

	scope, err := h.scopeSvc.FindScope(reqScope)
	if err != nil {
		h.sessionSvc.SetFlashMessage(err.Error())
		redirectSelf(w, r)
		return
	}

	at, rt, err := h.tokenSvc.Login(cli.ID, user.ID, scope)
	if err != nil {
		h.sessionSvc.SetFlashMessage(err.Error())
		redirectSelf(w, r)
		return
	}

	if err := h.sessionSvc.SetUserData(&session.UserData{
		ClientID:        cli.ID,
		ClientUUID:      cli.ClientID,
		Email:           user.Email,
		AccessToken:     at.AccessToken,
		RefreshToken:    rt.RefreshToken,
		AuthenticatedAt: time.Now(),
	}); err != nil {
		h.sessionSvc.SetFlashMessage(err.Error())
		redirectSelf(w, r)
		return
	}

	loginRedirectURI := r.URL.Query().Get("login_redirect_uri")
	if loginRedirectURI == "" {
		loginRedirectURI = defaultLoginRedirectURI
	}

	if loginRedirectURI == "/v1/oauth2/authorize" {
		redirectAuthorize(w, r)
		return
	}

	redirectToPath(w, r, loginRedirectURI)
}
