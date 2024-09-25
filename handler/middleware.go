package handler

import (
	"context"
	"net/http"

	"github.com/hrz8/simpath/session"
)

func (h *Handler) authenticate(userSession *session.UserData) error {
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

// Use to setup global session like csrf token
func (h *Handler) UseSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		h.sessionSvc.SetSessionService(w, r)
		if err := h.sessionSvc.StartSession(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := h.sessionSvc.SetCSRFToken(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Use to start user data session, so in the handler we can check if the user data already set or not
func (h *Handler) UseUserSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		h.sessionSvc.SetSessionService(w, r)
		if err := h.sessionSvc.StartUserSession(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Use to parse form so we can perform r.Form.Get()
func (h *Handler) UseForm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) CheckClientID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
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

		// fetch user's data session
		userData, err := h.sessionSvc.GetUserData()
		if err != nil {
			redirectLogin(w, r)
			return
		}
		ctx = context.WithValue(ctx, userDataKey, userData) // set user session to context

		// authenticate (possibly mutate the userSession if refresh token expired)
		err = h.authenticate(userData)
		if err != nil {
			redirectLogin(w, r)
			return
		}
		h.sessionSvc.SetUserData(userData)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) GuestOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// prevent logged-in user to access the login page again
		_, err := h.sessionSvc.GetUserData()
		if err == nil {
			redirectAuthorize(w, r)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) CheckCSRFToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		csrfToken := r.Form.Get("csrf_token")
		csrfTokenFromSession, err := h.sessionSvc.GetCSRFToken()
		if err != nil || csrfTokenFromSession == "" {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if csrfToken != csrfTokenFromSession {
			h.sessionSvc.SetFlashMessage("Invalid or expired csrf token")
			redirectSelf(w, r)
			return
		}
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
