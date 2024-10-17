package handlers

import (
	"net/http"

	"github.com/gorilla/csrf"
)

// LandingPage handles GET /
func LandingPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"csrfToken": csrf.Token(r),
	}
	RenderTemplate(w, "landing", data)
}
