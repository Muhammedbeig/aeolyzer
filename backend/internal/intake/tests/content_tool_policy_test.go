package tests

import (
	"aeolyzer/internal/intake/contracts"
	"aeolyzer/internal/intake/middleware"
	"testing"
)

func TestContentToolPolicy(t *testing.T) {
	decision := contracts.IntakeDecision{
		Intent: contracts.IntentDraftArticle,
		Mode:   contracts.ModeWrite,
	}

	t.Run("safe tool", func(t *testing.T) {
		req := contracts.ProposedToolRequest{
			ToolName: "readBrandContext",
			Params:   map[string]interface{}{},
		}
		err := middleware.ValidateContentToolPolicy(req, decision)
		if err != nil {
			t.Errorf("Expected nil error for safe tool, got %v", err)
		}
	})

	t.Run("unsafe output path", func(t *testing.T) {
		req := contracts.ProposedToolRequest{
			ToolName: "canvasWrite",
			Params:   map[string]interface{}{"output_path": "/etc/passwd"},
		}
		err := middleware.ValidateContentToolPolicy(req, decision)
		if err == nil {
			t.Errorf("Expected error for unsafe output path")
		}
	})

	t.Run("unknown tool", func(t *testing.T) {
		req := contracts.ProposedToolRequest{
			ToolName: "unknownTool",
			Params:   map[string]interface{}{},
		}
		err := middleware.ValidateContentToolPolicy(req, decision)
		if err == nil {
			t.Errorf("Expected error for unknown tool")
		}
	})

	t.Run("canvas edit without edit mode", func(t *testing.T) {
		req := contracts.ProposedToolRequest{
			ToolName: "canvasEdit",
		}
		err := middleware.ValidateContentSurfaceMutation(req, decision)
		if err == nil {
			t.Errorf("Expected error for canvas edit without edit mode")
		}
	})
}
