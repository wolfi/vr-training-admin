// internal/handlers/dashboard.go
package handlers

import (
	"log"
	"net/http"

	"github.com/saladinomario/vr-training-admin/templates/pages"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure we're only handling the root path
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Set content type header
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Render the dashboard
	component := pages.Dashboard()
	err := component.Render(r.Context(), w)
	if err != nil {
		log.Printf("Error rendering dashboard: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// DashboardContentHandler handles the AJAX request for dashboard content
func DashboardContentHandler(w http.ResponseWriter, r *http.Request) {
	// Get recent sessions and render dashboard content
	recentSessions := SessionStore.GetRecent(5)

	// Ensure we're using the proper capitalized store variables
	_ = ScenarioStore
	_ = AvatarStore
	_ = ObserverStore
	component := pages.DashboardContent(recentSessions)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := component.Render(r.Context(), w)
	if err != nil {
		log.Printf("Error rendering dashboard content: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
