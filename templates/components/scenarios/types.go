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
