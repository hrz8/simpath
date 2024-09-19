package handler

import (
	"net/http"
)

func (h *Handler) LoginFormHandler(w http.ResponseWriter, r *http.Request) {
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
