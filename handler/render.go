package handler

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
)

func TemplateRender(w http.ResponseWriter, r *http.Request, baseTemplate string, contentTemplate string, data interface{}) {
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
