package requests_test

import (
	"testing"
	"aeolyzer/internal/orchestrator"
	"aeolyzer/internal/orchestrator/requests"
)

// TestMemoryUpdateRequiresApproval verifies the hard stop against silent memory mutation.
// Tone rules and long-term context cannot be rewritten autonomously. A proposal is generated,
// but the state update remains locked until a matching ApprovedAction is received from Layer 2.
func TestMemoryUpdateRequiresApproval(t *testing.T) {
	decision := orchestrator.IntakeDecision{
		Intent:          "update_memory",
		ApprovedActions: []string{}, // Missing "memoryUpdate" explicit approval
	}
	
	err := requests.ValidateApprovalResult(decision, "memoryUpdate")
	if err == nil {
		t.Fatal("expected ErrMemoryApprovalMissing but got nil")
	}
	if err.Error() != "MEMORY_APPROVAL_MISSING" {
		t.Fatalf("expected MEMORY_APPROVAL_MISSING, got %v", err)
	}
}

// TestDeepResearchRequiresApproval validates that heavy, unconstrained compute tasks
// (like broad scraping or iterative deep research) require explicit hitl checkpoints.
func TestDeepResearchRequiresApproval(t *testing.T) {
	decision := orchestrator.IntakeDecision{
		Intent:          "content_research",
		ApprovedActions: []string{"someOtherAction"},
	}
	
	err := requests.ValidateApprovalResult(decision, "deepResearch")
	if err == nil {
		t.Fatal("expected ErrDeepResearchApprovalMissing but got nil")
	}
}
