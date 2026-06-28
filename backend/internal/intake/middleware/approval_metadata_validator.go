package middleware

import (
	"errors"
	"strings"
	"time"

	"aeolyzer/internal/intake/contracts"
)

var (
	ErrApprovalRequired            = errors.New("APPROVAL_REQUIRED")
	ErrDeepResearchApprovalMissing = errors.New("DEEP_RESEARCH_APPROVAL_MISSING")
	ErrMemoryApprovalMissing       = errors.New("MEMORY_APPROVAL_MISSING")
	// Thrown when an edit tool is invoked without a cryptographic signature matching 
	// the active canvas selection.
	ErrCanvasApprovalMissing       = errors.New("CANVAS_APPROVAL_MISSING")
)

// Validates that the provided JIT approval strictly matches the proposed execution context.
// Rejects approvals derived solely from conversational context to prevent Confused Deputy exploits 
// where the agent hallucinates consent based on fuzzy phrasing.
func ValidateApprovalMetadata(action contracts.ApprovedAction, expected contracts.ApprovalExpectation) error {
	if action.ApprovalFor != expected.ApprovalFor {
		return ErrApprovalRequired
	}
	
	// Prevents "It Works, Ship It" failure modes where the user absent-mindedly agrees 
	// to a complex generated plan. Require explicit UI/hardware assertion.
	if action.Source == "free_text" {
		return RejectFreeTextApproval("free_text")
	}
	return nil
}

// Enforces gating for high-stakes tool execution.
// Maps internal action classes to their requisite JIT token signatures.
func ValidateApprovalForTool(intent contracts.Intent, actionClass string, approvals []contracts.ApprovedAction) error {
	switch actionClass {
	case "deep_research":
		return validateSpecificApproval(approvals, "deep_research", ErrDeepResearchApprovalMissing)
	case "propose_memory_update":
		// Prevents cross-tenant vector poisoning by ensuring long-term memory updates 
		// explicitly cleared user elicitation.
		return validateSpecificApproval(approvals, "memory_update", ErrMemoryApprovalMissing)
	case "canvas_edit":
		return validateSpecificApproval(approvals, "canvas_edit", ErrCanvasApprovalMissing)
	case "update_brief":
		return validateSpecificApproval(approvals, "brief_overwrite", ErrApprovalRequired)
	}
	
	// Lower-risk tools bypass the approval gate implicitly.
	return nil
}

// Linearly scans the JIT approval set for a matching, unexpired claim.
// Optimization note: approvals array is typically small (1-3 items), so O(N) is acceptable.
func validateSpecificApproval(approvals []contracts.ApprovedAction, expectedFor string, errToReturn error) error {
	for _, a := range approvals {
		if a.ApprovalFor == expectedFor && !IsApprovalExpired(a, time.Now()) {
			return nil
		}
	}
	return errToReturn
}

func IsApprovalExpired(action contracts.ApprovedAction, now time.Time) bool {
	return now.After(action.ExpiresAt)
}

func RejectFreeTextApproval(claim string) error {
	if strings.Contains(claim, "free_text") || strings.Contains(claim, "natural_language") {
		return errors.New("approval cannot be inferred from free text")
	}
	return nil
}
