// Package evaluationengine evaluates recorded agent behavior without executing
// tools, workflows, connectors, or remediation actions.
package evaluationengine

import (
	"errors"
	"fmt"
)

// TrajectoryMode controls how required action classes are matched.
type TrajectoryMode string

const (
	// TrajectoryExact requires the complete action-class sequence to match.
	TrajectoryExact TrajectoryMode = "EXACT"
	// TrajectoryInOrder requires an ordered subsequence and permits other actions.
	TrajectoryInOrder TrajectoryMode = "IN_ORDER"
	// TrajectoryAnyOrder requires all action classes in any order.
	TrajectoryAnyOrder TrajectoryMode = "ANY_ORDER"
	// TrajectoryForbiddenAbsent checks only that forbidden actions are absent.
	TrajectoryForbiddenAbsent TrajectoryMode = "FORBIDDEN_ABSENT"
)

// ActionEvent is one sanitized action-class observation from a trace.
type ActionEvent struct {
	Class string `json:"class"`
}

// ActionSpec identifies an expected or forbidden action class.
type ActionSpec struct {
	Class string `json:"class"`
}

// AgentTrace is the sanitized action sequence being evaluated.
type AgentTrace struct {
	TraceID string        `json:"trace_id"`
	Actions []ActionEvent `json:"actions"`
}

// TrajectorySpec defines required and forbidden behavior.
type TrajectorySpec struct {
	Mode      TrajectoryMode `json:"mode"`
	Required  []ActionSpec   `json:"required,omitempty"`
	Forbidden []ActionSpec   `json:"forbidden,omitempty"`
}

// TrajectoryScore records deterministic trajectory evaluation.
type TrajectoryScore struct {
	TraceID         string         `json:"trace_id"`
	Mode            TrajectoryMode `json:"mode"`
	RequiredMatched bool           `json:"required_matched"`
	ForbiddenAbsent bool           `json:"forbidden_absent"`
	Pass            bool           `json:"pass"`
	ObservedClasses []string       `json:"observed_classes"`
	FailureReasons  []string       `json:"failure_reasons,omitempty"`
}

// EvaluateTrajectory evaluates one sanitized trace against a trajectory spec.
func EvaluateTrajectory(trace AgentTrace, expected TrajectorySpec) (TrajectoryScore, error) {
	if trace.TraceID == "" {
		return TrajectoryScore{}, errors.New("trace id is required")
	}
	if err := validateTrajectorySpec(expected); err != nil {
		return TrajectoryScore{}, err
	}
	for _, action := range trace.Actions {
		if action.Class == "" {
			return TrajectoryScore{}, errors.New("observed action class is required")
		}
	}

	requiredMatched := false
	switch expected.Mode {
	case TrajectoryExact:
		requiredMatched = MatchExact(trace.Actions, expected.Required)
	case TrajectoryInOrder:
		requiredMatched = MatchInOrder(trace.Actions, expected.Required)
	case TrajectoryAnyOrder:
		requiredMatched = MatchAnyOrder(trace.Actions, expected.Required)
	case TrajectoryForbiddenAbsent:
		requiredMatched = true
	}
	forbiddenAbsent := VerifyForbiddenAbsent(trace.Actions, expected.Forbidden)

	score := TrajectoryScore{
		TraceID:         trace.TraceID,
		Mode:            expected.Mode,
		RequiredMatched: requiredMatched,
		ForbiddenAbsent: forbiddenAbsent,
		Pass:            requiredMatched && forbiddenAbsent,
		ObservedClasses: observedClasses(trace.Actions),
	}
	if !requiredMatched {
		score.FailureReasons = append(score.FailureReasons, "required action trajectory did not match")
	}
	if !forbiddenAbsent {
		score.FailureReasons = append(score.FailureReasons, "forbidden action class was observed")
	}
	return score, nil
}

// MatchExact reports whether actual and expected classes match exactly.
func MatchExact(actual []ActionEvent, expected []ActionSpec) bool {
	if len(actual) != len(expected) {
		return false
	}
	for i := range expected {
		if actual[i].Class != expected[i].Class {
			return false
		}
	}
	return true
}

// MatchInOrder reports whether expected is an ordered subsequence of actual.
func MatchInOrder(actual []ActionEvent, expected []ActionSpec) bool {
	if len(expected) == 0 {
		return true
	}
	next := 0
	for _, action := range actual {
		if action.Class == expected[next].Class {
			next++
			if next == len(expected) {
				return true
			}
		}
	}
	return false
}

// MatchAnyOrder reports whether actual contains every expected class with the
// required multiplicity.
func MatchAnyOrder(actual []ActionEvent, expected []ActionSpec) bool {
	counts := make(map[string]int, len(actual))
	for _, action := range actual {
		counts[action.Class]++
	}
	for _, action := range expected {
		if counts[action.Class] == 0 {
			return false
		}
		counts[action.Class]--
	}
	return true
}

// VerifyForbiddenAbsent reports whether no forbidden class appears in actual.
func VerifyForbiddenAbsent(actual []ActionEvent, forbidden []ActionSpec) bool {
	if len(forbidden) == 0 {
		return true
	}
	denied := make(map[string]struct{}, len(forbidden))
	for _, action := range forbidden {
		denied[action.Class] = struct{}{}
	}
	for _, action := range actual {
		if _, found := denied[action.Class]; found {
			return false
		}
	}
	return true
}

func validateTrajectorySpec(spec TrajectorySpec) error {
	switch spec.Mode {
	case TrajectoryExact, TrajectoryInOrder, TrajectoryAnyOrder:
		if len(spec.Required) == 0 {
			return fmt.Errorf("%s trajectory requires expected action classes", spec.Mode)
		}
	case TrajectoryForbiddenAbsent:
		if len(spec.Forbidden) == 0 {
			return errors.New("forbidden-absent trajectory requires forbidden action classes")
		}
	default:
		return fmt.Errorf("unsupported trajectory mode %q", spec.Mode)
	}
	for _, action := range append(append([]ActionSpec(nil), spec.Required...), spec.Forbidden...) {
		if action.Class == "" {
			return errors.New("trajectory action class is required")
		}
	}
	return nil
}

func observedClasses(actions []ActionEvent) []string {
	classes := make([]string, len(actions))
	for i, action := range actions {
		classes[i] = action.Class
	}
	return classes
}
