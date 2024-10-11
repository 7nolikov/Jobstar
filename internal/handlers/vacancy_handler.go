package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/7nolikov/Jobstar/internal/db"
	"github.com/7nolikov/Jobstar/internal/models"
	"github.com/gorilla/mux"
)

func CreateVacancy(w http.ResponseWriter, r *http.Request) {
	var vacancy models.Vacancy
	err := json.NewDecoder(r.Body).Decode(&vacancy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vacancy)
}

func ListVacancies(w http.ResponseWriter, r *http.Request) {
	var vacancies []models.Vacancy
	err := db.DB.Select(&vacancies, "SELECT * FROM vacancies")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vacancies)
}

func UpdateVacancy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var vacancy models.Vacancy
	err := json.NewDecoder(r.Body).Decode(&vacancy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vacancy.ID, _ = strconv.Atoi(id)

	query := `UPDATE vacancies SET title=:title, description=:description, location=:location, updated_at=NOW() WHERE id=:id`
	_, err = db.DB.NamedExec(query, vacancy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteVacancy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := db.DB.Exec("DELETE FROM vacancies WHERE id=$1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
