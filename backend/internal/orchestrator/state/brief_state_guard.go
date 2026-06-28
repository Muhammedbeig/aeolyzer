package state

import (
	"errors"
	"aeolyzer/internal/orchestrator"
)

var (
	ErrBriefOverwriteRejected = errors.New("BRIEF_OVERWRITE_REJECTED")
)

// ValidateBriefRequiredFields verifies that the sanitized payload contains essential planning signals.
// By strictly matching keys, we eliminate the risk of executing partial content generation 
// that leads to hallucinatory interpolation by downstream LLMs.
func ValidateBriefRequiredFields(ctx map[string]string, required []string) error {
	for _, req := range required {
		if _, exists := ctx[req]; !exists {
			return errors.New("MISSING_REQUIRED_FIELD_" + req)
		}
	}
	return nil
}

// BuildBriefUpdateProposal packages incremental brief modifications into a safe request.
// Layer 3 never directly accesses the filesystem to write brief files. It merely requests 
// that the Layer 6 execution environment perform the state mutation.
func BuildBriefUpdateProposal(plan orchestrator.DAGPlan, task orchestrator.TaskNode, fields map[string]string) (orchestrator.ProposedToolRequest, error) {
	payload := make(map[string]interface{})
	for k, v := range fields {
		payload[k] = v
	}
	return orchestrator.ProposedToolRequest{
		TaskID:  task.ID,
		Tool:    "updateBrief",
		Payload: payload,
	}, nil
}

// ValidateNoBriefOverwrite prevents destructive modifications to already-approved strategic boundaries.
// By checking against locked keys, we protect against prompt-injection attacks aiming to quietly pivot 
// the topic or brand constraints mid-workflow.
func ValidateNoBriefOverwrite(fields map[string]string, existingKeys []string) error {
	existingMap := make(map[string]bool)
	for _, k := range existingKeys {
		existingMap[k] = true
	}
	for k := range fields {
		if existingMap[k] {
			return ErrBriefOverwriteRejected
		}
	}
	return nil
}
