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
		"Service Standards Focus",
		"Procedure Compliance",
		"Citizen Experience",
		"Communication Quality",
		"Problem-Solving Process",
		"Policy Implementation",
	}
}

// FeedbackTones returns available feedback tones
func FeedbackTones() []string {
	return []string{
		"Professional",
		"Constructive",
		"Service-Oriented",
		"Objective",
		"Instructive",
		"Supportive",
	}
}

// CommonTriggers returns common intervention triggers
func CommonTriggers() []string {
	return []string{
		"Incorrect procedure application",
		"Privacy/confidentiality breach",
		"Citizen distress indicators",
		"Service standards deviation",
		"Documentation errors",
		"Communication barriers",
		"Policy misinterpretation",
		"Escalation requirements missed",
		"Successful resolution achieved",
		"Proper referral needed",
	}
}
