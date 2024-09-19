package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
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
		// set session error
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	scope, err := h.scopeSvc.FindScope(reqScope)
	if err != nil {
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	at, err := h.tokenSvc.GrantAccessToken(cli, user, scope, 3600)
	if err != nil {
		fmt.Println(at)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, fmt.Sprintf("logged in: %s", user.Email))
}
