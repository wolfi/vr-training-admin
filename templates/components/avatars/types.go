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
		"Friendly",
		"Professional",
		"Challenging",
		"Analytical",
		"Empathetic",
		"Technical",
		"Difficult",
		"Supportive",
	}
}

// CommunicationStyles returns available communication styles
func CommunicationStyles() []string {
	return []string{
		"Formal",
		"Casual",
		"Direct",
		"Indirect",
		"Detailed",
		"Concise",
	}
}

// VoiceTypes returns available voice options
func VoiceTypes() []string {
	return []string{
		"Natural",
		"Robotic",
		"Deep",
		"High-pitched",
		"Accented",
	}
}
