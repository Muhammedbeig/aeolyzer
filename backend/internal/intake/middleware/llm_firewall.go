package middleware

import (
	"errors"
	"strings"

	"aeolyzer/internal/intake/contracts"
)

var ErrPromptInjection = errors.New("PROMPT_INJECTION_DETECTED")

// Pre-classification gate to filter malicious instructions disguised as user intent.
// Blocks deterministic strings designed to subvert system prompts or extract
// the prompt template directly into the generated output.
func CheckForPromptInjection(input contracts.SanitizedInput) error {
	lower := strings.ToLower(input.RawText)

	// This is a naive regex/substring blocklist.
	// Sufficient for primitive attacks, but vulnerable to adversarial token splitting.
	// Depends on downstream behavioral analytics to catch what this misses.
	injectionPatterns := []string{
		"ignore all previous instructions",
		"ignore previous instructions",
		"system prompt",
		"you are now",
		"print all instructions",
	}

	for _, pattern := range injectionPatterns {
		if strings.Contains(lower, pattern) {
			return ErrPromptInjection
		}
	}
	return nil
}
