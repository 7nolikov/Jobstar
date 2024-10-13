package handlers

import (
	"net/http"
	"strconv"

	"github.com/7nolikov/Jobstar/internal/db"
	"github.com/7nolikov/Jobstar/internal/models"
	"github.com/gorilla/mux"
)

// ListVacancies handles GET /vacancies
func ListVacancies(w http.ResponseWriter, r *http.Request) {
    var vacancies []models.Vacancy
    err := db.DB.Select(&vacancies, "SELECT * FROM vacancies ORDER BY id DESC")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Check if request is from HTMX
    if r.Header.Get("HX-Request") == "true" {
        // Return the entire vacancies list
        RenderTemplate(w, "vacancies", map[string]interface{}{
            "Vacancies": vacancies,
        })
    } else {
        // Render the base template with vacancies
        RenderTemplate(w, "vacancies", map[string]interface{}{
            "Vacancies": vacancies,
        })
    }
}

// NewVacancyForm handles GET /vacancies/new
func NewVacancyForm(w http.ResponseWriter, r *http.Request) {
    // Render the add vacancy form partial
    RenderTemplate(w, "add_vacancy_form", nil)
}

// CreateVacancy handles POST /vacancies
func CreateVacancy(w http.ResponseWriter, r *http.Request) {
    var vacancy models.Vacancy
    err := r.ParseForm()
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    vacancy.Title = r.FormValue("title")
    vacancy.Description = r.FormValue("description")
    vacancy.Location = r.FormValue("location")

    query := `INSERT INTO vacancies (title, description, location, created_at, updated_at)
              VALUES (:title, :description, :location, NOW(), NOW()) RETURNING id`

    stmt, err := db.DB.PrepareNamed(query)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    err = stmt.QueryRowx(vacancy).Scan(&vacancy.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Fetch the newly created vacancy
    err = db.DB.Get(&vacancy, "SELECT * FROM vacancies WHERE id=$1", vacancy.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Check if request is from HTMX
    if r.Header.Get("HX-Request") == "true" {
        // Return the vacancy item partial
        RenderTemplate(w, "vacancy_item", vacancy)
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

    // Render the edit vacancy form partial
    RenderTemplate(w, "edit_vacancy_form", vacancy)
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
        http.Error(w, "Vacancy not found", http.StatusNotFound)
        return
    }

    // Check if request is from HTMX
    if r.Header.Get("HX-Request") == "true" {
        // Return the updated vacancy item partial
        RenderTemplate(w, "vacancy_item", vacancy)
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
