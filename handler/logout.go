package handler

import (
	"net/http"
)

func (h *Handler) LogoutPage(w http.ResponseWriter, r *http.Request) {
	// prevent non-logged-in user to access the page
	userSession, err := getUserDataFromSession(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.tokenSvc.ClearUserTokens(userSession.RefreshToken, userSession.AccessToken)
	h.sessionSvc.ClearUserData()

	redirectLogin(w, r)
}
