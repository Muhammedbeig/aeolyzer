package extensions_test

import (
	"testing"
	"aeolyzer/internal/extensions"
)

// TestNoExecutableBoundary ensures that Layer 5 immediately drops any presentation intent 
// attempting to smuggle raw HTML or executable scripts, honoring the declarative-only mandate.
func TestNoExecutableBoundary(t *testing.T) {
	intent := extensions.PresentationIntent{
		Surface:      "chat",
		FallbackText: "Hello <script>alert(1)</script>",
	}
	
	err := extensions.ValidatePresentationIntent(intent)
	if err != extensions.ErrUnsafePayload {
		t.Fatalf("expected ErrUnsafePayload for script injection, got %v", err)
	}
}

// TestSurfaceValidation ensures that unknown surfaces are structurally rejected, 
// preventing unsupported state mappings.
func TestSurfaceValidation(t *testing.T) {
	intent := extensions.PresentationIntent{
		Surface:      "", // Missing surface
	}
	
	err := extensions.ValidatePresentationIntent(intent)
	if err != extensions.ErrUnknownSurface {
		t.Fatalf("expected ErrUnknownSurface, got %v", err)
	}
}
