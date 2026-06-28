package mcp_transport_plane

import (
	"aeolyzer/internal/interop"
	"errors"
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
