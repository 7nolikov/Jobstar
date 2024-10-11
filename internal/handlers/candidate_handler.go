package handlers

import (
    "encoding/json"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "github.com/7nolikov/Jobstar/internal/db"
    "github.com/7nolikov/Jobstar/internal/models"
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

func CreateCandidate(w http.ResponseWriter, r *http.Request) {
    var candidate models.Candidate
    err := json.NewDecoder(r.Body).Decode(&candidate)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    candidate.State = "Applied" // Initial state

    query := `INSERT INTO candidates (name, email, vacancy_id, state, created_at, updated_at)
              VALUES (:name, :email, :vacancy_id, :state, NOW(), NOW()) RETURNING id`

    stmt, err := db.DB.PrepareNamed(query)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    err = stmt.QueryRowx(candidate).Scan(&candidate.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(candidate)
}

func ListCandidates(w http.ResponseWriter, r *http.Request) {
    var candidates []models.Candidate
    err := db.DB.Select(&candidates, "SELECT * FROM candidates")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(candidates)
}

func UpdateCandidateState(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    var input struct {
        NewState string `json:"new_state"`
    }
    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    var candidate models.Candidate
    err = db.DB.Get(&candidate, "SELECT * FROM candidates WHERE id=$1", id)
    if err != nil {
        http.Error(w, "Candidate not found", http.StatusNotFound)
        return
    }

    if !CanTransition(candidate.State, input.NewState) {
        http.Error(w, "Invalid state transition", http.StatusBadRequest)
        return
    }

    candidate.State = input.NewState
    candidate.UpdatedAt = time.Now()

    query := `UPDATE candidates SET state=:state, updated_at=NOW() WHERE id=:id`
    _, err = db.DB.NamedExec(query, candidate)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}
