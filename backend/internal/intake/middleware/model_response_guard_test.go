package middleware

import (
	"strings"
	"testing"

	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

func TestGuardModelResponseRemovesThoughtsAndProtectedMetadata(t *testing.T) {
	t.Parallel()

	response := &model.LLMResponse{
		Content: &genai.Content{
			Role: genai.RoleModel,
			Parts: []*genai.Part{
				{Thought: true, Text: "private reasoning"},
				{Text: "The workflow_id is internal."},
			},
		},
	}
	guarded, err := GuardModelResponse(nil, response, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(guarded.Content.Parts) != 1 {
		t.Fatalf("GuardModelResponse() parts = %d, want 1", len(guarded.Content.Parts))
	}
	if strings.Contains(guarded.Content.Parts[0].Text, "workflow_id") {
		t.Fatal("GuardModelResponse() retained protected metadata")
	}
}
