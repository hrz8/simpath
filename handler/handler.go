package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/hrz8/simpath/internal/authcode"
	"github.com/hrz8/simpath/internal/client"
	"github.com/hrz8/simpath/internal/consent"
	"github.com/hrz8/simpath/internal/introspect"
	"github.com/hrz8/simpath/internal/scope"
	"github.com/hrz8/simpath/internal/token"
	"github.com/hrz8/simpath/internal/tokengrant"
	"github.com/hrz8/simpath/internal/user"
	"github.com/hrz8/simpath/session"
)

type contextKey int

const (
	clientKey contextKey = iota
	userDataKey
)

var (
	ErrClientNotPresent      = errors.New("Client not present in the request context")
	ErrUserSessionNotPresent = errors.New("User session not present in the request context")
)

type Handler struct {
	DB            *sql.DB
	sessionSvc    *session.Service
	userSvc       *user.Service
	clientSvc     *client.Service
	scopeSvc      *scope.Service
	tokenSvc      *token.Service
	authCodeSvc   *authcode.Service
	tokenGrantSvc *tokengrant.Service
	introspectSvc *introspect.Service
	consentSvc    *consent.Service
}

func NewHandler(
	db *sql.DB,
	sesSvc *session.Service,
	uSvc *user.Service,
	cSvc *client.Service,
	sSvc *scope.Service,
	tSvc *token.Service,
	acSvc *authcode.Service,
	tgSvc *tokengrant.Service,
	iSvc *introspect.Service,
	conSvc *consent.Service,
) *Handler {
	return &Handler{
		db,
		sesSvc,
		uSvc,
		cSvc,
		sSvc,
		tSvc,
		acSvc,
		tgSvc,
		iSvc,
		conSvc,
	}
}

func getQueryString(query url.Values) string {
	encoded := query.Encode()
	if len(encoded) > 0 {
		encoded = fmt.Sprintf("?%s", encoded)
	}
	return encoded
}

func templateRender(w http.ResponseWriter, _ *http.Request, baseTemplate string, contentTemplate string, data any) {
	baseTemplatePath := filepath.Join("templates", baseTemplate)
	contentTemplatePath := filepath.Join("templates", "partials", contentTemplate)

	tmpl, err := template.ParseFiles(baseTemplatePath, contentTemplatePath)
	if err != nil {
		http.Error(w, "unable to load template", http.StatusInternalServerError)
		log.Println("error parsing templates:", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "unable to render template", http.StatusInternalServerError)
		log.Println("error executing template:", err)
	}
}

func TemplateRenderNoBase(w http.ResponseWriter, _ *http.Request, templateName string, data any) {
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

func getClient(ctx context.Context) (*client.OauthClient, error) {
	cli, ok := ctx.Value(clientKey).(*client.OauthClient)
	if !ok {
		return nil, ErrClientNotPresent
	}

	return cli, nil
}

func getUserDataFromSession(ctx context.Context) (*session.UserData, error) {
	cli, ok := ctx.Value(userDataKey).(*session.UserData)
	if !ok {
		return nil, ErrUserSessionNotPresent
	}

	return cli, nil
}
