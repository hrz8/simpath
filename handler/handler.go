package handler

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	"github.com/hrz8/simpath/internal/client"
	"github.com/hrz8/simpath/internal/scope"
	"github.com/hrz8/simpath/internal/token"
	"github.com/hrz8/simpath/internal/user"
	"github.com/hrz8/simpath/session"
)

type Handler struct {
	DB         *sql.DB
	sessionSvc *session.Service
	userSvc    *user.Service
	clientSvc  *client.Service
	scopeSvc   *scope.Service
	tokenSvc   *token.Service
}

func NewHandler(
	db *sql.DB,
	sesSvc *session.Service,
	uSvc *user.Service,
	cSvc *client.Service,
	sSvc *scope.Service,
	tSvc *token.Service,
) *Handler {
	return &Handler{
		db,
		sesSvc,
		uSvc,
		cSvc,
		sSvc,
		tSvc,
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
