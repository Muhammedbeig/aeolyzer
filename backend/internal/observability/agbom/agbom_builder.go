package agbom

import "errors"

// RuntimeAgentBillOfMaterials tracks everything a session touches (Section 4.6).
type RuntimeAgentBillOfMaterials struct {
	TraceID    string   `json:"trace_id"`
	SkillsUsed []string `json:"skills_used"`
	ToolsUsed  []string `json:"tools_used"`
	Connectors []string `json:"connectors"`
}

// BuildAgBOM asserts that a valid AgBOM must be grounded to a Trace.
func BuildAgBOM(traceID string) (*RuntimeAgentBillOfMaterials, error) {
	if traceID == "" {
		return nil, errors.New("TRACE_REQUIRED_FOR_AGBOM")
	}
	return &RuntimeAgentBillOfMaterials{
		TraceID: traceID,
	}, nil
}
