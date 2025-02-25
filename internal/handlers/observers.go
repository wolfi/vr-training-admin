// internal/handlers/observers.go
package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/saladinomario/vr-training-admin/internal/models"
	"github.com/saladinomario/vr-training-admin/templates/components/observers"
	"github.com/saladinomario/vr-training-admin/templates/pages"
)

var observerStore *models.ObserverStore

func init() {
	observerStore = models.NewObserverStore()
}

// ObserversHandler handles the observers index page
func ObserversHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/observers" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	allObservers := observerStore.GetAll()
	component := pages.ObserversIndex(allObservers)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering observers page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ObserverNewHandler handles the new observer form
func ObserverNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	component := pages.ObserverNew()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering new observer form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ObserverEditHandler handles the edit observer form
func ObserverEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/observers/edit/")
	if idStr == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	// Get observer by ID
	observer, err := observerStore.GetByID(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	component := pages.ObserverEdit(observer)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering edit observer form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ObserverCreateHandler handles observer creation
func ObserverCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse form values
	observer := parseObserverForm(r)

	// Create observer
	_, err := observerStore.Create(observer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If this is an HTMX request, return the main content
	if r.Header.Get("HX-Request") == "true" {
		allObservers := observerStore.GetAll()
		component := pages.ObserversContent(allObservers)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering observers content: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Standard redirect for non-HTMX requests
	http.Redirect(w, r, "/observers", http.StatusSeeOther)
}

// ObserverUpdateHandler handles observer updates
func ObserverUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/observers/")
	if idStr == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse form values
	observer := parseObserverForm(r)

	// Update observer
	err := observerStore.Update(idStr, observer)
	if err != nil {
		if err == models.ErrObserverNotFound {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// If this is an HTMX request, return the main content
	if r.Header.Get("HX-Request") == "true" {
		allObservers := observerStore.GetAll()
		component := pages.ObserversContent(allObservers)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering observers content: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Standard redirect for non-HTMX requests
	http.Redirect(w, r, "/observers", http.StatusSeeOther)
}

// ObserverDeleteHandler handles observer deletion
func ObserverDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/observers/")
	if idStr == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	// Delete observer
	err := observerStore.Delete(idStr)
	if err != nil {
		if err == models.ErrObserverNotFound {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// If this is an HTMX request, return the updated observer list
	if r.Header.Get("HX-Request") == "true" {
		allObservers := observerStore.GetAll()
		component := observers.ObserverList(allObservers)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering observer list: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Standard redirect for non-HTMX requests
	http.Redirect(w, r, "/observers", http.StatusSeeOther)
}

// ObserverSearchHandler handles observer search
func ObserverSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	foundObservers := observerStore.Search(query)

	component := observers.ObserverList(foundObservers)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering observer search results: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Helper function to parse observer form data
func parseObserverForm(r *http.Request) observers.Observer {
	interventionLevel, _ := strconv.Atoi(r.FormValue("intervention_level"))
	detailLevel, _ := strconv.Atoi(r.FormValue("detail_level"))
	active := r.FormValue("active") == "on"

	// Process triggers
	triggerValues := r.Form["triggers"]
	customTriggers := r.FormValue("custom_triggers")

	var interventionTriggers []string
	// Add selected common triggers
	interventionTriggers = append(interventionTriggers, triggerValues...)

	// Add custom triggers
	if customTriggers != "" {
		customLines := strings.Split(customTriggers, "\n")
		for _, line := range customLines {
			line = strings.TrimSpace(line)
			if line != "" {
				interventionTriggers = append(interventionTriggers, line)
			}
		}
	}

	return observers.Observer{
		Name:                 r.FormValue("name"),
		Description:          r.FormValue("description"),
		FeedbackStyle:        r.FormValue("feedback_style"),
		InterventionLevel:    interventionLevel,
		DetailLevel:          detailLevel,
		FeedbackTone:         r.FormValue("feedback_tone"),
		SuccessMetrics:       r.FormValue("success_metrics"),
		InterventionTriggers: interventionTriggers,
		Active:               active,
	}
}

// SetupObserverRoutes registers all observer-related routes
func SetupObserverRoutes(mux *http.ServeMux) {
	log.Println("Setting up observer routes...")

	// List and Create
	mux.HandleFunc("/observers", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Handling observer request: %s %s", r.Method, r.URL.Path)
		switch r.Method {
		case http.MethodGet:
			ObserversHandler(w, r)
		case http.MethodPost:
			ObserverCreateHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// New form
	mux.HandleFunc("/observers/new", ObserverNewHandler)

	// Search
	mux.HandleFunc("/observers/search", ObserverSearchHandler)

	// Edit form
	mux.HandleFunc("/observers/edit/", ObserverEditHandler)

	// Update and Delete
	mux.HandleFunc("/observers/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut, http.MethodPost:
			ObserverUpdateHandler(w, r)
		case http.MethodDelete:
			ObserverDeleteHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Observer routes registered successfully")
}
