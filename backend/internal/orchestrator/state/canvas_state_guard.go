package state

import (
	"errors"
	"aeolyzer/internal/orchestrator"
)

var (
	ErrCanvasWriteRejected   = errors.New("CANVAS_WRITE_REJECTED")
	ErrEditSelectionRequired = errors.New("EDIT_SELECTION_REQUIRED")
	ErrWriteModeRequired     = errors.New("WRITE_MODE_REQUIRED")
)

// ValidateCanvasWriteMode acts as the primary firewall for content execution.
// If Layer 2 has not explicitly flagged the intent as 'write', all canvas mutation 
// nodes fail closed immediately. This protects against hallucinatory planning agents 
// generating unauthorized drafts.
func ValidateCanvasWriteMode(decision orchestrator.IntakeDecision) error {
	if decision.Mode != string(orchestrator.ModeWrite) {
		return ErrWriteModeRequired
	}
	return nil
}

// ValidateSelectedTextForEdit requires surgical targeting before any edit workflow begins.
// Without an exact text lock, the orchestrator refuses to proceed. This prevents
// catastrophic full-document overwrites by an unconstrained LLM.
func ValidateSelectedTextForEdit(decision orchestrator.IntakeDecision) error {
	if _, exists := decision.SanitizedContext["selected_text"]; !exists {
		return ErrEditSelectionRequired
	}
	return nil
}

// BuildCanvasChangeProposal delegates the actual content mutation to Layer 6 execution.
// By wrapping it in a ProposedToolRequest, Layer 3 avoids ever touching the raw canvas
// memory state or handling real-time hydration conflicts.
func BuildCanvasChangeProposal(plan orchestrator.DAGPlan, task orchestrator.TaskNode, change orchestrator.ContentGenerationTask) (orchestrator.ProposedToolRequest, error) {
	payload := make(map[string]interface{})
	for k, v := range change.Inputs {
		payload[k] = v
	}
	
	// Enforce that we only write to the targeted surface layer requested
	payload["target_surface"] = change.SurfaceHint
	
	return orchestrator.ProposedToolRequest{
		TaskID:  task.ID,
		Tool:    "writeCanvas", // Translates to the internal capability binding later
		Payload: payload,
	}, nil
}
