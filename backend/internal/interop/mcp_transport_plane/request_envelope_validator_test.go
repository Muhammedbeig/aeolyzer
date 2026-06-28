package mcptransportplane

import (
	"errors"
	"testing"

	"aeolyzer/internal/interop"
)

func TestValidateEnvelope(t *testing.T) {
	request := interop.InteropRequest{
		RequestID: "request-1",
		TenantID:  "tenant-1",
	}
	if err := ValidateEnvelope(request); !errors.Is(err, ErrMissingContext) {
		t.Fatalf("ValidateEnvelope() error = %v, want %v", err, ErrMissingContext)
	}
}
