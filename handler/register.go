package handler

import "net/http"

func (h *Handler) RegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	// start cookie session | handle guest
	h.sessionSvc.SetSessionService(w, r)
	h.sessionSvc.StartSession()

	data := map[string]any{
		"error":       r.URL.Query().Get("error"),
		"queryString": r.URL.RawQuery,
	}
	templateRender(w, r, "landing.html", "register.html", data)
}
