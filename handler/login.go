package handler

import (
	"net/http"

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

	if err := h.sessionSvc.SetUserSession(&session.UserSession{
		ClientID:     cli.ID,
		ClientUUID:   cli.ClientID,
		Email:        user.Email,
		AccessToken:  at.AccessToken,
		RefreshToken: rt.RefreshToken,
	}); err != nil {
		h.sessionSvc.SetFlashMessage(err.Error())
		redirectSelf(w, r)
		return
	}

	redirectAuthorize(w, r)
}
