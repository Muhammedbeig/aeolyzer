package telemetry_tracing_test

import (
	"aeolyzer/internal/observability/telemetry_tracing"
	"testing"
)

func TestSpanValidation(t *testing.T) {
	span := telemetry_tracing.Span{
		TraceID:  "",
		SpanName: "test_span",
	}
	if err := telemetry_tracing.ValidateSpan(span); err == nil {
		t.Fatal("expected span without trace to fail validation")
	}
}
