package extensions

// PresentationIntent represents a sanitized outcome from Layer 3 that is ready for display.
// Structural separation guarantees Layer 5 does not execute the actual business logic;
// it merely maps state into declarative UI frames.
type PresentationIntent struct {
	TraceID          string                 `json:"trace_id"`
	WorkflowID       string                 `json:"workflow_id,omitempty"`
	NodeID           string                 `json:"node_id,omitempty"`
	Surface          string                 `json:"surface"`
	EventKind        string                 `json:"event_kind"`
	Mode             string                 `json:"mode,omitempty"`
	Priority         string                 `json:"priority,omitempty"`
	Payload          map[string]interface{} `json:"payload"`
	OutputContracts  []string               `json:"output_contracts,omitempty"`
	ApprovalRequired bool                   `json:"approval_required,omitempty"`
	FallbackText     string                 `json:"fallback_text,omitempty"`
	Metadata         map[string]string      `json:"metadata,omitempty"`
}

// A2UIFrame is a strictly typed declarative UI model passed to the client.
// Executable boundaries: No javascript, HTML, or CSS is allowed here.
type A2UIFrame struct {
	FrameID        string     `json:"frame_id"`
	TraceID        string     `json:"trace_id,omitempty"`
	Surface        string     `json:"surface"`
	CatalogID      string     `json:"catalog_id"`
	CatalogVersion string     `json:"catalog_version"`
	SchemaVersion  string     `json:"schema_version"`
	RootID         string     `json:"root_id"`
	Nodes          []A2UINode `json:"nodes"`
	FallbackText   string     `json:"fallback_text,omitempty"`
	ExpiresAt      string     `json:"expires_at,omitempty"`
	Signature      string     `json:"signature"`
}

type A2UINode struct {
	ID       string                 `json:"id"`
	Type     string                 `json:"type"`
	Props    map[string]interface{} `json:"props,omitempty"`
	Children []string               `json:"children,omitempty"`
	Slot     string                 `json:"slot,omitempty"`
}
