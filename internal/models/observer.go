package models

import (
	"errors"
	"strings"
	"sync"
	"time"
)

// Observer represents the configuration for the "Observer LLM", i.e., the evaluator/coach.
type Observer struct {
	ID                          string
	Name                        string
	InterruptionTriggers        string
	FeedbackStyle               string
	InterventionFrequency       int
	KeyPerformanceIndicators    string
	ScoringSystem               string
	RequiredAchievements        string
	FailedStateDefinitions      string
	RealTimeFeedbackRules       string
	PostSessionReportTemplate   string
	LearningPointsFocus         string
	ImprovementSuggestionsStyle string
}

// Errors for the Observer domain.
var (
	ErrObserverNotFound = errors.New("observer not found")
	ErrInvalidObserver  = errors.New("invalid observer data")
)

// ObserverStore implements an in-memory storage for Observers.
type ObserverStore struct {
	observers map[string]Observer
	mu        sync.RWMutex
}

// NewObserverStore creates a new observer store (in-memory).
func NewObserverStore() *ObserverStore {
	store := &ObserverStore{
		observers: make(map[string]Observer),
	}

	// (Optional) seed with sample Observers:
	sampleObservers := []Observer{
		{
			ID:                          "obs1",
			Name:                        "Default Sales Observer",
			InterruptionTriggers:        "Off-topic, Negative language",
			FeedbackStyle:               "Encouraging",
			InterventionFrequency:       2,
			KeyPerformanceIndicators:    "Client persuasion, Objection handling",
			ScoringSystem:               "Points-based",
			RequiredAchievements:        "Close sale, Address objections",
			FailedStateDefinitions:      "Insufficient responses, client dissatisfaction",
			RealTimeFeedbackRules:       "Interrupt if user is silent > 10s",
			PostSessionReportTemplate:   "Summary, tips, and recommended reading",
			LearningPointsFocus:         "Soft skills, active listening",
			ImprovementSuggestionsStyle: "Bulleted list",
		},
		{
			ID:                          "obs2",
			Name:                        "Conflict Resolution Coach",
			InterruptionTriggers:        "Escalating tension, personal attacks",
			FeedbackStyle:               "Direct but constructive",
			InterventionFrequency:       3,
			KeyPerformanceIndicators:    "Resolution speed, Tone management",
			ScoringSystem:               "Score from 1-5 on empathy & clarity",
			RequiredAchievements:        "Agreement or compromise",
			FailedStateDefinitions:      "Aggressive or insulting language used",
			RealTimeFeedbackRules:       "Immediate feedback on emotional triggers",
			PostSessionReportTemplate:   "Detailed transcript analysis, highlight missteps",
			LearningPointsFocus:         "De-escalation, empathic listening",
			ImprovementSuggestionsStyle: "Step-by-step",
		},
	}

	for _, obs := range sampleObservers {
		store.observers[obs.ID] = obs
	}

	return store
}

// GetAll returns all Observers in the store.
func (s *ObserverStore) GetAll() []Observer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	results := make([]Observer, 0, len(s.observers))
	for _, obs := range s.observers {
		results = append(results, obs)
	}
	return results
}

// GetByID returns a single Observer by ID.
func (s *ObserverStore) GetByID(id string) (Observer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	obs, ok := s.observers[id]
	if !ok {
		return Observer{}, ErrObserverNotFound
	}
	return obs, nil
}

// Create adds a new observer to the store.
func (s *ObserverStore) Create(obs Observer) (Observer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if obs.Name == "" {
		return Observer{}, ErrInvalidObserver
	}

	obs.ID = generateObserverID() // generate a new ID
	s.observers[obs.ID] = obs
	return obs, nil
}

// Update modifies an existing observer.
func (s *ObserverStore) Update(id string, obs Observer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.observers[id]; !ok {
		return ErrObserverNotFound
	}
	if obs.Name == "" {
		return ErrInvalidObserver
	}

	// Preserve the original ID.
	obs.ID = id
	s.observers[id] = obs
	return nil
}

// Delete removes an observer from the store.
func (s *ObserverStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.observers[id]; !ok {
		return ErrObserverNotFound
	}
	delete(s.observers, id)
	return nil
}

// Search looks for observers matching the query string in name or key fields.
func (s *ObserverStore) Search(query string) []Observer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if query == "" {
		return s.GetAll()
	}

	query = strings.ToLower(query)
	result := make([]Observer, 0)

	for _, obs := range s.observers {
		// Expand matching logic as needed:
		if strings.Contains(strings.ToLower(obs.Name), query) ||
			strings.Contains(strings.ToLower(obs.FeedbackStyle), query) ||
			strings.Contains(strings.ToLower(obs.KeyPerformanceIndicators), query) {
			result = append(result, obs)
		}
	}
	return result
}

// Helper to generate a simple ID.
func generateObserverID() string {
	return "observer_" + time.Now().Format("20060102150405")
}
