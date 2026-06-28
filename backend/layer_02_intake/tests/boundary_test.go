package tests

import (
	"testing"
	"aeolyzer/layer_02_intake/middleware"
	"aeolyzer/layer_02_intake/contracts"
)

// Boundary test checks that no Layer 3+ concepts are returned by Layer 2 functions
func TestBoundaryViolations(t *testing.T) {
	input := contracts.SanitizedInput{RawText: "write the article using workflow-id-123"}
	intent, _, err := middleware.ClassifyContentIntent(input)
	
	if err != nil {
		t.Fatalf("Expected classification to succeed, got %v", err)
	}
	
	// Layer 2 should only emit the intent, NOT workflow-id-123. 
	// The intent enum should be clean and strict.
	if string(intent) == "workflow-id-123" {
		t.Errorf("Layer 2 leaked a workflow ID as an intent. Boundary violation.")
	}
	
	if intent != contracts.IntentDraftArticle {
		t.Errorf("Expected draft article intent")
	}
}
