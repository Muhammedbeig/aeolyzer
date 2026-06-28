package telemetry_tracing

import "errors"

// Span schema structure (internal telemetry logic).
type Span struct {
	TraceID  string `json:"trace_id"`
	SpanName string `json:"span_name"`
	TenantID string `json:"tenant_id,omitempty"`
}

var ErrMissingContext = errors.New("MISSING_SPAN_CONTEXT")

// ValidateSpan ensures that no trace is committed to storage without proper attribution (Section 7).
func ValidateSpan(span Span) error {
	if span.TraceID == "" || span.SpanName == "" {
		return ErrMissingContext
	}
	return nil
}
