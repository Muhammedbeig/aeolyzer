package middleware

import (
	"errors"
	"strings"

	"aeolyzer/internal/intake/contracts"
)

var (
	ErrUnknownIntent       = errors.New("UNKNOWN_INTENT")
	// Emitted when intent probabilities fall below threshold. Used to trigger 
	// deterministic disambiguation flows rather than guessing intent.
	ErrLowConfidenceIntent = errors.New("LOW_CONFIDENCE_INTENT")
	// Blocks probes designed to exfiltrate internal system boundaries.
	ErrProtectedDisclosure = errors.New("PROTECTED_DISCLOSURE")
)

// Maps raw, untrusted user strings into strict, statically typed intent enums.
// Eliminates the risk of orchestrating on hallucinated or adversarial workflow identifiers.
func ClassifyContentIntent(input contracts.SanitizedInput) (contracts.Intent, float64, error) {
	text := strings.ToLower(input.RawText)
	
	// Fast-path heuristic classification. 
	// Production path should swap this for a small, fast classifier model (e.g. fast-text).
	if strings.Contains(text, "find topic ideas") {
		return contracts.IntentTopicDiscovery, 0.9, nil
	} else if strings.Contains(text, "build a content brief") {
		return contracts.IntentContentBrief, 0.9, nil
	} else if strings.Contains(text, "research sources for") {
		return contracts.IntentContentResearch, 0.9, nil
	} else if strings.Contains(text, "plan keywords/internal links") {
		return contracts.IntentSEOPlanning, 0.9, nil
	} else if strings.Contains(text, "audit this url") {
		return contracts.IntentPageAnalysis, 0.9, nil
	} else if strings.Contains(text, "make an outline") {
		return contracts.IntentArticlePlanning, 0.9, nil
	} else if strings.Contains(text, "write the article") {
		return contracts.IntentDraftArticle, 0.9, nil
	} else if strings.Contains(text, "improve this article") {
		return contracts.IntentOptimizeContent, 0.9, nil
	} else if strings.Contains(text, "turn this into linkedin post") {
		return contracts.IntentRepurposeContent, 0.9, nil
	} else if strings.Contains(text, "edit selected paragraph") {
		return contracts.IntentEditExisting, 0.9, nil
	} else if strings.Contains(text, "remember this tone rule") {
		return contracts.IntentUpdateMemory, 0.9, nil
	} else if strings.Contains(text, "what tools do you use exactly") {
		return contracts.IntentProtectedDisclosure, 1.0, ErrProtectedDisclosure
	}

	// Default to clarification rather than hallucinating an intent.
	return contracts.IntentFallbackClarification, 0.4, ErrLowConfidenceIntent
}

// Optimization: used to quickly branch logic between content-generation pipelines 
// and raw technical SEO diagnostics without pulling in the full enum schema.
func IsContentIntent(intent contracts.Intent) bool {
	switch intent {
	case contracts.IntentTopicDiscovery, contracts.IntentContentBrief,
		contracts.IntentContentResearch, contracts.IntentSEOPlanning,
		contracts.IntentPageAnalysis, contracts.IntentArticlePlanning,
		contracts.IntentDraftArticle, contracts.IntentOptimizeContent,
		contracts.IntentRepurposeContent, contracts.IntentSwitchContentType,
		contracts.IntentEditExisting, contracts.IntentMemoryToneManagement,
		contracts.IntentUpdateMemory:
		return true
	default:
		return false
	}
}

// Verifies that a derived intent maps perfectly to the strictly typed schema boundary.
func ValidateIntentEnum(intent contracts.Intent) error {
	if IsContentIntent(intent) {
		return nil
	}
	switch intent {
	case contracts.IntentAuditSEO, contracts.IntentSiteHealth, contracts.IntentAnalyzePage,
		contracts.IntentCrawlSite, contracts.IntentAnalyzeAIVisibility, contracts.IntentAnalyzeGSC,
		contracts.IntentAnalyzeGA4, contracts.IntentAnalyzeCitations, contracts.IntentAnalyzeBrandFacts,
		contracts.IntentAnalyzeSentiment, contracts.IntentGenerateSchema, contracts.IntentGenerateLLMsTxt,
		contracts.IntentGenerateRobotsTxt, contracts.IntentGenerateSitemap, contracts.IntentOptimizeMetadata,
		contracts.IntentAnalyzeInternalLinks, contracts.IntentAnalyzeCoreWebVitals, contracts.IntentContentStrategy,
		contracts.IntentCapabilityExplanation, contracts.IntentDocumentationLookup, contracts.IntentFallbackClarification,
		contracts.IntentProtectedDisclosure, contracts.IntentOutOfBounds:
		return nil
	}
	return ErrUnknownIntent
}
