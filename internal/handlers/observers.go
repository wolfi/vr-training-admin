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

// ObserverNewHandler displays the form to create a new observer
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

	// Extract ID from URL: e.g. /observers/edit/observer_20250101010101
	idStr := strings.TrimPrefix(r.URL.Path, "/observers/edit/")
	if idStr == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	obs, err := observerStore.GetByID(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	component := pages.ObserverEdit(obs)

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

	observer := parseObserverForm(r)

	_, err := observerStore.Create(observer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Redirect to the observers page
	http.Redirect(w, r, "/observers", http.StatusSeeOther)
}

// ObserverUpdateHandler handles observer updates
func ObserverUpdateHandler(w http.ResponseWriter, r *http.Request) {
	// We allow PUT or POST for "update" to accommodate HTMX or plain forms
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL: e.g. /observers/observer_20250101010101
	idStr := strings.TrimPrefix(r.URL.Path, "/observers/")
	if idStr == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	observer := parseObserverForm(r)

	err := observerStore.Update(idStr, observer)
	if err != nil {
		if err == models.ErrObserverNotFound {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// Redirect to the observers page
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

// ObserverSearchHandler handles searching for observers
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

// parseObserverForm extracts Observer fields from a request form.
func parseObserverForm(r *http.Request) models.Observer {
	interventionFreq, _ := strconv.Atoi(r.FormValue("intervention_frequency"))

	return models.Observer{
		Name:                        r.FormValue("name"),
		InterruptionTriggers:        r.FormValue("interruption_triggers"),
		FeedbackStyle:               r.FormValue("feedback_style"),
		InterventionFrequency:       interventionFreq,
		KeyPerformanceIndicators:    r.FormValue("kpis"),
		ScoringSystem:               r.FormValue("scoring_system"),
		RequiredAchievements:        r.FormValue("required_achievements"),
		FailedStateDefinitions:      r.FormValue("failed_states"),
		RealTimeFeedbackRules:       r.FormValue("realtime_rules"),
		PostSessionReportTemplate:   r.FormValue("report_template"),
		LearningPointsFocus:         r.FormValue("learning_focus"),
		ImprovementSuggestionsStyle: r.FormValue("improvement_style"),
	}
}

// SetupObserverRoutes registers all observer-related routes
func SetupObserverRoutes(mux *http.ServeMux) {
	// List and Create
	mux.HandleFunc("/observers", func(w http.ResponseWriter, r *http.Request) {
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
}
