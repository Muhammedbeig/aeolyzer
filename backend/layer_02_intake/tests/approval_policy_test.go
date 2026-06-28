package tests

import (
	"testing"
	"time"
	"aeolyzer/layer_02_intake/contracts"
	"aeolyzer/layer_02_intake/middleware"
)

func TestApprovalValidation(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name        string
		intent      contracts.Intent
		actionClass string
		approvals   []contracts.ApprovedAction
		expectError bool
	}{
		{
			name:        "deep_research without approval",
			intent:      contracts.IntentContentResearch,
			actionClass: "deep_research",
			approvals:   []contracts.ApprovedAction{},
			expectError: true,
		},
		{
			name:        "deep_research with stale approval",
			intent:      contracts.IntentContentResearch,
			actionClass: "deep_research",
			approvals: []contracts.ApprovedAction{
				{ApprovalFor: "deep_research", ExpiresAt: now.Add(-time.Hour)},
			},
			expectError: true,
		},
		{
			name:        "deep_research with valid approval",
			intent:      contracts.IntentContentResearch,
			actionClass: "deep_research",
			approvals: []contracts.ApprovedAction{
				{ApprovalFor: "deep_research", ExpiresAt: now.Add(time.Hour)},
			},
			expectError: false,
		},
		{
			name:        "memory_update without approval",
			intent:      contracts.IntentUpdateMemory,
			actionClass: "propose_memory_update",
			approvals:   []contracts.ApprovedAction{},
			expectError: true,
		},
		{
			name:        "memory_update with free-text approval",
			intent:      contracts.IntentUpdateMemory,
			actionClass: "propose_memory_update",
			approvals: []contracts.ApprovedAction{
				{ApprovalFor: "memory_update", Source: "free_text", ExpiresAt: now.Add(time.Hour)},
			},
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := middleware.ValidateApprovalForTool(tc.intent, tc.actionClass, tc.approvals)
			// check free text explicitly if needed
			if len(tc.approvals) > 0 && tc.approvals[0].Source == "free_text" {
				err = middleware.RejectFreeTextApproval(tc.approvals[0].Source)
			}
			if (err != nil) != tc.expectError {
				t.Errorf("Expected error: %v, got: %v", tc.expectError, err)
			}
		})
	}
}
