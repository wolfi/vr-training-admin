// templates/components/avatars/types.go
package avatars

type Avatar struct {
	ID                  string
	Name                string
	Description         string
	PersonalityType     string
	CommunicationStyle  string
	KnowledgeLevel      int
	AggressivenessLevel int
	PatienceLevel       int
	EmotionalReactivity int
	VoiceType           string
	SpeakingSpeed       int // 1-5 scale
	ImageURL            string
	Keywords            string
}

// PersonalityTypes returns available personality types
func PersonalityTypes() []string {
	return []string{
		"Frustrated Citizen",
		"Confused Citizen",
		"Demanding Citizen",
		"Cooperative Citizen",
		"Elderly Citizen",
		"Language Barrier",
		"Emotional Citizen",
		"Time-Pressured Citizen",
	}
}

// CommunicationStyles returns available communication styles
func CommunicationStyles() []string {
	return []string{
		"Official",
		"Approachable",
		"Clear and Simple",
		"Detailed Explanation",
		"Patient and Supportive",
		"Firm but Respectful",
	}
}

// VoiceTypes returns available voice options
func VoiceTypes() []string {
	return []string{
		"Clear Native Speaker",
		"Non-Native Speaker",
		"Elderly Voice",
		"Emotional Tone",
		"Soft-Spoken",
	}
}
