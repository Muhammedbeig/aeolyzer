package drifttrust

import (
	"errors"
	"fmt"

	observabilityconfig "aeolyzer/internal/observability/config"
)

// LoopObservation contains sanitized node and request fingerprints.
type LoopObservation struct {
	NodeIDs            []string `json:"node_ids"`
	ActionFingerprints []string `json:"action_fingerprints"`
	ReplanCount        int      `json:"replan_count"`
}

// LoopResult identifies bounded loop-policy violations.
type LoopResult struct {
	Detected            bool     `json:"detected"`
	RepeatedNodes       []string `json:"repeated_nodes,omitempty"`
	RepeatedActions     []string `json:"repeated_actions,omitempty"`
	ReplanLimitExceeded bool     `json:"replan_limit_exceeded"`
}

// DetectLoops applies the configured repetition limits.
func DetectLoops(
	observation LoopObservation,
	policy observabilityconfig.DriftPolicy,
) (LoopResult, error) {
	if observation.ReplanCount < 0 {
		return LoopResult{}, errors.New("replan count must not be negative")
	}
	if policy.Loops.MaxRepeatedNode < 1 ||
		policy.Loops.MaxRepeatedToolSameParams < 1 ||
		policy.Loops.MaxReplanCount < 1 {
		return LoopResult{}, errors.New("loop policy limits are invalid")
	}

	repeatedNodes, err := repeatedAbove(
		observation.NodeIDs,
		policy.Loops.MaxRepeatedNode,
	)
	if err != nil {
		return LoopResult{}, fmt.Errorf("inspect repeated nodes: %w", err)
	}
	repeatedActions, err := repeatedAbove(
		observation.ActionFingerprints,
		policy.Loops.MaxRepeatedToolSameParams,
	)
	if err != nil {
		return LoopResult{}, fmt.Errorf("inspect repeated actions: %w", err)
	}
	replanExceeded := observation.ReplanCount > policy.Loops.MaxReplanCount
	return LoopResult{
		Detected:            len(repeatedNodes) > 0 || len(repeatedActions) > 0 || replanExceeded,
		RepeatedNodes:       repeatedNodes,
		RepeatedActions:     repeatedActions,
		ReplanLimitExceeded: replanExceeded,
	}, nil
}

func repeatedAbove(values []string, limit int) ([]string, error) {
	counts := make(map[string]int, len(values))
	for _, value := range values {
		if value == "" {
			return nil, errors.New("loop fingerprint must not be empty")
		}
		counts[value]++
	}
	var repeated []string
	for _, value := range values {
		if counts[value] <= limit {
			continue
		}
		alreadyAdded := false
		for _, existing := range repeated {
			if existing == value {
				alreadyAdded = true
				break
			}
		}
		if !alreadyAdded {
			repeated = append(repeated, value)
		}
	}
	return repeated, nil
}
