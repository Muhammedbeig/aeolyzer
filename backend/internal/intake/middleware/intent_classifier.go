package middleware

import (
	"strings"

	"aeolyzer/internal/intake/contracts"
)

func ClassifyIntent(input contracts.SanitizedInput) (contracts.Intent, float64, error) {
	// Try content intent first
	intent, conf, err := ClassifyContentIntent(input)
	if err == nil && conf > 0.5 {
		return intent, conf, nil
	}

	text := strings.ToLower(input.RawText)

	if strings.Contains(text, "audit seo") || strings.Contains(text, "seo audit") {
		return contracts.IntentAuditSEO, 0.9, nil
	} else if strings.Contains(text, "site health") {
		return contracts.IntentSiteHealth, 0.9, nil
	} else if strings.Contains(text, "analyze page") {
		return contracts.IntentAnalyzePage, 0.9, nil
	} else if strings.Contains(text, "crawl site") {
		return contracts.IntentCrawlSite, 0.9, nil
	}

	if intent == contracts.IntentProtectedDisclosure {
		return intent, 1.0, ErrProtectedDisclosure
	}

	return contracts.IntentFallbackClarification, 0.4, ErrLowConfidenceIntent
}
