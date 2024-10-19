package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/7nolikov/Jobstar/internal/db"
	"github.com/7nolikov/Jobstar/internal/models"
	"github.com/7nolikov/Jobstar/internal/templates"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

var ValidTransitions = map[string][]string{
	"Applied":      {"Interviewing", "Rejected"},
	"Interviewing": {"Offered", "Rejected"},
	"Offered":      {"Hired", "Rejected"},
	"Hired":        {},
	"Rejected":     {},
}

func CanTransition(currentState, newState string) bool {
	for _, s := range ValidTransitions[currentState] {
		if s == newState {
			return true
		}
	}
	return false
}

// ListCandidates handles GET /candidates
func ListCandidates(w http.ResponseWriter, r *http.Request) {
	var candidates []models.Candidate
	err := db.DB.Select(&candidates, "SELECT * FROM candidates ORDER BY id DESC")
	if err != nil {
		log.Printf("Error fetching candidates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	templates.Render(w, "candidates", candidates)
}

// NewCandidateForm handles GET /candidates/new
func NewCandidateForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"csrfField": csrf.TemplateField(r),
	}
	templates.Render(w, "partials/add_candidate_form", data)
}

// CreateCandidate handles POST /candidates
func CreateCandidate(w http.ResponseWriter, r *http.Request) {
	var candidate models.Candidate
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	candidate.FirstName = r.FormValue("first_name")
	candidate.LastName = r.FormValue("last_name")
	candidate.Email = r.FormValue("email")
	candidate.Phone = r.FormValue("phone")
	candidate.Resume = r.FormValue("resume")

	query := `INSERT INTO candidates (first_name, last_name, email, phone, resume, created_at, updated_at)
	          VALUES (:first_name, :last_name, :email, :phone, :resume, NOW(), NOW()) RETURNING id`

	stmt, err := db.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = stmt.QueryRowx(candidate).Scan(&candidate.ID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch the newly created candidate
	err = db.DB.Get(&candidate, "SELECT * FROM candidates WHERE id=$1", candidate.ID)
	if err != nil {
		log.Printf("Error fetching candidate after creation: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Candidate created with ID: %d", candidate.ID)

	// Check if request is from HTMX
	if r.Header.Get("HX-Request") == "true" {
		// Set a custom HTTP header to trigger modal closure
		w.Header().Set("HX-Trigger", "closeModal")

		// Return the candidate item partial
		templates.Render(w, "partials/candidate_item", candidate)
	} else {
		// Redirect to candidates page
		http.Redirect(w, r, "/candidates", http.StatusSeeOther)
	}
}

// EditCandidateForm handles GET /candidates/{id}/edit
func EditCandidateForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var candidate models.Candidate
	err := db.DB.Get(&candidate, "SELECT * FROM candidates WHERE id=$1", id)
	if err != nil {
		http.Error(w, "Candidate not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Candidate": candidate,
		"csrfField": csrf.TemplateField(r),
	}

	templates.Render(w, "partials/edit_candidate_form", data)
}

// UpdateCandidate handles PUT /candidates/{id}
func UpdateCandidate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var candidate models.Candidate
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	candidateID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid candidate ID", http.StatusBadRequest)
		return
	}
	candidate.ID = candidateID
	candidate.FirstName = r.FormValue("first_name")
	candidate.LastName = r.FormValue("last_name")
	candidate.Email = r.FormValue("email")
	candidate.Phone = r.FormValue("phone")
	candidate.Resume = r.FormValue("resume")

	query := `UPDATE candidates SET first_name=:first_name, last_name=:last_name, email=:email, phone=:phone, resume=:resume, updated_at=NOW() WHERE id=:id`

	_, err = db.DB.NamedExec(query, candidate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the updated candidate
	err = db.DB.Get(&candidate, "SELECT * FROM candidates WHERE id=$1", id)
	if err != nil {
		http.Error(w, "Candidate not found after update", http.StatusNotFound)
		return
	}

	// Check if request is from HTMX
	if r.Header.Get("HX-Request") == "true" {
		// Return the updated candidate item partial
		templates.Render(w, "partials/candidate_item", candidate)
	} else {
		// Redirect to candidates page
		http.Redirect(w, r, "/candidates", http.StatusSeeOther)
	}
}

// DeleteCandidate handles DELETE /candidates/{id}
func DeleteCandidate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.DB.Exec("DELETE FROM candidates WHERE id=$1", id)
	if err != nil {
		log.Printf("Error deleting candidate with ID %s: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Candidate deleted with ID: %s", id)

	// Check if request is from HTMX
	if r.Header.Get("HX-Request") == "true" {
		// Return an empty response to remove the candidate item
		w.WriteHeader(http.StatusNoContent)
	} else {
		// Redirect to candidates page
		http.Redirect(w, r, "/candidates", http.StatusSeeOther)
	}
}
