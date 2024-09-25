package handler

import (
	"net/http"

	"github.com/hrz8/simpath/internal/role"
)

func (h *Handler) RegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	flashMsg, _ := h.sessionSvc.GetFlashMessage()
	csrfToken, _ := h.sessionSvc.GetCSRFToken()

	data := map[string]any{
		"csrf_token":  csrfToken,
		"error":       flashMsg,
		"queryString": getQueryString(r.URL.Query()),
	}
	templateRender(w, r, "landing.html", "register.html", data)
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	if h.userSvc.IsUserExists(email) {
		h.sessionSvc.SetFlashMessage("Email already taken")
		redirectSelf(w, r)
		return
	}

	_, err := h.userSvc.CreateUser(
		role.User,
		email,
		password,
	)
	if err != nil {
		h.sessionSvc.SetFlashMessage(err.Error())
		redirectSelf(w, r)
		return
	}

	redirectLogin(w, r)
}
