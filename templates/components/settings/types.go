// templates/components/settings/types.go
package settings

type LLMSettings struct {
	ID                string
	Provider          string
	APIKey            string
	Model             string
	MaxTokens         int
	Temperature       float64
	TopP              float64
	FrequencyPenalty  float64
	PresencePenalty   float64
	ProjectID         string // For Google LLM services
	Location          string // For Google LLM services
	Endpoint          string // For Google LLM services
	ServiceAccountKey string // For Google LLM services
}

type GeneralSettings struct {
	ApplicationName       string
	LogLevel              string
	MaxConcurrentSessions int
	SessionTimeout        int // minutes
	RecordSessions        bool
	StoreSessionData      bool
	DataRetentionDays     int
}

// Providers returns available LLM providers
func Providers() []string {
	return []string{
		"Google Vertex AI",
		"Google PaLM API",
		"OpenAI",
		"Anthropic",
		"Custom Endpoint",
	}
}

// GoogleModels returns available Google LLM models
func GoogleModels() []string {
	return []string{
		"gemini-pro",
		"gemini-pro-vision",
		"text-bison",
		"text-unicorn",
		"chat-bison",
		"codechat-bison",
	}
}

// Locations returns available Google Cloud regions
func Locations() []string {
	return []string{
		"us-central1",
		"us-east1",
		"us-west1",
		"europe-west1",
		"europe-west2",
		"europe-west4",
		"asia-east1",
		"asia-northeast1",
		"asia-southeast1",
	}
}

// LogLevels returns available logging levels
func LogLevels() []string {
	return []string{
		"DEBUG",
		"INFO",
		"WARNING",
		"ERROR",
		"CRITICAL",
	}
}
