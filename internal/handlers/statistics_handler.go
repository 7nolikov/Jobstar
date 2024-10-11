package handlers

import (
    "bytes"
	"html/template"
	"net/http"

	"github.com/7nolikov/Jobstar/internal/db"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func GetPipelineStatistics() (map[string]int, error) {
	stats := make(map[string]int)
	rows, err := db.DB.Queryx("SELECT state, COUNT(*) as count FROM candidates GROUP BY state")
	if err != nil {
		return stats, err
	}
	defer rows.Close()

	for rows.Next() {
		var state string
		var count int
		err = rows.Scan(&state, &count)
		if err != nil {
			return stats, err
		}
		stats[state] = count
	}
	return stats, nil
}

func ShowStatistics(w http.ResponseWriter, r *http.Request) {
	stats, err := GetPipelineStatistics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pie := charts.NewPie()
	pie.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: "Candidate Pipeline Statistics",
	}))

	var items []opts.PieData
	for state, count := range stats {
		items = append(items, opts.PieData{Name: state, Value: count})
	}

	pie.AddSeries("States", items)

	tmpl, err := template.ParseFiles("templates/statistics.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a buffer to hold the rendered pie chart
    var buf bytes.Buffer
    err = pie.Render(&buf)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Execute the template with the rendered pie chart
    err = tmpl.Execute(w, template.HTML(buf.String()))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
