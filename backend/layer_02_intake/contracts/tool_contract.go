package contracts

type ProposedToolRequest struct {
	TraceID   string                 `json:"trace_id"`
	TaskID    string                 `json:"task_id"`
	ToolName  string                 `json:"tool_name"`
	Params    map[string]interface{} `json:"params"`
}

type AuthorizedToolRequest struct {
	TraceID     string                 `json:"trace_id"`
	TaskID      string                 `json:"task_id"`
	ToolName    string                 `json:"tool_name"`
	ActionClass string                 `json:"action_class"`
	Params      map[string]interface{} `json:"params"`
}

type SanitizedInput struct {
	RawText string `json:"raw_text"`
}

type SanitizedSelectedText struct {
	Text string `json:"text"`
	Hash string `json:"hash"`
}
