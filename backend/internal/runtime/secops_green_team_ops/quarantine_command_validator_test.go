package secopsgreenteamops

import (
	"errors"
	"testing"
	"time"

	"aeolyzer/internal/runtime"
)

func TestCommandValidatorAcceptsSignedCommand(t *testing.T) {
	now := time.Date(2026, 6, 28, 12, 0, 0, 0, time.UTC)
	key := []byte("01234567890123456789012345678901")
	command := validCommand(now)
	signature, err := SignCommand(key, command)
	if err != nil {
		t.Fatalf("SignCommand() failed: %v", err)
	}
	command.Signature = signature
	validator, err := NewCommandValidator(key, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewCommandValidator() failed: %v", err)
	}
	if err := validator.Validate(command); err != nil {
		t.Fatalf("CommandValidator.Validate() failed: %v", err)
	}
}

func TestCommandValidatorRejectsTamperingExpiryAndUnknownAction(t *testing.T) {
	now := time.Date(2026, 6, 28, 12, 0, 0, 0, time.UTC)
	key := []byte("01234567890123456789012345678901")
	validator, err := NewCommandValidator(key, func() time.Time { return now })
	if err != nil {
		t.Fatalf("NewCommandValidator() failed: %v", err)
	}
	command := validCommand(now)
	signature, err := SignCommand(key, command)
	if err != nil {
		t.Fatalf("SignCommand() failed: %v", err)
	}
	command.Signature = signature

	tampered := command
	tampered.TargetScope = "other-session"
	if err := validator.Validate(tampered); !errors.Is(err, ErrQuarantineSignature) {
		t.Fatalf("tampered Validate() error = %v, want signature failure", err)
	}
	expired := command
	expired.ExpiresAt = now.Format(time.RFC3339)
	if err := validator.Validate(expired); !errors.Is(err, ErrInvalidQuarantineCommand) {
		t.Fatalf("expired Validate() error = %v, want invalid command", err)
	}
	unknown := command
	unknown.RequestedActions = []string{"delete_everything"}
	if err := validator.Validate(unknown); !errors.Is(err, ErrInvalidQuarantineCommand) {
		t.Fatalf("unknown action Validate() error = %v, want invalid command", err)
	}
}

func validCommand(now time.Time) runtime.QuarantineCommand {
	return runtime.QuarantineCommand{
		TraceID:          "trace-1",
		SessionID:        "session-1",
		TargetScope:      "session-1",
		TriggerReason:    "critical trust threshold crossed",
		Severity:         "critical",
		RequestedActions: []string{"stop_new_executions", "revoke_jit_tokens", "block_egress"},
		PreserveState:    true,
		DecisionRef:      "decision-1",
		ExpiresAt:        now.Add(time.Minute).Format(time.RFC3339),
	}
}
