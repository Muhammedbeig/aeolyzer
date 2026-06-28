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

// TestProtectedMetadata validates that internal registry data is not exposed raw.
func TestProtectedMetadata(t *testing.T) {
	err := skills.VerifyChecksum("hash1", "hash2")
	if err != skills.ErrChecksumMismatch {
		t.Fatalf("VerifyChecksum() error = %v, want %v", err, skills.ErrChecksumMismatch)
	}
}
