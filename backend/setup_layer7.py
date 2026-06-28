import os
import shutil

base = r"C:\Users\Muham\AEOlyzer\backend"
old_l7 = os.path.join(base, "layer_07_interop")
new_l7 = os.path.join(base, "internal", "interop")

dirs = [
    "config",
    "mcp_transport_plane",
    "connector_registry",
    "enterprise_mcp_servers",
    "data_security_mesh",
    "vector_rag_store",
    "resource_packaging",
    "interop_events",
    "docs"
]

os.makedirs(new_l7, exist_ok=True)
for d in dirs:
    os.makedirs(os.path.join(new_l7, d), exist_ok=True)

# Move existing files if present
if os.path.exists(old_l7):
    for f in os.listdir(old_l7):
        src = os.path.join(old_l7, f)
        dst = os.path.join(new_l7, f)
        if os.path.isfile(src):
            shutil.move(src, dst)
    shutil.rmtree(old_l7)

schemas = {
    "config/connector-registry.yaml": "version: 2\n",
    "config/mcp-server-manifest.schema.json": "{}",
    "config/source-contracts.yaml": "version: 2\n",
    "docs/connector-onboarding.md": "# Connector Onboarding\n"
}
for name, content in schemas.items():
    with open(os.path.join(new_l7, name), "w") as f:
        f.write(content)

# 1. Types
types_go = """package interop

// InteropRequest defines the strict envelope required to invoke any data connector (Section 7.1).
// Layer 7 requires explicit context (like TenantID and PolicyDecisionID) to enforce zero-ambient authority.
type InteropRequest struct {
	RequestID           string `json:"request_id"`
	SessionID           string `json:"session_id"`
	TenantID            string `json:"tenant_id"`
	ConnectorID         string `json:"connector_id"`
	ToolID              string `json:"tool_id"`
	PolicyDecisionID    string `json:"policy_decision_id"`
	JITCredentialRef    string `json:"jit_credential_ref,omitempty"`
}

// InteropResult ensures that raw connector output is properly tainted and attributed (Section 7.2).
type InteropResult struct {
	RequestID          string   `json:"request_id"`
	TenantID           string   `json:"tenant_id"`
	Status             string   `json:"status"`
	DataClassification string   `json:"data_classification"`
	TaintFlags         []string `json:"taint_flags,omitempty"`
}
"""
with open(os.path.join(new_l7, "types.go"), "w") as f:
    f.write(types_go)

# 2. mcp_transport_plane/request_envelope_validator.go
validator_go = """package mcp_transport_plane

import (
	"errors"
	"aeolyzer/internal/interop"
)

var ErrMissingContext = errors.New("MISSING_REQUIRED_CONTEXT")

// ValidateEnvelope enforces the fails-closed policy for incomplete connector requests (Section 3.1).
// Ensures no arbitrary execution can occur without explicit orchestration and policy lineage.
func ValidateEnvelope(req interop.InteropRequest) error {
	if req.RequestID == "" || req.TenantID == "" || req.ConnectorID == "" || req.PolicyDecisionID == "" {
		return ErrMissingContext
	}
	return nil
}
"""
with open(os.path.join(new_l7, "mcp_transport_plane", "request_envelope_validator.go"), "w") as f:
    f.write(validator_go)

# 3. data_security_mesh/tenant_context_enforcer.go
tenant_go = """package data_security_mesh

import "errors"

var ErrCrossTenantLeak = errors.New("CROSS_TENANT_LEAK_PREVENTED")

// EnforceTenantBoundary prevents multi-tenant data bleed during vector or MCP retrieval (Section 3.3).
// If a JIT credential belongs to Tenant A, it strictly cannot execute a query for Tenant B.
func EnforceTenantBoundary(requestTenantID, credentialTenantID string) error {
	if requestTenantID != credentialTenantID {
		return ErrCrossTenantLeak
	}
	return nil
}
"""
with open(os.path.join(new_l7, "data_security_mesh", "tenant_context_enforcer.go"), "w") as f:
    f.write(tenant_go)

# 4. Tests
validator_test_go = """package mcp_transport_plane_test

import (
	"testing"
	"aeolyzer/internal/interop"
	"aeolyzer/internal/interop/mcp_transport_plane"
)

func TestValidateEnvelope(t *testing.T) {
	req := interop.InteropRequest{
		RequestID: "req-1",
		TenantID: "tenant-1",
		// Missing ConnectorID and PolicyDecisionID
	}
	if err := mcp_transport_plane.ValidateEnvelope(req); err == nil {
		t.Fatal("expected missing context error")
	}
}
"""
with open(os.path.join(new_l7, "mcp_transport_plane", "request_envelope_validator_test.go"), "w") as f:
    f.write(validator_test_go)

tenant_test_go = """package data_security_mesh_test

import (
	"testing"
	"aeolyzer/internal/interop/data_security_mesh"
)

func TestEnforceTenantBoundary(t *testing.T) {
	err := data_security_mesh.EnforceTenantBoundary("tenant-A", "tenant-B")
	if err == nil {
		t.Fatal("expected cross-tenant leak to be prevented")
	}
}
"""
with open(os.path.join(new_l7, "data_security_mesh", "tenant_context_enforcer_test.go"), "w") as f:
    f.write(tenant_test_go)

# Fix imports in Go files
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
            replace_in_file(os.path.join(root, file), "aeolyzer/layer_07_interop", "aeolyzer/internal/interop")

print("Layer 7 scaffolded successfully.")
