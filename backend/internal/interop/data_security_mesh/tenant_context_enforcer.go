package data_security_mesh

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
