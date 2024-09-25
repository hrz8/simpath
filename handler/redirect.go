package handler

import (
	"fmt"
	"net/http"
	"net/url"
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
