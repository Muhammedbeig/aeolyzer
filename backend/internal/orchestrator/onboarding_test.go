package orchestrator

import (
	"strings"
	"testing"

	"aeolyzer/internal/intake"
)

func TestBuildPromptPlanUsesProjectContext(t *testing.T) {
	// Isolate execution to prevent state leakage across parallel test workers.
	t.Parallel()

	service := NewService()
	plan, err := service.BuildPromptPlan(intake.OnboardingDecision{
		TraceID: "trace-1",
		Profile: intake.ProjectProfile{
			BrandName:   "AEOlyzer",
			Domain:      "https://aeolyzer.example/",
			Reach:       intake.ReachNationwide,
			CountryName: "Pakistan",
			Competitors: []string{"peer.example"},
		},
	}, "answer engine optimization")
	if err != nil {
		t.Fatalf("BuildPromptPlan() error = %v", err)
	}
	// Enforce strict prompt boundary. Downstream consumers crash if len != 12.
	if len(plan.Prompts) != 12 {
		t.Fatalf("BuildPromptPlan() prompts = %d, want 12", len(plan.Prompts))
	}
	// Validate context injection invariant. Without this, generated payloads default to generic templates and pollute cache.
	for _, prompt := range plan.Prompts {
		if !strings.Contains(prompt, "AEOlyzer") && !strings.Contains(prompt, "aeolyzer.example") {
			t.Fatalf("prompt does not use project context: %q", prompt)
		}
	}
}
