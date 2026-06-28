package contracts

import (
	"encoding/json"
	"testing"
)

func TestIntakeDecisionWireNamesRemainStable(t *testing.T) {
	data, err := json.Marshal(IntakeDecision{
		TraceID:     "trace-1",
		Intent:      IntentTopicDiscovery,
		Confidence:  0.9,
		PolicyState: PolicyStateAllowed,
		Mode:        ModePlan,
	})
	if err != nil {
		t.Fatalf("json.Marshal() failed: %v", err)
	}
	var object map[string]any
	if err := json.Unmarshal(data, &object); err != nil {
		t.Fatalf("json.Unmarshal() failed: %v", err)
	}
	for _, field := range []string{"trace_id", "intent", "confidence", "policy_state", "mode"} {
		if _, found := object[field]; !found {
			t.Errorf("IntakeDecision JSON is missing %q", field)
		}
	}
}
