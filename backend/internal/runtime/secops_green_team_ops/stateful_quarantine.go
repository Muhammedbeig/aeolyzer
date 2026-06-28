package secopsgreenteamops

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"aeolyzer/internal/runtime"
)

// RuntimeFreezer freezes a running sandbox or session.
type RuntimeFreezer interface {
	Freeze(context.Context, string) error
}

// TokenRevoker revokes all JIT credentials for a trace.
type TokenRevoker interface {
	RevokeTrace(context.Context, string) error
}

// EgressBlocker blocks all outbound traffic for a target scope.
type EgressBlocker interface {
	Block(context.Context, string) error
}

// ToolAccessRevoker denies further tool execution for a target scope.
type ToolAccessRevoker interface {
	RevokeToolAccess(context.Context, string) error
}

// ForensicPreserver creates a restricted forensic snapshot.
type ForensicPreserver interface {
	Preserve(context.Context, string, string) error
}

// QuarantineState is safe status metadata without forensic contents.
type QuarantineState struct {
	SessionID        string   `json:"session_id"`
	DecisionRef      string   `json:"decision_ref"`
	Status           string   `json:"status"`
	CompletedActions []string `json:"completed_actions"`
}

// QuarantineExecutor executes only signed, allowlisted quarantine actions.
type QuarantineExecutor struct {
	mu        sync.RWMutex
	validator *CommandValidator
	freezer   RuntimeFreezer
	tokens    TokenRevoker
	egress    EgressBlocker
	tools     ToolAccessRevoker
	forensics ForensicPreserver
	states    map[string]QuarantineState
	decisions map[string]struct{}
}

// NewQuarantineExecutor constructs a stateful executor. Dependencies are
// required because silently skipping a requested action is unsafe.
func NewQuarantineExecutor(
	validator *CommandValidator,
	freezer RuntimeFreezer,
	tokens TokenRevoker,
	egress EgressBlocker,
	tools ToolAccessRevoker,
	forensics ForensicPreserver,
) (*QuarantineExecutor, error) {
	if validator == nil ||
		freezer == nil ||
		tokens == nil ||
		egress == nil ||
		tools == nil ||
		forensics == nil {
		return nil, errors.New("quarantine executor dependencies are required")
	}
	return &QuarantineExecutor{
		validator: validator,
		freezer:   freezer,
		tokens:    tokens,
		egress:    egress,
		tools:     tools,
		forensics: forensics,
		states:    make(map[string]QuarantineState),
		decisions: make(map[string]struct{}),
	}, nil
}

// Execute validates and applies a quarantine command idempotently.
func (e *QuarantineExecutor) Execute(
	ctx context.Context,
	command runtime.QuarantineCommand,
) (QuarantineState, error) {
	if e == nil || e.validator == nil {
		return QuarantineState{}, errors.New("quarantine executor is not configured")
	}
	if err := e.validator.Validate(command); err != nil {
		return QuarantineState{}, err
	}

	e.mu.Lock()
	if _, duplicate := e.decisions[command.DecisionRef]; duplicate {
		state := e.states[command.SessionID]
		e.mu.Unlock()
		return cloneState(state), nil
	}
	e.decisions[command.DecisionRef] = struct{}{}
	state := QuarantineState{
		SessionID:   command.SessionID,
		DecisionRef: command.DecisionRef,
		Status:      "quarantining",
	}
	e.states[command.SessionID] = state
	e.mu.Unlock()

	for _, action := range command.RequestedActions {
		if err := e.executeAction(ctx, command, action); err != nil {
			e.mu.Lock()
			state.Status = "quarantine_failed"
			e.states[command.SessionID] = cloneState(state)
			e.mu.Unlock()
			return cloneState(state), fmt.Errorf("execute quarantine action %s: %w", action, err)
		}
		state.CompletedActions = append(state.CompletedActions, action)
		e.mu.Lock()
		e.states[command.SessionID] = cloneState(state)
		e.mu.Unlock()
	}
	state.Status = "quarantined"
	e.mu.Lock()
	e.states[command.SessionID] = cloneState(state)
	e.mu.Unlock()
	return cloneState(state), nil
}

// AllowExecution reports whether new execution remains permitted.
func (e *QuarantineExecutor) AllowExecution(sessionID string) bool {
	if e == nil || sessionID == "" {
		return false
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	state, found := e.states[sessionID]
	return !found || (state.Status != "quarantining" &&
		state.Status != "quarantined" &&
		state.Status != "quarantine_failed")
}

// Status returns safe quarantine metadata.
func (e *QuarantineExecutor) Status(sessionID string) (QuarantineState, bool) {
	if e == nil {
		return QuarantineState{}, false
	}
	e.mu.RLock()
	defer e.mu.RUnlock()
	state, found := e.states[sessionID]
	return cloneState(state), found
}

func (e *QuarantineExecutor) executeAction(
	ctx context.Context,
	command runtime.QuarantineCommand,
	action string,
) error {
	switch action {
	case "freeze_runtime", "stop_new_executions":
		return e.freezer.Freeze(ctx, command.TargetScope)
	case "revoke_jit_tokens":
		return e.tokens.RevokeTrace(ctx, command.TraceID)
	case "block_egress":
		return e.egress.Block(ctx, command.TargetScope)
	case "revoke_tool_access":
		return e.tools.RevokeToolAccess(ctx, command.TargetScope)
	case "preserve_forensic_snapshot":
		if !command.PreserveState {
			return errors.New("forensic preservation was not authorized")
		}
		return e.forensics.Preserve(ctx, command.TargetScope, command.DecisionRef)
	case "allow_read_only_status", "start_repair_scope":
		return nil
	default:
		return ErrInvalidQuarantineCommand
	}
}

func cloneState(state QuarantineState) QuarantineState {
	result := state
	result.CompletedActions = append([]string(nil), state.CompletedActions...)
	return result
}
