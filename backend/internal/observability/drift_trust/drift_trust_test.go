package drifttrust

import (
	"testing"

	observabilityconfig "aeolyzer/internal/observability/config"
)

func TestDetectIntentDriftEscalatesModeAndUnauthorizedActions(t *testing.T) {
	policies := testPolicies(t)
	result, err := DetectIntentDrift(DriftObservation{
		Layer2Intent:            "draft_article",
		PlannedIntent:           "draft_article",
		ActiveMode:              "write",
		ObservedMode:            "act",
		ApprovedGoalSummaryHash: "hash-a",
		ObservedGoalSummaryHash: "hash-b",
		AuthorizedActionClasses: []string{"brief_read", "draft_section"},
		ObservedActionClasses:   []string{"draft_section", "external_publish"},
	}, policies.Drift)
	if err != nil {
		t.Fatalf("DetectIntentDrift() failed: %v", err)
	}
	if result.Level != DriftQuarantine {
		t.Fatalf("DetectIntentDrift().Level = %s, want %s", result.Level, DriftQuarantine)
	}
	if len(result.UnauthorizedActions) != 1 ||
		result.UnauthorizedActions[0] != "external_publish" {
		t.Fatalf(
			"DetectIntentDrift().UnauthorizedActions = %v, want external_publish",
			result.UnauthorizedActions,
		)
	}
}

func TestDetectIntentDriftAcceptsAlignedObservation(t *testing.T) {
	policies := testPolicies(t)
	result, err := DetectIntentDrift(DriftObservation{
		Layer2Intent:            "topic_discovery",
		PlannedIntent:           "topic_discovery",
		ActiveMode:              "plan",
		ObservedMode:            "plan",
		ApprovedGoalSummaryHash: "hash",
		ObservedGoalSummaryHash: "hash",
		AuthorizedActionClasses: []string{"source_read"},
		ObservedActionClasses:   []string{"source_read"},
	}, policies.Drift)
	if err != nil {
		t.Fatalf("DetectIntentDrift() failed: %v", err)
	}
	if result.Level != DriftNormal || result.Score != 0 {
		t.Fatalf("DetectIntentDrift() = %+v, want normal score zero", result)
	}
}

func TestScoreTrustCapsRecoveryAndEscalatesCriticalSignal(t *testing.T) {
	policies := testPolicies(t)
	result, err := ScoreTrust(TrustInput{
		PreviousScore: 1,
		DecaySignals: []TrustSignal{
			{Class: "cross_tenant_signal", Count: 1},
		},
		RecoverySignals: []TrustSignal{
			{Class: "eval_pass_after_repair", Count: 10},
		},
	}, policies.Drift, policies.Trust)
	if err != nil {
		t.Fatalf("ScoreTrust() failed: %v", err)
	}
	if result.RecoveryApplied != policies.Trust.Score.RecoveryCapPerTrace {
		t.Fatalf(
			"ScoreTrust().RecoveryApplied = %f, want cap %f",
			result.RecoveryApplied,
			policies.Trust.Score.RecoveryCapPerTrace,
		)
	}
	if !result.RecommendQuarantine {
		t.Fatal("ScoreTrust().RecommendQuarantine = false, want true")
	}
}

func TestDetectLoopsFindsNodeActionAndReplanLoops(t *testing.T) {
	policies := testPolicies(t)
	result, err := DetectLoops(LoopObservation{
		NodeIDs:            []string{"a", "a", "a", "a"},
		ActionFingerprints: []string{"hash", "hash", "hash"},
		ReplanCount:        5,
	}, policies.Drift)
	if err != nil {
		t.Fatalf("DetectLoops() failed: %v", err)
	}
	if !result.Detected ||
		len(result.RepeatedNodes) != 1 ||
		len(result.RepeatedActions) != 1 ||
		!result.ReplanLimitExceeded {
		t.Fatalf("DetectLoops() = %+v, want all loop signals", result)
	}
}

func testPolicies(t *testing.T) observabilityconfig.Policies {
	t.Helper()
	policies, err := observabilityconfig.LoadEmbeddedPolicies()
	if err != nil {
		t.Fatalf("LoadEmbeddedPolicies() failed: %v", err)
	}
	return policies
}
