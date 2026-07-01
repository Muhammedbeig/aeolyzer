# Layer 2 Production-Ready Specs v2
## Request Intake, Context, Safety, Policy, and Protected Disclosure Guard for SEO/AEO Auditor + Content Agent

**Status:** Production-ready upgrade  
**Supersedes:** Layer 2 v1 where routing enums and policy gates were SEO/AEO-first  
**Primary rule:** Layer 2 is still only the intake and safety gate. It blocks, classifies, sanitizes, authorizes, redacts, validates approval metadata, and emits sanitized safety facts. It never orchestrates, executes, retrieves skills, renders UI, manages runtime, connects MCP, stores telemetry, updates memory, or decides workflows.

---

## 1. Upgrade Decision

Upgrade is needed.

The old Layer 2 spec is correct as a safety boundary, but it is no longer complete after Layer 3 v2. The existing routing schema only covers the original SEO/AEO intent set, such as audit, analytics, schema, metadata, internal links, Core Web Vitals, protected disclosure, fallback clarification, and out-of-bounds. Layer 3 v2 now depends on additional content-agent intents, mode flags, approval metadata, selected-text validation, and content-surface safety gates.

Layer 2 must therefore be upgraded to support:

```text
topic_discovery
content_brief
content_research
seo_planning
page_analysis
article_planning
draft_article
optimize_content
repurpose_content
switch_content_type
edit_existing
memory_tone_management
update_memory
```

The rules stay the same:

```text
Layer 2 owns raw intake, firewall, disclosure guard, intent enum classification, context sanitization, tool authorization, approval metadata validation, and outbound redaction.
Layer 3 owns workflow/profile routing and DAG planning.
Layer 4 owns skills.
Layer 5 owns UI/canvas/brief/chat rendering.
Layer 6 owns execution and runtime enforcement.
Layer 7 owns MCP/connectors/RAG/data plumbing.
Layer 8 owns telemetry, evals, drift, SecOps, quarantine, and long-term audit storage.
```

---

## 2. Layer Objective

Layer 2 receives raw user input and proposed action payloads before any agentic execution.

It must:

1. Normalize raw intake.
2. Detect and block prompt injection.
3. Detect protected internal-disclosure requests.
4. Classify sanitized user intent into strict enums.
5. Extract trusted context fields into sanitized context.
6. Validate content-agent modes.
7. Validate approval metadata for high-impact actions.
8. Validate selected text before edit flows.
9. Validate content-surface mutation requests before execution.
10. Validate tool names and parameter shape.
11. Sanitize and serialize proposed tool payloads.
12. Authorize tools using policy.
13. Run semantic pre-execution checks.
14. Redact outbound responses before user release.
15. Emit sanitized safety events to Layer 8.
16. Fail closed on unknown, malformed, unsafe, unapproved, over-limit, protected, or boundary-violating requests.

Layer 2 must output only safe contracts to downstream layers. It must never perform orchestration or execution.

---

## 3. Strict Layer Boundary

### 3.1 Layer 2 owns

```text
raw intake normalization
prompt-injection detection
zero-width/control-character stripping
protected-disclosure classification
intent enum classification
routing-schema validation
sanitized context extraction
mode validation
approval metadata validation
selected_text validation
tool name validation
tool parameter shape validation
tool payload sanitation
tool authorization
semantic pre-execution checks
outbound response redaction
safe product-level capability responses
sanitized intake events
```

### 3.2 Layer 2 does not own

| Responsibility | Owning layer |
|---|---|
| Workflow selection, DAG construction, profile selection, task sequencing | Layer 3 |
| Skill registry, SKILL.md loading, skill references/assets/scripts | Layer 4 |
| Canvas, brief, chat, dashboard, A2UI, approval UI rendering | Layer 5 |
| Tool execution, sandboxing, filesystem mounts, egress, JIT credentials, quarantine execution | Layer 6 |
| MCP transport, connectors, RAG, vector retrieval, API clients | Layer 7 |
| Telemetry storage, drift scoring, eval scoring, AgBOM, SecOps triad, quarantine decisions | Layer 8 |

### 3.3 Content-agent non-overlap rule

Layer 2 may classify and authorize content-agent requests, but it must not:

```text
choose content workflow
assign content profile
build article plan
write article sections
edit canvas text
save brief directly
read memory documents
write memory documents
render approval card
render canvas or brief
call research tools
call writing tools
call memory tools
```

Instead, Layer 2 emits:

```text
IntakeDecision -> Layer 3
AuthorizedToolRequest -> Layer 6
SafetyEvent -> Layer 8
SafeTextResponse -> Layer 1/5 release path
```

---

## 4. Required Directory Upgrade

Add these files to Layer 2 v1.

```text
/layer_02_intake
  ├── /config
  │   ├── policies.yaml
  │   ├── routing-schema.json
  │   ├── context-schema.json
  │   ├── approval-policy.yaml
  │   └── outbound-redaction-policy.yaml
  ├── /middleware
  │   ├── llm_firewall.go
  │   ├── disclosure_guard.go
  │   ├── intent_classifier.go
  │   ├── content_intent_classifier.go
  │   ├── context_resolver.go
  │   ├── context_sanitizer.go
  │   ├── mode_gate.go
  │   ├── approval_metadata_validator.go
  │   ├── selected_text_validator.go
  │   ├── tool_policy_engine.go
  │   ├── content_tool_policy_engine.go
  │   ├── outbound_response_guard.go
  │   └── protected_metadata_redactor.go
  ├── /contracts
  │   ├── intake_decision.go
  │   ├── approval_contract.go
  │   ├── tool_contract.go
  │   └── safety_event.go
  ├── /tests
  │   ├── routing_schema_test.go
  │   ├── content_intent_test.go
  │   ├── mode_gate_test.go
  │   ├── approval_policy_test.go
  │   ├── selected_text_test.go
  │   ├── content_tool_policy_test.go
  │   ├── outbound_redaction_test.go
  │   └── boundary_test.go
  └── /intake_events
      └── safety_emitter.go
```

Do not add `/workflows`, `/capability_profiles`, `/skills`, `/mcp`, `/runtime`, `/a2ui`, `/memory`, or `/observability` ownership to Layer 2.

---

## 5. Configuration Contracts

Layer 2 owns these policy files:

```text
/config/policies.yaml
/config/routing-schema.json
/config/context-schema.json
/config/approval-policy.yaml
/config/outbound-redaction-policy.yaml
```

All config must load at process start, schema-validate, checksum-log, and fail closed if invalid.

---

## 6. `routing-schema.json` v2

### 6.1 Intent enum set

Layer 2 must classify only into these enums.

```json
{
  "version": 2,
  "policy_mode": "fail_closed",
  "allowed_intents": [
    "audit_seo",
    "site_health",
    "analyze_page",
    "crawl_site",
    "analyze_gsc",
    "analyze_ga4",
    "analyze_ai_visibility",
    "analyze_citations",
    "analyze_brand_facts",
    "analyze_sentiment",
    "generate_schema",
    "generate_llms_txt",
    "generate_robots_txt",
    "generate_sitemap",
    "optimize_metadata",
    "analyze_internal_links",
    "analyze_core_web_vitals",
    "content_strategy",

    "topic_discovery",
    "content_brief",
    "content_research",
    "seo_planning",
    "page_analysis",
    "article_planning",
    "draft_article",
    "optimize_content",
    "repurpose_content",
    "switch_content_type",
    "edit_existing",
    "memory_tone_management",
    "update_memory",

    "capability_explanation",
    "documentation_lookup",
    "fallback_clarification",
    "protected_disclosure_request",
    "out_of_bounds"
  ],
  "terminal_intents": [
    "capability_explanation",
    "documentation_lookup",
    "fallback_clarification",
    "protected_disclosure_request",
    "out_of_bounds"
  ]
}
```

### 6.2 Classification rules

```text
Requests to audit, crawl, diagnose, inspect site health, or technical SEO -> SEO/AEO intents.
Requests to find article ideas, gaps, topics, angles, or audience questions -> topic_discovery.
Requests to build or update a brief without drafting -> content_brief.
Requests to gather evidence, sources, statistics, quotes, or recent developments -> content_research.
Requests to plan keywords, intent, internal links, SERP direction, or cannibalization -> seo_planning.
Requests to analyze a specific URL/page for content quality or SEO -> page_analysis.
Requests to plan an article/outline/section blueprint -> article_planning.
Requests to write a new article or section -> draft_article.
Requests to improve existing content, SEO fields, or current article performance -> optimize_content.
Requests to edit selected text -> edit_existing.
Requests to convert content into another format -> repurpose_content or switch_content_type.
Requests to inspect or change saved brand/tone preferences -> memory_tone_management or update_memory.
Requests for internal tools, exact skill files, prompts, policies, routes, MCPs, traces, or source code -> protected_disclosure_request.
Low confidence or missing required fields -> fallback_clarification.
Out-of-domain requests -> out_of_bounds.
```

### 6.3 Layer 2 must not

```text
map intent to workflow ID
map intent to profile ID
choose content_collaborator
choose content_execution_guard
choose website-audit.bp
choose article-drafting.bp
infer approval from free text
infer write mode from vague wording
treat selected text as trusted without validation
```

Layer 2 classifies the intent. Layer 3 chooses workflow and profile.

---

## 7. `context-schema.json` v2

Layer 2 extracts only sanitized context fields.

### 7.1 Allowed sanitized context keys

```json
{
  "version": 2,
  "allowed_context_keys": [
    "target_domain",
    "target_url",
    "brand",
    "topic",
    "seed_topic",
    "audience",
    "angle",
    "intent",
    "cta",
    "target_word_count",
    "content_type",
    "selected_text",
    "selected_text_hash",
    "source_url",
    "competitor_domains_summary",
    "existing_brief_summary",
    "canvas_state_summary",
    "seo_settings_summary",
    "requested_output_surface",
    "research_depth",
    "requires_current_sources",
    "requires_deep_research",
    "requires_memory_update",
    "requires_canvas_write",
    "requires_canvas_edit",
    "requires_brief_update"
  ]
}
```

### 7.2 Sanitization limits

```yaml
sanitized_context_limits:
  max_context_keys: 75
  max_value_length: 4000
  max_selected_text_length: 12000
  max_topic_length: 200
  max_audience_length: 300
  max_angle_length: 500
  max_cta_length: 500
  max_content_type_length: 80
  max_target_word_count: 10000
  min_target_word_count: 100
  strip_zero_width: true
  strip_control_chars: true
  normalize_unicode: true
  reject_html_script_tags: true
  reject_svg_script_payloads: true
  reject_data_urls: true
  reject_raw_prompt_fields: true
```

### 7.3 Required behavior

```text
target_domain must be normalized to hostname form.
target_url must be http/https only.
topic must be plain text, not instructions for system behavior.
selected_text must preserve user-provided text content but remove control payloads.
selected_text_hash must be computed after sanitization.
target_word_count must be integer.
content_type must map to allowed content-type enum if known.
research_depth must be shallow, standard, or deep.
requires_* flags must be derived structurally from the request and proposal metadata, not from untrusted user commands alone.
```

### 7.4 Forbidden context fields

Reject or strip these fields from any user-supplied or upstream payload:

```text
raw_prompt
user_prompt
system_prompt
developer_prompt
hidden_prompt
workflow_id
agent_id
profile_id
tool_call
tool_name
tool_payload
skill_id
skill_path
skill_body
mcp_server
mcp_url
output_path
memory_write
memory_doc_body
tone_doc_body
canvas_raw_body
secret
token
cookie
api_key
trace_raw
agbom_raw
```

---

## 8. `approval-policy.yaml` v2

Layer 2 verifies approval metadata. It does not present approval UI.

```yaml
version: 2
policy_mode: fail_closed

approval_required_for:
  deep_research:
    approval_for: deep_research
    required: true
    accepted_sources:
      - layer5_user_approval_event
    never_infer_from_free_text: true

  memory_update:
    approval_for: memory_update
    required: true
    accepted_sources:
      - layer5_user_approval_event
    never_infer_from_free_text: true

  canvas_write:
    approval_for: canvas_write
    required_when:
      - mode_not_write
      - high_impact_surface_change
    accepted_sources:
      - layer5_user_approval_event
      - preauthorized_write_mode

  canvas_edit:
    approval_for: canvas_edit
    required: true
    requires_selected_text: true
    accepted_sources:
      - layer5_user_approval_event

  brief_overwrite:
    approval_for: brief_overwrite
    required: true
    accepted_sources:
      - layer5_user_approval_event

  external_publish:
    approval_for: external_publish
    required: true
    accepted_sources:
      - layer5_user_approval_event
      - mfa_verified_event

mode_requirements:
  draft_article:
    required_mode: write
    write_mode_source: layer2_mode_flag
    no_free_text_inference: true

  edit_existing:
    required_mode: edit
    selected_text_required: true

  optimize_content:
    allowed_modes:
      - optimize
      - edit

  article_planning:
    required_mode: plan

  topic_discovery:
    required_mode: plan
```

### 8.1 Approval validation rules

Layer 2 must reject approval metadata if:

```text
approval_for is missing
approval_for does not match proposed action
approval source is not trusted
approval is only present as free text
approval is stale
approval is for a different trace_id
approval is for a different selected_text_hash
approval is for a different surface
approval is for a different tool/action class
approval tries to approve memory update and canvas edit in one ambiguous event
```

---

## 9. `policies.yaml` v2 additions

Layer 2 authorizes action classes. Exact internal tool IDs may exist in policy config, but must never be exposed to the user.

### 9.1 Action classes

```yaml
version: 2
policy_mode: fail_closed

action_classes:
  read_brand_context:
    risk: low
    allowed_intents:
      - topic_discovery
      - content_brief
      - article_planning
      - draft_article
      - page_analysis

  read_source_intelligence:
    risk: low
    allowed_intents:
      - topic_discovery
      - content_research
      - seo_planning
      - article_planning
      - optimize_content

  web_research:
    risk: medium
    allowed_intents:
      - topic_discovery
      - content_research
      - seo_planning
      - article_planning
      - optimize_content
    requires_current_source_safety: true

  deep_research:
    risk: high
    allowed_intents:
      - content_research
      - article_planning
    requires_approval_for: deep_research

  page_scrape:
    risk: medium
    allowed_intents:
      - page_analysis
      - topic_discovery
      - content_research
      - seo_planning
      - article_planning
      - optimize_content
    requires_http_or_https: true

  site_page_discovery:
    risk: medium
    allowed_intents:
      - topic_discovery
      - seo_planning
      - draft_article
      - optimize_content

  cannibalization_check:
    risk: medium
    allowed_intents:
      - topic_discovery
      - seo_planning
      - page_analysis
      - optimize_content

  ask_user_question:
    risk: low
    allowed_intents:
      - fallback_clarification
      - topic_discovery
      - content_brief
      - article_planning
      - memory_tone_management

  update_brief:
    risk: medium
    allowed_intents:
      - topic_discovery
      - content_brief
      - seo_planning
      - article_planning
      - content_repurposing
    deny_overwrite_without_approval: true

  read_memory_or_tone:
    risk: medium
    allowed_intents:
      - content_brief
      - article_planning
      - draft_article
      - repurpose_content
      - memory_tone_management
    return_summary_only: true

  propose_memory_update:
    risk: high
    allowed_intents:
      - content_brief
      - optimize_content
      - memory_tone_management
      - update_memory
    requires_approval_for: memory_update

  set_content_type:
    risk: medium
    allowed_intents:
      - repurpose_content
      - switch_content_type
    require_allowed_content_type: true

  canvas_write:
    risk: high
    allowed_intents:
      - draft_article
      - repurpose_content
    requires_mode: write
    deny_direct_output_path: true

  canvas_edit:
    risk: high
    allowed_intents:
      - edit_existing
      - optimize_content
      - repurpose_content
    requires_selected_text: true
    requires_approval_for: canvas_edit
    require_exact_match_patch: true

  seo_support_update:
    risk: medium
    allowed_intents:
      - optimize_content
      - seo_planning
    deny_overwrite_existing_fields_without_approval: true
```

### 9.2 Semantic checks

`CheckActionSemantic` must block:

```text
tool payload contains prompt-injection strings
tool payload contains system/developer prompt text
tool payload contains secrets or credentials
tool payload contains raw memory or tone document bodies
tool payload contains raw internal routing/workflow/profile IDs from user
tool payload attempts to write arbitrary output paths
tool payload attempts to bypass Layer 2 authorization
tool payload attempts direct network or MCP connection
tool payload attempts filesystem operations outside approved contract
canvas write while mode != write
canvas edit without selected_text_hash
memory update without approved memory_update event
deep research without approved deep_research event
brief overwrite without approved brief_overwrite event
SEO field overwrite without explicit permission
```

---

## 10. Intake Contract v2

```go
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
```

### 10.1 Intent type

```go
type Intent string

const (
    IntentAuditSEO              Intent = "audit_seo"
    IntentSiteHealth            Intent = "site_health"
    IntentAnalyzePage           Intent = "analyze_page"
    IntentCrawlSite             Intent = "crawl_site"
    IntentAnalyzeAIVisibility   Intent = "analyze_ai_visibility"
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
```

### 10.2 Mode type

```go
type OrchestrationMode string

const (
    ModePlan     OrchestrationMode = "plan"
    ModeWrite    OrchestrationMode = "write"
    ModeEdit     OrchestrationMode = "edit"
    ModeOptimize OrchestrationMode = "optimize"
    ModeAudit    OrchestrationMode = "audit"
)
```

### 10.3 ApprovedAction type

```go
type ApprovedAction struct {
    ApprovalID       string            `json:"approval_id"`
    ApprovalFor      string            `json:"approval_for"`
    TraceID          string            `json:"trace_id"`
    Source           string            `json:"source"`
    Surface          string            `json:"surface,omitempty"`
    SelectedTextHash string            `json:"selected_text_hash,omitempty"`
    ExpiresAt        time.Time         `json:"expires_at"`
    Constraints      map[string]string `json:"constraints,omitempty"`
}
```

---

## 11. Middleware File Specs

### 11.1 `content_intent_classifier.go`

Purpose: extend intent classification for content-agent workflows while still emitting only strict enums.

Required functions:

```go
func ClassifyContentIntent(input SanitizedInput) (Intent, float64, error)
func IsContentIntent(intent Intent) bool
func ValidateIntentEnum(intent Intent) error
```

Required behavior:

```text
Never output free-text intent.
Never output workflow ID.
Never output profile ID.
Low confidence -> fallback_clarification.
Internal disclosure probe -> protected_disclosure_request.
```

### 11.2 `mode_gate.go`

Purpose: validate mode before Layer 3 receives the decision.

Required functions:

```go
func DeriveMode(input SanitizedInput, intent Intent, metadata map[string]interface{}) (OrchestrationMode, error)
func ValidateModeForIntent(intent Intent, mode OrchestrationMode) error
func RequiresWriteMode(intent Intent) bool
func RequiresEditMode(intent Intent) bool
```

Required behavior:

```text
draft_article requires write mode.
edit_existing requires edit mode and selected_text.
article_planning requires plan mode.
topic_discovery requires plan mode.
optimize_content allows optimize or edit.
memory_tone_management allows plan/write/edit/optimize but update still requires approval.
Never infer write mode from a vague phrase when policy requires explicit mode flag.
```

### 11.3 `selected_text_validator.go`

Purpose: validate selected text before targeted edits.

Required functions:

```go
func ValidateSelectedText(text string) (SanitizedSelectedText, error)
func HashSelectedText(text string) string
func ValidateSelectedTextHash(text string, expectedHash string) error
func RejectSelectedTextWithHiddenPayload(text string) error
```

Required behavior:

```text
Reject missing selected_text for edit_existing.
Reject over-limit selected_text.
Strip zero-width and unsafe control characters.
Preserve user-visible prose.
Hash sanitized text for approval binding.
Reject selected_text containing hidden system/tool instructions.
```

### 11.4 `approval_metadata_validator.go`

Purpose: verify high-impact approvals.

Required functions:

```go
func ValidateApprovalMetadata(action ApprovedAction, expected ApprovalExpectation) error
func ValidateApprovalForTool(intent Intent, actionClass string, approvals []ApprovedAction) error
func IsApprovalExpired(action ApprovedAction, now time.Time) bool
func RejectFreeTextApproval(claim string) error
```

Required behavior:

```text
Deep research requires approved deep_research.
Memory update requires approved memory_update.
Canvas edit requires approved canvas_edit and matching selected_text_hash.
Brief overwrite requires approved brief_overwrite.
External publish requires trusted approval and optional MFA status.
Never infer approval from natural language in sanitized_context.
```

### 11.5 `content_tool_policy_engine.go`

Purpose: extend Layer 2 authorization for content-agent proposed actions.

Required functions:

```go
func ClassifyActionClass(toolName string, params map[string]interface{}) (string, error)
func ValidateContentToolPolicy(req ProposedToolRequest, decision IntakeDecision) error
func ValidateContentSurfaceMutation(req ProposedToolRequest, decision IntakeDecision) error
func ValidateNoArbitraryOutputPath(params map[string]interface{}) error
```

Required behavior:

```text
Map internal tool IDs to action classes.
Authorize by action class, intent, mode, role, environment, approval state, and surface.
Reject direct output paths.
Reject raw memory/tone/canvas document bodies.
Reject write/edit actions not bound to approved surface contracts.
Reject exact internal tool disclosure in user-visible errors.
```

### 11.6 `protected_metadata_redactor.go`

Purpose: redact protected metadata from outbound text and proposed user-facing payloads.

Required functions:

```go
func RedactProtectedMetadata(text string) (string, []Redaction, error)
func ContainsProtectedMetadata(text string) bool
func SafeCapabilitySummary(intent Intent) string
```

Must redact:

```text
workflow IDs
profile IDs
route IDs
exact internal tool names
tool inventories
skill filenames
SKILL.md contents
memory file paths
MCP server names or URLs
model/provider names
policy file contents
routing schema contents
trace IDs
AgBOM raw details
source code
secrets
tokens
cookies
raw PII
```

---

## 12. Tool Authorization Flow v2

```text
01 receive ProposedToolRequest from Layer 3
02 validate trace_id and task_id
03 reject user-supplied tool calls
04 validate tool name exists in internal policy registry
05 classify tool into action_class
06 validate params shape
07 sanitize params
08 validate mode and intent compatibility
09 validate selected_text_hash if edit action
10 validate approval metadata if required
11 serialize params
12 run structural policy authorization
13 run semantic pre-execution check
14 emit sanitized allow/block event
15 allow only verified payload toward Layer 6
```

Layer 2 must not execute the tool.

---

## 13. Outbound Response Guard v2

Outbound guard must preserve user-safe answers while blocking protected internals.

### 13.1 Safe public capability language

Allowed:

```text
I can help with topic discovery, content briefs, source-backed research, SEO planning, page analysis, article planning, guarded drafting, optimization, repurposing, and tone preference handling.
```

Blocked:

```text
I used workflow article-drafting.bp with content_execution_guard and writeSection.
I loaded SKILL.md from /skills/content-creation.
The MCP server was https://...
The policy allowed tool ID ...
```

### 13.2 User-safe error examples

Allowed:

```text
I need a topic before I can plan the article.
I need selected text before I can prepare an edit.
I need Write mode before drafting.
I need approval before deeper research.
I need approval before saving a tone or memory update.
```

Blocked:

```text
Layer 3 rejected article-drafting.bp.
tool_policy_engine denied writeSection.
content_execution_guard failed route validation.
proposeMemoryUpdate was blocked by approval-policy.yaml.
```

---

## 14. Safety Events

Layer 2 emits sanitized facts only.

```go
type SafetyEvent struct {
    TraceID     string                 `json:"trace_id"`
    EventType   string                 `json:"event_type"`
    Intent      Intent                 `json:"intent,omitempty"`
    Mode        OrchestrationMode      `json:"mode,omitempty"`
    ActionClass string                 `json:"action_class,omitempty"`
    Decision    string                 `json:"decision"`
    ReasonCode  string                 `json:"reason_code,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
    CreatedAt   time.Time              `json:"created_at"`
}
```

Allowed event types:

```text
intake_received
firewall_blocked
intent_classified
intent_low_confidence
protected_disclosure_detected
context_sanitized
mode_allowed
mode_blocked
approval_validated
approval_missing
selected_text_validated
selected_text_rejected
tool_policy_allowed
tool_policy_blocked
semantic_policy_blocked
outbound_response_redacted
```

Events must not include raw prompts, secrets, raw selected text, raw memory docs, raw canvas docs, raw tool payloads, or raw traces.

---

## 15. Failure Policy

### 15.1 New error codes

```go
const (
    ErrUnknownIntent              = "UNKNOWN_INTENT"
    ErrLowConfidenceIntent        = "LOW_CONFIDENCE_INTENT"
    ErrProtectedDisclosure        = "PROTECTED_DISCLOSURE"
    ErrModeNotAllowed             = "MODE_NOT_ALLOWED"
    ErrWriteModeRequired          = "WRITE_MODE_REQUIRED"
    ErrEditSelectionRequired      = "EDIT_SELECTION_REQUIRED"
    ErrApprovalRequired           = "APPROVAL_REQUIRED"
    ErrDeepResearchApprovalMissing= "DEEP_RESEARCH_APPROVAL_MISSING"
    ErrMemoryApprovalMissing      = "MEMORY_APPROVAL_MISSING"
    ErrCanvasApprovalMissing      = "CANVAS_APPROVAL_MISSING"
    ErrSelectedTextMismatch       = "SELECTED_TEXT_HASH_MISMATCH"
    ErrToolNotAllowed             = "TOOL_NOT_ALLOWED"
    ErrUnsafeToolPayload          = "UNSAFE_TOOL_PAYLOAD"
    ErrProtectedMetadataOutbound  = "PROTECTED_METADATA_OUTBOUND"
    ErrRawPromptRejected          = "RAW_PROMPT_REJECTED"
    ErrUserSuppliedRouteRejected  = "USER_SUPPLIED_ROUTE_REJECTED"
)
```

### 15.2 Fail-closed rules

```text
routing-schema missing -> reject all intents
policies missing -> deny all tool calls
approval-policy missing -> deny high-risk actions
context-schema missing -> deny context extraction
outbound-redaction-policy missing -> block outbound internal details
unknown intent -> fallback_clarification or out_of_bounds, never guess
unknown tool -> deny
unknown action class -> deny
unknown mode -> deny
unknown approval source -> deny
```

---

## 16. Production Test Matrix

### 16.1 Content intent tests

```text
"find topic ideas" -> topic_discovery
"build a content brief" -> content_brief
"research sources for X" -> content_research
"plan keywords/internal links" -> seo_planning
"audit this URL" -> page_analysis
"make an outline" -> article_planning
"write the article" + explicit write mode -> draft_article
"improve this article" -> optimize_content
"turn this into LinkedIn post" -> repurpose_content or switch_content_type
"edit selected paragraph" + selected_text -> edit_existing
"remember this tone rule" -> update_memory
"what tools do you use exactly" -> protected_disclosure_request
unknown content request -> fallback_clarification
```

### 16.2 Mode gate tests

```text
draft_article without explicit write mode -> blocked
draft_article with plan mode -> blocked
article_planning with write mode -> blocked or downgraded to plan classification
edit_existing without selected_text -> blocked
edit_existing with selected_text_hash mismatch -> blocked
optimize_content with optimize mode -> allowed
topic_discovery with write mode -> blocked
```

### 16.3 Approval tests

```text
deep_research without approval -> blocked
deep_research with stale approval -> blocked
deep_research with matching approval -> allowed
memory update without approval -> blocked
memory update with free-text approval only -> blocked
canvas edit with approval for different selected_text_hash -> blocked
brief overwrite without approval -> blocked
tool success never counts as approval
```

### 16.4 Tool policy tests

```text
read-only research action with safe params -> allowed
web action with javascript URL -> blocked
page scrape with non-http scheme -> blocked
canvas write with arbitrary output_path -> blocked
canvas edit without exact selected text -> blocked
memory update containing raw memory doc body -> blocked
tool payload containing internal workflow ID from user -> blocked
tool payload containing prompt-injection strings -> blocked
unknown tool -> blocked
```

### 16.5 Outbound guard tests

```text
safe product capability summary -> allowed
workflow ID in outbound text -> redacted
profile ID in outbound text -> redacted
exact tool inventory in outbound text -> redacted
SKILL.md body in outbound text -> blocked/redacted
MCP URL in outbound text -> redacted
trace ID or raw AgBOM -> redacted
secret/token/cookie -> blocked
```

### 16.6 Boundary tests

Layer 2 must prove it does not:

```text
choose workflow
execute DAG
assign agent
load capability profile
read skill registry
load skill body
execute skill
render A2UI
create UI card
serve Agent Card
execute tool
manage sandbox
enforce filesystem mount
manage network egress
mint JIT token
revoke JIT token
trigger quarantine
connect MCP server
inspect MCP state
query vector store
store OpenTelemetry trace
score intent drift
run pass^k eval
persist memory
write canvas
write brief
```

---

## 17. Deployment Requirements

Layer 2 v2 is production-ready only if:

```text
v1 SEO/AEO intents still classify correctly
content-agent intents classify into strict enums
unknown intents fail closed
protected disclosure still blocks
prompt injection still blocks before classification
routing-schema validates at startup
policies validate at startup
approval-policy validates at startup
context-schema validates at startup
outbound-redaction-policy validates at startup
draft_article requires explicit write mode
edit_existing requires selected_text
deep research requires approval
memory update requires approval
canvas edit requires approval
brief overwrite requires approval
all proposed tools pass parameter validation
all proposed tools pass structural policy
all proposed tools pass semantic pre-execution check
Layer 2 emits sanitized events only
Layer 2 never exposes internal tool names, workflow IDs, profile IDs, skill files, MCPs, trace internals, memory paths, or policy contents to users
all boundary tests pass
```

### 17.1 Executable enforcement baseline

The repository implementation must retain these executable controls:

```text
config.LoadEmbedded strictly decodes routing-schema.json and policies.yaml and rejects unknown fields
config validation cross-checks every allowed intent and action class and requires fail_closed semantic denials
cmd/api validates the embedded Layer 2 configuration before opening a listener
middleware/llm_firewall normalizes Unicode with NFKC, enforces input limits, detects split role markers and bounded encoded injection, and fails closed
middleware/protected_metadata_redactor removes protected identifiers and value-level credential patterns without returning the matched value
middleware/content_intent_classifier returns a closed intent enum or clarification for ambiguous and low-confidence input
intake contracts preserve trace, tenant, approval, and authorization wire names
```

Required tests:

```text
strict config decode and unknown-field rejection
unsafe semantic-policy rejection
plain, split-token, full-width, role-marker, and base64 prompt-injection rejection
safe-input false-positive tests
credential and protected-metadata redaction tests
classifier ambiguity and unsupported-intent tests
wire-contract and sanitized-event tests
```

These controls are a baseline, not proof that all Layer 2 authorization paths
or adversarial corpora are complete.

### 17.2 Repository readiness evidence

The repository-level production check is:

```text
go run ./cmd/readiness -root .
```

For Layer 2, the check must fail closed when required policy or routing-schema
artifacts are missing, unreadable, or placeholder-only. It must also report
explicit prototype markers in Layer 2 production sources.

This check proves repository evidence only. It does not replace startup schema
validation, semantic policy validation, prompt-injection tests, authorization
tests, redaction tests, or the production test matrix above.

`cmd/readiness` and `internal/releasegate` are read-only platform CI tooling.
They must not classify intent, inspect or store raw prompts, authorize tools,
modify policies, emit runtime safety decisions, or otherwise perform Layer 2
runtime behavior.

---

## 18. Acceptance Criteria

Layer 2 v2 is accepted when:

1. It keeps all v1 zero-overlap boundaries.
2. It adds content-agent intent enums.
3. It produces `IntakeDecision` objects compatible with Layer 3 v2.
4. It validates mode for plan/write/edit/optimize/audit flows.
5. It blocks drafting without explicit write mode.
6. It blocks editing without selected text.
7. It validates approval metadata for deep research, memory updates, canvas edits, brief overwrites, and external publishing.
8. It sanitizes content-agent context without storing raw prompts, raw memory docs, raw canvas docs, or secrets.
9. It authorizes content-agent proposed actions by action class, intent, mode, role, environment, surface, and approval state.
10. It redacts protected metadata from outbound responses.
11. It emits only sanitized safety events.
12. It does not choose workflows, assign agents, load skills, execute tools, render UI, connect MCP, persist memory, store telemetry, run evals, or quarantine.

---

## 19. Final Non-Goals

```text
Layer 2 must not choose topic-discovery.bp, article-drafting.bp, or any workflow.
Layer 2 must not select content_collaborator or content_execution_guard.
Layer 2 must not build DAG nodes.
Layer 2 must not create content-generation task contracts.
Layer 2 must not write article sections.
Layer 2 must not edit canvas text.
Layer 2 must not update the brief directly.
Layer 2 must not read memory.md or tone.md.
Layer 2 must not write memory.md or tone.md.
Layer 2 must not load skill-registry.yaml.
Layer 2 must not read SKILL.md.
Layer 2 must not execute skill scripts.
Layer 2 must not run web research.
Layer 2 must not run deep research.
Layer 2 must not scrape pages.
Layer 2 must not query site pages.
Layer 2 must not check cannibalization directly.
Layer 2 must not render approval cards.
Layer 2 must not render canvas, brief, chat, tables, dashboards, or A2UI.
Layer 2 must not open MCP connections.
Layer 2 must not call connectors.
Layer 2 must not manage sandbox, filesystem, egress, packages, or credentials.
Layer 2 must not store OpenTelemetry traces.
Layer 2 must not calculate evals.
Layer 2 must not score drift.
Layer 2 must not trigger quarantine.
Layer 2 must not expose workflow IDs, profile IDs, route IDs, exact tool names, skill files, memory paths, MCP endpoints, policy contents, or trace internals to end users.
```

---

## 20. One-Line Architecture Summary

Layer 2 v2 is the fail-closed intake and policy gate that converts raw SEO/AEO and content-agent requests into sanitized intent decisions, validated modes, trusted context, approved-action metadata, and authorized tool payloads while leaving workflow planning, skill loading, execution, rendering, MCP/data access, telemetry, evaluation, memory persistence, and recovery to their owning layers.

---

## 21. Secure Chat and Attachment Intake Addendum

Layer 2 owns validation of conversation input before it can reach either agent.

Required request limits:

```yaml
chat_intake:
  text_max_runes: 12000
  attachment_count_max: 5
  idempotency_key_min_chars: 16
  idempotency_key_max_chars: 128
  require_text_or_attachment: true
```

Required behavior:

```text
validate the explicit agent enum
reject unknown multipart fields
reject duplicate agent or text fields
reject invalid UTF-8
scan chat text for prompt injection
scan text, Markdown, CSV, HTML, source, and JSON attachments in bounded overlapping chunks
pass only validated bytes to Layer 6 file processing
require the allowed browser origin
bind anonymous guest history to a signed HttpOnly SameSite cookie
never treat the guest cookie as tool authorization
never place raw message or attachment content in telemetry
```

Tests must cover empty input, text-only input, attachment-only input, excessive text, excessive attachment count, direct prompt injection, attachment-borne prompt injection, invalid agent values, malformed multipart bodies, tampered guest cookies, and disallowed origins.
