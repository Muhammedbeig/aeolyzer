package middleware

import (
	"errors"
	"strings"

	"aeolyzer/layer_02_intake/contracts"
)

var (
	ErrToolNotAllowed    = errors.New("TOOL_NOT_ALLOWED")
	// Thrown when an output path attempts directory traversal or absolute pathing.
	// Used to prevent accidental overwriting of host configuration or secrets via
	// dynamically generated tool parameters.
	ErrUnsafeToolPayload = errors.New("UNSAFE_TOOL_PAYLOAD")
)

// Maps raw, dynamically generated tool names into canonical internal action classes.
// Decouples the external LLM/MCP tool definition from internal policy evaluation,
// preventing policy bypass if an attacker hallucinates an undocumented tool name.
func ClassifyActionClass(toolName string, params map[string]interface{}) (string, error) {
	switch toolName {
	case "readBrandContext":
		return "read_brand_context", nil
	case "webResearch":
		return "web_research", nil
	case "deepResearch":
		return "deep_research", nil
	case "pageScrape":
		return "page_scrape", nil
	case "sitePageDiscovery":
		return "site_page_discovery", nil
	case "cannibalizationCheck":
		return "cannibalization_check", nil
	case "askUserQuestion":
		return "ask_user_question", nil
	case "updateBrief":
		return "update_brief", nil
	case "readMemory":
		return "read_memory_or_tone", nil
	case "proposeMemoryUpdate":
		return "propose_memory_update", nil
	case "setContent":
		return "set_content_type", nil
	case "canvasWrite":
		return "canvas_write", nil
	case "canvasEdit":
		return "canvas_edit", nil
	case "seoUpdate":
		return "seo_support_update", nil
	default:
		// Default-deny for unknown tools. Prevents newly added un-audited tools 
		// from running without explicit policy mappings.
		return "", ErrToolNotAllowed
	}
}

// Intercepts tool execution requests before they reach the sandbox.
// Fails the request if the derived action class lacks explicit JIT approval
// or contains dangerous path primitives.
func ValidateContentToolPolicy(req contracts.ProposedToolRequest, decision contracts.IntakeDecision) error {
	class, err := ClassifyActionClass(req.ToolName, req.Params)
	if err != nil {
		return err
	}
	
	if err := ValidateApprovalForTool(decision.Intent, class, decision.ApprovedActions); err != nil {
		return err
	}

	if err := ValidateNoArbitraryOutputPath(req.Params); err != nil {
		return err
	}

	return nil
}

// Ensures state mutations strictly align with the current orchestration mode.
// Rejects write/edit requests originating from read-only (plan) mode sessions,
// preventing speculative plans from persisting side effects.
func ValidateContentSurfaceMutation(req contracts.ProposedToolRequest, decision contracts.IntakeDecision) error {
	class, _ := ClassifyActionClass(req.ToolName, req.Params)
	
	if class == "canvas_write" && decision.Mode != contracts.ModeWrite {
		return ErrWriteModeRequired
	}
	if class == "canvas_edit" && decision.Mode != contracts.ModeEdit {
		return ErrEditModeRequired
	}
	return nil
}

// Blocks path traversal (..) and absolute paths (/) in output_path parameters.
// Only relative, sandboxed paths are permitted.
func ValidateNoArbitraryOutputPath(params map[string]interface{}) error {
	if path, ok := params["output_path"].(string); ok && path != "" {
		if strings.Contains(path, "..") || strings.HasPrefix(path, "/") {
			return ErrUnsafeToolPayload
		}
	}
	return nil
}
