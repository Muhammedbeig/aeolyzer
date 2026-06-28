package mcp_transport_plane_test

import (
	"aeolyzer/internal/interop"
	"aeolyzer/internal/interop/mcp_transport_plane"
	"testing"
)

func TestValidateEnvelope(t *testing.T) {
	req := interop.InteropRequest{
		RequestID: "req-1",
		TenantID:  "tenant-1",
		// Missing ConnectorID and PolicyDecisionID
	}
	if err := mcp_transport_plane.ValidateEnvelope(req); err == nil {
		t.Fatal("expected missing context error")
	}
}
