package datasecuritymesh

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

// TokenClaims is a verified token's security-relevant claim set.
type TokenClaims struct {
	TokenID     string
	TenantID    string
	ConnectorID string
	Audience    string
	Scopes      []string
	NotBefore   time.Time
	ExpiresAt   time.Time
}

// TokenExpectation is the exact authorized connector binding.
type TokenExpectation struct {
	TenantID    string
	ConnectorID string
	Audience    string
	Scopes      []string
	MaxTTL      time.Duration
}

// ValidateTokenAudience rejects tenant, connector, audience, scope, and time
// confusion after cryptographic token verification by the credential provider.
func ValidateTokenAudience(
	claims TokenClaims,
	expected TokenExpectation,
	now time.Time,
) error {
	if claims.TokenID == "" ||
		claims.TenantID == "" ||
		claims.ConnectorID == "" ||
		claims.Audience == "" ||
		expected.TenantID == "" ||
		expected.ConnectorID == "" ||
		expected.Audience == "" ||
		expected.MaxTTL <= 0 ||
		now.IsZero() {
		return errors.New("token validation context is incomplete")
	}
	if err := EnforceTenantBoundary(expected.TenantID, claims.TenantID); err != nil {
		return err
	}
	if claims.ConnectorID != expected.ConnectorID ||
		claims.Audience != expected.Audience {
		return errors.New("token connector or audience mismatch")
	}
	if claims.NotBefore.IsZero() ||
		claims.ExpiresAt.IsZero() ||
		claims.ExpiresAt.Sub(claims.NotBefore) > expected.MaxTTL ||
		now.Before(claims.NotBefore) ||
		!now.Before(claims.ExpiresAt) {
		return errors.New("token validity window is invalid")
	}
	actualScopes := append([]string(nil), claims.Scopes...)
	expectedScopes := append([]string(nil), expected.Scopes...)
	sort.Strings(actualScopes)
	sort.Strings(expectedScopes)
	if len(actualScopes) != len(expectedScopes) {
		return errors.New("token scope set mismatch")
	}
	for i := range actualScopes {
		if actualScopes[i] != expectedScopes[i] {
			return fmt.Errorf("token scope %q is not authorized", actualScopes[i])
		}
	}
	return nil
}
