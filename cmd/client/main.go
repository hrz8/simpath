package main

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/hrz8/simpath/handler"
)

func main() {
	mux := chi.NewRouter()
	mux.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		code := queryParams.Get("code")
		if code != "" {
			data := map[string]any{
				"isCode":         true,
				"code":           code,
				"clientID":       "600ef080-d02c-426d-bf79-64247ba0fc90",
				"clientSecret":   "test_secret",
				"redirectURIRaw": "http://localhost:8088/signin",
			}
			handler.TemplateRenderNoBase(w, r, "client_login.html", data)
			return
		}

		data := map[string]any{
			"isCode":           false,
			"loginURL":         "http://localhost:5001/v1/oauth2/authorize",
			"clientID":         "600ef080-d02c-426d-bf79-64247ba0fc90",
			"clientSecret":     "test_secret",
			"loginRedirectURI": url.QueryEscape("/v1/oauth2/authorize"),
			"redirectURIRaw":   "http://localhost:8088/signin",
			"redirectURI":      url.QueryEscape("http://localhost:8088/signin"),
			"scope":            url.QueryEscape("read_write openid"),
			"state":            "somestate",
			"btnLabel":         "Login dengan Simpath",
		}
		handler.TemplateRenderNoBase(w, r, "client_login.html", data)
	})
	srv := &http.Server{Addr: ":8088", Handler: mux}
	srv.ListenAndServe()
}
