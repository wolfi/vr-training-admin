// internal/models/avatar.go
package models

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/saladinomario/vr-training-admin/templates/components/avatars"
)

var (
	ErrAvatarNotFound = errors.New("avatar not found")
	ErrInvalidAvatar  = errors.New("invalid avatar data")
)

// AvatarStore implements an in-memory storage for avatars
type AvatarStore struct {
	avatars map[string]avatars.Avatar
	mu      sync.RWMutex
}

// NewAvatarStore creates a new avatar store with some sample data
func NewAvatarStore() *AvatarStore {
	store := &AvatarStore{
		avatars: make(map[string]avatars.Avatar),
	}

	// Add some sample avatars
	sampleAvatars := []avatars.Avatar{
		{
			ID:                  "1",
			Name:                "Sales Manager",
			Description:         "A challenging sales manager who will test your negotiation skills.",
			PersonalityType:     "Challenging",
			CommunicationStyle:  "Direct",
			KnowledgeLevel:      8,
			AggressivenessLevel: 7,
			PatienceLevel:       4,
			EmotionalReactivity: 6,
			VoiceType:           "Deep",
			SpeakingSpeed:       3,
			Keywords:            "negotiation, sales targets, performance, objections",
		},
		{
			ID:                  "2",
			Name:                "Customer Support",
			Description:         "A patient customer support representative for practicing service scenarios.",
			PersonalityType:     "Supportive",
			CommunicationStyle:  "Empathetic",
			KnowledgeLevel:      6,
			AggressivenessLevel: 2,
			PatienceLevel:       9,
			EmotionalReactivity: 4,
			VoiceType:           "Natural",
			SpeakingSpeed:       3,
			Keywords:            "service, resolution, listening, empathy",
		},
	}

	for _, avatar := range sampleAvatars {
		store.avatars[avatar.ID] = avatar
	}

	return store
}

// GetAll returns all avatars
func (s *AvatarStore) GetAll() []avatars.Avatar {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]avatars.Avatar, 0, len(s.avatars))
	for _, avatar := range s.avatars {
		result = append(result, avatar)
	}
	return result
}

// GetByID returns an avatar by its ID
func (s *AvatarStore) GetByID(id string) (avatars.Avatar, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	avatar, ok := s.avatars[id]
	if !ok {
		return avatars.Avatar{}, ErrAvatarNotFound
	}
	return avatar, nil
}

// Create adds a new avatar
func (s *AvatarStore) Create(avatar avatars.Avatar) (avatars.Avatar, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Very basic validation
	if avatar.Name == "" {
		return avatars.Avatar{}, ErrInvalidAvatar
	}

	// Generate a simple ID based on timestamp
	avatar.ID = generateAvatarID()
	s.avatars[avatar.ID] = avatar
	return avatar, nil
}

// Update modifies an existing avatar
func (s *AvatarStore) Update(id string, avatar avatars.Avatar) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.avatars[id]; !ok {
		return ErrAvatarNotFound
	}

	// Basic validation
	if avatar.Name == "" {
		return ErrInvalidAvatar
	}

	// Preserve the ID
	avatar.ID = id
	s.avatars[id] = avatar
	return nil
}

// Delete removes an avatar
func (s *AvatarStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.avatars[id]; !ok {
		return ErrAvatarNotFound
	}

	delete(s.avatars, id)
	return nil
}

// Search looks for avatars matching the query
func (s *AvatarStore) Search(query string) []avatars.Avatar {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if query == "" {
		return s.GetAll()
	}

	query = strings.ToLower(query)
	result := make([]avatars.Avatar, 0)

	for _, avatar := range s.avatars {
		if strings.Contains(strings.ToLower(avatar.Name), query) ||
			strings.Contains(strings.ToLower(avatar.Description), query) ||
			strings.Contains(strings.ToLower(avatar.PersonalityType), query) {
			result = append(result, avatar)
		}
	}

	return result
}

// Helper to generate a simple ID
func generateAvatarID() string {
	return "avatar_" + time.Now().Format("20060102150405")
}
