// internal/models/observer.go
package models

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/saladinomario/vr-training-admin/templates/components/observers"
)

var (
	ErrObserverNotFound = errors.New("observer not found")
	ErrInvalidObserver  = errors.New("invalid observer data")
)

// ObserverStore implements an in-memory storage for observers
type ObserverStore struct {
	observers map[string]observers.Observer
	mu        sync.RWMutex
}

// NewObserverStore creates a new observer store with some sample data
func NewObserverStore() *ObserverStore {
	store := &ObserverStore{
		observers: make(map[string]observers.Observer),
	}

	// Add some sample observers
	sampleObservers := []observers.Observer{
		{
			ID:                "1",
			Name:              "Service Standards Supervisor",
			Description:       "Focuses on ensuring proper procedure implementation and citizen service quality in complex situations.",
			FeedbackStyle:     "Service Standards Focus",
			InterventionLevel: 3,
			DetailLevel:       4,
			FeedbackTone:      "Professional",
			SuccessMetrics:    "Adherence to service procedures, clear communication of requirements, and proper documentation handling while maintaining citizen satisfaction.",
			InterventionTriggers: []string{
				"Incorrect procedure application",
				"Documentation errors",
				"Privacy/confidentiality breach",
				"Service standards deviation",
				"Policy misinterpretation",
			},
			Active: true,
		},
		{
			ID:                "2",
			Name:              "Citizen Experience Specialist",
			Description:       "Evaluates communication effectiveness and citizen interaction quality, especially in challenging situations.",
			FeedbackStyle:     "Citizen Experience",
			InterventionLevel: 2,
			DetailLevel:       5,
			FeedbackTone:      "Constructive",
			SuccessMetrics:    "Effective citizen communication, proper handling of emotional situations, and successful service delivery.",
			InterventionTriggers: []string{
				"Citizen distress indicators",
				"Communication barriers",
				"Escalation requirements missed",
				"Proper referral needed",
			},
			Active: true,
		},
	}

	for _, observer := range sampleObservers {
		store.observers[observer.ID] = observer
	}

	return store
}

// GetAll returns all observers
func (s *ObserverStore) GetAll() []observers.Observer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]observers.Observer, 0, len(s.observers))
	for _, observer := range s.observers {
		result = append(result, observer)
	}
	return result
}

// GetByID returns an observer by its ID
func (s *ObserverStore) GetByID(id string) (observers.Observer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	observer, ok := s.observers[id]
	if !ok {
		return observers.Observer{}, ErrObserverNotFound
	}
	return observer, nil
}

// Create adds a new observer
func (s *ObserverStore) Create(observer observers.Observer) (observers.Observer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Very basic validation
	if observer.Name == "" {
		return observers.Observer{}, ErrInvalidObserver
	}

	// Generate a simple ID based on timestamp
	observer.ID = generateObserverID()
	s.observers[observer.ID] = observer
	return observer, nil
}

// Update modifies an existing observer
func (s *ObserverStore) Update(id string, observer observers.Observer) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.observers[id]; !ok {
		return ErrObserverNotFound
	}

	// Basic validation
	if observer.Name == "" {
		return ErrInvalidObserver
	}

	// Preserve the ID
	observer.ID = id
	s.observers[id] = observer
	return nil
}

// Delete removes an observer
func (s *ObserverStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.observers[id]; !ok {
		return ErrObserverNotFound
	}

	delete(s.observers, id)
	return nil
}

// Search looks for observers matching the query
func (s *ObserverStore) Search(query string) []observers.Observer {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if query == "" {
		return s.GetAll()
	}

	query = strings.ToLower(query)
	result := make([]observers.Observer, 0)

	for _, observer := range s.observers {
		if strings.Contains(strings.ToLower(observer.Name), query) ||
			strings.Contains(strings.ToLower(observer.Description), query) ||
			strings.Contains(strings.ToLower(observer.FeedbackStyle), query) {
			result = append(result, observer)
		}
	}

	return result
}

// Helper to generate a simple ID
func generateObserverID() string {
	return "observer_" + time.Now().Format("20060102150405")
}
