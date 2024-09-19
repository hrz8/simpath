package handler

import (
	"net/http"
)

func (h *Handler) LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clientID := r.Form.Get("client_id")
	_, err := h.clientSvc.FindClientByClientID(clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"error":       r.URL.Query().Get("error"),
		"queryString": r.URL.RawQuery,
	}

	baseTemplate := "landing.html"
	TemplateRender(w, r, baseTemplate, "login.html", data)
}

func (h *Handler) RegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"error":       r.URL.Query().Get("error"),
		"queryString": r.URL.RawQuery,
	}

	baseTemplate := "landing.html"
	TemplateRender(w, r, baseTemplate, "register.html", data)
}
