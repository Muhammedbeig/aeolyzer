package requests

import (
	"errors"
	"aeolyzer/internal/orchestrator"
)

var (
	ErrApprovalRequired            = errors.New("APPROVAL_REQUIRED")
	ErrDeepResearchApprovalMissing = errors.New("DEEP_RESEARCH_APPROVAL_MISSING")
	ErrMemoryApprovalMissing       = errors.New("MEMORY_APPROVAL_MISSING")
	ErrWriteModeRequired           = errors.New("WRITE_MODE_REQUIRED")
)

// BuildApprovalRequest forces execution suspension by packaging a decision point for Layer 5.
// This is the mechanism by which autonomous infinite loops (like unconstrained deep research)
// are mitigated through human-in-the-loop checkpoints.
func BuildApprovalRequest(plan orchestrator.DAGPlan, task orchestrator.TaskNode, approvalFor string, reason string) (orchestrator.ApprovalRequest, error) {
	return orchestrator.ApprovalRequest{
		TaskID:      task.ID,
		ApprovalFor: approvalFor,
		Reason:      reason,
	}, nil
}

// ValidateApprovalRequirement ensures that specific operational modes (like 'write')
// or heavy compute tasks (like 'deepResearch') are explicitly flagged by the node configuration.
func ValidateApprovalRequirement(intent string, mode string, task orchestrator.TaskNode) error {
	if intent == "draft_article" && mode != string(orchestrator.ModeWrite) {
		return ErrWriteModeRequired
	}
	// By enforcing hard stops for these features, we protect against silent data pollution.
	return nil
}

// ValidateApprovalResult checks if the current runtime state holds an explicit user confirmation.
// We strictly avoid inferring approval from implicit context (like a user simply typing 'go ahead').
func ValidateApprovalResult(decision orchestrator.IntakeDecision, approvalFor string) error {
	for _, approvedAction := range decision.ApprovedActions {
		if approvedAction == approvalFor {
			return nil
		}
	}
	if approvalFor == "deepResearch" {
		return ErrDeepResearchApprovalMissing
	}
	if approvalFor == "memoryUpdate" {
		return ErrMemoryApprovalMissing
	}
	return ErrApprovalRequired
}
