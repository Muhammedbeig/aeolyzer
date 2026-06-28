package intake_events

import (
	"testing"
	"time"

	"aeolyzer/internal/intake/contracts"
)

func TestSafetyEmitterAddsTimestampAndPreservesSafeFields(t *testing.T) {
	var received contracts.SafetyEvent
	emitter := NewSafetyEmitter(func(event contracts.SafetyEvent) {
		received = event
	})
	emitter.Emit(contracts.SafetyEvent{
		TraceID:   "trace-1",
		EventType: "policy_block",
		Decision:  "blocked",
	})
	if received.CreatedAt.IsZero() {
		t.Fatal("SafetyEmitter.Emit() did not add timestamp")
	}
	if received.TraceID != "trace-1" || received.Decision != "blocked" {
		t.Fatalf("SafetyEmitter.Emit() = %+v", received)
	}
	if received.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Fatal("SafetyEmitter.Emit() produced future timestamp")
	}
}
