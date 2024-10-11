package main

import (
	"log"
	"net/http"

	"github.com/7nolikov/Jobstar/internal/db"
	"github.com/7nolikov/Jobstar/internal/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database connection
	db.InitDB()

	// Set up the router
	r := mux.NewRouter()

	// Vacancy routes
	r.HandleFunc("/vacancies", handlers.ListVacancies).Methods("GET")
	r.HandleFunc("/vacancies", handlers.CreateVacancy).Methods("POST")
	r.HandleFunc("/vacancies/{id}", handlers.UpdateVacancy).Methods("PUT")
	r.HandleFunc("/vacancies/{id}", handlers.DeleteVacancy).Methods("DELETE")

	// Candidate routes
	r.HandleFunc("/candidates", handlers.ListCandidates).Methods("GET")
	r.HandleFunc("/candidates", handlers.CreateCandidate).Methods("POST")
	r.HandleFunc("/candidates/{id}/state", handlers.UpdateCandidateState).Methods("POST")

	// Statistics route
	r.HandleFunc("/statistics", handlers.ShowStatistics).Methods("GET")

	// Serve static files (e.g., templates, static assets)
	// r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// Start the server
	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
