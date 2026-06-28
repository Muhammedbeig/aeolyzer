package telemetrytracing_test

import (
	"testing"

	telemetrytracing "aeolyzer/internal/observability/telemetry_tracing"
)

func TestSpanValidation(t *testing.T) {
	span := telemetrytracing.Span{
		TraceID:  "",
		SpanName: "test_span",
	}
	if err := telemetrytracing.ValidateSpan(span); err == nil {
		t.Fatal("expected span without trace to fail validation")
	}
}
