package orchestrator

// OrchestrationMode enforces the isolation boundary between distinct cognitive phases.
// By strictly typing modes, we prevent state leakage where a planning agent inadvertently
// mutates execution state (e.g., writing article prose during topic discovery).
type OrchestrationMode string

const (
	ModePlan     OrchestrationMode = "plan"
	ModeWrite    OrchestrationMode = "write"
	ModeEdit     OrchestrationMode = "edit"
	ModeOptimize OrchestrationMode = "optimize"
	ModeAudit    OrchestrationMode = "audit"
)

type NodeType string

const (
	NodeTypeContentGeneration NodeType = "content_generation_task"
	NodeTypeApprovalRequest   NodeType = "approval_request"
)

type WorkflowID string

const (
	WorkflowTopicDiscovery      WorkflowID = "topic-discovery.bp"
	WorkflowContentBrief        WorkflowID = "content-brief.bp"
	WorkflowContentResearch     WorkflowID = "content-research.bp"
	WorkflowSEOPlanning         WorkflowID = "seo-planning.bp"
	WorkflowPageAnalysis        WorkflowID = "page-analysis.bp"
	WorkflowArticlePlanning     WorkflowID = "article-planning.bp"
	WorkflowArticleDrafting     WorkflowID = "article-drafting.bp"
	WorkflowContentOptimization WorkflowID = "content-optimization.bp"
	WorkflowContentRepurposing  WorkflowID = "content-repurposing.bp"
	WorkflowMemoryTone          WorkflowID = "memory-tone-management.bp"
)

type AgentID string

const (
	AgentContentCollaborator   AgentID = "content_collaborator"
	AgentContentExecutionGuard AgentID = "content_execution_guard"
)

// IntakeDecision defines the read-only contract handed off by Layer 2.
// Layer 3 relies entirely on these pre-validated fields, avoiding raw intent parsing
// to maintain the prompt-injection firewall established in Layer 2.
type IntakeDecision struct {
	TraceID          string                 `json:"trace_id"`
	Intent           string                 `json:"intent"`
	Confidence       float64                `json:"confidence"`
	SanitizedContext map[string]string      `json:"sanitized_context"`
	DisclosureStatus string                 `json:"disclosure_status,omitempty"`
	PolicyState      string                 `json:"policy_state,omitempty"`
	Mode             string                 `json:"mode,omitempty"`
	ApprovedActions  []string               `json:"approved_actions,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}

// TaskID securely references a specific node in the DAG plan.
type TaskID string

// ContentGenerationTask represents a structural delegation contract.
// Notice that this structure completely avoids containing raw LLM prompts or prompt templates.
// It is designed strictly for capability matching and constrained surface targeting.
type ContentGenerationTask struct {
	TraceID              string            `json:"trace_id"`
	TaskID               TaskID            `json:"task_id"`
	Intent               string            `json:"intent"`
	Mode                 string            `json:"mode"`
	SurfaceHint          string            `json:"surface_hint"`
	RequiredCapabilities []string          `json:"required_capabilities"`
	Inputs               map[string]string `json:"inputs"`
	OutputContract       string            `json:"output_contract"`
	Constraints          []string          `json:"constraints"`
}

// ApprovalRequest halts automation by elevating a specific decision back to Layer 5/A2UI.
// This is critical for preventing unprompted recursive deep research or unauthorized memory mutations.
type ApprovalRequest struct {
	TraceID     string                 `json:"trace_id"`
	TaskID      TaskID                 `json:"task_id"`
	ApprovalFor string                 `json:"approval_for"`
	Reason      string                 `json:"reason"`
	Options     []string               `json:"options,omitempty"`
	Payload     map[string]interface{} `json:"payload,omitempty"`
}

// ProposedToolRequest acts as a request envelope for Layer 6 execution. Layer 3 lacks authorization to execute directly.
type ProposedToolRequest struct {
	TraceID string                 `json:"trace_id"`
	TaskID  TaskID                 `json:"task_id"`
	Tool    string                 `json:"tool"`
	Payload map[string]interface{} `json:"payload"`
}

// DAGPlan is a placeholder representation of a loaded workflow blueprint.
type DAGPlan struct {
	WorkflowID WorkflowID
}

// TaskNode represents a node in the executing DAG.
type TaskNode struct {
	ID TaskID
}

// PlanningContext is the localized state context available to a task node.
type PlanningContext struct {
	SanitizedInputs map[string]string
}
