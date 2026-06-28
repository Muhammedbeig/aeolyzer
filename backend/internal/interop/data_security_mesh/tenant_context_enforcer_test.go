package data_security_mesh_test

import (
	"aeolyzer/internal/interop/data_security_mesh"
	"testing"
)

func TestEnforceTenantBoundary(t *testing.T) {
	err := data_security_mesh.EnforceTenantBoundary("tenant-A", "tenant-B")
	if err == nil {
		t.Fatal("expected cross-tenant leak to be prevented")
	}
}
