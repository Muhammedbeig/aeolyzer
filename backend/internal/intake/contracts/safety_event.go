package contracts

import "time"

type SafetyEvent struct {
	TraceID     string                 `json:"trace_id"`
	EventType   string                 `json:"event_type"`
	Intent      Intent                 `json:"intent,omitempty"`
	Mode        OrchestrationMode      `json:"mode,omitempty"`
	ActionClass string                 `json:"action_class,omitempty"`
	Decision    string                 `json:"decision"`
	ReasonCode  string                 `json:"reason_code,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}
