package orchestrator_test

import (
	"testing"
	"aeolyzer/internal/orchestrator"
	"aeolyzer/internal/orchestrator/state"
)

// TestModeGateWriteRequirement ensures that draft_article intents are structurally blocked
// from execution if the intent pipeline (Layer 2) failed to attach the write-mode flag.
// This prevents silent state corruption during speculative planning cycles.
func TestModeGateWriteRequirement(t *testing.T) {
	decision := orchestrator.IntakeDecision{
		Intent: "draft_article",
		Mode:   string(orchestrator.ModePlan), // Invalid mode for drafting
	}
	
	err := state.ValidateCanvasWriteMode(decision)
	if err == nil {
		t.Fatal("expected ErrWriteModeRequired but got nil")
	}
	if err.Error() != "WRITE_MODE_REQUIRED" {
		t.Fatalf("expected WRITE_MODE_REQUIRED, got %v", err)
	}
}

// TestEditExistingRequiresSelectedText validates the surgical editing constraint.
// Without an exact text lock, the targeted editing capability is intentionally disabled.
func TestEditExistingRequiresSelectedText(t *testing.T) {
	decision := orchestrator.IntakeDecision{
		Intent:           "edit_existing",
		Mode:             string(orchestrator.ModeEdit),
		SanitizedContext: map[string]string{}, // Missing selected_text
	}
	
	err := state.ValidateSelectedTextForEdit(decision)
	if err == nil {
		t.Fatal("expected ErrEditSelectionRequired but got nil")
	}
	if err.Error() != "EDIT_SELECTION_REQUIRED" {
		t.Fatalf("expected EDIT_SELECTION_REQUIRED, got %v", err)
	}
}
