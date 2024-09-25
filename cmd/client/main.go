package main

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"text/template"

	"github.com/go-chi/chi/v5"
)

func templateRender(w http.ResponseWriter, _ *http.Request, templateName string, data any) {
	templatePath := filepath.Join("templates", templateName)
	funcMap := template.FuncMap{
		"sprintf": fmt.Sprintf,
	}

	tmpl, err := template.New("client.html").Funcs(funcMap).ParseFiles(templatePath)
	if err != nil {
		http.Error(w, "unable to load template", http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "unable to render template", http.StatusInternalServerError)
	}
}

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
			templateRender(w, r, "client_login.html", data)
			return
		}

		data := map[string]any{
			"isCode":         false,
			"loginURL":       "http://localhost:5001/v1/authorize",
			"clientID":       "600ef080-d02c-426d-bf79-64247ba0fc90",
			"clientSecret":   "test_secret",
			"redirectURIRaw": "http://localhost:8088/signin",
			"redirectURI":    url.QueryEscape("http://localhost:8088/signin"),
			"scope":          "read_write",
			"state":          "somestate",
			"btnLabel":       "Login dengan Simpath",
		}
		templateRender(w, r, "client_login.html", data)
	})
	srv := &http.Server{Addr: ":8088", Handler: mux}
	srv.ListenAndServe()
}
