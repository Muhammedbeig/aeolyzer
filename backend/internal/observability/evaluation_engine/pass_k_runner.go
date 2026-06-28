package evaluationengine

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const maxPassKRuns = 100

var (
	// ErrPassKFailed blocks release when any required run fails.
	ErrPassKFailed = errors.New("pass k reliability gate failed")
	// ErrFlakeRateExceeded blocks release when instability exceeds policy.
	ErrFlakeRateExceeded = errors.New("evaluation flake rate exceeds threshold")
	// ErrSafetyEvaluationFailed blocks release after any safety failure.
	ErrSafetyEvaluationFailed = errors.New("safety evaluation failed")
)

// EvalRunResult is one independent evaluation execution.
type EvalRunResult struct {
	Run       int           `json:"run"`
	Pass      bool          `json:"pass"`
	Safety    bool          `json:"safety"`
	Duration  time.Duration `json:"duration"`
	ErrorCode string        `json:"error_code,omitempty"`
}

// EvalCase executes one evaluation run. Implementations must be deterministic
// where possible and must honor context cancellation.
type EvalCase interface {
	ID() string
	Run(context.Context, int) EvalRunResult
}

// PassKResult contains repeated reliability evidence.
type PassKResult struct {
	EvalID       string          `json:"eval_id"`
	K            int             `json:"k"`
	Runs         []EvalRunResult `json:"runs"`
	AllPassed    bool            `json:"all_passed"`
	FlakeRate    float64         `json:"flake_rate"`
	SafetyPassed bool            `json:"safety_passed"`
}

// RunPassK runs an evaluation sequentially k times with cancellation.
func RunPassK(ctx context.Context, eval EvalCase, k int) (PassKResult, error) {
	if eval == nil {
		return PassKResult{}, errors.New("evaluation case is required")
	}
	if eval.ID() == "" {
		return PassKResult{}, errors.New("evaluation id is required")
	}
	if k < 1 || k > maxPassKRuns {
		return PassKResult{}, fmt.Errorf("pass k must be between 1 and %d", maxPassKRuns)
	}

	results := make([]EvalRunResult, 0, k)
	for run := 1; run <= k; run++ {
		if err := ctx.Err(); err != nil {
			return PassKResult{}, fmt.Errorf("run pass k: %w", err)
		}
		started := time.Now()
		result := eval.Run(ctx, run)
		result.Run = run
		if result.Duration <= 0 {
			result.Duration = time.Since(started)
		}
		results = append(results, result)
	}
	return PassKResult{
		EvalID:       eval.ID(),
		K:            k,
		Runs:         results,
		AllPassed:    RequireAllPass(results),
		FlakeRate:    CalculateFlakeRate(results),
		SafetyPassed: requireSafetyPass(results),
	}, nil
}

// RequireAllPass reports whether every run passed.
func RequireAllPass(results []EvalRunResult) bool {
	if len(results) == 0 {
		return false
	}
	for _, result := range results {
		if !result.Pass {
			return false
		}
	}
	return true
}

// CalculateFlakeRate returns the fraction of failed runs.
func CalculateFlakeRate(results []EvalRunResult) float64 {
	if len(results) == 0 {
		return 1
	}
	failures := 0
	for _, result := range results {
		if !result.Pass {
			failures++
		}
	}
	return float64(failures) / float64(len(results))
}

// BlockIfFlaky applies safety, all-pass, and flake-rate release gates.
func BlockIfFlaky(result PassKResult, threshold float64) error {
	if threshold < 0 || threshold > 1 {
		return errors.New("flake threshold must be between zero and one")
	}
	if !result.SafetyPassed {
		return ErrSafetyEvaluationFailed
	}
	if !result.AllPassed {
		return ErrPassKFailed
	}
	if result.FlakeRate > threshold {
		return ErrFlakeRateExceeded
	}
	return nil
}

func requireSafetyPass(results []EvalRunResult) bool {
	if len(results) == 0 {
		return false
	}
	for _, result := range results {
		if !result.Safety {
			return false
		}
	}
	return true
}
