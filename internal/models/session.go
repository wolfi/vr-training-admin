// internal/models/session.go
package models

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/saladinomario/vr-training-admin/templates/components/avatars"
	"github.com/saladinomario/vr-training-admin/templates/components/observers"
	"github.com/saladinomario/vr-training-admin/templates/components/scenarios"
	"github.com/saladinomario/vr-training-admin/templates/components/sessions"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrInvalidSession  = errors.New("invalid session data")
)

// TrainingSessionStore manages VR training sessions
type SessionStore struct {
	sessions map[string]*sessions.Session
	mu       sync.RWMutex
	filePath string
}

// NewSessionStore creates a new session store
func NewSessionStore(filePath string) *SessionStore {
	store := &SessionStore{
		sessions: make(map[string]*sessions.Session),
		filePath: filePath,
	}

	// Create the directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Error creating directory for sessions: %v", err)
	}

	// Load existing sessions if file exists
	if _, err := os.Stat(filePath); err == nil {
		store.loadSessions()
	}

	return store
}

// loadSessions loads sessions from disk
func (s *SessionStore) loadSessions() {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		log.Printf("Error reading sessions file: %v", err)
		return
	}

	var sessions []*sessions.Session
	if err := json.Unmarshal(data, &sessions); err != nil {
		log.Printf("Error unmarshaling sessions: %v", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, session := range sessions {
		s.sessions[session.ID] = session
	}

	log.Printf("Loaded %d sessions from disk", len(sessions))
}

// saveSessions saves sessions to disk

func (s *SessionStore) saveSessions() error {
	s.mu.RLock()
	sessionsList := make([]*sessions.Session, 0, len(s.sessions))
	for _, session := range s.sessions {
		sessionsList = append(sessionsList, session)
	}
	s.mu.RUnlock()

	data, err := json.MarshalIndent(sessionsList, "", "  ")
	if err != nil {
		log.Printf("Error marshaling sessions: %v", err)
		return err
	}

	// Make sure the directory exists
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Printf("Error creating directory for sessions: %v", err)
		return err
	}

	// Write file with explicit sync
	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		log.Printf("Error writing sessions to file: %v", err)
		return err
	}

	log.Printf("Saved %d sessions to %s", len(sessionsList), s.filePath)
	return nil
}

// GetAll returns all sessions
func (s *SessionStore) GetAll() []*sessions.Session {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*sessions.Session, 0, len(s.sessions))
	for _, session := range s.sessions {
		result = append(result, session)
	}

	// Sort sessions by start time, newest first
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].StartTime.Before(result[j].StartTime) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result
}

// GetRecent returns the n most recent sessions
func (s *SessionStore) GetRecent(n int) []*sessions.Session {
	allSessions := s.GetAll()
	if len(allSessions) <= n {
		return allSessions
	}
	return allSessions[:n]
}

// GetByID returns a session by ID
func (s *SessionStore) GetByID(id string) (*sessions.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	session, ok := s.sessions[id]
	if !ok {
		return nil, ErrSessionNotFound
	}
	return session, nil
}

// Update the Create method:
func (s *SessionStore) Create(scenarioID, avatarID, observerID string) (*sessions.Session, error) {
	// Generate ID based on timestamp
	id := "session_" + time.Now().Format("20060102150405")

	// Create the session
	session := &sessions.Session{
		ID:         id,
		ScenarioID: scenarioID,
		AvatarID:   avatarID,
		ObserverID: observerID,
		Status:     sessions.StatusPending,
		StartTime:  time.Now(),
		UpdateTime: time.Now(),
	}

	s.mu.Lock()
	s.sessions[id] = session
	s.mu.Unlock()

	// Save to disk synchronously
	err := s.saveSessions()
	if err != nil {
		log.Printf("Error saving sessions: %v", err)
	}

	return session, nil
}

// Update modifies an existing session
// internal/models/session.go
// Update the Update method:

func (s *SessionStore) Update(id string, status string) error {
	s.mu.Lock()

	session, ok := s.sessions[id]
	if !ok {
		return ErrSessionNotFound
	}

	// Update status and timestamp
	session.Status = status
	session.UpdateTime = time.Now()

	// If completed, set end time
	if status == sessions.StatusCompleted {
		now := time.Now()
		session.EndTime = &now
	}

	s.mu.Unlock()

	// Save to disk synchronously
	err := s.saveSessions()
	if err != nil {
		log.Printf("Error saving sessions: %v", err)
	}

	return nil
}

// Delete removes a session
func (s *SessionStore) Delete(id string) error {
	s.mu.RLock()
	if _, ok := s.sessions[id]; !ok {
		return ErrSessionNotFound
	}
	s.mu.RUnlock()

	delete(s.sessions, id)

	// Save to disk
	go s.saveSessions()

	return nil
}

// GetSessionDetails retrieves the scenario, avatar, and observer details for a session
func (s *SessionStore) GetSessionDetails(id string, scenarioStore *ScenarioStore, avatarStore *AvatarStore, observerStore *ObserverStore) (*sessions.SessionDetails, error) {
	session, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Get associated entities
	scenario, err := scenarioStore.GetByID(session.ScenarioID)
	if err != nil {
		log.Printf("Warning: Session %s references non-existent scenario %s", id, session.ScenarioID)
	}

	avatar, err := avatarStore.GetByID(session.AvatarID)
	if err != nil {
		log.Printf("Warning: Session %s references non-existent avatar %s", id, session.AvatarID)
	}

	observer, err := observerStore.GetByID(session.ObserverID)
	if err != nil {
		log.Printf("Warning: Session %s references non-existent observer %s", id, session.ObserverID)
	}

	return &sessions.SessionDetails{
		Session:  *session,
		Scenario: scenario,
		Avatar:   avatar,
		Observer: observer,
	}, nil
}

// CreateURESessionPayload creates the payload to send to Unreal Engine
func (s *SessionStore) CreateURESessionPayload(id string, scenarioStore *ScenarioStore, avatarStore *AvatarStore, observerStore *ObserverStore) ([]byte, error) {
	details, err := s.GetSessionDetails(id, scenarioStore, avatarStore, observerStore)
	if err != nil {
		return nil, err
	}

	// Create a payload structure for Unreal Engine
	type UREPayload struct {
		SessionID string             `json:"sessionId"`
		Status    string             `json:"status"`
		Scenario  scenarios.Scenario `json:"scenario"`
		Avatar    avatars.Avatar     `json:"avatar"`
		Observer  observers.Observer `json:"observer"`
		Timestamp string             `json:"timestamp"`
	}

	payload := UREPayload{
		SessionID: details.Session.ID,
		Status:    details.Session.Status,
		Scenario:  details.Scenario,
		Avatar:    details.Avatar,
		Observer:  details.Observer,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	return json.Marshal(payload)
}
