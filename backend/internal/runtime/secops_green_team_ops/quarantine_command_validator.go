// Package secopsgreenteamops executes signed Layer 8 quarantine decisions.
package secopsgreenteamops

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"aeolyzer/internal/runtime"
)

var (
	// ErrInvalidQuarantineCommand indicates malformed, expired, or unsigned
	// quarantine input.
	ErrInvalidQuarantineCommand = errors.New("invalid quarantine command")
	// ErrQuarantineSignature indicates a signature mismatch.
	ErrQuarantineSignature = errors.New("quarantine signature verification failed")
)

var allowedQuarantineActions = map[string]struct{}{
	"freeze_runtime":             {},
	"revoke_jit_tokens":          {},
	"block_egress":               {},
	"revoke_tool_access":         {},
	"preserve_forensic_snapshot": {},
	"stop_new_executions":        {},
	"allow_read_only_status":     {},
	"start_repair_scope":         {},
}

// CommandValidator verifies structure, expiry, allowed actions, and HMAC.
type CommandValidator struct {
	key []byte
	now func() time.Time
}

// NewCommandValidator constructs a validator with a minimum 256-bit key.
func NewCommandValidator(
	key []byte,
	now func() time.Time,
) (*CommandValidator, error) {
	if len(key) < 32 || now == nil {
		return nil, errors.New("quarantine validator is not configured")
	}
	return &CommandValidator{
		key: append([]byte(nil), key...),
		now: now,
	}, nil
}

// SignCommand signs a complete command for tests and trusted Layer 8 adapters.
// Production signing keys must be held by an external secret manager.
func SignCommand(key []byte, command runtime.QuarantineCommand) (string, error) {
	if len(key) < 32 {
		return "", ErrInvalidQuarantineCommand
	}
	command.Signature = ""
	if err := validateCommandFields(command, time.Time{}); err != nil {
		return "", err
	}
	payload, err := canonicalCommand(command)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write(payload)
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

// Validate verifies a signed quarantine command.
func (v *CommandValidator) Validate(command runtime.QuarantineCommand) error {
	if v == nil || len(v.key) < 32 || v.now == nil {
		return errors.New("quarantine validator is not configured")
	}
	if command.Signature == "" {
		return ErrInvalidQuarantineCommand
	}
	if err := validateCommandFields(command, v.now()); err != nil {
		return err
	}
	provided, err := base64.RawURLEncoding.DecodeString(command.Signature)
	if err != nil {
		return ErrQuarantineSignature
	}
	unsigned := command
	unsigned.Signature = ""
	payload, err := canonicalCommand(unsigned)
	if err != nil {
		return err
	}
	mac := hmac.New(sha256.New, v.key)
	_, _ = mac.Write(payload)
	if !hmac.Equal(provided, mac.Sum(nil)) {
		return ErrQuarantineSignature
	}
	return nil
}

// ValidateQuarantineCommand performs structural validation for compatibility.
// Runtime execution must use CommandValidator.Validate for cryptographic proof.
func ValidateQuarantineCommand(command runtime.QuarantineCommand) error {
	return validateCommandFields(command, time.Time{})
}

func validateCommandFields(
	command runtime.QuarantineCommand,
	now time.Time,
) error {
	if command.TraceID == "" ||
		command.SessionID == "" ||
		command.TargetScope == "" ||
		command.TriggerReason == "" ||
		command.DecisionRef == "" ||
		command.ExpiresAt == "" ||
		len(command.RequestedActions) == 0 ||
		len(command.RequestedActions) > 8 {
		return ErrInvalidQuarantineCommand
	}
	switch command.Severity {
	case "medium", "high", "critical":
	default:
		return ErrInvalidQuarantineCommand
	}
	seen := make(map[string]struct{}, len(command.RequestedActions))
	for _, action := range command.RequestedActions {
		if _, allowed := allowedQuarantineActions[action]; !allowed {
			return ErrInvalidQuarantineCommand
		}
		if _, duplicate := seen[action]; duplicate {
			return ErrInvalidQuarantineCommand
		}
		seen[action] = struct{}{}
	}
	expiresAt, err := time.Parse(time.RFC3339, command.ExpiresAt)
	if err != nil {
		return ErrInvalidQuarantineCommand
	}
	if !now.IsZero() && !now.Before(expiresAt) {
		return ErrInvalidQuarantineCommand
	}
	return nil
}

func canonicalCommand(command runtime.QuarantineCommand) ([]byte, error) {
	data, err := json.Marshal(command)
	if err != nil {
		return nil, fmt.Errorf("encode quarantine command: %w", err)
	}
	return data, nil
}
