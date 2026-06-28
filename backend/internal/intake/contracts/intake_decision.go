package contracts

type Intent string

const (
	IntentAuditSEO              Intent = "audit_seo"
	IntentSiteHealth            Intent = "site_health"
	IntentAnalyzePage           Intent = "analyze_page"
	IntentCrawlSite             Intent = "crawl_site"
	IntentAnalyzeAIVisibility   Intent = "analyze_ai_visibility"
	IntentAnalyzeGSC            Intent = "analyze_gsc"
	IntentAnalyzeGA4            Intent = "analyze_ga4"
	IntentAnalyzeCitations      Intent = "analyze_citations"
	IntentAnalyzeBrandFacts     Intent = "analyze_brand_facts"
	IntentAnalyzeSentiment      Intent = "analyze_sentiment"
	IntentGenerateSchema        Intent = "generate_schema"
	IntentGenerateLLMsTxt       Intent = "generate_llms_txt"
	IntentGenerateRobotsTxt     Intent = "generate_robots_txt"
	IntentGenerateSitemap       Intent = "generate_sitemap"
	IntentOptimizeMetadata      Intent = "optimize_metadata"
	IntentAnalyzeInternalLinks  Intent = "analyze_internal_links"
	IntentAnalyzeCoreWebVitals  Intent = "analyze_core_web_vitals"
	IntentContentStrategy       Intent = "content_strategy"

	IntentTopicDiscovery        Intent = "topic_discovery"
	IntentContentBrief          Intent = "content_brief"
	IntentContentResearch       Intent = "content_research"
	IntentSEOPlanning           Intent = "seo_planning"
	IntentPageAnalysis          Intent = "page_analysis"
	IntentArticlePlanning       Intent = "article_planning"
	IntentDraftArticle          Intent = "draft_article"
	IntentOptimizeContent       Intent = "optimize_content"
	IntentRepurposeContent      Intent = "repurpose_content"
	IntentSwitchContentType     Intent = "switch_content_type"
	IntentEditExisting          Intent = "edit_existing"
	IntentMemoryToneManagement  Intent = "memory_tone_management"
	IntentUpdateMemory          Intent = "update_memory"

	IntentCapabilityExplanation Intent = "capability_explanation"
	IntentDocumentationLookup   Intent = "documentation_lookup"
	IntentFallbackClarification Intent = "fallback_clarification"
	IntentProtectedDisclosure   Intent = "protected_disclosure_request"
	IntentOutOfBounds           Intent = "out_of_bounds"
)

type OrchestrationMode string

const (
	ModePlan     OrchestrationMode = "plan"
	ModeWrite    OrchestrationMode = "write"
	ModeEdit     OrchestrationMode = "edit"
	ModeOptimize OrchestrationMode = "optimize"
	ModeAudit    OrchestrationMode = "audit"
)

type DisclosureStatus string

const (
	DisclosureStatusNone     DisclosureStatus = "none"
	DisclosureStatusDetected DisclosureStatus = "detected"
)

type PolicyState string

const (
	PolicyStateAllowed PolicyState = "allowed"
	PolicyStateBlocked PolicyState = "blocked"
)

type IntakeDecision struct {
	TraceID          string                 `json:"trace_id"`
	Intent           Intent                 `json:"intent"`
	Confidence       float64                `json:"confidence"`
	SanitizedContext map[string]string      `json:"sanitized_context"`
	DisclosureStatus DisclosureStatus       `json:"disclosure_status,omitempty"`
	PolicyState      PolicyState            `json:"policy_state,omitempty"`
	Mode             OrchestrationMode      `json:"mode,omitempty"`
	ApprovedActions  []ApprovedAction       `json:"approved_actions,omitempty"`
	SafetyClasses    []string               `json:"safety_classes,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
}
