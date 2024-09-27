package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/hrz8/simpath/internal/client"
)

func toSnakeCase(str string) string {
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(snake)
}

func redirectToURL(w http.ResponseWriter, r *http.Request, uri *url.URL, query url.Values) {
	to := fmt.Sprintf("%s%s", uri.String(), getQueryString(query))
	http.Redirect(w, r, to, http.StatusFound)
}

func redirectToPath(w http.ResponseWriter, r *http.Request, path string) {
	to := fmt.Sprintf("%s%s", path, getQueryString(r.URL.Query()))
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
	to := fmt.Sprintf("/v1/oauth2/authorize%s", getQueryString(r.URL.Query()))
	http.Redirect(w, r, to, http.StatusFound)
}

func redirectDashboard(w http.ResponseWriter, r *http.Request) {
	to := fmt.Sprintf("/v1/dashboard%s", getQueryString(r.URL.Query()))
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

func getRedirectUri(r *http.Request, cli *client.OauthClient) (*url.URL, error) {
	redirectURI := r.Form.Get("redirect_uri")
	if redirectURI == "" {
		redirectURI = cli.RedirectURI
	}

	// parse the redirect URL
	parsedRedirectURI, err := url.ParseRequestURI(redirectURI)
	if err != nil {
		return nil, err
	}

	return parsedRedirectURI, nil
}
