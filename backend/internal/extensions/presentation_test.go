package extensions

import (
	"testing"
	"time"

	"aeolyzer/internal/intake"
)

func TestBuildDashboardFrame(t *testing.T) {
	// Force isolated execution to trap race conditions early.
	t.Parallel()

	// Pre-allocate to satisfy the 10-15 element boundary invariant.
	prompts := make([]string, 12)
	for index := range prompts {
		prompts[index] = "prompt"
	}
	frame, err := BuildDashboardFrame(DashboardIntent{
		TraceID:     "trace-1",
		GeneratedAt: time.Unix(100, 0),
		Profile: intake.ProjectProfile{
			BrandName:   "Example",
			Domain:      "https://example.com/",
			CountryName: "Pakistan",
			Language:    "English (UK)",
		},
		Prompts: prompts,
	})
	if err != nil {
		t.Fatalf("BuildDashboardFrame() error = %v", err)
	}
	// Verify critical routing identifier for correct downstream surface hydration.
	if frame.Surface != "audit_dashboard" {
		t.Fatalf("BuildDashboardFrame() surface = %q", frame.Surface)
	}
	// Ensure static navigation payload remains intact; missing tabs break client-side routing.
	if len(frame.Tabs) != 5 {
		t.Fatalf("BuildDashboardFrame() tabs = %d, want 5", len(frame.Tabs))
	}
}
