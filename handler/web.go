package handler

import (
	"fmt"
	"net/http"
	"net/url"
)

func getQueryString(query url.Values) string {
	encoded := query.Encode()
	if len(encoded) > 0 {
		encoded = fmt.Sprintf("?%s", encoded)
	}
	return encoded
}

func (h *Handler) LoginFormHandler(w http.ResponseWriter, r *http.Request) {
	// start cookie session | handle guest
	h.sessionSvc.SetSessionService(w, r)
	h.sessionSvc.StartSession()
	flashMsg, _ := h.sessionSvc.GetFlashMessage()

	// parse form input to perform r.Form.Method()
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

	data := map[string]any{
		"error":       flashMsg,
		"queryString": r.URL.RawQuery,
	}
	TemplateRender(w, r, "landing.html", "login.html", data)
}

func (h *Handler) RegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	// start cookie session | handle guest
	h.sessionSvc.SetSessionService(w, r)
	h.sessionSvc.StartSession()

	data := map[string]any{
		"error":       r.URL.Query().Get("error"),
		"queryString": r.URL.RawQuery,
	}
	TemplateRender(w, r, "landing.html", "register.html", data)
}

func (h *Handler) AuthorizeFormHandler(w http.ResponseWriter, r *http.Request) {
	// start cookie session | handle guest
	h.sessionSvc.SetSessionService(w, r)
	h.sessionSvc.StartSession()

	// parse form input to perform r.Form.Method()
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clientID := r.Form.Get("client_id")
	client, err := h.clientSvc.FindClientByClientID(clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]any{
		"error":       r.URL.Query().Get("error"),
		"queryString": getQueryString(r.URL.Query()),
		"token":       false,
		"clientID":    client.AppName,
	}
	TemplateRender(w, r, "base.html", "authorize.html", data)
}
