// Package datasecuritymesh enforces tenant, identity, projection, provenance,
// and taint boundaries for Layer 7 data access.
package datasecuritymesh

import (
	"crypto/subtle"
	"errors"
)

var ErrCrossTenantLeak = errors.New("cross-tenant access is blocked")

// EnforceTenantBoundary requires non-empty, exact tenant binding.
func EnforceTenantBoundary(requestTenantID, credentialTenantID string) error {
	if requestTenantID == "" ||
		credentialTenantID == "" ||
		len(requestTenantID) != len(credentialTenantID) ||
		subtle.ConstantTimeCompare(
			[]byte(requestTenantID),
			[]byte(credentialTenantID),
		) != 1 {
		return ErrCrossTenantLeak
	}
	return nil
}
