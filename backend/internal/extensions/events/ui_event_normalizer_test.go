package events_test

import (
	"aeolyzer/internal/extensions/events"
	"testing"
	"time"
)

func TestStaleEventRejection(t *testing.T) {
	evt := events.UserInteractionEvent{
		Signature: "valid",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	if err := events.NormalizeInteraction(evt); err == nil {
		t.Fatal("expected stale event to be rejected")
	}
}
