package datasecuritymesh

import (
	"errors"
	"testing"
)

func TestEnforceTenantBoundary(t *testing.T) {
	for name, tenants := range map[string][2]string{
		"different":        {"tenant-a", "tenant-b"},
		"request empty":    {"", "tenant-b"},
		"credential empty": {"tenant-a", ""},
	} {
		t.Run(name, func(t *testing.T) {
			if err := EnforceTenantBoundary(
				tenants[0],
				tenants[1],
			); !errors.Is(err, ErrCrossTenantLeak) {
				t.Fatalf("EnforceTenantBoundary() error = %v, want %v", err, ErrCrossTenantLeak)
			}
		})
	}
	if err := EnforceTenantBoundary("tenant-a", "tenant-a"); err != nil {
		t.Fatalf("EnforceTenantBoundary(equal) failed: %v", err)
	}
}
