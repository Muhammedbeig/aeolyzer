package evaluationengine

import (
	"context"
	"encoding/json"
	"testing"
)

type skillSelectionStub struct {
	unstable bool
	calls    int
}

func (s *skillSelectionStub) SelectSkills(
	_ context.Context,
	request SkillSelectionRequest,
) ([]byte, error) {
	s.calls++
	targetAlias := ""
	for _, candidate := range request.Candidates {
		if candidate.Description == "Target skill description for routing." {
			targetAlias = candidate.Alias
		}
	}
	results := make([]map[string]any, 0, len(request.Cases))
	for _, test := range request.Cases {
		selected := targetAlias
		if test.Input == "reject this request" {
			selected = "none"
		}
		if s.unstable && s.calls == 2 {
			selected = "none"
		}
		results = append(results, map[string]any{
			"case_id":            test.CaseID,
			"selected_candidate": selected,
			"confidence":         0.97,
			"summary":            "The request matches the bounded candidate description.",
		})
	}
	return json.Marshal(map[string]any{"results": results})
}

func TestEvaluateSkillTriggersPassesStableBatch(t *testing.T) {
	report, err := EvaluateSkillTriggers(
		context.Background(),
		&skillSelectionStub{},
		testSkillEvalBatch(),
		3,
		0.8,
	)
	if err != nil {
		t.Fatalf("EvaluateSkillTriggers() error = %v", err)
	}
	if !report.Passed ||
		report.PassedCases != 2 ||
		report.FailedCases != 0 ||
		report.FlakyCases != 0 {
		t.Fatalf("unexpected report: %+v", report)
	}
	if report.Cases[0].CaseHash == "" {
		t.Fatal("report omitted case hash")
	}
}

func TestEvaluateSkillTriggersDetectsInstability(t *testing.T) {
	report, err := EvaluateSkillTriggers(
		context.Background(),
		&skillSelectionStub{unstable: true},
		testSkillEvalBatch(),
		3,
		0.8,
	)
	if err != nil {
		t.Fatalf("EvaluateSkillTriggers() error = %v", err)
	}
	if report.Passed || report.FlakyCases == 0 || report.FailedCases == 0 {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestValidateSkillSelectionOutputRejectsUnknownFields(t *testing.T) {
	request, _ := buildSkillSelectionRequest(testSkillEvalBatch(), 0)
	raw := []byte(`{"results":[` +
		`{"case_id":"positive","selected_candidate":"candidate_1",` +
		`"confidence":0.9,"summary":"ok","reasoning":"hidden"}` +
		`]}`)
	if _, err := validateSkillSelectionOutput(raw, request); err == nil {
		t.Fatal("validateSkillSelectionOutput() error = nil, want error")
	}
}

func testSkillEvalBatch() SkillEvalBatch {
	return SkillEvalBatch{
		EvalID:         "skill-eval",
		CorpusChecksum: "sha256:test",
		TargetSkillID:  "target",
		Model:          "test-model",
		PromptVersion:  "v1",
		Candidates: []SkillEvalCandidate{
			{
				SkillID:      "target",
				Description:  "Target skill description for routing.",
				AntiTriggers: []string{"unrelated request"},
			},
			{
				SkillID:      "neighbor",
				Description:  "Neighbor skill description for routing.",
				AntiTriggers: []string{"target request"},
			},
			{
				SkillID:      "third",
				Description:  "Third skill description for routing.",
				AntiTriggers: []string{"target request"},
			},
			{
				SkillID:      "fourth",
				Description:  "Fourth skill description for routing.",
				AntiTriggers: []string{"target request"},
			},
			{
				SkillID:      "fifth",
				Description:  "Fifth skill description for routing.",
				AntiTriggers: []string{"target request"},
			},
		},
		Cases: []SkillEvalCase{
			{
				CaseID:        "positive",
				Group:         "positive",
				Input:         "use the target skill",
				ExpectedSkill: "target",
			},
			{
				CaseID:         "negative",
				Group:          "negative",
				Input:          "reject this request",
				ForbiddenSkill: "target",
				SafetyCritical: true,
			},
		},
	}
}
