package event_ingestion_test

import (
	"aeolyzer/internal/observability/event_ingestion"
	"strings"
	"testing"
)

func TestEventRedaction(t *testing.T) {
	out := event_ingestion.RedactEvent("this contains raw_system_prompt inside")
	if strings.Contains(out, "raw_system_prompt") {
		t.Fatal("expected sensitive context to be redacted")
	}
}
