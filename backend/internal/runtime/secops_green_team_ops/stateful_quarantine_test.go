package secopsgreenteamops

import (
	"context"
	"testing"
	"time"
)

type quarantineDependencies struct {
	frozen    int
	revoked   int
	blocked   int
	tools     int
	preserved int
}

func (d *quarantineDependencies) Freeze(context.Context, string) error {
	d.frozen++
	return nil
}

func (d *quarantineDependencies) RevokeTrace(context.Context, string) error {
	d.revoked++
	return nil
}

func (d *quarantineDependencies) Block(context.Context, string) error {
	d.blocked++
	return nil
}

func (d *quarantineDependencies) RevokeToolAccess(context.Context, string) error {
	d.tools++
	return nil
}

func (d *quarantineDependencies) Preserve(context.Context, string, string) error {
	d.preserved++
	return nil
}

func TestQuarantineExecutorAppliesSignedActionsAndBlocksExecution(t *testing.T) {
	now := time.Date(2026, 6, 28, 12, 0, 0, 0, time.UTC)
	key := []byte("01234567890123456789012345678901")
	validator, err := NewCommandValidator(key, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewCommandValidator() failed: %v", err)
	}
	dependencies := &quarantineDependencies{}
	executor, err := NewQuarantineExecutor(
		validator,
		dependencies,
		dependencies,
		dependencies,
		dependencies,
		dependencies,
	)
	if err != nil {
		t.Fatalf("NewQuarantineExecutor() failed: %v", err)
	}
	command := validCommand(now)
	command.RequestedActions = []string{
		"stop_new_executions",
		"revoke_jit_tokens",
		"block_egress",
		"revoke_tool_access",
		"preserve_forensic_snapshot",
	}
	signature, err := SignCommand(key, command)
	if err != nil {
		t.Fatalf("SignCommand() failed: %v", err)
	}
	command.Signature = signature

	state, err := executor.Execute(context.Background(), command)
	if err != nil {
		t.Fatalf("QuarantineExecutor.Execute() failed: %v", err)
	}
	if state.Status != "quarantined" || executor.AllowExecution(command.SessionID) {
		t.Fatalf("quarantine state = %+v, allow = %t", state, executor.AllowExecution(command.SessionID))
	}
	if dependencies.frozen != 1 ||
		dependencies.revoked != 1 ||
		dependencies.blocked != 1 ||
		dependencies.tools != 1 ||
		dependencies.preserved != 1 {
		t.Fatalf("quarantine dependencies = %+v, want each action once", dependencies)
	}

	if _, err := executor.Execute(context.Background(), command); err != nil {
		t.Fatalf("idempotent QuarantineExecutor.Execute() failed: %v", err)
	}
	if dependencies.frozen != 1 {
		t.Fatal("idempotent execution repeated side effects")
	}
}
