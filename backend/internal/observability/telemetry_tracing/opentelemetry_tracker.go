package telemetrytracing

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const instrumentationName = "aeolyzer/internal/observability"

var allowedSpanNames = map[string]struct{}{
	"agent.session":       {},
	"intake.decision":     {},
	"orchestration.plan":  {},
	"skill.activation":    {},
	"tool.proposal":       {},
	"tool.authorization":  {},
	"runtime.execution":   {},
	"connector.call":      {},
	"presentation.intent": {},
	"output.final":        {},
	"evaluation.run":      {},
	"governance.decision": {},
}

// SessionContext contains sanitized trace attributes. Hash fields must already
// be non-reversible.
type SessionContext struct {
	TenantHash  string
	SessionHash string
	Intent      string
	Mode        string
}

// OperationContext contains safe categorical child-span attributes.
type OperationContext struct {
	SourceLayer string
	ActionClass string
	Outcome     string
}

// Tracker creates policy-constrained OpenTelemetry spans.
type Tracker struct {
	tracer trace.Tracer
}

// NewTracker creates instrumentation from an injected provider. Applications
// must configure an SDK/exporter; this package never installs a global provider.
func NewTracker(provider trace.TracerProvider, version string) (*Tracker, error) {
	if provider == nil || version == "" {
		return nil, errors.New("opentelemetry provider and version are required")
	}
	return &Tracker{
		tracer: provider.Tracer(
			instrumentationName,
			trace.WithInstrumentationVersion(version),
		),
	}, nil
}

// StartSession starts the required root span.
func (t *Tracker) StartSession(
	ctx context.Context,
	session SessionContext,
) (context.Context, trace.Span, error) {
	if t == nil || t.tracer == nil {
		return ctx, nil, errors.New("opentelemetry tracker is not configured")
	}
	if !safeAttribute(session.TenantHash) ||
		!safeAttribute(session.SessionHash) ||
		!safeAttribute(session.Intent) ||
		!safeAttribute(session.Mode) {
		return ctx, nil, errors.New("session trace context is invalid")
	}
	ctx, span := t.tracer.Start(
		ctx,
		"agent.session",
		trace.WithSpanKind(trace.SpanKindServer),
		trace.WithAttributes(
			attribute.String("tenant.hash", session.TenantHash),
			attribute.String("session.hash", session.SessionHash),
			attribute.String("agent.intent", session.Intent),
			attribute.String("agent.mode", session.Mode),
		),
	)
	return ctx, span, nil
}

// StartOperation starts an allowlisted child span.
func (t *Tracker) StartOperation(
	ctx context.Context,
	spanName string,
	operation OperationContext,
) (context.Context, trace.Span, error) {
	if t == nil || t.tracer == nil {
		return ctx, nil, errors.New("opentelemetry tracker is not configured")
	}
	if _, allowed := allowedSpanNames[spanName]; !allowed || spanName == "agent.session" {
		return ctx, nil, fmt.Errorf("span name %q is not allowed", spanName)
	}
	if !safeAttribute(operation.SourceLayer) ||
		!safeAttribute(operation.ActionClass) ||
		!safeAttribute(operation.Outcome) {
		return ctx, nil, errors.New("operation trace context is invalid")
	}
	ctx, span := t.tracer.Start(
		ctx,
		spanName,
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(
			attribute.String("source.layer", operation.SourceLayer),
			attribute.String("action.class", operation.ActionClass),
			attribute.String("operation.outcome", operation.Outcome),
		),
	)
	return ctx, span, nil
}

// RecordError records a safe error class, never a raw error message.
func RecordError(span trace.Span, errorClass string) error {
	if span == nil {
		return errors.New("span is required")
	}
	if !safeAttribute(errorClass) {
		return errors.New("safe error class is required")
	}
	span.SetStatus(codes.Error, errorClass)
	span.SetAttributes(attribute.String("error.class", errorClass))
	return nil
}

// RecordSuccess marks a span successful.
func RecordSuccess(span trace.Span) error {
	if span == nil {
		return errors.New("span is required")
	}
	span.SetStatus(codes.Ok, "success")
	return nil
}

func safeAttribute(value string) bool {
	if value == "" || len(value) > 128 || strings.ContainsAny(value, "\r\n\t") {
		return false
	}
	lower := strings.ToLower(value)
	for _, fragment := range []string{
		"bearer ",
		"api_key",
		"password",
		"private key",
		"raw_prompt",
		"chain_of_thought",
	} {
		if strings.Contains(lower, fragment) {
			return false
		}
	}
	return true
}
