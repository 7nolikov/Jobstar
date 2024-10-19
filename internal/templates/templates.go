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
		var err error

		// Create a new blocks engine.
		templateEngine := blocks.New(templatesFS)
				

		if err != nil {
			log.Fatalf("Failed to create blocks engine: %v", err)
		}

		// Add embedded templates to the engine.
		err = templateEngine.Load()
		if err != nil {
			log.Fatalf("Failed to parse embedded templates: %v", err)
		}

		log.Println("All templates successfully loaded and cached.")
	})
}

// Render renders a specified template with provided data.
func Render(w http.ResponseWriter, tmpl string, data interface{}) {

	// Set the response header content type.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Render the template.
	err := templateEngine.ExecuteTemplate(w, tmpl, "", data)
	if err != nil {
		log.Printf("Error rendering template '%s': %v", tmpl, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
