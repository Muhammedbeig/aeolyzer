package skills

import "time"

type SkillActivationRequest struct {
	TraceID             string            `json:"trace_id"`
	WorkflowID          string            `json:"workflow_id"`
	NodeID              string            `json:"node_id"`
	Intent              string            `json:"intent"`
	Mode                string            `json:"mode"`
	ProfileID           string            `json:"profile_id"`
	RequestedSkillIDs   []string          `json:"requested_skill_ids"`
	RequiredTags        []string          `json:"required_tags,omitempty"`
	OutputContracts     []string          `json:"output_contracts,omitempty"`
	MaxTokenBudget      int               `json:"max_token_budget"`
	SanitizedContextRef string            `json:"sanitized_context_ref,omitempty"`
	ResourceHints       []string          `json:"resource_hints,omitempty"`
	EvalMode            bool              `json:"eval_mode,omitempty"`
	Metadata            map[string]string `json:"metadata,omitempty"`
}

type SkillActivationResponse struct {
	TraceID             string        `json:"trace_id"`
	LoadedSkills        []SkillBundle `json:"loaded_skills"`
	OmittedSkills       []string      `json:"omitted_skills,omitempty"`
	TokenEstimate       int           `json:"token_estimate"`
	ResourceHandles     []string      `json:"resource_handles,omitempty"`
	CompatibilityStatus string        `json:"compatibility_status"`
	Warnings            []string      `json:"warnings,omitempty"`
	SafeSummary         string        `json:"safe_summary,omitempty"`
}

type SkillBundle struct {
	SkillID         string            `json:"skill_id"`
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Tier            string            `json:"tier"`
	BodyMarkdown    string            `json:"body_markdown"`
	LoadedResources []string          `json:"loaded_resources,omitempty"`
	TokenEstimate   int               `json:"token_estimate"`
	Checksums       map[string]string `json:"checksums"`
	OutputContracts []string          `json:"output_contracts"`
}

type SkillEvent struct {
	TraceID       string            `json:"trace_id,omitempty"`
	EventType     string            `json:"event_type"`
	SkillID       string            `json:"skill_id,omitempty"`
	SkillVersion  string            `json:"skill_version,omitempty"`
	Decision      string            `json:"decision"`
	ReasonCode    string            `json:"reason_code,omitempty"`
	TokenEstimate int               `json:"token_estimate,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	CreatedAt     time.Time         `json:"created_at"`
}
