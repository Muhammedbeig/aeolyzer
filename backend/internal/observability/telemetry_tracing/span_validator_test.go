package telemetry_tracing_test

import (
	"testing"
	"aeolyzer/internal/observability/telemetry_tracing"
)

func TestSpanValidation(t *testing.T) {
	span := telemetry_tracing.Span{
		TraceID: "",
		SpanName: "test_span",
	}
	if err := telemetry_tracing.ValidateSpan(span); err == nil {
		t.Fatal("expected span without trace to fail validation")
	}
}
