// Package drifttrust detects sanitized intent drift, trust decay, and loops.
package drifttrust

import (
	"errors"
	"math"

	observabilityconfig "aeolyzer/internal/observability/config"
)

// DriftLevel is the policy outcome for a drift score.
type DriftLevel string

const (
	// DriftNormal is below the warning threshold.
	DriftNormal DriftLevel = "normal"
	// DriftWarning requires observation but does not recommend a block.
	DriftWarning DriftLevel = "warning"
	// DriftBlock recommends blocking or replanning through the owning layers.
	DriftBlock DriftLevel = "block_recommendation"
	// DriftQuarantine recommends a Layer 8 quarantine decision for Layer 6.
	DriftQuarantine DriftLevel = "quarantine_decision"
)

// DriftObservation contains only sanitized categorical and hashed context.
type DriftObservation struct {
	Layer2Intent            string   `json:"layer2_intent"`
	PlannedIntent           string   `json:"planned_intent"`
	ActiveMode              string   `json:"active_mode"`
	ObservedMode            string   `json:"observed_mode"`
	ApprovedGoalSummaryHash string   `json:"approved_goal_summary_hash"`
	ObservedGoalSummaryHash string   `json:"observed_goal_summary_hash"`
	AuthorizedActionClasses []string `json:"authorized_action_classes"`
	ObservedActionClasses   []string `json:"observed_action_classes"`
}

// DriftResult is a deterministic drift assessment.
type DriftResult struct {
	Score               float64    `json:"score"`
	Level               DriftLevel `json:"level"`
	IntentMismatch      bool       `json:"intent_mismatch"`
	ModeMismatch        bool       `json:"mode_mismatch"`
	GoalMismatch        bool       `json:"goal_mismatch"`
	UnauthorizedActions []string   `json:"unauthorized_actions,omitempty"`
}

// DetectIntentDrift scores categorical differences without reading raw prompts
// or hidden reasoning.
func DetectIntentDrift(
	observation DriftObservation,
	policy observabilityconfig.DriftPolicy,
) (DriftResult, error) {
	if observation.Layer2Intent == "" ||
		observation.PlannedIntent == "" ||
		observation.ActiveMode == "" ||
		observation.ObservedMode == "" {
		return DriftResult{}, errors.New("drift intent and mode context is required")
	}
	if observation.ApprovedGoalSummaryHash == "" ||
		observation.ObservedGoalSummaryHash == "" {
		return DriftResult{}, errors.New("drift goal hashes are required")
	}
	if len(observation.AuthorizedActionClasses) == 0 {
		return DriftResult{}, errors.New("authorized action classes are required")
	}

	intentMismatch := observation.Layer2Intent != observation.PlannedIntent
	modeMismatch := observation.ActiveMode != observation.ObservedMode
	goalMismatch := observation.ApprovedGoalSummaryHash != observation.ObservedGoalSummaryHash
	unauthorized := unauthorizedActions(
		observation.AuthorizedActionClasses,
		observation.ObservedActionClasses,
	)

	score := 0.0
	if intentMismatch {
		score += 0.25
	}
	if modeMismatch {
		score += 0.35
	}
	if goalMismatch {
		score += 0.20
	}
	if len(observation.ObservedActionClasses) > 0 {
		ratio := float64(len(unauthorized)) / float64(len(observation.ObservedActionClasses))
		score += math.Min(0.40, ratio*0.40)
	}
	score = math.Min(1, score)

	return DriftResult{
		Score:               score,
		Level:               classifyDrift(score, policy),
		IntentMismatch:      intentMismatch,
		ModeMismatch:        modeMismatch,
		GoalMismatch:        goalMismatch,
		UnauthorizedActions: unauthorized,
	}, nil
}

func classifyDrift(score float64, policy observabilityconfig.DriftPolicy) DriftLevel {
	thresholds := policy.IntentDrift.Thresholds
	switch {
	case score >= thresholds.QuarantineDecision:
		return DriftQuarantine
	case score >= thresholds.BlockRecommendation:
		return DriftBlock
	case score >= thresholds.Warn:
		return DriftWarning
	default:
		return DriftNormal
	}
}

func unauthorizedActions(allowed, observed []string) []string {
	allowlist := make(map[string]struct{}, len(allowed))
	for _, action := range allowed {
		allowlist[action] = struct{}{}
	}
	seen := make(map[string]struct{})
	var result []string
	for _, action := range observed {
		if _, allowed := allowlist[action]; allowed {
			continue
		}
		if _, duplicate := seen[action]; duplicate {
			continue
		}
		seen[action] = struct{}{}
		result = append(result, action)
	}
	return result
}
