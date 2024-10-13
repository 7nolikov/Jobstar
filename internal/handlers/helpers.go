package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// Template cache
var templates *template.Template

func init() {
	// Parse base templates and partials
	templates = template.Must(template.ParseGlob(filepath.Join("templates", "*.html")))
	templates = template.Must(templates.ParseGlob(filepath.Join("templates", "partials", "*.html")))
}

// RenderTemplate renders a specified template with data
func RenderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    err := templates.ExecuteTemplate(w, tmpl, data)
    if err != nil {
        log.Printf("Error executing template '%s': %v", tmpl, err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}
