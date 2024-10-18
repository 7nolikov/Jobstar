package handlers

import (
    "github.com/7nolikov/Jobstar/internal/templates"
    "github.com/gorilla/csrf"
    "net/http"
)

// LandingPage handles GET /
func LandingPage(w http.ResponseWriter, r *http.Request) {
    data := map[string]interface{}{
        "csrfToken": csrf.Token(r),
    }
    templates.RenderTemplate(w, "base", data)
}
