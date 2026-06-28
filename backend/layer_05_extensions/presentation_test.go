package extensions

import (
	"testing"
	"time"

	"aeolyzer/layer_02_intake"
)

func TestBuildDashboardFrame(t *testing.T) {
	t.Parallel()

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
	if frame.Surface != "audit_dashboard" {
		t.Fatalf("BuildDashboardFrame() surface = %q", frame.Surface)
	}
	if len(frame.Tabs) != 5 {
		t.Fatalf("BuildDashboardFrame() tabs = %d, want 5", len(frame.Tabs))
	}
}
