package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hrz8/simpath/session"
)

func redirectToURL(w http.ResponseWriter, r *http.Request, uri *url.URL, query url.Values) {
	to := fmt.Sprintf("%s%s", uri.String(), getQueryString(query))
	http.Redirect(w, r, to, http.StatusFound)
}

func redirectSelf(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, r.RequestURI, http.StatusFound)
}

func redirectLogin(w http.ResponseWriter, r *http.Request) {
	to := fmt.Sprintf("/v1/login%s", getQueryString(r.URL.Query()))
	http.Redirect(w, r, to, http.StatusFound)
}

func redirectAuthorize(w http.ResponseWriter, r *http.Request) {
	to := fmt.Sprintf("/v1/authorize%s", getQueryString(r.URL.Query()))
	http.Redirect(w, r, to, http.StatusFound)
}

func redirectError(w http.ResponseWriter, r *http.Request, redirectURI *url.URL, err, state string) {
	query := redirectURI.Query()
	query.Set("error", err)
	if state != "" {
		query.Set("state", state)
	}

	to := redirectURI.String()
	http.Redirect(w, r, fmt.Sprintf("%s%s", to, getQueryString(query)), http.StatusFound)
}

func (h *Handler) authenticate(userSession *session.UserSession) error {
	// try to authenticate with access token
	_, err := h.tokenSvc.Authenticate(userSession.AccessToken)
	if err == nil {
		return nil
	}

	// got error okay maybe it's refresh token
	// fetch user's client first
	client, err := h.clientSvc.FindClientByClientUUID(userSession.ClientUUID)
	if err != nil {
		return err
	}

	// check refresh token validity (might be expired already)
	refreshToken, err := h.tokenSvc.GetValidRefreshToken(userSession.RefreshToken, client.ID)
	if err != nil {
		return err
	}

	// login to create access token and refresh token for user
	at, rt, err := h.tokenSvc.Login(
		refreshToken.ClientID,
		refreshToken.UserID,
		refreshToken.Scope,
	)
	if err != nil {
		return err
	}

	// mutate userSession with latest token
	userSession.AccessToken = at.AccessToken
	userSession.RefreshToken = rt.RefreshToken

	return nil
}

func (h *Handler) ShouldHaveClientID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// parse so we can perform r.Form.Get()
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clientID := r.Form.Get("client_id")
		client, err := h.clientSvc.FindClientByClientUUID(clientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ctx = context.WithValue(ctx, clientKey, client)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) LoggedInOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// start cookie session
		h.sessionSvc.SetSessionService(w, r)
		if err := h.sessionSvc.StartSession(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// fetch user's session
		userSession, err := h.sessionSvc.GetUserSession()
		if err != nil {
			redirectLogin(w, r)
			return
		}
		ctx = context.WithValue(ctx, userSessionKey, userSession) // set user session to context

		// authenticate (possibly mutate the userSession if refresh token expired)
		err = h.authenticate(userSession)
		if err != nil {
			redirectLogin(w, r)
			return
		}

		h.sessionSvc.SetUserSession(userSession)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) GuestOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// start cookie session
		h.sessionSvc.SetSessionService(w, r)
		if err := h.sessionSvc.StartSession(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// prevent logged-in user to access the login page again
		_, err := h.sessionSvc.GetUserSession()
		if err == nil {
			redirectAuthorize(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
