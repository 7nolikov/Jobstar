package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/7nolikov/Jobstar/internal/db"
	"github.com/7nolikov/Jobstar/internal/handlers"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Continuing with environment variables.")
	}

	// Initialize the database connection
	db.InitDB()

	// Run database migrations
	db.RunMigrations()

	// Set up the router
	r := mux.NewRouter()

	// Serve static files
	staticDir := "/static/"
	r.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir(filepath.Join(".", "static")))))

	// Define your routes
	r.HandleFunc("/vacancies", handlers.ListVacancies).Methods("GET")
	r.HandleFunc("/vacancies/new", handlers.NewVacancyForm).Methods("GET")
	r.HandleFunc("/vacancies", handlers.CreateVacancy).Methods("POST")
	r.HandleFunc("/vacancies/{id}/edit", handlers.EditVacancyForm).Methods("GET")
	r.HandleFunc("/vacancies/{id}", handlers.UpdateVacancy).Methods("PUT")
	r.HandleFunc("/vacancies/{id}", handlers.DeleteVacancy).Methods("DELETE")

	// Define the root route to redirect to /vacancies
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/vacancies", http.StatusSeeOther)
	})

	csrfKey := os.Getenv("CSRF_KEY")
	// Initialize CSRF protection
	csrfMiddleware := csrf.Protect(
		[]byte(csrfKey), // Replace with your secure key
		csrf.Secure(false),              // Set to true in production (requires HTTPS)
	)

	// Start the server
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", csrfMiddleware(r)))
}
