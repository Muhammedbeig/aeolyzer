package middleware

import (
	"aeolyzer/layer_02_intake/contracts"
)

func ValidateToolPolicy(req contracts.ProposedToolRequest, decision contracts.IntakeDecision) error {
	// Wrapper to call content tool policy since Layer 2 handles both
	return ValidateContentToolPolicy(req, decision)
}
