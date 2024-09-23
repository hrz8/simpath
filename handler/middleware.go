package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hrz8/simpath/session"
)

func redirectLogin(w http.ResponseWriter, r *http.Request) {
	to := fmt.Sprintf("/v1/login%s", getQueryString(r.URL.Query()))
	http.Redirect(w, r, to, http.StatusFound)
}

func redirectAuthorize(w http.ResponseWriter, r *http.Request) {
	to := fmt.Sprintf("/v1/authorize%s", getQueryString(r.URL.Query()))
	http.Redirect(w, r, to, http.StatusFound)
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
	at, rt, err := h.login(
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
