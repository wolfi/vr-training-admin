// internal/handlers/dashboard.go
package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/saladinomario/vr-training-admin/templates/pages"
)

// DashboardHandler renders the main dashboard page
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Only handle the root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Only handle GET requests
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
		return
	}

	// Generate the dashboard component
	component := pages.Dashboard()

	// Set content type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Render the template
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering dashboard: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
