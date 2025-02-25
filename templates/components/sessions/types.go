// templates/components/sessions/types.go
package sessions

import (
	"fmt"
	"time"

	"github.com/saladinomario/vr-training-admin/templates/components/avatars"
	"github.com/saladinomario/vr-training-admin/templates/components/observers"
	"github.com/saladinomario/vr-training-admin/templates/components/scenarios"
)

// Session status constants
const (
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusPaused    = "paused"
	StatusCompleted = "completed"
	StatusFailed    = "failed"
)

// Session represents a VR training session
type Session struct {
	ID         string     `json:"id"`
	ScenarioID string     `json:"scenarioId"`
	AvatarID   string     `json:"avatarId"`
	ObserverID string     `json:"observerId"`
	Status     string     `json:"status"`
	StartTime  time.Time  `json:"startTime"`
	EndTime    *time.Time `json:"endTime,omitempty"`
	UpdateTime time.Time  `json:"updateTime"`
	Score      *int       `json:"score,omitempty"`
	Notes      string     `json:"notes,omitempty"`
}

// SessionDetails contains all details for a session including related entities
type SessionDetails struct {
	Session  Session            `json:"session"`
	Scenario scenarios.Scenario `json:"scenario"`
	Avatar   avatars.Avatar     `json:"avatar"`
	Observer observers.Observer `json:"observer"`
}

// GetDuration returns the duration of the session
func (s *Session) GetDuration() time.Duration {
	if s.EndTime == nil {
		return time.Since(s.StartTime)
	}
	return s.EndTime.Sub(s.StartTime)
}

// GetFormattedDuration returns the duration as a string
func (s *Session) GetFormattedDuration() string {
	duration := s.GetDuration()
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}

	return fmt.Sprintf("%dm %ds", minutes, seconds)
}

// GetStatusClass returns the CSS class for the status badge
func (s *Session) GetStatusClass() string {
	switch s.Status {
	case StatusRunning:
		return "badge-primary"
	case StatusPaused:
		return "badge-warning"
	case StatusCompleted:
		return "badge-success"
	case StatusFailed:
		return "badge-error"
	default:
		return "badge-ghost"
	}
}
