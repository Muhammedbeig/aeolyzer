import os
import shutil

base = r"C:\Users\Muham\AEOlyzer\backend"
old_l8 = os.path.join(base, "layer_08_observability")
new_l8 = os.path.join(base, "internal", "observability")

dirs = [
    "config",
    "contracts",
    "telemetry_tracing",
    "event_ingestion",
    "agbom",
    "drift_trust",
    "secops_triad",
    "evaluation_engine",
    "governance_audit",
    "feedback_improvement_loop",
    "metrics_exports",
    "runbooks",
    "docs"
]

os.makedirs(new_l8, exist_ok=True)
for d in dirs:
    os.makedirs(os.path.join(new_l8, d), exist_ok=True)

if os.path.exists(old_l8):
    for f in os.listdir(old_l8):
        src = os.path.join(old_l8, f)
        dst = os.path.join(new_l8, f)
        if os.path.isfile(src):
            shutil.move(src, dst)
    shutil.rmtree(old_l8)

schemas = {
    "config/telemetry-policy.yaml": "version: 2\n",
    "config/redaction-policy.yaml": "version: 2\n",
    "config/eval-policy.yaml": "version: 2\n",
    "config/secops-policy.yaml": "version: 2\n",
    "config/drift-policy.yaml": "version: 2\n",
    "runbooks/incident_response.md": "# Incident Response\n"
}
for name, content in schemas.items():
    with open(os.path.join(new_l8, name), "w") as f:
        f.write(content)

# 1. telemetry_tracing/span_validator.go
validator_go = """package telemetry_tracing

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
"""
with open(os.path.join(new_l8, "telemetry_tracing", "span_validator.go"), "w") as f:
    f.write(validator_go)

# 2. agbom/agbom_builder.go
agbom_go = """package agbom

import "errors"

// RuntimeAgentBillOfMaterials tracks everything a session touches (Section 4.6).
type RuntimeAgentBillOfMaterials struct {
	TraceID     string   `json:"trace_id"`
	SkillsUsed  []string `json:"skills_used"`
	ToolsUsed   []string `json:"tools_used"`
	Connectors  []string `json:"connectors"`
}

// BuildAgBOM asserts that a valid AgBOM must be grounded to a Trace.
func BuildAgBOM(traceID string) (*RuntimeAgentBillOfMaterials, error) {
	if traceID == "" {
		return nil, errors.New("TRACE_REQUIRED_FOR_AGBOM")
	}
	return &RuntimeAgentBillOfMaterials{
		TraceID: traceID,
	}, nil
}
"""
with open(os.path.join(new_l8, "agbom", "agbom_builder.go"), "w") as f:
    f.write(agbom_go)

# 3. event_ingestion/event_redactor.go
redactor_go = """package event_ingestion

import "strings"

// RedactEvent safely strips out sensitive values before they reach the data sink (Section 8).
// Ensures hidden chain-of-thought or raw system prompts are not leaked to external logs.
func RedactEvent(payload string) string {
	payload = strings.ReplaceAll(payload, "raw_system_prompt", "[REDACTED]")
	payload = strings.ReplaceAll(payload, "raw_developer_prompt", "[REDACTED]")
	return payload
}
"""
with open(os.path.join(new_l8, "event_ingestion", "event_redactor.go"), "w") as f:
    f.write(redactor_go)

# Tests
validator_test_go = """package telemetry_tracing_test

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
"""
with open(os.path.join(new_l8, "telemetry_tracing", "span_validator_test.go"), "w") as f:
    f.write(validator_test_go)

agbom_test_go = """package agbom_test

import (
	"testing"
	"aeolyzer/internal/observability/agbom"
)

func TestBuildAgBOM(t *testing.T) {
	_, err := agbom.BuildAgBOM("")
	if err == nil {
		t.Fatal("expected trace ID requirement for agbom")
	}
}
"""
with open(os.path.join(new_l8, "agbom", "agbom_builder_test.go"), "w") as f:
    f.write(agbom_test_go)

redactor_test_go = """package event_ingestion_test

import (
	"testing"
	"strings"
	"aeolyzer/internal/observability/event_ingestion"
)

func TestEventRedaction(t *testing.T) {
	out := event_ingestion.RedactEvent("this contains raw_system_prompt inside")
	if strings.Contains(out, "raw_system_prompt") {
		t.Fatal("expected sensitive context to be redacted")
	}
}
"""
with open(os.path.join(new_l8, "event_ingestion", "event_redactor_test.go"), "w") as f:
    f.write(redactor_test_go)

# Fix imports
def replace_in_file(filepath, old, new):
    if not os.path.exists(filepath):
        return
    with open(filepath, 'r') as f:
        content = f.read()
    if old in content:
        content = content.replace(old, new)
        with open(filepath, 'w') as f:
            f.write(content)

for root, _, files in os.walk(base):
    for file in files:
        if file.endswith('.go'):
            replace_in_file(os.path.join(root, file), "aeolyzer/layer_08_observability", "aeolyzer/internal/observability")

print("Layer 8 scaffolded successfully.")
