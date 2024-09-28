package main

import (
	"context"
	"encoding/json"
	"errors"

	"log"
	"net/http"
	"net/url"

	"github.com/hrz8/simpath/handler"
	"golang.org/x/oauth2"
)

func addLoginRedirectURI(authURL string) (string, error) {
	parsedURL, err := url.Parse(authURL)
	if err != nil {
		return "", errors.New("Error parsing URL")
	}
	query := parsedURL.Query()
	query.Set("login_redirect_uri", "/v1/oauth2/authorize")
	parsedURL.RawQuery = query.Encode()

	return parsedURL.String(), nil
}

func SimpathLogin(simpathOauth *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authURL := simpathOauth.AuthCodeURL("somestate")
		uri, err := addLoginRedirectURI(authURL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, uri, http.StatusSeeOther)
	}
}

func SimpathCallback(simpathOauth *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")
		if state != "somestate" {
			http.Error(w, "States don't Match!!", http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")

		token, err := simpathOauth.Exchange(context.Background(), code)
		if err != nil {
			http.Error(w, "Code-Token Exchange Failed", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(token); err != nil {
			http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	simpathOauth := &oauth2.Config{
		RedirectURL:  "http://localhost:8089/simpath/callback",
		ClientID:     "600ef080-d02c-426d-bf79-64247ba0fc90",
		ClientSecret: "test_secret",
		Scopes:       []string{"read_write", "openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://localhost:5001/v1/oauth2/authorize",
			TokenURL: "http://localhost:5001/v1/oauth2/token",
		},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler.TemplateRenderNoBase(w, r, "oauthclient_login.html", map[string]any{})
	})
	http.HandleFunc("/simpath/login", SimpathLogin(simpathOauth))
	http.HandleFunc("/simpath/callback", SimpathCallback(simpathOauth))

	log.Println("Server is running on http://localhost:8089")
	log.Fatal(http.ListenAndServe(":8089", nil))
}
