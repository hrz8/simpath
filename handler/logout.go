package handler

import (
	"net/http"
)

func (h *Handler) LogoutPage(w http.ResponseWriter, r *http.Request) {
	// prevent non-logged-in user to access the page
	userSession, err := getUserSession(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.tokenSvc.ClearUserTokens(userSession.RefreshToken, userSession.AccessToken)
	h.sessionSvc.ClearUserSession()

	redirectLogin(w, r)
}
