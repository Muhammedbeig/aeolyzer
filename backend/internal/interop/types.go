package interop

// InteropRequest defines the strict envelope required to invoke any data connector (Section 7.1).
// Layer 7 requires explicit context (like TenantID and PolicyDecisionID) to enforce zero-ambient authority.
type InteropRequest struct {
	RequestID           string `json:"request_id"`
	SessionID           string `json:"session_id"`
	TenantID            string `json:"tenant_id"`
	ConnectorID         string `json:"connector_id"`
	ToolID              string `json:"tool_id"`
	PolicyDecisionID    string `json:"policy_decision_id"`
	JITCredentialRef    string `json:"jit_credential_ref,omitempty"`
}

// InteropResult ensures that raw connector output is properly tainted and attributed (Section 7.2).
type InteropResult struct {
	RequestID          string   `json:"request_id"`
	TenantID           string   `json:"tenant_id"`
	Status             string   `json:"status"`
	DataClassification string   `json:"data_classification"`
	TaintFlags         []string `json:"taint_flags,omitempty"`
}
