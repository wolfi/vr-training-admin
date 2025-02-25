// internal/models/scenario.go
package models

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/saladinomario/vr-training-admin/templates/components/scenarios"
)

var (
	ErrScenarioNotFound = errors.New("scenario not found")
	ErrInvalidScenario  = errors.New("invalid scenario data")
)

// ScenarioStore implements an in-memory storage for scenarios
type ScenarioStore struct {
	scenarios map[string]scenarios.Scenario
	mu        sync.RWMutex
}

// NewScenarioStore creates a new scenario store with some sample data
func NewScenarioStore() *ScenarioStore {
	store := &ScenarioStore{
		scenarios: make(map[string]scenarios.Scenario),
	}

	// Add some sample scenarios
	sampleScenarios := []scenarios.Scenario{
		{
			ID:              "1",
			Name:            "Digital Service Assistance",
			Description:     "Help an elderly citizen navigate online form submission and digital identity verification.",
			Category:        "Document Processing",
			Difficulty:      4,
			Duration:        20,
			Scene:           "Information Desk",
			BackgroundNoise: true,
			SuccessCriteria: "Clear Communication and Proper Procedure Following",
			Keywords:        "digital services, online forms, identity verification, patience, step-by-step guidance",
		},
		{
			ID:              "2",
			Name:            "Language Barrier Resolution",
			Description:     "Assist a non-native speaker with completing essential service applications while maintaining clear communication.",
			Category:        "Special Needs Assistance",
			Difficulty:      3,
			Duration:        25,
			Scene:           "Private Consultation Room",
			BackgroundNoise: false,
			SuccessCriteria: "Effective Communication and Accurate Information Provided",
			Keywords:        "translation, clear speech, visual aids, documentation requirements, verification",
		},
		{
			ID:              "3",
			Name:            "Emergency Service Request",
			Description:     "Handle an urgent situation with a distressed citizen requiring immediate assistance with social services.",
			Category:        "Emergency Assistance",
			Difficulty:      5,
			Duration:        15,
			Scene:           "Quick Service Counter",
			BackgroundNoise: true,
			SuccessCriteria: "Effective Conflict Resolution and Appropriate Referral Made",
			Keywords:        "urgency, empathy, crisis management, service coordination, immediate response",
		},
	}

	for _, scenario := range sampleScenarios {
		store.scenarios[scenario.ID] = scenario
	}

	return store
}

// GetAll returns all scenarios
func (s *ScenarioStore) GetAll() []scenarios.Scenario {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]scenarios.Scenario, 0, len(s.scenarios))
	for _, scenario := range s.scenarios {
		result = append(result, scenario)
	}
	return result
}

// GetByID returns a scenario by its ID
func (s *ScenarioStore) GetByID(id string) (scenarios.Scenario, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	scenario, ok := s.scenarios[id]
	if !ok {
		return scenarios.Scenario{}, ErrScenarioNotFound
	}
	return scenario, nil
}

// Create adds a new scenario
func (s *ScenarioStore) Create(scenario scenarios.Scenario) (scenarios.Scenario, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Very basic validation
	if scenario.Name == "" {
		return scenarios.Scenario{}, ErrInvalidScenario
	}

	// Generate a simple ID based on timestamp
	scenario.ID = generateID()
	s.scenarios[scenario.ID] = scenario
	return scenario, nil
}

// Update modifies an existing scenario
func (s *ScenarioStore) Update(id string, scenario scenarios.Scenario) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.scenarios[id]; !ok {
		return ErrScenarioNotFound
	}

	// Basic validation
	if scenario.Name == "" {
		return ErrInvalidScenario
	}

	// Preserve the ID
	scenario.ID = id
	s.scenarios[id] = scenario
	return nil
}

// Delete removes a scenario
func (s *ScenarioStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.scenarios[id]; !ok {
		return ErrScenarioNotFound
	}

	delete(s.scenarios, id)
	return nil
}

// Search looks for scenarios matching the query
func (s *ScenarioStore) Search(query string) []scenarios.Scenario {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if query == "" {
		return s.GetAll()
	}

	query = strings.ToLower(query)
	result := make([]scenarios.Scenario, 0)

	for _, scenario := range s.scenarios {
		if strings.Contains(strings.ToLower(scenario.Name), query) ||
			strings.Contains(strings.ToLower(scenario.Description), query) ||
			strings.Contains(strings.ToLower(scenario.Category), query) {
			result = append(result, scenario)
		}
	}

	return result
}

// Helper to generate a simple ID
func generateID() string {
	return "scenario_" + time.Now().Format("20060102150405")
}
