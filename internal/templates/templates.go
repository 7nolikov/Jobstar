// internal/templates/templates.go

package templates

import (
	"embed"
	"log"
	"net/http"
	"sync"

	"github.com/kataras/blocks"
)

// Embed all HTML templates.
//go:embed views/*
var templatesFS embed.FS

var (
	templateEngine *blocks.Blocks
	once           sync.Once
)

// Init initializes the blocks engine and parses all templates.
// It ensures that initialization happens only once.
func Init() {
	once.Do(func() {
		templateEngine = blocks.New(templatesFS).RootDir("views").Reload(true)

		err := templateEngine.Load()
		if err != nil {
			log.Fatalf("Failed to parse embedded templates: %v", err)
		}
	})
}

// Render renders a specified template with provided data.
func Render(w http.ResponseWriter, tmpl string, data interface{}) {

	// Set the response header content type.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Render the template.
	err := templateEngine.ExecuteTemplate(w, tmpl, "main", data)
	if err != nil {
		log.Printf("Error rendering template '%s': %v", tmpl, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
