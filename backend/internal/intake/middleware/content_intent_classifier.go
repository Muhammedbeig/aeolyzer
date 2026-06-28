package middleware

import (
	"errors"
	"strings"

	"aeolyzer/internal/intake/contracts"
)

var (
	ErrUnknownIntent = errors.New("unknown intent")
	// Emitted when intent probabilities fall below threshold. Used to trigger
	// deterministic disambiguation flows rather than guessing intent.
	ErrLowConfidenceIntent = errors.New("intent confidence is too low")
	// Blocks probes designed to exfiltrate internal system boundaries.
	ErrProtectedDisclosure = errors.New("protected disclosure request detected")
)

type intentRule struct {
	intent  contracts.Intent
	signals []string
}

var contentIntentRules = []intentRule{
	{intent: contracts.IntentTopicDiscovery, signals: []string{"topic ideas", "content ideas", "what should we write", "topic gaps"}},
	{intent: contracts.IntentContentBrief, signals: []string{"content brief", "article brief", "build a brief", "prepare a brief"}},
	{intent: contracts.IntentContentResearch, signals: []string{"research sources", "collect evidence", "source backed research", "find credible sources"}},
	{intent: contracts.IntentSEOPlanning, signals: []string{"plan keywords", "seo plan", "internal link plan", "keyword strategy"}},
	{intent: contracts.IntentPageAnalysis, signals: []string{"audit this url", "analyze this page", "review this page", "page content audit"}},
	{intent: contracts.IntentArticlePlanning, signals: []string{"article outline", "make an outline", "article plan", "plan the sections"}},
	{intent: contracts.IntentDraftArticle, signals: []string{"write the article", "draft the article", "write this section", "draft a blog post"}},
	{intent: contracts.IntentOptimizeContent, signals: []string{"improve this article", "optimize this content", "refresh this page", "improve existing content"}},
	{intent: contracts.IntentRepurposeContent, signals: []string{"turn this into", "repurpose this", "convert this article", "make a linkedin post"}},
	{intent: contracts.IntentSwitchContentType, signals: []string{"switch content type", "change this to a", "convert format"}},
	{intent: contracts.IntentEditExisting, signals: []string{"edit selected", "rewrite selected", "change this paragraph", "edit this passage"}},
	{intent: contracts.IntentMemoryToneManagement, signals: []string{"use our tone", "brand voice", "tone preferences", "writing voice"}},
	{intent: contracts.IntentUpdateMemory, signals: []string{"remember this tone", "save this preference", "update memory", "remember this rule"}},
}

// Maps raw, untrusted user strings into strict, statically typed intent enums.
// Eliminates the risk of orchestrating on hallucinated or adversarial workflow identifiers.
func ClassifyContentIntent(input contracts.SanitizedInput) (contracts.Intent, float64, error) {
	if len(input.RawText) == 0 || len(input.RawText) > 64<<10 {
		return contracts.IntentFallbackClarification, 0, ErrLowConfidenceIntent
	}
	text := normalizeSecurityText(input.RawText)
	if CheckForProtectedDisclosure(input) == contracts.DisclosureStatusDetected {
		return contracts.IntentProtectedDisclosure, 1.0, ErrProtectedDisclosure
	}

	bestScore := 0
	var bestIntent contracts.Intent
	tied := false
	for _, rule := range contentIntentRules {
		score := 0
		for _, signal := range rule.signals {
			if strings.Contains(text, normalizeSecurityText(signal)) {
				score++
			}
		}
		if score > bestScore {
			bestScore = score
			bestIntent = rule.intent
			tied = false
		} else if score > 0 && score == bestScore {
			tied = true
		}
	}
	if bestScore == 0 || tied {
		return contracts.IntentFallbackClarification, 0.4, ErrLowConfidenceIntent
	}
	confidence := 0.75 + float64(bestScore-1)*0.1
	if confidence > 0.95 {
		confidence = 0.95
	}
	return bestIntent, confidence, nil
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
