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

// ListVacancies handles GET /vacancies
func ListVacancies(w http.ResponseWriter, r *http.Request) {
	var vacancies []models.Vacancy
	err := db.DB.Select(&vacancies, "SELECT * FROM vacancies ORDER BY id DESC")
	if err != nil {
		log.Printf("Error fetching vacancies: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Vacancies": vacancies,
		"csrfToken": csrf.Token(r),
	}

	templates.RenderTemplate(w, "vacancies.html", data)
}

// NewVacancyForm handles GET /vacancies/new
func NewVacancyForm(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"csrfField": csrf.TemplateField(r),
	}
	templates.RenderTemplate(w, "add_vacancy_form", data)
}

// CreateVacancy handles POST /vacancies
func CreateVacancy(w http.ResponseWriter, r *http.Request) {
	var vacancy models.Vacancy
	err := r.ParseForm()
	if err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	vacancy.Title = r.FormValue("title")
	vacancy.Description = r.FormValue("description")
	vacancy.Location = r.FormValue("location")

	query := `INSERT INTO vacancies (title, description, location, created_at, updated_at)
              VALUES (:title, :description, :location, NOW(), NOW()) RETURNING id`

	stmt, err := db.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = stmt.QueryRowx(vacancy).Scan(&vacancy.ID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch the newly created vacancy
	err = db.DB.Get(&vacancy, "SELECT * FROM vacancies WHERE id=$1", vacancy.ID)
	if err != nil {
		log.Printf("Error fetching vacancy after creation: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	log.Printf("Vacancy created with ID: %d", vacancy.ID)

	// Check if request is from HTMX
	if r.Header.Get("HX-Request") == "true" {
		// Set a custom HTTP header to trigger modal closure
		w.Header().Set("HX-Trigger", "closeModal")

		// Return the vacancy item partial
		templates.RenderTemplate(w, "vacancy_item", vacancy)
	} else {
		// Redirect to vacancies page
		http.Redirect(w, r, "/vacancies", http.StatusSeeOther)
	}
}

// EditVacancyForm handles GET /vacancies/{id}/edit
func EditVacancyForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var vacancy models.Vacancy
	err := db.DB.Get(&vacancy, "SELECT * FROM vacancies WHERE id=$1", id)
	if err != nil {
		http.Error(w, "Vacancy not found", http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"Vacancy":   vacancy,
		"csrfField": csrf.TemplateField(r),
	}

	templates.RenderTemplate(w, "edit_vacancy_form", data)
}

// UpdateVacancy handles PUT /vacancies/{id}
func UpdateVacancy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var vacancy models.Vacancy
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vacancyID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid vacancy ID", http.StatusBadRequest)
		return
	}
	vacancy.ID = vacancyID
	vacancy.Title = r.FormValue("title")
	vacancy.Description = r.FormValue("description")
	vacancy.Location = r.FormValue("location")

	query := `UPDATE vacancies SET title=:title, description=:description, location=:location, updated_at=NOW() WHERE id=:id`

	_, err = db.DB.NamedExec(query, vacancy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the updated vacancy
	err = db.DB.Get(&vacancy, "SELECT * FROM vacancies WHERE id=$1", id)
	if err != nil {
		http.Error(w, "Vacancy not found after update", http.StatusNotFound)
		return
	}

	// Check if request is from HTMX
	if r.Header.Get("HX-Request") == "true" {
		// Return the updated vacancy item partial
		templates.RenderTemplate(w, "vacancy_item", vacancy)
	} else {
		// Redirect to vacancies page
		http.Redirect(w, r, "/vacancies", http.StatusSeeOther)
	}
}

// DeleteVacancy handles DELETE /vacancies/{id}
func DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.DB.Exec("DELETE FROM vacancies WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if request is from HTMX
	if r.Header.Get("HX-Request") == "true" {
		// Return an empty response to remove the vacancy item
		w.WriteHeader(http.StatusNoContent)
	} else {
		// Redirect to vacancies page
		http.Redirect(w, r, "/vacancies", http.StatusSeeOther)
	}
}
