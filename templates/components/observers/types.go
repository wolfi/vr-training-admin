// templates/components/observers/types.go
package observers

type Observer struct {
	ID                   string
	Name                 string
	Description          string
	FeedbackStyle        string
	InterventionLevel    int // 1-5 scale (1: Minimal, 5: Frequent)
	DetailLevel          int // 1-5 scale (1: Brief, 5: Comprehensive)
	FeedbackTone         string
	SuccessMetrics       string
	InterventionTriggers []string
	Active               bool
}

// FeedbackStyles returns available feedback styles
func FeedbackStyles() []string {
	return []string{
		"Supportive",
		"Challenging",
		"Analytical",
		"Instructional",
		"Socratic",
		"Direct",
	}
}

// FeedbackTones returns available feedback tones
func FeedbackTones() []string {
	return []string{
		"Formal",
		"Casual",
		"Encouraging",
		"Neutral",
		"Authoritative",
		"Gentle",
	}
}

// CommonTriggers returns common intervention triggers
func CommonTriggers() []string {
	return []string{
		"Off-topic conversation",
		"Inappropriate communication",
		"Silence over 10 seconds",
		"Missed opportunity",
		"Incorrect information shared",
		"Customer frustration detected",
		"Talking over the customer",
		"Success criteria met",
	}
}
