// templates/components/scenarios/types.go
package scenarios

type Scenario struct {
	ID              string
	Name            string
	Description     string
	Category        string
	Difficulty      int
	Duration        int
	Scene           string
	BackgroundNoise bool
	SuccessCriteria string
	Keywords        string
}

// ScenarioCategories returns available scenario categories for public service training
func ScenarioCategories() []string {
	return []string{
		"Document Processing",
		"Identity Verification",
		"Service Application",
		"Complaint Handling",
		"Emergency Assistance",
		"Information Request",
		"Payment Processing",
		"Language Barrier",
		"Special Needs Assistance",
	}
}

// SceneTypes returns available scene settings for service counter interactions
func SceneTypes() []string {
	return []string{
		"Main Service Counter",
		"Private Consultation Room",
		"Waiting Area",
		"Information Desk",
		"Document Processing Station",
		"Quick Service Counter",
		"Special Assistance Desk",
	}
}

// SuccessCriteriaTypes returns predefined success criteria for scenarios
func SuccessCriteriaTypes() []string {
	return []string{
		"Correct Document Processing",
		"Clear Communication",
		"Proper Procedure Following",
		"Effective Conflict Resolution",
		"Accurate Information Provided",
		"Appropriate Referral Made",
		"Customer Satisfaction Achieved",
		"Privacy Guidelines Followed",
		"Service Standards Met",
	}
}
