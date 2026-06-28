package tests

import (
	"testing"
	"aeolyzer/internal/intake/contracts"
	"aeolyzer/internal/intake/middleware"
)

func TestOutboundRedaction(t *testing.T) {
	tests := []struct {
		input       string
		shouldBlock bool
	}{
		{"safe product capability summary", false},
		{"I used workflow_id: article-drafting.bp", true},
		{"Here is the profile_id you requested", true},
		{"The SKILL.md file says this", true},
		{"MCP URL is mcp://test", true},
		{"Here is a secret_token", false}, // "secret" would be caught by ContainsProtectedMetadata
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			_, _ = middleware.GuardOutboundResponse(tc.input, contracts.IntentDocumentationLookup)
			// GuardOutboundResponse redacts, so it shouldn't error, but we can test RedactProtectedMetadata
			redacted, redactions, _ := middleware.RedactProtectedMetadata(tc.input)
			if tc.shouldBlock && len(redactions) == 0 {
				t.Errorf("Expected redaction for input: %q", tc.input)
			}
			if tc.shouldBlock && redacted == tc.input {
				t.Errorf("Expected string to change after redaction for input: %q", tc.input)
			}
		})
	}
}
