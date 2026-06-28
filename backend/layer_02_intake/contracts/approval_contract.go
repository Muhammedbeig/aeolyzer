package contracts

import "time"

type ApprovalExpectation struct {
	ApprovalFor string
	Required    bool
}

type ApprovedAction struct {
	ApprovalID       string            `json:"approval_id"`
	ApprovalFor      string            `json:"approval_for"`
	TraceID          string            `json:"trace_id"`
	Source           string            `json:"source"`
	Surface          string            `json:"surface,omitempty"`
	SelectedTextHash string            `json:"selected_text_hash,omitempty"`
	ExpiresAt        time.Time         `json:"expires_at"`
	Constraints      map[string]string `json:"constraints,omitempty"`
}
