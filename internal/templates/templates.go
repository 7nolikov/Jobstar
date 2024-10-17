package templates

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// TemplateCache holds the parsed templates
var TemplateCache *template.Template
var once sync.Once

// Init parses all templates and caches them with detailed logging
func Init() {
	once.Do(func() {
		var err error
		// Initialize a new template
		TemplateCache = template.New("")

		// Walk through the templates directory and parse each .html file
		err = filepath.Walk("templates", func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".html" {
				relPath, err := filepath.Rel("templates", path)
				if err != nil {
					return err
				}
				// Read the file content
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				// Parse the template file
				_, err = TemplateCache.New(relPath).Parse(string(content))
				if err != nil {
					log.Printf("Error parsing template %s: %v", relPath, err)
					return err
				}
				log.Printf("Loaded template: %s", relPath)
			}
			return nil
		})

		if err != nil {
			log.Fatalf("Error walking through templates: %v", err)
		}

		log.Println("All templates successfully loaded and cached.")
	})
}

// RenderTemplate executes the specified template
func RenderTemplate(w http.ResponseWriter, name string, data interface{}) {
	tmpl := TemplateCache.Lookup(name)
	if tmpl == nil {
		log.Printf("Template '%s' not found", name)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err := tmpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template '%s': %v", name, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
