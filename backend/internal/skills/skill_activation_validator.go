package skills

import (
	"errors"
)

var (
	ErrSkillUnknown        = errors.New("SKILL_UNKNOWN")
	ErrSkillBlocked        = errors.New("SKILL_BLOCKED")
	ErrIntentIncompatible  = errors.New("INTENT_INCOMPATIBLE")
	ErrModeIncompatible    = errors.New("MODE_INCOMPATIBLE")
	ErrProfileIncompatible = errors.New("PROFILE_INCOMPATIBLE")
	ErrTokenBudgetExceeded = errors.New("TOKEN_BUDGET_EXCEEDED")
	ErrProtectedMetadata   = errors.New("PROTECTED_METADATA")
)

// ValidateActivationRequest enforces Layer 4 strict progressive disclosure boundaries.
// Note: We only validate compatibility; Layer 4 does NOT modify intents or choose workflows.
// This design forces all state routing logic to remain strictly in Layer 3, leaving this
// layer as a pure procedural memory store without orchestration side-effects.
func ValidateActivationRequest(req SkillActivationRequest) error {
	if req.Intent == "" {
		return ErrIntentIncompatible
	}
	if req.Mode == "" {
		return ErrModeIncompatible
	}
	if req.ProfileID == "" {
		return ErrProfileIncompatible
	}
	if req.MaxTokenBudget <= 0 {
		return ErrTokenBudgetExceeded
	}
	return nil
}

// ChecksumVerifier acts as a supply-chain firewall ensuring resources are untampered.
func VerifyChecksum(expected, actual string) error {
	if expected != actual {
		return errors.New("CHECKSUM_MISMATCH")
	}
	return nil
}

// TokenEstimator provides heuristics for context-window budgeting without injecting payload.
func EstimateTokens(bodyLength int, resourceLengths []int) int {
	total := bodyLength
	for _, l := range resourceLengths {
		total += l
	}
	return total
}
