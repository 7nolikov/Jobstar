package handlers

import (
	"log"
	"net/http"

	"github.com/7nolikov/Jobstar/internal/db"
	"github.com/7nolikov/Jobstar/internal/templates"
	"github.com/gorilla/csrf"
)

// StatisticsData holds the data for the Statistics page
type StatisticsData struct {
	TotalVacancies    int
	TotalCandidates   int
	FilledVacancies   int
	PendingCandidates int
}

// StatisticsHandler handles GET /statistics
func StatisticsHandler(w http.ResponseWriter, r *http.Request) {
	var data StatisticsData

	// Fetch Total Vacancies
	err := db.DB.Get(&data.TotalVacancies, "SELECT COUNT(*) FROM vacancies")
	if err != nil {
		log.Printf("Error fetching total vacancies: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch Total Candidates
	err = db.DB.Get(&data.TotalCandidates, "SELECT COUNT(*) FROM candidates")
	if err != nil {
		log.Printf("Error fetching total candidates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch Filled Vacancies (assuming a 'filled' boolean column)
	err = db.DB.Get(&data.FilledVacancies, "SELECT COUNT(*) FROM vacancies WHERE filled = TRUE")
	if err != nil {
		log.Printf("Error fetching filled vacancies: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Fetch Pending Candidates (assuming a 'status' column)
	err = db.DB.Get(&data.PendingCandidates, "SELECT COUNT(*) FROM candidates WHERE status = 'pending'")
	if err != nil {
		log.Printf("Error fetching pending candidates: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Pass statistics data to the template
	dataMap := map[string]interface{}{
		"Statistics": data,
		"csrfToken":  csrf.Token(r),
	}

	templates.Render(w, "statistics.html", dataMap)
}
