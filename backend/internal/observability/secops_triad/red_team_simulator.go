// Package secopstriad produces bounded security findings and recovery plans.
package secopstriad

import (
	"errors"
	"fmt"
)

// TestEnvironment is a safe Red Team execution environment.
type TestEnvironment string

const (
	// EnvironmentCI is an isolated CI run.
	EnvironmentCI TestEnvironment = "ci"
	// EnvironmentShadow is a no-side-effect shadow run.
	EnvironmentShadow TestEnvironment = "shadow"
	// EnvironmentStaging is an isolated staging run.
	EnvironmentStaging TestEnvironment = "staging"
)

// RedTeamFixture is a synthetic adversarial input with no live target.
type RedTeamFixture struct {
	ID           string `json:"id"`
	PayloadClass string `json:"payload_class"`
	Input        string `json:"input"`
}

// RedTeamResult records expected detection without retaining raw production
// data.
type RedTeamResult struct {
	FixtureID      string `json:"fixture_id"`
	PayloadClass   string `json:"payload_class"`
	Detected       bool   `json:"detected"`
	DetectionClass string `json:"detection_class,omitempty"`
}

// Detector evaluates one synthetic fixture.
type Detector interface {
	Detect(RedTeamFixture) (string, bool)
}

// RunRedTeamSuite runs only static synthetic fixtures in safe environments.
func RunRedTeamSuite(
	environment TestEnvironment,
	detector Detector,
	fixtures []RedTeamFixture,
) ([]RedTeamResult, error) {
	switch environment {
	case EnvironmentCI, EnvironmentShadow, EnvironmentStaging:
	default:
		return nil, errors.New("red team suite is forbidden in this environment")
	}
	if detector == nil || len(fixtures) == 0 || len(fixtures) > 1000 {
		return nil, errors.New("red team suite is not configured")
	}
	results := make([]RedTeamResult, 0, len(fixtures))
	seen := make(map[string]struct{}, len(fixtures))
	for _, fixture := range fixtures {
		if fixture.ID == "" || fixture.PayloadClass == "" || fixture.Input == "" {
			return nil, errors.New("red team fixture is invalid")
		}
		if _, duplicate := seen[fixture.ID]; duplicate {
			return nil, fmt.Errorf("duplicate red team fixture %q", fixture.ID)
		}
		seen[fixture.ID] = struct{}{}
		detectionClass, detected := detector.Detect(fixture)
		results = append(results, RedTeamResult{
			FixtureID:      fixture.ID,
			PayloadClass:   fixture.PayloadClass,
			Detected:       detected,
			DetectionClass: detectionClass,
		})
	}
	return results, nil
}
