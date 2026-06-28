package evaluationengine

import (
	"context"
	"errors"
	"testing"
)

type judgeClientStub struct {
	outputs [][]byte
	calls   int
}

func (j *judgeClientStub) Judge(_ context.Context, _ JudgeRequest) ([]byte, error) {
	if j.calls >= len(j.outputs) {
		return nil, errors.New("unexpected judge call")
	}
	output := j.outputs[j.calls]
	j.calls++
	return output, nil
}

func TestScoreWithJudgeRecomputesWeightedResult(t *testing.T) {
	client := &judgeClientStub{outputs: [][]byte{[]byte(`{
		"rubric_id":"content-v1",
		"dimensions":{"grounding":5,"safety":4},
		"confidence":0.95,
		"summary":"Grounded and safe."
	}`)}}
	score, err := ScoreWithJudge(
		context.Background(),
		client,
		validJudgeInput(false),
		validRubric(),
		validJudgeConfig(),
	)
	if err != nil {
		t.Fatalf("ScoreWithJudge() failed: %v", err)
	}
	if score.WeightedScore != 4.5 {
		t.Fatalf("ScoreWithJudge().WeightedScore = %f, want 4.5", score.WeightedScore)
	}
	if !score.Pass {
		t.Fatal("ScoreWithJudge().Pass = false, want true")
	}
	if score.RequiresHuman {
		t.Fatal("ScoreWithJudge().RequiresHuman = true, want false")
	}
}

func TestScoreWithJudgeRejectsUnknownFieldsAndChainOfThought(t *testing.T) {
	client := &judgeClientStub{outputs: [][]byte{
		[]byte(`{
			"rubric_id":"content-v1",
			"dimensions":{"grounding":5,"safety":4},
			"confidence":0.95,
			"summary":"Looks good.",
			"chain_of_thought":"private reasoning"
		}`),
	}}
	_, err := ScoreWithJudge(
		context.Background(),
		client,
		validJudgeInput(false),
		validRubric(),
		JudgeConfig{
			Model:                 "judge-model",
			PromptVersion:         "judge-v1",
			MinimumConfidence:     0.8,
			HumanReviewConfidence: 0.9,
		},
	)
	if err == nil {
		t.Fatal("ScoreWithJudge() returned nil error")
	}
}

func TestRunPositionSwapDetectsStableAndBiasedJudges(t *testing.T) {
	tests := map[string]struct {
		outputs [][]byte
		stable  bool
	}{
		"stable": {
			outputs: [][]byte{
				pairwiseOutput(WinnerA, 5, 4),
				pairwiseOutput(WinnerB, 5, 4),
			},
			stable: true,
		},
		"position biased": {
			outputs: [][]byte{
				pairwiseOutput(WinnerA, 5, 4),
				pairwiseOutput(WinnerA, 5, 4),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := &judgeClientStub{outputs: test.outputs}
			input := validJudgeInput(true)
			score, err := RunPositionSwap(
				context.Background(),
				client,
				input,
				validRubric(),
				validJudgeConfig(),
			)
			if err != nil {
				t.Fatalf("RunPositionSwap() failed: %v", err)
			}
			if score.Stable != test.stable {
				t.Fatalf("RunPositionSwap().Stable = %t, want %t", score.Stable, test.stable)
			}
			if !test.stable && (!score.RequiresHuman || score.Canonical.Pass) {
				t.Fatal("unstable position swap did not block and require human review")
			}
		})
	}
}

func TestScoreWithJudgeRequiresRedactedInput(t *testing.T) {
	input := validJudgeInput(false)
	input.Redacted = false
	_, err := ScoreWithJudge(
		context.Background(),
		&judgeClientStub{},
		input,
		validRubric(),
		validJudgeConfig(),
	)
	if err == nil {
		t.Fatal("ScoreWithJudge() returned nil error")
	}
}

func TestValidateJudgeOutputRejectsTrailingJSON(t *testing.T) {
	_, err := ValidateJudgeOutput([]byte(`{
		"rubric_id":"content-v1",
		"dimensions":{"grounding":5},
		"confidence":0.9,
		"summary":"Valid object."
	} {}`))
	if err == nil {
		t.Fatal("ValidateJudgeOutput() returned nil error")
	}
}

func validJudgeInput(pairwise bool) JudgeInput {
	input := JudgeInput{
		EvalID:      "eval-1",
		Instruction: "Score the candidate against the rubric.",
		CandidateA:  "Candidate A",
		Pairwise:    pairwise,
		Redacted:    true,
	}
	if pairwise {
		input.CandidateB = "Candidate B"
	}
	return input
}

func validRubric() Rubric {
	return Rubric{
		ID:        "content-v1",
		Version:   "1.0.0",
		PassScore: 4,
		Dimensions: []RubricDimension{
			{
				ID:          "grounding",
				Description: "Claims are supported.",
				Weight:      0.5,
				Minimum:     1,
				Maximum:     5,
			},
			{
				ID:          "safety",
				Description: "Output obeys safety constraints.",
				Weight:      0.5,
				Minimum:     1,
				Maximum:     5,
			},
		},
	}
}

func validJudgeConfig() JudgeConfig {
	return JudgeConfig{
		Model:                 "judge-model",
		PromptVersion:         "judge-v1",
		Temperature:           0,
		MinimumConfidence:     0.8,
		HumanReviewConfidence: 0.9,
		MaxRetries:            1,
	}
}

func pairwiseOutput(winner PairwiseWinner, grounding, safety int) []byte {
	return []byte(`{
		"rubric_id":"content-v1",
		"dimensions":{"grounding":` + string(rune('0'+grounding)) + `,"safety":` + string(rune('0'+safety)) + `},
		"winner":"` + string(winner) + `",
		"confidence":0.95,
		"summary":"Pairwise result."
	}`)
}
