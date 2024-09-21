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
	flashMsg, _ := h.sessionSvc.GetFlashMessage()

	// prevent non-logged-in user to access the page
	userSession, err := h.sessionSvc.GetUserSession()
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/v1/login%s", getQueryString(r.URL.Query())), http.StatusFound)
		return
	}

	// parse form input to perform r.Form.Method()
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// check if client_id exists
	clientID := r.Form.Get("client_id")
	client, err := h.clientSvc.FindClientByClientID(clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check token
	_, err = h.tokenSvc.Authenticate(userSession.AccessToken)
	if err != nil {
		fmt.Println(err)
		// got error okay may be it's refresh token
		// check refresh token validity (might be expired already)
		refreshToken, err := h.tokenSvc.GetValidRefreshToken(userSession.RefreshToken, client.ID)
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("/v1/login%s", getQueryString(r.URL.Query())), http.StatusFound)
			return
		}

		// login to create access token and refresh token for user
		at, rt, err := h.login(
			refreshToken.ClientID,
			refreshToken.UserID,
			refreshToken.Scope,
		)
		if err != nil {
			http.Redirect(w, r, fmt.Sprintf("/v1/login%s", getQueryString(r.URL.Query())), http.StatusFound)
			return
		}

		userSession.AccessToken = at.AccessToken
		userSession.RefreshToken = rt.RefreshToken
	}

	h.sessionSvc.SetUserSession(userSession)

	data := map[string]any{
		"error":       flashMsg,
		"queryString": getQueryString(r.URL.Query()),
		"token":       false,
		"clientID":    client.AppName,
	}
	TemplateRender(w, r, "base.html", "authorize.html", data)
}
