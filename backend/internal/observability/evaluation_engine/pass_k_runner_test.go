package evaluationengine

import (
	"context"
	"errors"
	"testing"
)

type evalCaseStub struct {
	id      string
	results []EvalRunResult
}

func (e evalCaseStub) ID() string {
	return e.id
}

func (e evalCaseStub) Run(_ context.Context, run int) EvalRunResult {
	return e.results[run-1]
}

func TestRunPassKRequiresEveryRun(t *testing.T) {
	result, err := RunPassK(context.Background(), evalCaseStub{
		id: "reliability",
		results: []EvalRunResult{
			{Pass: true, Safety: true},
			{Pass: false, Safety: true, ErrorCode: "output_mismatch"},
			{Pass: true, Safety: true},
		},
	}, 3)
	if err != nil {
		t.Fatalf("RunPassK() failed: %v", err)
	}
	if result.AllPassed {
		t.Fatal("RunPassK().AllPassed = true, want false")
	}
	if result.FlakeRate != 1.0/3.0 {
		t.Fatalf("RunPassK().FlakeRate = %f, want %f", result.FlakeRate, 1.0/3.0)
	}
	if err := BlockIfFlaky(result, 0.4); !errors.Is(err, ErrPassKFailed) {
		t.Fatalf("BlockIfFlaky() error = %v, want %v", err, ErrPassKFailed)
	}
}

func TestBlockIfFlakyBlocksSingleSafetyFailure(t *testing.T) {
	result, err := RunPassK(context.Background(), evalCaseStub{
		id: "safety",
		results: []EvalRunResult{
			{Pass: true, Safety: true},
			{Pass: true, Safety: false},
		},
	}, 2)
	if err != nil {
		t.Fatalf("RunPassK() failed: %v", err)
	}
	if err := BlockIfFlaky(result, 0); !errors.Is(err, ErrSafetyEvaluationFailed) {
		t.Fatalf("BlockIfFlaky() error = %v, want %v", err, ErrSafetyEvaluationFailed)
	}
}

func TestRunPassKHonorsCancellationAndBounds(t *testing.T) {
	eval := evalCaseStub{
		id:      "cancelled",
		results: []EvalRunResult{{Pass: true, Safety: true}},
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := RunPassK(ctx, eval, 1); !errors.Is(err, context.Canceled) {
		t.Fatalf("RunPassK() error = %v, want %v", err, context.Canceled)
	}
	if _, err := RunPassK(context.Background(), eval, 0); err == nil {
		t.Fatal("RunPassK(k=0) returned nil error")
	}
	if _, err := RunPassK(context.Background(), eval, maxPassKRuns+1); err == nil {
		t.Fatal("RunPassK(k>max) returned nil error")
	}
}

func TestBlockIfFlakyAcceptsStableRuns(t *testing.T) {
	result, err := RunPassK(context.Background(), evalCaseStub{
		id: "stable",
		results: []EvalRunResult{
			{Pass: true, Safety: true},
			{Pass: true, Safety: true},
			{Pass: true, Safety: true},
		},
	}, 3)
	if err != nil {
		t.Fatalf("RunPassK() failed: %v", err)
	}
	if err := BlockIfFlaky(result, 0); err != nil {
		t.Fatalf("BlockIfFlaky() failed: %v", err)
	}
}
