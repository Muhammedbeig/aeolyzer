package skills

import (
	"errors"
)

var (
	ErrSkillUnknown         = errors.New("skill is unknown")
	ErrSkillBlocked         = errors.New("skill is blocked")
	ErrIntentIncompatible   = errors.New("intent is incompatible")
	ErrModeIncompatible     = errors.New("mode is incompatible")
	ErrProfileIncompatible  = errors.New("profile is incompatible")
	ErrTokenBudgetExceeded  = errors.New("token budget exceeded")
	ErrProtectedMetadata    = errors.New("protected metadata detected")
	ErrSchemasNotConfigured = errors.New("skill schemas are not configured")
	ErrChecksumMismatch     = errors.New("checksum mismatch")
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
		return ErrChecksumMismatch
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
