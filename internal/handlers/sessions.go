// internal/handlers/sessions.go
package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/saladinomario/vr-training-admin/internal/models"
	"github.com/saladinomario/vr-training-admin/templates/components/sessions"
	"github.com/saladinomario/vr-training-admin/templates/pages"
)

var SessionStore *models.SessionStore

func init() {
	// Create data directory if it doesn't exist
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("Error creating data directory: %v", err)
	}

	// Initialize session store
	sessionFilePath := dataDir + "/sessions.json"
	SessionStore = models.NewSessionStore(sessionFilePath)
}

// StartSessionHandler handles the session creation form submission
func StartSessionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Extract form values
	scenarioID := r.FormValue("scenario_id")
	avatarID := r.FormValue("avatar_id")
	observerID := r.FormValue("observer_id")

	// Validate required fields
	if scenarioID == "" || avatarID == "" || observerID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Create session
	session, err := SessionStore.Create(scenarioID, avatarID, observerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Start the session in Unreal Engine
	go startUnrealEngineSession(session.ID)

	// Return success response
	if r.Header.Get("HX-Request") == "true" {
		// Get updated dashboard content
		recentSessions := SessionStore.GetRecent(5)
		component := pages.DashboardContent(recentSessions)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering dashboard content: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Redirect to dashboard for non-HTMX requests
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// SessionStatusHandler handles updating session status
func SessionStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract session ID from URL
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	sessionID := parts[2]

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Extract status from form
	status := r.FormValue("status")
	if status == "" {
		http.Error(w, "Status is required", http.StatusBadRequest)
		return
	}

	// Only allow valid status values
	validStatus := map[string]bool{
		sessions.StatusRunning:   true,
		sessions.StatusPaused:    true,
		sessions.StatusCompleted: true,
	}

	if !validStatus[status] {
		http.Error(w, "Invalid status value", http.StatusBadRequest)
		return
	}

	// Update session status
	err := SessionStore.Update(sessionID, status)
	if err != nil {
		if err == models.ErrSessionNotFound {
			http.NotFound(w, r)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Update session in Unreal Engine
	go updateUnrealEngineSession(sessionID, status)

	// Return success response
	if r.Header.Get("HX-Request") == "true" {
		// Get updated dashboard content
		recentSessions := SessionStore.GetRecent(5)
		component := pages.RecentActivity(recentSessions)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("HX-Trigger", "closeModal")
		if err := component.Render(r.Context(), w); err != nil {
			log.Printf("Error rendering recent activity: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// Redirect to dashboard for non-HTMX requests
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// SessionFormHandler handles serving the new session form
func SessionFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get all scenarios, avatars, and observers for form dropdowns
	allScenarios := ScenarioStore.GetAll()
	allAvatars := AvatarStore.GetAll()
	allObservers := ObserverStore.GetAll()

	// Render form page
	component := pages.SessionNew(allScenarios, allAvatars, allObservers)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering session form: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// SetupSessionRoutes registers all session-related routes
func SetupSessionRoutes(mux *http.ServeMux) {
	log.Println("Setting up session routes...")

	// Session form
	mux.HandleFunc("/sessions/new", SessionFormHandler)

	// Start session
	mux.HandleFunc("/sessions/start", StartSessionHandler)

	// Update session status
	mux.HandleFunc("/sessions/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/sessions/") && r.Method == http.MethodPost {
			SessionStatusHandler(w, r)
			return
		}
		http.NotFound(w, r)
	})

	log.Println("Session routes registered successfully")
}

// Unreal Engine integration functions

// startUnrealEngineSession sends a request to start a session in Unreal Engine
func startUnrealEngineSession(sessionID string) {
	// Update status to "running"
	err := SessionStore.Update(sessionID, sessions.StatusRunning)
	if err != nil {
		log.Printf("Error updating session status: %v", err)
		return
	}

	// Create payload for Unreal Engine
	payload, err := SessionStore.CreateURESessionPayload(sessionID, ScenarioStore, AvatarStore, ObserverStore)
	if err != nil {
		log.Printf("Error creating UE payload: %v", err)
		return
	}

	// Send to Unreal Engine
	sendToUnrealEngine(payload)
}

// updateUnrealEngineSession sends a request to update a session in Unreal Engine
func updateUnrealEngineSession(sessionID, status string) {
	// Create payload for Unreal Engine
	payload, err := SessionStore.CreateURESessionPayload(sessionID, ScenarioStore, AvatarStore, ObserverStore)
	if err != nil {
		log.Printf("Error creating UE payload: %v", err)
		return
	}

	// Send to Unreal Engine
	sendToUnrealEngine(payload)
}

// sendToUnrealEngine sends data to the Unreal Engine endpoint
func sendToUnrealEngine(payload []byte) {
	// TODO: Get this from settings
	unrealEndpoint := "http://localhost:8081/api/vr-session"

	// For now, just log the payload
	log.Printf("Would send to Unreal Engine: %s", unrealEndpoint)

	// Create a nicely formatted version of the payload for logging
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, payload, "", "  "); err == nil {
		log.Printf("Payload: %s", prettyJSON.String())
	}

	// In production, uncomment this code to actually send the request
	/*
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		req, err := http.NewRequest("POST", unrealEndpoint, bytes.NewBuffer(payload))
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error sending request to Unreal Engine: %v", err)
			return
		}
		defer resp.Body.Close()

		log.Printf("Unreal Engine response status: %s", resp.Status)
	*/
}
