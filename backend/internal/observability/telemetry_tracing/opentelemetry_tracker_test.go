package telemetrytracing

import (
	"context"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestTrackerEmitsRootAndChildSpansWithSafeAttributes(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	defer func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			t.Errorf("TracerProvider.Shutdown() failed: %v", err)
		}
	}()
	tracker, err := NewTracker(provider, "1.0.0")
	if err != nil {
		t.Fatalf("NewTracker() failed: %v", err)
	}
	ctx, root, err := tracker.StartSession(context.Background(), SessionContext{
		TenantHash:  "sha256-tenant",
		SessionHash: "sha256-session",
		Intent:      "site_audit",
		Mode:        "audit",
	})
	if err != nil {
		t.Fatalf("Tracker.StartSession() failed: %v", err)
	}
	_, child, err := tracker.StartOperation(ctx, "runtime.execution", OperationContext{
		SourceLayer: "layer_6",
		ActionClass: "read_only_fetch",
		Outcome:     "success",
	})
	if err != nil {
		t.Fatalf("Tracker.StartOperation() failed: %v", err)
	}
	if err := RecordSuccess(child); err != nil {
		t.Fatalf("RecordSuccess() failed: %v", err)
	}
	child.End()
	root.End()

	spans := recorder.Ended()
	if len(spans) != 2 {
		t.Fatalf("ended spans = %d, want 2", len(spans))
	}
	if spans[0].Name() != "runtime.execution" ||
		spans[1].Name() != "agent.session" {
		t.Fatalf("span names = %q, %q", spans[0].Name(), spans[1].Name())
	}
	if spans[0].Parent().SpanID() != spans[1].SpanContext().SpanID() {
		t.Fatal("child span is not linked to root")
	}
}

func TestTrackerRejectsUnknownSpanAndProtectedAttribute(t *testing.T) {
	provider := sdktrace.NewTracerProvider()
	defer func() {
		_ = provider.Shutdown(context.Background())
	}()
	tracker, err := NewTracker(provider, "1.0.0")
	if err != nil {
		t.Fatalf("NewTracker() failed: %v", err)
	}
	if _, _, err := tracker.StartOperation(
		context.Background(),
		"arbitrary.span",
		OperationContext{
			SourceLayer: "layer_6",
			ActionClass: "read_only_fetch",
			Outcome:     "success",
		},
	); err == nil {
		t.Fatal("Tracker.StartOperation() accepted unknown span")
	}
	if _, _, err := tracker.StartSession(context.Background(), SessionContext{
		TenantHash:  "tenant",
		SessionHash: "session",
		Intent:      "raw_prompt=secret",
		Mode:        "audit",
	}); err == nil {
		t.Fatal("Tracker.StartSession() accepted protected attribute")
	}
}
