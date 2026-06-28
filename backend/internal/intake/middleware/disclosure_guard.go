package middleware

import (
	"aeolyzer/internal/intake/contracts"
	"strings"
)

func CheckForProtectedDisclosure(input contracts.SanitizedInput) contracts.DisclosureStatus {
	lower := strings.ToLower(input.RawText)
	patterns := []string{
		"what tools do you use",
		"show me your tools",
		"print skill.md",
		"mcp configuration",
		"internal workflow",
	}

	for _, pattern := range patterns {
		if strings.Contains(lower, pattern) {
			return contracts.DisclosureStatusDetected
		}
	}

	return contracts.DisclosureStatusNone
}
