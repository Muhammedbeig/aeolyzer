package middleware

import (
	"errors"
	"strings"

	"aeolyzer/internal/intake/contracts"
)

var (
	ErrProtectedMetadataOutbound = errors.New("PROTECTED_METADATA_OUTBOUND")
)

type Redaction struct {
	Pattern  string
	Replaced bool
}

// Scans outbound response buffers for sensitive internal configuration, 
// traces, or tool details that the agent may have hallucinated or leaked 
// during the generation process. 
func RedactProtectedMetadata(text string) (string, []Redaction, error) {
	if ContainsProtectedMetadata(text) {
		redactions := []Redaction{}
		
		// Hardcoded suppression of MCP tooling endpoints and internal workflow routing.
		patterns := []string{
			"workflow_id", "profile_id", "route_id", "trace_id",
			"SKILL.md", "mcp://", "https://mcp",
		}

		// Replaces leaking constants with [REDACTED] to maintain response structure 
		// without exposing the exact internal taxonomy.
		for _, p := range patterns {
			if strings.Contains(text, p) {
				text = strings.ReplaceAll(text, p, "[REDACTED]")
				redactions = append(redactions, Redaction{Pattern: p, Replaced: true})
			}
		}

		return text, redactions, nil
	}
	
	return text, nil, nil
}

func ContainsProtectedMetadata(text string) bool {
	// Secret detection (token/cookie) is naive here; assumes preceding nodes 
	// emit specific strings. Should ideally be supplemented with entropy analysis.
	patterns := []string{
		"workflow_id", "profile_id", "route_id", "trace_id",
		"SKILL.md", "mcp://", "https://mcp", "secret", "token", "cookie",
	}
	
	for _, p := range patterns {
		if strings.Contains(text, p) {
			return true
		}
	}
	return false
}

// Fallback capability response when the agent attempts to reveal exact toolsets.
// Provides a consumer-safe abstraction of system capabilities instead of a raw tool dump.
func SafeCapabilitySummary(intent contracts.Intent) string {
	return "I can help with topic discovery, content briefs, source-backed research, SEO planning, page analysis, article planning, guarded drafting, optimization, repurposing, and tone preference handling."
}
