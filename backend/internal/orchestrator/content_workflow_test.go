package orchestrator_test

import (
	"testing"
	"aeolyzer/internal/orchestrator"
	"aeolyzer/internal/orchestrator/handoff"
)

// TestValidateContentGenerationTask guards against execution delegation missing required context.
// Layer 3 must marshal an explicit 'mode' and 'required_capabilities' before transferring
// the payload to Layer 6 to prevent runaway capability invocation.
func TestValidateContentGenerationTask(t *testing.T) {
	// Task missing Mode
	invalidTask := orchestrator.ContentGenerationTask{
		Intent:               "article_planning",
		RequiredCapabilities: []string{"content_strategy"},
	}
	
	err := handoff.ValidateContentGenerationTask(invalidTask)
	if err == nil {
		t.Fatal("expected ErrContentTaskInvalid for missing mode")
	}

	// Task missing capabilities
	invalidTask2 := orchestrator.ContentGenerationTask{
		Intent: "article_planning",
		Mode:   "plan",
	}
	
	err2 := handoff.ValidateContentGenerationTask(invalidTask2)
	if err2 == nil {
		t.Fatal("expected ErrContentTaskInvalid for missing capabilities")
	}
}
