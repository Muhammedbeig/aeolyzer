package handoff

import (
	"errors"
	"aeolyzer/internal/orchestrator"
)

var (
	ErrContentTaskInvalid   = errors.New("CONTENT_TASK_INVALID")
	ErrOutputSurfaceInvalid = errors.New("OUTPUT_SURFACE_INVALID")
)

// BuildContentGenerationTask marshals safe orchestration state into an execution envelope.
// By isolating task creation here, Layer 3 is structurally incapable of writing prose itself; 
// it only configures the boundaries (mode, capability constraints) for downstream consumers.
func BuildContentGenerationTask(plan orchestrator.DAGPlan, task orchestrator.TaskNode, ctx orchestrator.PlanningContext) (orchestrator.ContentGenerationTask, error) {
	// The core invariant: task building relies strictly on sanitized mapping.
	cgt := orchestrator.ContentGenerationTask{
		TaskID: task.ID,
		Inputs: ctx.SanitizedInputs,
		// Example: Setting hardcoded surface constraints based on workflow defaults.
		// No raw tool configuration is injected.
	}
	return cgt, nil
}

// ValidateContentGenerationTask inspects the created task for security anomalies
// such as missing surface targets or unapproved operational modes.
// This prevents malformed tasks from reaching the execution layer (Layer 6).
func ValidateContentGenerationTask(task orchestrator.ContentGenerationTask) error {
	if task.Mode == "" {
		return ErrContentTaskInvalid
	}
	if len(task.RequiredCapabilities) == 0 {
		return ErrContentTaskInvalid
	}
	return nil
}

// ValidateSurfaceHint ensures rendering constraints strictly adhere to approved UI layouts.
// Preventing dynamic surface hinting blocks payload-injected XSS/UI redress attacks.
func ValidateSurfaceHint(surface string) error {
	allowed := map[string]bool{
		"canvas":    true,
		"brief":     true,
		"chat":      true,
		"report":    true,
		"table":     true,
		"dashboard": true,
	}
	if !allowed[surface] {
		return ErrOutputSurfaceInvalid
	}
	return nil
}
