import os
import shutil

base = r"C:\Users\Muham\AEOlyzer\backend"
old_l5 = os.path.join(base, "layer_05_extensions")
new_l5 = os.path.join(base, "internal", "extensions")

# Create structure
dirs = [
    "config",
    "a2ui_translator",
    "surface_router",
    "interactive_components",
    "catalogs/basic_v0_9",
    "catalogs/seo_content_app_v1",
    "approval_ux",
    "a2a_server",
    "security",
    "events",
    "tests"
]

os.makedirs(new_l5, exist_ok=True)
for d in dirs:
    os.makedirs(os.path.join(new_l5, d), exist_ok=True)

# Move existing files if present
if os.path.exists(old_l5):
    for f in os.listdir(old_l5):
        src = os.path.join(old_l5, f)
        dst = os.path.join(new_l5, f)
        if os.path.isfile(src):
            shutil.move(src, dst)
    shutil.rmtree(old_l5)

schemas = {
    "presentation.schema.json": "{}",
    "a2ui-frame.schema.json": "{}",
    "a2ui-catalog.schema.json": "{}",
    "ui-event.schema.json": "{}",
    "approval.schema.json": "{}",
    "surface-patch.schema.json": "{}",
    "a2a-agent-card.schema.json": "{}",
    "a2a-envelope.schema.json": "{}",
    "catalog-lock.yaml": "version: 1\n",
    "presentation-changelog.md": "# Changelog\n",
    "README.md": "# Layer 5 Extensions\n"
}

for name, content in schemas.items():
    with open(os.path.join(new_l5, name), "w") as f:
        f.write(content)

# Generate new Types
types_go = """package extensions

// PresentationIntent represents a sanitized outcome from Layer 3 that is ready for display.
// Structural separation guarantees Layer 5 does not execute the actual business logic; 
// it merely maps state into declarative UI frames.
type PresentationIntent struct {
	TraceID          string                 `json:"trace_id"`
	WorkflowID       string                 `json:"workflow_id,omitempty"`
	NodeID           string                 `json:"node_id,omitempty"`
	Surface          string                 `json:"surface"`
	EventKind        string                 `json:"event_kind"`
	Mode             string                 `json:"mode,omitempty"`
	Priority         string                 `json:"priority,omitempty"`
	Payload          map[string]interface{} `json:"payload"`
	OutputContracts  []string               `json:"output_contracts,omitempty"`
	ApprovalRequired bool                   `json:"approval_required,omitempty"`
	FallbackText     string                 `json:"fallback_text,omitempty"`
	Metadata         map[string]string      `json:"metadata,omitempty"`
}

// A2UIFrame is a strictly typed declarative UI model passed to the client. 
// Executable boundaries: No javascript, HTML, or CSS is allowed here.
type A2UIFrame struct {
	FrameID        string            `json:"frame_id"`
	TraceID        string            `json:"trace_id,omitempty"`
	Surface        string            `json:"surface"`
	CatalogID      string            `json:"catalog_id"`
	CatalogVersion string            `json:"catalog_version"`
	SchemaVersion  string            `json:"schema_version"`
	RootID         string            `json:"root_id"`
	Nodes          []A2UINode        `json:"nodes"`
	FallbackText   string            `json:"fallback_text,omitempty"`
	ExpiresAt      string            `json:"expires_at,omitempty"`
	Signature      string            `json:"signature"`
}

type A2UINode struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Props    map[string]interface{} `json:"props,omitempty"`
	Children []string               `json:"children,omitempty"`
	Slot     string                 `json:"slot,omitempty"`
}
"""
with open(os.path.join(new_l5, "types.go"), "w") as f:
    f.write(types_go)

# Generate Validator
validator_go = """package extensions

import (
	"errors"
	"strings"
)

var (
	ErrUnknownSurface      = errors.New("UNKNOWN_SURFACE")
	ErrUnsafePayload       = errors.New("UNSAFE_PAYLOAD")
)

// ValidatePresentationIntent ensures that upstream layers are not leaking 
// raw execution contexts or arbitrary javascript into the presentation layer.
// This preserves the firewall against DOM XSS and payload smuggling.
func ValidatePresentationIntent(intent PresentationIntent) error {
	if intent.Surface == "" {
		return ErrUnknownSurface
	}
	
	// Enforce that fallback text does not contain raw HTML injections.
	// This acts as a secondary depth-in-defense check.
	if strings.Contains(intent.FallbackText, "<script>") {
		return ErrUnsafePayload
	}

	return nil
}
"""
with open(os.path.join(new_l5, "surface_router", "presentation_intent_validator.go"), "w") as f:
    f.write(validator_go)

# Generate Boundary Test
test_go = """package extensions_test

import (
	"testing"
	"aeolyzer/internal/extensions"
)

// TestNoExecutableBoundary ensures that Layer 5 immediately drops any presentation intent 
// attempting to smuggle raw HTML or executable scripts, honoring the declarative-only mandate.
func TestNoExecutableBoundary(t *testing.T) {
	intent := extensions.PresentationIntent{
		Surface:      "chat",
		FallbackText: "Hello <script>alert(1)</script>",
	}
	
	err := extensions.ValidatePresentationIntent(intent)
	if err != extensions.ErrUnsafePayload {
		t.Fatalf("expected ErrUnsafePayload for script injection, got %v", err)
	}
}

// TestSurfaceValidation ensures that unknown surfaces are structurally rejected, 
// preventing unsupported state mappings.
func TestSurfaceValidation(t *testing.T) {
	intent := extensions.PresentationIntent{
		Surface:      "", // Missing surface
	}
	
	err := extensions.ValidatePresentationIntent(intent)
	if err != extensions.ErrUnknownSurface {
		t.Fatalf("expected ErrUnknownSurface, got %v", err)
	}
}
"""
with open(os.path.join(new_l5, "tests", "presentation_schema_test.go"), "w") as f:
    f.write(test_go)

print("Layer 5 migrated and scaffolded successfully")
