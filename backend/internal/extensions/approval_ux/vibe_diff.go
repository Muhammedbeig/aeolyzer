package approval_ux

import "errors"

// VibeDiff represents a strictly constrained summary of a state change (Section 12.2).
// It decouples the approval presentation from the underlying system state.
type VibeDiff struct {
	Summary        string      `json:"summary"`
	ChangeType     string      `json:"change_type"`
	Before         interface{} `json:"before,omitempty"`
	After          interface{} `json:"after,omitempty"`
	RiskNotes      []string    `json:"risk_notes,omitempty"`
}

// ValidateVibeDiff ensures the diff does not violate data leakage boundaries.
// By statically rejecting internal identifiers, it maintains the abstraction 
// between the user-facing approval card and the orchestrator's state machine.
func ValidateVibeDiff(diff VibeDiff) error {
	if diff.Summary == "" {
		return errors.New("MISSING_SUMMARY")
	}
	
	// Prevent leaking internal system policies into the user approval screen.
	for _, note := range diff.RiskNotes {
		if note == "policy.yaml" || note == "trace_id" {
			return errors.New("INTERNAL_METADATA_LEAK")
		}
	}
	
	return nil
}
