package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) AuthorizeFormHandler(w http.ResponseWriter, r *http.Request) {
	// check if client_id exists
	client, err := getClient(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	flashMsg, _ := h.sessionSvc.GetFlashMessage()
	data := map[string]any{
		"error":       flashMsg,
		"queryString": getQueryString(r.URL.Query()),
		"clientID":    client.AppName,
	}
	templateRender(w, r, "base.html", "authorize.html", data)
}

func (h *Handler) AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "kuy")
}
