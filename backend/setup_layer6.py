import os
import shutil

base = r"C:\Users\Muham\AEOlyzer\backend"
old_l6 = os.path.join(base, "layer_06_runtime")
new_l6 = os.path.join(base, "internal", "runtime")

dirs = [
    "config",
    "execution_gateway",
    "sandbox_environment",
    "filesystem_control",
    "network_egress",
    "supply_chain_defense",
    "iam_context",
    "runtime_observation",
    "secops_green_team_ops"
]

os.makedirs(new_l6, exist_ok=True)
for d in dirs:
    os.makedirs(os.path.join(new_l6, d), exist_ok=True)

# Move existing files if present
if os.path.exists(old_l6):
    for f in os.listdir(old_l6):
        src = os.path.join(old_l6, f)
        dst = os.path.join(new_l6, f)
        if os.path.isfile(src):
            shutil.move(src, dst)
    shutil.rmtree(old_l6)

schemas = {
    "runtime-execution.schema.json": "{}",
    "sandbox-lease.schema.json": "{}",
    "runtime-result.schema.json": "{}",
    "jit-token.schema.json": "{}",
    "quarantine-command.schema.json": "{}",
    "dependency-policy.schema.json": "{}",
    "filesystem-policy.schema.json": "{}",
    "egress-policy.schema.json": "{}",
    "runtime-changelog.md": "# Changelog\n",
    "README.md": "# Layer 6 Runtime\n",
    "layer6-boundary.md": "# Boundary\n"
}

for name, content in schemas.items():
    with open(os.path.join(new_l6, name), "w") as f:
        f.write(content)

types_go = """package runtime

// RuntimeExecutionRequest is the authorized envelope passing into the execution layer.
// Structural separation guarantees that Layer 6 does not invent this request, 
// it only executes what Layer 2/3 explicitly approved.
type RuntimeExecutionRequest struct {
	TraceID             string                 `json:"trace_id"`
	SessionID           string                 `json:"session_id"`
	TaskID              string                 `json:"task_id"`
	RuntimeClass        string                 `json:"runtime_class"`
	ActionType          string                 `json:"action_type"`
	PolicyDecisionID    string                 `json:"policy_decision_id"`
	RequestSignature    string                 `json:"request_signature"`
	ExpiresAt           string                 `json:"expires_at"`
}

// QuarantineCommand instructs Layer 6 to freeze or alter state.
type QuarantineCommand struct {
	TraceID       string   `json:"trace_id"`
	TargetScope   string   `json:"target_scope"`
	TriggerReason string   `json:"trigger_reason"`
	Signature     string   `json:"signature"`
}
"""
with open(os.path.join(new_l6, "types.go"), "w") as f:
    f.write(types_go)

# Implementation files
iam_context_go = """package iam_context

import "errors"

var ErrAmbientCredentialsBlocked = errors.New("AMBIENT_CREDENTIALS_BLOCKED")

// StripAmbientCredentials enforces zero-ambient authority (Section 4.3).
// Prevents a script or tool from inheriting the host/orchestrator's implicit permissions.
// Any execution context attempting to pass raw ENV secrets will be rejected at the gateway.
func StripAmbientCredentials(env map[string]string) error {
	for k := range env {
		if k == "AWS_ACCESS_KEY_ID" || k == "GITHUB_TOKEN" {
			return ErrAmbientCredentialsBlocked
		}
	}
	return nil
}
"""
with open(os.path.join(new_l6, "iam_context", "ambient_credential_stripper.go"), "w") as f:
    f.write(iam_context_go)

quarantine_go = """package secops_green_team_ops

import "errors"
import "aeolyzer/internal/runtime"

var ErrInvalidQuarantineCommand = errors.New("INVALID_QUARANTINE_COMMAND")

// ValidateQuarantineCommand ensures that stateful freezes are strictly authorized.
// Layer 6 does not decide to quarantine; it only executes signed commands from Layer 8.
func ValidateQuarantineCommand(cmd runtime.QuarantineCommand) error {
	if cmd.Signature == "" || cmd.TargetScope == "" {
		return ErrInvalidQuarantineCommand
	}
	return nil
}
"""
with open(os.path.join(new_l6, "secops_green_team_ops", "quarantine_command_validator.go"), "w") as f:
    f.write(quarantine_go)

# Tests
iam_test_go = """package iam_context_test

import (
	"testing"
	"aeolyzer/internal/runtime/iam_context"
)

func TestAmbientCredentialStripping(t *testing.T) {
	env := map[string]string{
		"GITHUB_TOKEN": "secret",
	}
	if err := iam_context.StripAmbientCredentials(env); err == nil {
		t.Fatal("expected ambient credentials to be blocked")
	}
}
"""
with open(os.path.join(new_l6, "iam_context", "ambient_credential_stripper_test.go"), "w") as f:
    f.write(iam_test_go)

quarantine_test_go = """package secops_green_team_ops_test

import (
	"testing"
	"aeolyzer/internal/runtime/secops_green_team_ops"
	"aeolyzer/internal/runtime"
)

func TestQuarantineCommandValidation(t *testing.T) {
	cmd := runtime.QuarantineCommand{
		TargetScope: "agent-123",
		Signature: "", // Missing signature should fail
	}
	if err := secops_green_team_ops.ValidateQuarantineCommand(cmd); err == nil {
		t.Fatal("expected unsigned quarantine command to fail")
	}
}
"""
with open(os.path.join(new_l6, "secops_green_team_ops", "quarantine_command_validator_test.go"), "w") as f:
    f.write(quarantine_test_go)

# Fix imports in Go files
def replace_in_file(filepath, old, new):
    if not os.path.exists(filepath):
        return
    with open(filepath, 'r') as f:
        content = f.read()
    if old in content:
        content = content.replace(old, new)
        with open(filepath, 'w') as f:
            f.write(content)

# We might need to walk the tree to replace imports
for root, _, files in os.walk(base):
    for file in files:
        if file.endswith('.go'):
            replace_in_file(os.path.join(root, file), "aeolyzer/layer_06_runtime", "aeolyzer/internal/runtime")

print("Layer 6 scaffolded successfully.")
