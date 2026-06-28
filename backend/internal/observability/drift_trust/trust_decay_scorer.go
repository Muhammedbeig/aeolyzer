package drifttrust

import (
	"errors"
	"fmt"
	"math"

	observabilityconfig "aeolyzer/internal/observability/config"
)

// TrustSignal records one sanitized policy-defined trust event.
type TrustSignal struct {
	Class string `json:"class"`
	Count int    `json:"count"`
}

// TrustInput contains the prior score and bounded signals for one trace.
type TrustInput struct {
	PreviousScore   float64       `json:"previous_score"`
	DecaySignals    []TrustSignal `json:"decay_signals,omitempty"`
	RecoverySignals []TrustSignal `json:"recovery_signals,omitempty"`
}

// TrustResult is a deterministic trust update and decision recommendation.
type TrustResult struct {
	PreviousScore       float64  `json:"previous_score"`
	Score               float64  `json:"score"`
	DecayApplied        float64  `json:"decay_applied"`
	RecoveryApplied     float64  `json:"recovery_applied"`
	CriticalSignals     []string `json:"critical_signals,omitempty"`
	RecommendQuarantine bool     `json:"recommend_quarantine"`
	Warning             bool     `json:"warning"`
}

// ScoreTrust applies policy-defined decay and bounded recovery factors.
func ScoreTrust(
	input TrustInput,
	driftPolicy observabilityconfig.DriftPolicy,
	trustPolicy observabilityconfig.TrustPolicy,
) (TrustResult, error) {
	if input.PreviousScore < trustPolicy.Score.Minimum ||
		input.PreviousScore > trustPolicy.Score.Maximum {
		return TrustResult{}, errors.New("previous trust score is outside policy bounds")
	}

	decay, err := applySignals(input.DecaySignals, driftPolicy.TrustDecay.DecayFactors)
	if err != nil {
		return TrustResult{}, fmt.Errorf("apply trust decay: %w", err)
	}
	recovery, err := applySignals(
		input.RecoverySignals,
		driftPolicy.TrustDecay.RecoveryFactors,
	)
	if err != nil {
		return TrustResult{}, fmt.Errorf("apply trust recovery: %w", err)
	}
	recovery = math.Min(recovery, trustPolicy.Score.RecoveryCapPerTrace)

	criticalSet := make(map[string]struct{}, len(trustPolicy.CriticalSignals))
	for _, signal := range trustPolicy.CriticalSignals {
		criticalSet[signal] = struct{}{}
	}
	var critical []string
	for _, signal := range input.DecaySignals {
		if signal.Count < 1 {
			return TrustResult{}, errors.New("trust signal count must be positive")
		}
		if _, found := criticalSet[signal.Class]; found {
			critical = append(critical, signal.Class)
		}
	}

	score := clamp(
		input.PreviousScore-decay+recovery,
		trustPolicy.Score.Minimum,
		trustPolicy.Score.Maximum,
	)
	recommendQuarantine := score <= trustPolicy.Score.QuarantineDecision ||
		(len(critical) > 0 &&
			trustPolicy.Rules.CriticalSignalRequiresQuarantineDecision)
	return TrustResult{
		PreviousScore:       input.PreviousScore,
		Score:               score,
		DecayApplied:        decay,
		RecoveryApplied:     recovery,
		CriticalSignals:     critical,
		RecommendQuarantine: recommendQuarantine,
		Warning:             score <= trustPolicy.Score.Warning,
	}, nil
}

func applySignals(signals []TrustSignal, factors map[string]float64) (float64, error) {
	total := 0.0
	for _, signal := range signals {
		if signal.Class == "" || signal.Count < 1 || signal.Count > 100 {
			return 0, errors.New("trust signal is invalid")
		}
		factor, found := factors[signal.Class]
		if !found {
			return 0, fmt.Errorf("unknown trust signal %q", signal.Class)
		}
		total += factor * float64(signal.Count)
	}
	return total, nil
}

func clamp(value, minimum, maximum float64) float64 {
	return math.Max(minimum, math.Min(maximum, value))
}
