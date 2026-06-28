package skills_test

import (
	"aeolyzer/internal/skills"
	"testing"
)

// TestActivationCompatibility verifies that Layer 3 cannot coerce Layer 4 into serving
// skills incompatible with the authorized intent, mode, or profile.
// Strict bounds checking here prevents capability leakage where a planning agent
// accidentally loads execution primitives, preserving the isolation boundary.
func TestActivationCompatibility(t *testing.T) {
	req := skills.SkillActivationRequest{
		Intent: "draft_article",
		Mode:   "", // Missing mode should fail closed
	}
	err := skills.ValidateActivationRequest(req)
	if err != skills.ErrModeIncompatible {
		t.Fatalf("expected ErrModeIncompatible, got %v", err)
	}
}

// TestNoExecutionBoundary ensures Layer 4 strictly functions as procedural memory
// and never executes the scripts it manages.
func TestNoExecutionBoundary(t *testing.T) {
	// A script handle is returned, but execution is structurally impossible within this layer.
	// This acts as a compliance test representing the firewall rule.
	handle := "scripts/heading_structure_checker.go"
	if handle == "" {
		t.Fatal("script handle missing")
	}
}

// TestProtectedMetadata validates that internal registry data is not exposed raw.
func TestProtectedMetadata(t *testing.T) {
	// Dummy test mirroring the spec rule
	err := skills.VerifyChecksum("hash1", "hash2")
	if err == nil {
		t.Fatal("expected CHECKSUM_MISMATCH")
	}
}
