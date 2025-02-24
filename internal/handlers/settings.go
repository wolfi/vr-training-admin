// Copy the content from updated-settings-handlers artifact
// internal/handlers/settings.go
package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/a-h/templ"
	"github.com/saladinomario/vr-training-admin/internal/models"
	"github.com/saladinomario/vr-training-admin/templates/components/settings"
	"github.com/saladinomario/vr-training-admin/templates/pages"
)

var settingsStore *models.SettingsStore

func init() {
	log.Println("Initializing settings handler...")

	// Create data directory if it doesn't exist
	dataDir := "./data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("Error creating data directory: %v", err)
	} else {
		log.Printf("Data directory ensured: %s", dataDir)
	}

	// Initialize settings store
	settingsFilePath := dataDir + "/settings.json"
	log.Printf("Using settings file: %s", settingsFilePath)
	settingsStore = models.NewSettingsStore(settingsFilePath)

	log.Println("Settings handler initialized successfully")
}

// SettingsHandler handles the settings index page
func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/settings" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	llmSettings := settingsStore.GetLLMSettings()
	generalSettings := settingsStore.GetGeneralSettings()

	component := pages.SettingsIndex(llmSettings, generalSettings)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering settings page: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// UpdateLLMSettingsHandler handles updating the LLM settings
func UpdateLLMSettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse form values
	llmSettings := parseLLMSettingsForm(r)

	// Update settings
	err := settingsStore.UpdateLLMSettings(llmSettings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success message
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
    <div class="alert alert-success">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
        <span>LLM settings updated successfully!</span>
    </div>
    `))
}

// UpdateGeneralSettingsHandler handles updating the general settings
func UpdateGeneralSettingsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse form values
	generalSettings := parseGeneralSettingsForm(r)

	// Update settings
	err := settingsStore.UpdateGeneralSettings(generalSettings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a success message
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(`
    <div class="alert alert-success">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
        <span>General settings updated successfully!</span>
    </div>
    `))
}

// TestConnectionHandler handles testing the LLM API connection
func TestConnectionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get test prompt
	prompt := r.FormValue("test_prompt")
	if prompt == "" {
		prompt = "Hello! This is a test prompt."
	}

	// Test connection
	success, message, response := settingsStore.TestConnection(prompt)

	// Render test result
	component := settings.ConnectionResult(success, message, response)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering connection result: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// Helper function to parse LLM settings form
func parseLLMSettingsForm(r *http.Request) settings.LLMSettings {
	maxTokens, _ := strconv.Atoi(r.FormValue("max_tokens"))
	if maxTokens == 0 {
		maxTokens = 1024 // Default value
	}

	// Get current settings to preserve the ID
	currentSettings := settingsStore.GetLLMSettings()

	return settings.LLMSettings{
		ID:                currentSettings.ID,
		Provider:          r.FormValue("provider"),
		APIKey:            r.FormValue("api_key"),
		Model:             r.FormValue("model"),
		MaxTokens:         maxTokens,
		Temperature:       0.7,  // Default value
		TopP:              0.95, // Default value
		ProjectID:         r.FormValue("project_id"),
		Location:          r.FormValue("location"),
		Endpoint:          r.FormValue("endpoint"),
		ServiceAccountKey: r.FormValue("service_account_key"),
	}
}

// Helper function to parse general settings form
func parseGeneralSettingsForm(r *http.Request) settings.GeneralSettings {
	sessionTimeout, _ := strconv.Atoi(r.FormValue("session_timeout"))
	if sessionTimeout == 0 {
		sessionTimeout = 60 // Default value
	}

	// Default values for missing fields
	applicationName := r.FormValue("application_name")
	if applicationName == "" {
		applicationName = "VR Training Admin"
	}

	logLevel := r.FormValue("log_level")
	if logLevel == "" {
		logLevel = "INFO"
	}

	return settings.GeneralSettings{
		ApplicationName:   applicationName,
		LogLevel:          logLevel,
		SessionTimeout:    sessionTimeout,
		RecordSessions:    true, // Default value
		StoreSessionData:  true, // Default value
		DataRetentionDays: 90,   // Default value
	}
}

// SetupSettingsRoutes registers all settings-related routes
func SetupSettingsRoutes(mux *http.ServeMux) {
	log.Println("Setting up settings routes...")

	// Settings index
	log.Println("  Registering route: /settings")
	mux.HandleFunc("/settings", SettingsHandler)

	// LLM settings update
	log.Println("  Registering route: /settings/llm")
	mux.HandleFunc("/settings/llm", UpdateLLMSettingsHandler)

	// General settings update
	log.Println("  Registering route: /settings/general")
	mux.HandleFunc("/settings/general", UpdateGeneralSettingsHandler)

	// Test connection
	log.Println("  Registering route: /settings/test-connection")
	mux.HandleFunc("/settings/test-connection", TestConnectionHandler)

	// Provider fields
	log.Println("  Registering route: /settings/provider-fields")
	mux.HandleFunc("/settings/provider-fields", ProviderFieldsHandler)

	log.Println("Settings routes registered successfully")
}

// ProviderFieldsHandler handles returning provider-specific fields
func ProviderFieldsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the selected provider from query params
	provider := r.URL.Query().Get("provider")
	if provider == "" {
		http.Error(w, "Provider is required", http.StatusBadRequest)
		return
	}

	// Get current settings to preserve values
	llmSettings := settingsStore.GetLLMSettings()
	// Update the provider to match the selected one
	llmSettings.Provider = provider

	// Render the appropriate provider fields
	var component templ.Component
	if provider == "Google Vertex AI" || provider == "Google PaLM API" {
		component = settings.GoogleProviderFields(&llmSettings)
	} else {
		component = settings.GenericProviderFields(&llmSettings)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(r.Context(), w); err != nil {
		log.Printf("Error rendering provider fields: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
