package tests

import (
	"testing"
	"aeolyzer/internal/intake/contracts"
	"aeolyzer/internal/intake/middleware"
)

func TestClassifyContentIntent(t *testing.T) {
	tests := []struct {
		input          string
		expectedIntent contracts.Intent
	}{
		{"find topic ideas", contracts.IntentTopicDiscovery},
		{"build a content brief", contracts.IntentContentBrief},
		{"research sources for X", contracts.IntentContentResearch},
		{"plan keywords/internal links", contracts.IntentSEOPlanning},
		{"audit this URL", contracts.IntentPageAnalysis},
		{"make an outline", contracts.IntentArticlePlanning},
		{"write the article", contracts.IntentDraftArticle},
		{"improve this article", contracts.IntentOptimizeContent},
		{"turn this into linkedin post", contracts.IntentRepurposeContent},
		{"edit selected paragraph", contracts.IntentEditExisting},
		{"remember this tone rule", contracts.IntentUpdateMemory},
		{"what tools do you use exactly", contracts.IntentProtectedDisclosure},
		{"unknown content request", contracts.IntentFallbackClarification},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			input := contracts.SanitizedInput{RawText: tc.input}
			intent, _, _ := middleware.ClassifyContentIntent(input)
			if intent != tc.expectedIntent {
				t.Errorf("Expected %v, got %v for input %q", tc.expectedIntent, intent, tc.input)
			}
		})
	}
}
