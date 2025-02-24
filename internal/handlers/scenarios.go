// internal/handlers/scenarios.go
package handlers

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/saladinomario/vr-training-admin/internal/models"
	"github.com/saladinomario/vr-training-admin/templates/components/scenarios"
	"github.com/saladinomario/vr-training-admin/templates/pages"
)

var scenarioStore *models.ScenarioStore

func init() {
	scenarioStore = models.NewScenarioStore()
}

// ScenariosHandler handles the scenarios index page
func ScenariosHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/scenarios" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	allScenarios := scenarioStore.GetAll()
	component := pages.ScenariosIndex(allScenarios)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering scenarios page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ScenarioNewHandler handles the new scenario form
func ScenarioNewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	component := pages.ScenarioNew()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering new scenario form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ScenarioEditHandler handles the edit scenario form
func ScenarioEditHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/scenarios/edit/")
	if idStr == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	// Get scenario by ID
	scenario, err := scenarioStore.GetByID(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	component := pages.ScenarioEdit(scenario)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering edit scenario form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ScenarioCreateHandler handles scenario creation
func ScenarioCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse form values
	scenario := parseScenarioForm(r)

	// Create scenario
	_, err := scenarioStore.Create(scenario)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// If this is an HTMX request, return the updated content
	if r.Header.Get("HX-Request") == "true" {
		allScenarios := scenarioStore.GetAll()
		component := pages.ScenariosContent(allScenarios)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering scenarios content: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Standard redirect for non-HTMX requests
	http.Redirect(w, r, "/scenarios", http.StatusSeeOther)
}

// ScenarioUpdateHandler handles scenario updates
func ScenarioUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/scenarios/")
	if idStr == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse form values
	scenario := parseScenarioForm(r)

	// Update scenario
	err := scenarioStore.Update(idStr, scenario)
	if err != nil {
		if err == models.ErrScenarioNotFound {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	// If this is an HTMX request, return the updated content
	if r.Header.Get("HX-Request") == "true" {
		allScenarios := scenarioStore.GetAll()
		component := pages.ScenariosContent(allScenarios)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering scenarios content: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Standard redirect for non-HTMX requests
	http.Redirect(w, r, "/scenarios", http.StatusSeeOther)
}

// ScenarioDeleteHandler handles scenario deletion
func ScenarioDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL
	idStr := strings.TrimPrefix(r.URL.Path, "/scenarios/")
	if idStr == r.URL.Path {
		http.NotFound(w, r)
		return
	}

	// Delete scenario
	err := scenarioStore.Delete(idStr)
	if err != nil {
		if err == models.ErrScenarioNotFound {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// If this is an HTMX request, return the updated scenario list
	if r.Header.Get("HX-Request") == "true" {
		allScenarios := scenarioStore.GetAll()
		component := scenarios.ScenarioList(allScenarios)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering scenario list: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Standard redirect for non-HTMX requests
	http.Redirect(w, r, "/scenarios", http.StatusSeeOther)
}

// ScenarioSearchHandler handles scenario search
func ScenarioSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	foundScenarios := scenarioStore.Search(query)

	component := scenarios.ScenarioList(foundScenarios)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering scenario search results: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Helper function to parse scenario form data
func parseScenarioForm(r *http.Request) scenarios.Scenario {
	difficulty, _ := strconv.Atoi(r.FormValue("difficulty"))
	duration, _ := strconv.Atoi(r.FormValue("duration"))
	backgroundNoise := r.FormValue("background_noise") == "on"

	return scenarios.Scenario{
		Name:            r.FormValue("name"),
		Description:     r.FormValue("description"),
		Category:        r.FormValue("category"),
		Difficulty:      difficulty,
		Duration:        duration,
		Scene:           r.FormValue("scene"),
		BackgroundNoise: backgroundNoise,
		SuccessCriteria: r.FormValue("success_criteria"),
		Keywords:        r.FormValue("keywords"),
	}
}

// SetupScenarioRoutes registers all scenario-related routes
func SetupScenarioRoutes(mux *http.ServeMux) {
	// List and Create
	mux.HandleFunc("/scenarios", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ScenariosHandler(w, r)
		case http.MethodPost:
			ScenarioCreateHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// New form
	mux.HandleFunc("/scenarios/new", ScenarioNewHandler)

	// Search
	mux.HandleFunc("/scenarios/search", ScenarioSearchHandler)

	// Edit form
	mux.HandleFunc("/scenarios/edit/", ScenarioEditHandler)

	// Update and Delete
	mux.HandleFunc("/scenarios/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut, http.MethodPost:
			ScenarioUpdateHandler(w, r)
		case http.MethodDelete:
			ScenarioDeleteHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
