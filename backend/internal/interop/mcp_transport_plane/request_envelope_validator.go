// Package mcptransportplane implements bounded MCP JSON-RPC transports.
package mcptransportplane

import (
	"errors"

	"aeolyzer/internal/interop"
)

var ErrMissingContext = errors.New("mcp request context is incomplete")

// ValidateEnvelope enforces required policy and tenant lineage.
func ValidateEnvelope(request interop.InteropRequest) error {
	if request.RequestID == "" ||
		request.TenantID == "" ||
		request.ConnectorID == "" ||
		request.PolicyDecisionID == "" {
		return ErrMissingContext
	}
	return nil
}
