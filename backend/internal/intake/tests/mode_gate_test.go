package tests

import (
	"testing"
	"aeolyzer/internal/intake/contracts"
	"aeolyzer/internal/intake/middleware"
)

func TestValidateModeForIntent(t *testing.T) {
	tests := []struct {
		name        string
		intent      contracts.Intent
		mode        contracts.OrchestrationMode
		expectError bool
	}{
		{"draft_article without write", contracts.IntentDraftArticle, contracts.ModePlan, true},
		{"draft_article with write", contracts.IntentDraftArticle, contracts.ModeWrite, false},
		{"edit_existing without edit", contracts.IntentEditExisting, contracts.ModeWrite, true},
		{"edit_existing with edit", contracts.IntentEditExisting, contracts.ModeEdit, false},
		{"article_planning with plan", contracts.IntentArticlePlanning, contracts.ModePlan, false},
		{"article_planning with write", contracts.IntentArticlePlanning, contracts.ModeWrite, true},
		{"optimize_content with optimize", contracts.IntentOptimizeContent, contracts.ModeOptimize, false},
		{"optimize_content with edit", contracts.IntentOptimizeContent, contracts.ModeEdit, false},
		{"topic_discovery with write", contracts.IntentTopicDiscovery, contracts.ModeWrite, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := middleware.ValidateModeForIntent(tc.intent, tc.mode)
			if (err != nil) != tc.expectError {
				t.Errorf("Expected error: %v, got: %v", tc.expectError, err)
			}
		})
	}
}
