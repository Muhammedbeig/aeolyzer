package orchestrator_test

import (
	"testing"
	"aeolyzer/internal/orchestrator/handoff"
)

// TestValidateSurfaceHint ensures that orchestrator handoffs rigidly conform
// to the known presentation layers. An LLM cannot spontaneously invent a new surface
// (e.g., 'system_shell' or 'eval_override') as a vector for output injection.
func TestValidateSurfaceHint(t *testing.T) {
	// Valid surfaces
	if err := handoff.ValidateSurfaceHint("canvas"); err != nil {
		t.Errorf("expected valid for canvas, got %v", err)
	}
	if err := handoff.ValidateSurfaceHint("brief"); err != nil {
		t.Errorf("expected valid for brief, got %v", err)
	}

	// Invalid surface
	err := handoff.ValidateSurfaceHint("arbitrary_path_or_unknown_surface")
	if err == nil {
		t.Fatal("expected ErrOutputSurfaceInvalid for unknown surface")
	}
	if err.Error() != "OUTPUT_SURFACE_INVALID" {
		t.Fatalf("expected OUTPUT_SURFACE_INVALID, got %v", err)
	}
}
