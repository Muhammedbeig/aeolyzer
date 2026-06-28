package orchestration

import (
	"strings"
	"testing"

	"aeolyzer/layer_02_intake"
)

func TestBuildPromptPlanUsesProjectContext(t *testing.T) {
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
	if len(plan.Prompts) != 12 {
		t.Fatalf("BuildPromptPlan() prompts = %d, want 12", len(plan.Prompts))
	}
	for _, prompt := range plan.Prompts {
		if !strings.Contains(prompt, "AEOlyzer") && !strings.Contains(prompt, "aeolyzer.example") {
			t.Fatalf("prompt does not use project context: %q", prompt)
		}
	}
}
