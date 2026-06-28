package runtime

// RuntimeExecutionRequest is the authorized envelope passing into the execution layer.
// Structural separation guarantees that Layer 6 does not invent this request,
// it only executes what Layer 2/3 explicitly approved.
type RuntimeExecutionRequest struct {
	TraceID          string `json:"trace_id"`
	SessionID        string `json:"session_id"`
	TaskID           string `json:"task_id"`
	RuntimeClass     string `json:"runtime_class"`
	ActionType       string `json:"action_type"`
	PolicyDecisionID string `json:"policy_decision_id"`
	RequestSignature string `json:"request_signature"`
	ExpiresAt        string `json:"expires_at"`
}

// QuarantineCommand instructs Layer 6 to freeze or alter state.
type QuarantineCommand struct {
	TraceID          string            `json:"trace_id"`
	SessionID        string            `json:"session_id"`
	TargetScope      string            `json:"target_scope"`
	TriggerReason    string            `json:"trigger_reason"`
	Severity         string            `json:"severity"`
	RequestedActions []string          `json:"requested_actions"`
	PreserveState    bool              `json:"preserve_state"`
	DecisionRef      string            `json:"decision_ref"`
	ExpiresAt        string            `json:"expires_at"`
	Signature        string            `json:"signature"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}
