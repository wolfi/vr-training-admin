// internal/models/settings.go
package models

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"github.com/saladinomario/vr-training-admin/templates/components/settings"
)

// SettingsStore manages settings and persists them to a file
type SettingsStore struct {
	llmSettings     settings.LLMSettings
	generalSettings settings.GeneralSettings
	filePath        string
	mu              sync.RWMutex
}

// NewSettingsStore creates a new settings store with default values
func NewSettingsStore(filePath string) *SettingsStore {
	store := &SettingsStore{
		llmSettings: settings.LLMSettings{
			ID:                "default",
			Provider:          "Google Vertex AI",
			Model:             "gemini-pro",
			MaxTokens:         1024,
			Temperature:       0.7,
			TopP:              0.95,
			FrequencyPenalty:  0.0,
			PresencePenalty:   0.0,
			ProjectID:         "",
			Location:          "us-central1",
			Endpoint:          "",
			ServiceAccountKey: "",
		},
		generalSettings: settings.GeneralSettings{
			ApplicationName:       "VR Training Admin",
			LogLevel:              "INFO",
			MaxConcurrentSessions: 10,
			SessionTimeout:        60,
			RecordSessions:        true,
			StoreSessionData:      true,
			DataRetentionDays:     90,
		},
		filePath: filePath,
	}

	// Load existing settings if file exists
	if _, err := os.Stat(filePath); err == nil {
		store.loadFromFile()
	} else {
		// Save default settings
		store.saveToFile()
	}

	return store
}

// GetLLMSettings returns the current LLM settings
func (s *SettingsStore) GetLLMSettings() settings.LLMSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.llmSettings
}

// GetGeneralSettings returns the current general settings
func (s *SettingsStore) GetGeneralSettings() settings.GeneralSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.generalSettings
}

// UpdateLLMSettings updates the LLM settings
func (s *SettingsStore) UpdateLLMSettings(newSettings settings.LLMSettings) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.llmSettings = newSettings
	return s.saveToFile()
}

// UpdateGeneralSettings updates the general settings
func (s *SettingsStore) UpdateGeneralSettings(newSettings settings.GeneralSettings) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.generalSettings = newSettings
	return s.saveToFile()
}

// TestConnection simulates testing a connection to the LLM API
func (s *SettingsStore) TestConnection(prompt string) (bool, string, string) {
	// In a real implementation, this would make an actual API call
	// For now, we'll simulate a successful connection
	s.mu.RLock()
	provider := s.llmSettings.Provider
	s.mu.RUnlock()

	if provider == "" || prompt == "" {
		return false, "Invalid configuration. Please check your settings.", ""
	}

	// Simulate a successful connection
	response := "Connection successful! The API responded to your prompt."

	return true, "API connection test completed successfully.", response
}

// Private helper methods

// Combined settings for storage
type combinedSettings struct {
	LLM     settings.LLMSettings     `json:"llm"`
	General settings.GeneralSettings `json:"general"`
}

// loadFromFile loads settings from the JSON file
func (s *SettingsStore) loadFromFile() error {
	data, err := ioutil.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	var combined combinedSettings
	if err := json.Unmarshal(data, &combined); err != nil {
		return err
	}

	s.llmSettings = combined.LLM
	s.generalSettings = combined.General
	return nil
}

// saveToFile saves settings to the JSON file
func (s *SettingsStore) saveToFile() error {
	combined := combinedSettings{
		LLM:     s.llmSettings,
		General: s.generalSettings,
	}

	data, err := json.MarshalIndent(combined, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(s.filePath, data, 0644)
}
