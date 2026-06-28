package secops_green_team_ops_test

import (
	"aeolyzer/internal/runtime"
	"aeolyzer/internal/runtime/secops_green_team_ops"
	"testing"
)

func TestQuarantineCommandValidation(t *testing.T) {
	cmd := runtime.QuarantineCommand{
		TargetScope: "agent-123",
		Signature:   "", // Missing signature should fail
	}
	if err := secops_green_team_ops.ValidateQuarantineCommand(cmd); err == nil {
		t.Fatal("expected unsigned quarantine command to fail")
	}
}
