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

	user, err := h.UserSvc.AuthUser(
		r.Form.Get("email"),
		r.Form.Get("password"),
	)
	if err != nil {
		// set session error
		http.Redirect(w, r, r.RequestURI, http.StatusFound)
		return
	}

	fmt.Fprint(w, user.EncryptedPassword)
}
