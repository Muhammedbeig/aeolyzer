# Layer 5 Production-Ready Specs v2
## Extensions, Presentation Surfaces, A2UI Translation, Approval UX, and A2A Interop

**Status:** Production-ready upgrade  
**Supersedes:** Layer 5 v1 where A2UI translation, component mapping, audit dashboard, and A2A exposure were underspecified  
**Primary rule:** Layer 5 is the declarative presentation and extension boundary. It converts already-orchestrated, already-sanitized agent results and presentation intents into safe user-facing surfaces and external agent-facing envelopes. It never classifies raw user intent, chooses workflows, loads skills, executes tools, connects MCP, runs sandboxed code, stores telemetry, scores evals, persists memory, or performs SecOps recovery.

---

## 1. Upgrade Decision

Upgrade is required.

Layer 5 v1 correctly identified the core responsibilities:

```text
A2UI translation
interactive component mappings
audit dashboard UI
A2A server
Agent Card endpoint
```

But the rest of the system has evolved into a broader SEO/AEO auditor and content-writing agent with multiple surfaces:

```text
chat coordination
content brief
canvas writing surface
SEO settings surface
audit dashboard
issue and recommendation lists
source/evidence views
approval and Vibe Diff cards
memory update proposal cards
A2A external agent entrypoints
```

Layer 5 must therefore become a production-grade extension and presentation layer that supports:

```text
SEO/AEO audit UI
Content planning UI
Guarded writing and editing UI
Approval / rejection / clarification UX
A2A external agent interop
A2UI schema validation
Safe declarative component generation
```

The upgrade preserves the architecture:

```text
Layer 1: actual frontend shell, device UI, DOM, native app, CLI, browser, and hardware MFA control
Layer 2: intake, safety, redaction, policy authorization, protected disclosure guard
Layer 3: workflow selection, DAG orchestration, agent assignment, mode/profile selection
Layer 4: procedural skill bodies and resources
Layer 5: declarative presentation contracts, A2UI conversion, extension envelopes, UI event normalization
Layer 6: isolated runtime, sandbox, JIT tokens, filesystem and egress enforcement
Layer 7: MCP, connectors, transport, RAG, data plumbing, memory retrieval
Layer 8: telemetry, observability, evaluation, drift, SecOps triad, quarantine decisioning
```

Layer 5 must not become an orchestrator, frontend implementation, tool executor, policy engine, or telemetry backend.

---

## 2. Layer Objective

Layer 5 turns structured agent outcomes into safe, interactive, user-operable interfaces and extension envelopes.

It must:

1. Accept only structured `PresentationIntent`, `SurfacePatchIntent`, `ApprovalRequestIntent`, and `A2AEnvelope` objects from trusted upstream layers.
2. Validate every presentation object against strict schemas.
3. Convert presentation intents into safe A2UI component trees.
4. Maintain the approved A2UI catalog and catalog version metadata.
5. Provide deterministic UI mappings for common SEO/AEO and content-agent surfaces.
6. Support LLM-generated A2UI only inside strict component, schema, token, and safety limits.
7. Support hybrid output: human-readable text plus structured UI parts.
8. Render approval cards, Vibe Diff summaries, and memory-update proposal cards as declarative UI contracts.
9. Normalize user interaction events back into trusted event envelopes for Layer 2 and Layer 3.
10. Expose an A2A application-level endpoint and Agent Card without owning network transport, identity issuance, or workflow routing.
11. Keep all UI action payloads non-executable.
12. Prevent XSS, event-handler injection, component-smuggling, malformed JSON, hidden payloads, stale-event replay, and protected metadata leaks.
13. Emit sanitized UI and A2A events to Layer 8.
14. Never mutate persistent business state directly.
15. Fail closed when a component schema, action binding, approval token, surface version, or external A2A envelope is invalid.

Layer 5 exists to answer this question:

```text
Given a validated workflow result or presentation intent, what safe user interface or external agent envelope should represent it?
```

Layer 5 must not answer:

```text
What did the user ask?
Which workflow should run?
Which agent should act?
Which skill should load?
Which tool should execute?
Which connector should fetch data?
Which credential should be issued?
Was the agent correct?
Should the runtime be quarantined?
```

---

## 3. Strict Layer Boundary

### 3.1 Layer 5 owns

```text
presentation intent validation
A2UI catalog manifests
A2UI schema versions
A2UI component tree validation
A2UI part conversion
A2UI safety filtering
component allowlists
surface schema definitions
surface patch contracts
approval card contracts
Vibe Diff presentation schema
memory update proposal card schema
question card schema
audit dashboard schema
content brief surface schema
canvas surface schema
SEO settings surface schema
source/evidence table schema
content quality checklist schema
UI action binding contracts
UI event normalization
UI event idempotency tokens
stale event rejection
safe markdown rendering policy
safe URL rendering policy
A2A Agent Card schema
A2A extension capability advertisement
A2A inbound envelope validation
A2A outbound envelope formatting
A2A-to-internal request normalization
presentation catalog tests
A2UI regression fixtures
accessibility metadata requirements
sanitized UI/A2A event emission to Layer 8
```

### 3.2 Layer 5 does not own

| Responsibility | Owning layer |
|---|---|
| Actual browser DOM rendering, mobile/native rendering, CLI display, client styling | Layer 1 |
| Hardware MFA challenge execution | Layer 1, with Layer 2/6 policy/runtime support |
| Raw user-input normalization | Layer 2 |
| Prompt-injection detection | Layer 2 |
| Protected-disclosure detection and outbound redaction policy | Layer 2 |
| Intent classification | Layer 2 |
| Tool authorization | Layer 2 |
| Workflow selection | Layer 3 |
| DAG construction and sequencing | Layer 3 |
| Agent assignment | Layer 3 |
| Capability profile/mode selection | Layer 3 |
| Skill loading | Layer 4 |
| SKILL.md or asset ownership | Layer 4 |
| Script execution | Layer 6 |
| Tool execution | Layer 6 |
| Sandbox, filesystem, egress, JIT credentials | Layer 6 |
| MCP, API connectors, A2A transport, mTLS, RAG, memory retrieval | Layer 7 |
| Memory document persistence | Layer 7, governed by Layer 2/3 approval contracts |
| Telemetry storage | Layer 8 |
| Evaluation scoring | Layer 8 |
| Drift scoring and trust decay | Layer 8 |
| Red/blue/green SecOps loops | Layer 8 |
| Quarantine decisioning | Layer 8 |
| Quarantine execution | Layer 6 |

### 3.3 Non-overlap rule for Layer 5

Allowed:

```text
Build an approval card from an upstream approval request.
Render a diff as safe declarative UI.
Validate that a UI action has a known action_id and current surface_version.
Transform a structured audit report into a dashboard layout.
Transform an external A2A request envelope into a normalized internal request envelope.
Advertise public product capabilities through an Agent Card.
Normalize user clicks, choices, and form submissions into typed events.
Emit sanitized presentation events to Layer 8.
```

Forbidden:

```text
Classify the user's natural-language request.
Choose the website-audit workflow.
Choose the content-writing workflow.
Load or inspect SKILL.md bodies.
Call web search, crawl, GSC, GA4, or page scrape tools.
Execute scripts.
Open MCP connections.
Mint credentials.
Persist memory.
Store telemetry.
Score evals.
Start quarantine.
Generate arbitrary executable UI code.
Render raw HTML from the model.
Expose internal tool names, workflow IDs, file paths, SKILL.md bodies, trace internals, or connector details.
```

---

## 4. Whitepaper-Aligned Design Principles

### 4.1 A2UI is declarative sheet music, not executable UI

Layer 5 must follow the A2UI separation-of-concerns model:

```text
Agent or converter declares UI intent.
Layer 5 validates the declarative UI contract.
Layer 1 renders it using trusted components.
No untrusted executable code crosses the boundary.
```

An A2UI payload is not JavaScript, CSS, HTML, JSX, React code, or template code. It is a typed component tree using only approved component names, props, children references, and action bindings.

### 4.2 UI catalogs are trusted component libraries

The model or converter may choose from a catalog. It may not invent components at runtime.

Production rule:

```text
The A2UI basic catalog may be used for prototypes.
Production surfaces must use an application catalog that maps approved design-system components to A2UI types.
```

Layer 5 owns catalog metadata and validation. Layer 1 owns component implementation.

### 4.3 A2UI is for interaction, not data plumbing

Layer 5 must never fetch data to fill a chart, scrape a URL, call analytics, or query memory. It receives already-sanitized presentation data from Layer 3 and turns it into a user-facing surface.

### 4.4 Prefer deterministic UI mapping for known surfaces

Known system surfaces should use deterministic converters, not model-generated layouts:

```text
audit dashboard
technical issue list
content brief panel
canvas status panel
SEO settings panel
approval card
memory update proposal card
question card
source evidence table
content quality checklist
```

LLM-generated A2UI is allowed only when flexible composition is needed and strict validation is applied.

### 4.5 Hybrid output is the default

Every important UI response should support:

```text
short human-readable summary
structured A2UI payload
safe fallback text
machine-readable event bindings
```

If A2UI validation fails, the user still receives a safe textual fallback.

### 4.6 Approval UX must reduce fatigue, not create reflex clicking

Layer 5 must not create endless micro-approval cards. It should batch low-risk review items into structured summary cards, reserve blocking approval UI for meaningful state changes, and make high-stakes approval cards explicit, reviewable, and diff-based.

Layer 5 does not decide whether approval is required. It renders the approval experience required by Layer 2/3.

### 4.7 Canvas, brief, and dashboard are surfaces, not workflow owners

Layer 5 owns presentation schemas for canvas, brief, chat, SEO settings, source/evidence panels, and dashboards. It must not decide what content belongs there or when it should change. It receives patch intents from Layer 3 and user events from Layer 1, validates them, then sends normalized events back to the owning workflow and policy layers.

### 4.8 A2A is an application extension, not transport

Layer 5 owns the Agent Card, A2A application envelope schemas, extension negotiation, task schemas, and public capability descriptions. Layer 7 owns A2A/MCP transport and mTLS. Layer 2 owns policy. Layer 3 owns internal workflow routing.

The external agent contract is A2A Protocol 1.0 exposed through Google ADK Go. Implementations must use `google.golang.org/adk/server/adka2a/v2` and the canonical `github.com/a2aproject/a2a-go/v2/a2a` data model. Agent Card JSON must use the A2A camelCase fields (`supportedInterfaces`, `defaultInputModes`, `securityRequirements`, and `securitySchemes`), not AEOlyzer-only snake_case fields.

OpenAPI is not the A2A compliance contract. It may be added separately for ordinary REST endpoints owned by `internal/httpapi`, but it must not replace the A2A Agent Card, A2A JSON-RPC/HTTP bindings, or ADK executor wiring.

---

## 5. Required Directory Upgrade

### 5.1 Final Layer 5 tree

```text
/layer_05_extensions
  ├── README.md
  ├── layer5-boundary.md
  ├── presentation.schema.json
  ├── a2ui-frame.schema.json
  ├── a2ui-catalog.schema.json
  ├── ui-event.schema.json
  ├── approval.schema.json
  ├── surface-patch.schema.json
  ├── a2a-agent-card.schema.json
  ├── a2a-envelope.schema.json
  ├── catalog-lock.yaml
  ├── presentation-changelog.md
  │
  ├── /config
  │   ├── layer5_policy.yaml
  │   ├── surface_registry.yaml
  │   ├── component_catalog_registry.yaml
  │   ├── action_binding_registry.yaml
  │   ├── a2a_public_capability_policy.yaml
  │   ├── markdown_safety_policy.yaml
  │   ├── url_safety_policy.yaml
  │   └── accessibility_policy.yaml
  │
  ├── /a2ui_translator
  │   ├── a2ui_schema_manager.go
  │   ├── a2ui_catalog_loader.go
  │   ├── a2ui_part_converter.go
  │   ├── a2ui_llm_output_validator.go
  │   ├── a2ui_safety_filter.go
  │   ├── a2ui_fallback_renderer.go
  │   ├── a2ui_component_resolver.go
  │   ├── a2ui_token_budgeter.go
  │   └── a2ui_event_binding_builder.go
  │
  ├── /surface_router
  │   ├── presentation_intent_validator.go
  │   ├── surface_patch_validator.go
  │   ├── surface_router.go
  │   ├── surface_state_versioner.go
  │   ├── surface_action_normalizer.go
  │   ├── stale_event_rejector.go
  │   └── fallback_surface_builder.go
  │
  ├── /interactive_components
  │   ├── ask_question.ui.yaml
  │   ├── approval_vibe_diff.ui.yaml
  │   ├── propose_edit.ui.yaml
  │   ├── propose_memory_update.ui.yaml
  │   ├── content_brief.ui.yaml
  │   ├── canvas_editor.ui.yaml
  │   ├── seo_settings.ui.yaml
  │   ├── audit_dashboard.ui.yaml
  │   ├── issue_list.ui.yaml
  │   ├── recommendation_card.ui.yaml
  │   ├── source_evidence_table.ui.yaml
  │   ├── quality_gate_checklist.ui.yaml
  │   ├── progress_timeline.ui.yaml
  │   ├── export_artifact.ui.yaml
  │   └── external_agent_status.ui.yaml
  │
  ├── /catalogs
  │   ├── /basic_v0_9
  │   │   ├── catalog.yaml
  │   │   ├── component.schema.json
  │   │   └── examples.yaml
  │   └── /seo_content_app_v1
  │       ├── catalog.yaml
  │       ├── component.schema.json
  │       ├── component_map.yaml
  │       ├── design_system_binding.yaml
  │       └── examples.yaml
  │
  ├── /approval_ux
  │   ├── approval_request_validator.go
  │   ├── vibe_diff_builder.go
  │   ├── approval_card_builder.go
  │   ├── decision_event_normalizer.go
  │   ├── approval_batcher.go
  │   ├── approval_fatigue_guard.go
  │   └── approval_copy_policy.yaml
  │
  ├── /a2a_server
  │   ├── agent_card.yaml
  │   ├── agent_card_loader.go
  │   ├── a2a_envelope_validator.go
  │   ├── a2a_extension_negotiator.go
  │   ├── a2a_request_normalizer.go
  │   ├── a2a_response_formatter.go
  │   ├── a2a_public_disclosure_filter.go
  │   └── a2a_event_emitter.go
  │
  ├── /security
  │   ├── component_allowlist.go
  │   ├── prop_sanitizer.go
  │   ├── markdown_sanitizer.go
  │   ├── url_sanitizer.go
  │   ├── action_binding_sanitizer.go
  │   ├── hidden_payload_scanner.go
  │   ├── protected_metadata_scanner.go
  │   ├── idempotency_token_validator.go
  │   └── ui_payload_signer.go
  │
  ├── /events
  │   ├── ui_event_emitter.go
  │   ├── a2a_event_emitter.go
  │   └── event_redactor.go
  │
  └── /tests
      ├── presentation_schema_test.go
      ├── a2ui_schema_test.go
      ├── a2ui_catalog_test.go
      ├── component_allowlist_test.go
      ├── prop_sanitization_test.go
      ├── markdown_safety_test.go
      ├── url_safety_test.go
      ├── action_binding_test.go
      ├── approval_card_test.go
      ├── stale_event_test.go
      ├── a2a_agent_card_test.go
      ├── a2a_envelope_test.go
      ├── protected_metadata_boundary_test.go
      ├── no_intent_classification_boundary_test.go
      ├── no_workflow_selection_boundary_test.go
      ├── no_tool_execution_boundary_test.go
      ├── no_mcp_boundary_test.go
      ├── no_memory_persistence_boundary_test.go
      └── golden_surfaces_test.go
```

---

## 6. Canonical Layer 5 Surfaces

Layer 5 must support these first-class surfaces.

```yaml
surfaces:
  chat:
    purpose: "Short human coordination and fallback text."
    mutability: "append_only_by_presentation_intent"
    owner_of_content_generation: "layer_3_agent_output"
    layer5_role: "format and render"

  canvas:
    purpose: "Long-form content, report body, draft, outline, or audit narrative."
    mutability: "patch_intent_required"
    owner_of_persistent_state: "workflow/runtime storage outside Layer 5"
    layer5_role: "render, diff, collect edit decisions"

  brief:
    purpose: "Structured content contract: topic, angle, audience, intent, CTA, constraints."
    mutability: "patch_intent_or_user_event_required"
    owner_of_strategy_decisions: "layer_3"
    layer5_role: "render fields, collect explicit edits"

  seo_settings:
    purpose: "Metadata, slug, FAQs, schema, internal-link suggestions, crawler fields."
    mutability: "approval_required_for_overwrite"
    owner_of_generation: "layer_3_plus_layer4_skill"
    layer5_role: "render form and diff"

  audit_dashboard:
    purpose: "SEO/AEO site health, issue lists, evidence summaries, priorities, next actions."
    mutability: "workflow_patch_only"
    owner_of_analysis: "layer_3"
    layer5_role: "render prioritized dashboard"

  source_evidence:
    purpose: "Source quality, citations, unsupported claims, competitor exclusion notes."
    mutability: "read_only_by_default"
    owner_of_research: "layer_3_plus_layer7_data"
    layer5_role: "render table and filters"

  approval_stream:
    purpose: "Human-in-the-loop approval cards, Vibe Diff, memory update proposals."
    mutability: "decision_event_only"
    owner_of_authorization: "layer_2"
    layer5_role: "render request and normalize decision"

  a2a_external:
    purpose: "External agent-facing task status and extension response envelopes."
    mutability: "protocol_envelope_only"
    owner_of_execution: "layer_3_and_layer6"
    layer5_role: "format and validate envelope"
```

Layer 5 may keep ephemeral surface versions, selected UI state, collapsed/expanded rows, local sort/filter configuration, and pending UI card IDs. It must not store authoritative business state.

---

## 7. Core Input and Output Contracts

### 7.1 Input from Layer 3: `PresentationIntent`

Layer 3 sends a presentation intent after workflow state is known.

```go
type PresentationIntent struct {
    TraceID          string                 `json:"trace_id"`
    WorkflowID       string                 `json:"workflow_id,omitempty"`
    NodeID           string                 `json:"node_id,omitempty"`
    Surface          string                 `json:"surface"`
    EventKind        string                 `json:"event_kind"`
    Mode             string                 `json:"mode,omitempty"`
    Priority         string                 `json:"priority,omitempty"`
    Payload          map[string]interface{} `json:"payload"`
    OutputContracts  []string               `json:"output_contracts,omitempty"`
    ApprovalRequired bool                   `json:"approval_required,omitempty"`
    FallbackText     string                 `json:"fallback_text,omitempty"`
    Metadata         map[string]string      `json:"metadata,omitempty"`
}
```

Validation:

```text
trace_id required
surface must be in surface_registry.yaml
event_kind must be known
payload must match surface schema
fallback_text must be safe markdown
workflow_id/node_id are internal and must not be surfaced to users
payload must not contain raw tool payloads unless already normalized by Layer 3
payload must not contain secrets, tokens, cookies, raw logs, raw traces, exact internal tool IDs, SKILL.md bodies, or connector details
```

### 7.2 Input from Layer 3/2: `ApprovalRequestIntent`

Approval cards must be driven by structured approval requests.

```go
type ApprovalRequestIntent struct {
    TraceID             string                 `json:"trace_id"`
    ApprovalRequestID   string                 `json:"approval_request_id"`
    RequestedActionKind string                 `json:"requested_action_kind"`
    RiskLevel           string                 `json:"risk_level"`
    Surface             string                 `json:"surface"`
    PlainEnglishSummary string                 `json:"plain_english_summary"`
    VibeDiff            VibeDiff              `json:"vibe_diff"`
    Options             []ApprovalOption      `json:"options"`
    ExpiresAt           string                 `json:"expires_at"`
    RequiresMFA         bool                   `json:"requires_mfa"`
    PolicyRef           string                 `json:"policy_ref,omitempty"`
    Metadata            map[string]string      `json:"metadata,omitempty"`
}
```

Layer 5 renders this request. It does not decide `RiskLevel`, `RequiresMFA`, or policy outcome.

### 7.3 Output to Layer 1: `A2UIFrame`

Layer 5 outputs a signed declarative frame.

```go
type A2UIFrame struct {
    FrameID        string            `json:"frame_id"`
    TraceID        string            `json:"trace_id,omitempty"`
    Surface        string            `json:"surface"`
    CatalogID      string            `json:"catalog_id"`
    CatalogVersion string            `json:"catalog_version"`
    SchemaVersion  string            `json:"schema_version"`
    RootID         string            `json:"root_id"`
    Nodes          []A2UINode        `json:"nodes"`
    Actions        []UIActionBinding `json:"actions,omitempty"`
    FallbackText   string            `json:"fallback_text,omitempty"`
    ExpiresAt      string            `json:"expires_at,omitempty"`
    Signature      string            `json:"signature"`
}
```

The frame is safe for Layer 1 to render. Layer 1 still performs client-side validation before display.

### 7.4 Output from Layer 1 to Layer 5: `UserInteractionEvent`

Layer 1 sends typed interaction events.

```go
type UserInteractionEvent struct {
    FrameID        string                 `json:"frame_id"`
    Surface        string                 `json:"surface"`
    ActionID       string                 `json:"action_id"`
    InteractionID  string                 `json:"interaction_id"`
    SurfaceVersion int64                  `json:"surface_version"`
    UserDecision   string                 `json:"user_decision,omitempty"`
    Values         map[string]interface{} `json:"values,omitempty"`
    ClientTime     string                 `json:"client_time,omitempty"`
    Signature      string                 `json:"signature"`
}
```

Layer 5 validates, normalizes, and forwards to Layer 2/3. It must reject stale, unsigned, replayed, unknown, malformed, or overlarge events.

### 7.5 Output to Layer 3: `NormalizedSurfaceEvent`

```go
type NormalizedSurfaceEvent struct {
    TraceID          string                 `json:"trace_id"`
    Surface          string                 `json:"surface"`
    EventKind        string                 `json:"event_kind"`
    ActionID         string                 `json:"action_id"`
    IdempotencyKey   string                 `json:"idempotency_key"`
    SurfaceVersion   int64                  `json:"surface_version"`
    SafeValues       map[string]interface{} `json:"safe_values,omitempty"`
    ApprovalDecision *ApprovalDecision      `json:"approval_decision,omitempty"`
    RequiresPolicy   bool                   `json:"requires_policy"`
    Metadata         map[string]string      `json:"metadata,omitempty"`
}
```

Layer 3 decides how this event affects the workflow. Layer 2 validates policy for protected or mutating actions.

---

## 8. A2UI Frame Schema

### 8.1 Component node

```go
type A2UINode struct {
    ID          string                 `json:"id"`
    Type        string                 `json:"type"`
    Props       map[string]interface{} `json:"props,omitempty"`
    Children    []string               `json:"children,omitempty"`
    Slot        string                 `json:"slot,omitempty"`
    Visibility  *VisibilityRule        `json:"visibility,omitempty"`
    Accessibility AccessibilityProps   `json:"accessibility,omitempty"`
}
```

Validation:

```text
id must be unique in frame
type must exist in active catalog
props must match component prop schema
children must reference existing ids
no cycles unless explicit catalog component permits it
no raw HTML
no raw CSS
no inline JavaScript
no function strings
no event handler strings
no iframe
no script
no object/embed
no external image without safe URL policy
no unbounded markdown
```

### 8.2 UI action binding

```go
type UIActionBinding struct {
    ActionID           string            `json:"action_id"`
    ComponentID        string            `json:"component_id"`
    Event              string            `json:"event"`
    ActionKind         string            `json:"action_kind"`
    TargetSurface      string            `json:"target_surface,omitempty"`
    RequiresApproval   bool              `json:"requires_approval"`
    RequiresPolicy     bool              `json:"requires_policy"`
    RequiresMFA        bool              `json:"requires_mfa"`
    IdempotencyRequired bool             `json:"idempotency_required"`
    ValueSchemaRef     string            `json:"value_schema_ref,omitempty"`
    Metadata           map[string]string `json:"metadata,omitempty"`
}
```

Allowed action kinds:

```yaml
action_kinds:
  view_only:
    - expand
    - collapse
    - sort
    - filter
    - paginate
    - copy_to_clipboard
  workflow_event:
    - answer_question
    - request_clarification
    - accept_recommendation
    - reject_recommendation
    - request_revision
    - select_topic
    - update_brief_candidate
    - accept_canvas_patch
    - reject_canvas_patch
    - accept_seo_setting_candidate
    - reject_seo_setting_candidate
    - accept_memory_update_candidate
    - reject_memory_update_candidate
    - approve_high_risk_action
    - deny_high_risk_action
    - submit_a2a_task_response
```

Layer 5 may perform `view_only` locally when Layer 1 supports it. Every workflow or mutating event returns to Layer 3 and/or Layer 2.

---

## 9. A2UI Generation Patterns

### 9.1 Pattern A: deterministic converter

Use deterministic converters when the surface is known.

```text
PresentationIntent
-> validate payload
-> choose surface template
-> bind data to approved components
-> validate component props
-> sign A2UIFrame
-> send to Layer 1
```

Required deterministic mappings:

```yaml
mappings:
  ask_question:
    file: interactive_components/ask_question.ui.yaml
    components: [Card, Text, ChoicePicker, Button]

  approval_vibe_diff:
    file: interactive_components/approval_vibe_diff.ui.yaml
    components: [Card, Text, DiffView, Badge, Row, Button]

  propose_edit:
    file: interactive_components/propose_edit.ui.yaml
    components: [Card, Text, DiffView, ButtonRow]

  propose_memory_update:
    file: interactive_components/propose_memory_update.ui.yaml
    components: [Card, Text, FieldList, ButtonRow]

  content_brief:
    file: interactive_components/content_brief.ui.yaml
    components: [Surface, FieldList, StatusBadge, ButtonRow]

  canvas_editor:
    file: interactive_components/canvas_editor.ui.yaml
    components: [Surface, MarkdownViewer, DiffView, CommentThread, ButtonRow]

  seo_settings:
    file: interactive_components/seo_settings.ui.yaml
    components: [Surface, FieldList, TextArea, Badge, ButtonRow]

  audit_dashboard:
    file: interactive_components/audit_dashboard.ui.yaml
    components: [Surface, SummaryCard, DataGrid, SeverityBadge, IssueList]

  source_evidence_table:
    file: interactive_components/source_evidence_table.ui.yaml
    components: [Surface, DataGrid, Link, Badge, FilterBar]

  quality_gate_checklist:
    file: interactive_components/quality_gate_checklist.ui.yaml
    components: [Card, Checklist, Badge, Text]
```

Pattern A is the default for production.

### 9.2 Pattern B: LLM-generated A2UI

Use LLM-generated A2UI only when:

```text
the output shape is not known in advance
the user asks for a custom exploratory presentation
the surface is view-only or low-risk
deterministic mapping would produce a poor UX
Layer 3 explicitly requests A2UI generation
Layer 2 outbound/protected-disclosure policy permits the payload
```

Required controls:

```text
Layer 5 provides catalog schema only, not internal implementation details.
The generated output must be pure A2UI JSON.
The generated output must validate against schema.
The output must use only approved catalog components.
All props must validate.
All URLs must pass URL policy.
Markdown must pass markdown safety policy.
Actions must be known, typed, and non-executable.
Mutation actions must be blocked unless the approval contract already exists.
Invalid A2UI falls back to deterministic text or a safe static card.
```

### 9.3 Pattern C: hybrid output

Hybrid output is preferred when a response must work across chat, canvas, and dashboard contexts.

```go
type HybridPresentation struct {
    SummaryText string     `json:"summary_text"`
    A2UIFrame   A2UIFrame  `json:"a2ui_frame"`
    DataRef     string     `json:"data_ref,omitempty"`
    Fallback    string     `json:"fallback"`
}
```

Use hybrid output for:

```text
audit summaries
source evidence
keyword opportunity groups
content brief proposals
article outline options
memory update proposals
high-risk approval summaries
external A2A task status
```

---

## 10. Component Catalog Policy

### 10.1 Catalog manifest

```yaml
catalog_id: seo_content_app
version: 1.0.0
schema_version: a2ui-frame.v1
status: active
owner_team: frontend_platform
default_catalog: false
allowed_surfaces:
  - chat
  - canvas
  - brief
  - seo_settings
  - audit_dashboard
  - approval_stream
components:
  - type: Card
    prop_schema: schemas/card.props.schema.json
    allowed_children: [Text, Badge, ButtonRow, DataGrid, DiffView, Checklist]
    max_children: 12

  - type: MarkdownViewer
    prop_schema: schemas/markdown_viewer.props.schema.json
    allowed_markdown_profile: safe_markdown_v1
    max_chars: 50000

  - type: DiffView
    prop_schema: schemas/diff_view.props.schema.json
    allowed_diff_types: [inline, side_by_side, field_level]
    max_diff_chars: 25000

  - type: DataGrid
    prop_schema: schemas/datagrid.props.schema.json
    max_rows: 500
    max_columns: 20
```

### 10.2 Basic catalog policy

The basic catalog is allowed for:

```text
development
test fixtures
low-risk prototypes
fallback cards
non-branded demos
```

Production surfaces should use the app catalog to preserve design-system consistency, accessibility, analytics hooks, and safety constraints.

### 10.3 Component admission rules

A component can be added only if:

```text
prop schema exists
accessibility requirements exist
snapshot tests exist
security review passes
no arbitrary code execution
no untrusted HTML rendering
no external script loading
no dynamic import based on model output
no unrestricted URL/image props
no telemetry calls from model-controlled props
no hidden data exfiltration channels
```

---

## 11. Surface-Specific Specs

### 11.1 Chat surface

Purpose:

```text
Short coordination messages, fallback text, status updates, and plain-language summaries.
```

Rules:

```text
1-2 sentences for normal coordination
no raw JSON unless user explicitly asks for export/debug and policy allows it
no exact internal tool names
no internal file paths
no workflow IDs
no trace IDs
no protected metadata
safe markdown only
```

Layer 5 responsibilities:

```text
render chat text
attach simple A2UI cards when needed
fallback when complex A2UI fails
route user text back to Layer 2 intake
```

Layer 5 must not classify the message.

### 11.2 Canvas surface

Purpose:

```text
Long-form content, audit reports, article drafts, outlines, and revision patches.
```

Canonical capabilities:

```text
render markdown document
render section status
render word count from upstream state
render inline or side-by-side diffs
render accept/reject patch cards
render comments and revision requests
render quality gate summary
```

Layer 5 must not:

```text
generate article text
decide article structure
write sections
persist canvas
overwrite canvas without a patch intent
silently accept a patch
```

Canvas patch contract:

```go
type CanvasPatchIntent struct {
    TraceID        string `json:"trace_id"`
    SurfaceVersion int64  `json:"surface_version"`
    PatchID        string `json:"patch_id"`
    PatchKind      string `json:"patch_kind"`
    TargetRange    Range  `json:"target_range,omitempty"`
    BeforeMarkdown string `json:"before_markdown,omitempty"`
    AfterMarkdown  string `json:"after_markdown"`
    Summary        string `json:"summary"`
    RequiresApproval bool `json:"requires_approval"`
}
```

Layer 5 renders and normalizes the user decision. Actual patch application belongs outside Layer 5 via Layer 3-controlled flow.

### 11.3 Brief surface

Purpose:

```text
Structured content contract for planning and writing workflows.
```

Fields:

```yaml
brief_fields:
  - topic
  - angle
  - audience
  - intent
  - hidden_intent
  - CTA
  - content_type
  - target_keywords
  - internal_link_targets
  - source_requirements
  - competitor_exclusions
  - brand_constraints
  - tone_notes
  - approval_status
```

Layer 5 may render edit controls. User edits become normalized events, not direct persistent mutations.

### 11.4 SEO settings surface

Purpose:

```text
Render and approve metadata, slug, canonical URL, Open Graph, FAQs, schema, internal links, noindex/nofollow, and crawler directives.
```

Rules:

```text
never overwrite existing SEO fields without explicit approval
distinguish suggested, accepted, rejected, and already-set fields
show field-level diffs
validate display length client-side and server-side
route final acceptance through Layer 2/3
```

Layer 5 must not generate metadata or schema. It renders candidates produced by the workflow.

### 11.5 Audit dashboard surface

Purpose:

```text
Render SEO/AEO site health, prioritized issues, evidence, impact, fix recommendations, and progress.
```

Required dashboard areas:

```text
health summary
priority issue list
technical SEO issues
content quality issues
AEO readiness
structured data status
Core Web Vitals summary
internal linking opportunities
analytics/search trends
recommended next actions
evidence links
status timeline
```

Layer 5 must render issue data received from Layer 3. It must not crawl, analyze, score, or prioritize independently.

### 11.6 Source/evidence surface

Purpose:

```text
Render supporting sources, credibility status, claim mapping, competitor exclusion notes, and unsupported claims.
```

Required fields:

```yaml
source_row:
  title: string
  url_display: string
  domain: string
  source_type: institutional | official | publication | competitor_excluded | user_owned | unknown
  credibility_status: accepted | caution | rejected
  used_for_claims: array
  notes: string
```

URL safety:

```text
display URL may be visible after sanitization
clickable URL must pass safe URL policy
competitor/excluded sources may be shown as excluded without creating outbound link
```

### 11.7 Approval stream

Purpose:

```text
Render user-reviewable cards for actions requiring explicit approval.
```

Approval card must include:

```text
plain-English action summary
why approval is needed
what will change
where it will change
before/after diff when applicable
risk level
expiration
approve and reject options
optional revision/comment path
MFA-required indicator when required
```

Layer 5 must not:

```text
decide approval need
authorize action
start MFA challenge
mint credentials
execute the approved action
```

### 11.8 A2A external surface

Purpose:

```text
Expose the agent as a safe external service to other agents through a public Agent Card and A2A application envelopes.
```

Layer 5 owns:

```text
public capability descriptions
input/output schemas
extension advertisement
task envelope validation
safe response formatting
A2A error envelopes
```

Layer 5 must not:

```text
own network transport
own mTLS
issue external identity
select internal workflow
bypass Layer 2 intake/policy
execute external requests directly
```

---

## 12. Approval UX and Vibe Diff

### 12.1 Approval request types

```yaml
approval_request_types:
  canvas_patch:
    surface: canvas
    requires_diff: true
    default_mfa: false

  brief_update:
    surface: brief
    requires_diff: true
    default_mfa: false

  seo_settings_update:
    surface: seo_settings
    requires_diff: true
    default_mfa: false

  memory_update_proposal:
    surface: approval_stream
    requires_diff: true
    default_mfa: false

  external_publish:
    surface: approval_stream
    requires_diff: true
    default_mfa: true

  high_risk_external_action:
    surface: approval_stream
    requires_diff: true
    default_mfa: true

  a2a_external_task_acceptance:
    surface: a2a_external
    requires_diff: false
    default_mfa: policy_dependent
```

### 12.2 Vibe Diff schema

```go
type VibeDiff struct {
    Summary        string      `json:"summary"`
    ChangeType     string      `json:"change_type"`
    Before         interface{} `json:"before,omitempty"`
    After          interface{} `json:"after,omitempty"`
    FieldDiffs     []FieldDiff `json:"field_diffs,omitempty"`
    RiskNotes      []string    `json:"risk_notes,omitempty"`
    UserVisibleIDs []string    `json:"user_visible_ids,omitempty"`
}
```

Rules:

```text
Summary must be plain English.
Diff must be scoped to the requested action.
Risk notes must be user-facing and non-internal.
No trace IDs.
No internal tool names.
No exact workflow IDs.
No policy file names.
No hidden system prompts.
No secrets.
```

### 12.3 Approval fatigue guard

Layer 5 must include UX-level fatigue protections:

```text
batch related low-risk approvals
avoid repeated cards for the same decision
collapse routine read-only confirmations
show one clear diff per state-changing action
avoid micro-approval for non-mutating view changes
use expiration to prevent stale approvals
support reject-with-comment and revise paths
```

Layer 5 does not override Layer 2/3 approval requirements. It only prevents the UI from creating unnecessary friction when the policy allows batching.

---

## 13. A2A Agent Card

### 13.1 Agent Card schema

```json
{
  "name": "AEOlyzer",
  "description": "Provides guarded website visibility, AEO audit, and content-planning capabilities through the A2A protocol.",
  "supportedInterfaces": [
    {
      "url": "https://api.example.com/a2a",
      "protocolBinding": "JSONRPC",
      "protocolVersion": "1.0"
    }
  ],
  "capabilities": {
    "streaming": false,
    "pushNotifications": false,
    "extendedAgentCard": false
  },
  "defaultInputModes": ["text/plain"],
  "defaultOutputModes": ["text/plain"],
  "skills": [
    {
      "id": "site_visibility_guidance",
      "name": "Site visibility guidance",
      "description": "Explains safe, public website visibility and content improvement options without exposing internal topology.",
      "tags": ["aeo", "seo", "content"]
    }
  ],
  "securitySchemes": {
    "googleOidc": {
      "openIdConnectSecurityScheme": {
        "openIdConnectUrl": "https://accounts.google.com/.well-known/openid-configuration"
      }
    }
  },
  "securityRequirements": [
    {
      "schemes": {
        "googleOidc": []
      }
    }
  ],
  "version": "1.0.0"
}
```

Required implementation rules:

```text
serve the public Agent Card at /.well-known/agent-card.json
build the card from official a2a.AgentCard types
validate the serialized card against the Layer 5 schema at startup
reject cards that disclose workflow IDs, tool inventories, MCP endpoints, traces, sandbox details, policy file names, or skill file paths
advertise only public product capabilities as A2A skills
use Google ADK Go adka2a/v2 to bridge ADK agents to A2A execution
do not use OpenAPI as the A2A contract
```

### 13.2 Agent Card rules

The Agent Card may disclose:

```text
public product capabilities
public input/output schema summaries
supported extension types
authentication requirements at a high level
rate-limit and payload-size limits
contact/support metadata if applicable
```

The Agent Card must not disclose:

```text
exact internal tool inventory
MCP server names or URLs
connector implementation details
workflow IDs
DAG definitions
capability profile internals
SKILL.md bodies
skill inventory beyond public capability categories
policy files
trace formats
runtime sandbox internals
model/provider names
secrets
```

### 13.3 A2A inbound flow

```text
01 Layer 7 receives an A2A JSON-RPC request through the ADK/A2A transport mount
02 Layer 7 authenticates transport metadata and creates the ADK/A2A request context
03 Google ADK adka2a/v2 converts the A2A Message or Task request into ADK invocation input
04 Layer 5 validates public Agent Card and public A2A shape, including disclosure policy
05 Layer 5 normalizes safe public A2A content into an external request envelope
06 Layer 2 performs policy, prompt-injection, and outbound-disclosure checks on natural-language content
07 Layer 3 receives only Layer 2-permitted requests that require workflow routing
08 Layer 5 formats safe task/message output through official A2A response types
09 Layer 7 writes the A2A JSON-RPC response without exposing internal transport or workflow details
10 Layer 8 receives sanitized A2A event facts only
```

Layer 5 must not skip Layer 2.

The signed `a2a-envelope` schema remains an AEOlyzer extension payload for trusted interop scenarios. It is not the base A2A protocol shape and must not replace canonical A2A `Message`, `Task`, `Part`, or Agent Card objects.

---

## 14. Security Policy

### 14.1 No executable UI

Reject any A2UI or presentation payload containing:

```text
<script>
javascript:
data:text/html
raw HTML
inline CSS
event handler attributes
function strings
eval
new Function
import()
iframe
object
embed
form action to untrusted URL
style injection
SVG with scriptable content
base64 HTML payloads
zero-width hidden instruction text
prompt-injection instructions
```

### 14.2 Markdown safety

Allowed markdown:

```text
headings
paragraphs
bold/italic
lists
tables
blockquotes
inline code
code blocks when marked safe and not executable instructions
links that pass URL policy
```

Disallowed markdown:

```text
raw HTML
script tags
iframe
style tags
javascript links
data HTML links
hidden comments containing instructions
model/provider secret leakage
trace IDs or protected metadata
```

### 14.3 URL safety

Layer 5 must normalize and validate URLs before rendering them as clickable links.

Rules:

```text
allow only http and https
reject javascript:, data:, file:, blob:, chrome:, vscode:, ssh:, ftp:
canonicalize punycode domains
strip tracking parameters when policy requires
label competitor-excluded links as non-clickable if upstream marks them excluded
display domain separately from full URL when safer
never render credential-bearing URLs
never render signed URLs unless upstream policy explicitly permits and expiry is shown
```

### 14.4 Protected metadata scanner

Layer 5 must scan outgoing A2UI and A2A payloads for:

```text
internal tool names
workflow IDs
DAG node IDs
policy file names
exact skill inventory
SKILL.md body excerpts
file paths
MCP endpoints
connector URLs
trace IDs when user-facing
AgBOM details
raw logs
hidden chain-of-thought
secrets
tokens
cookies
raw PII
```

If detected:

```text
block frame
emit ui_payload_blocked
return safe fallback text or ask Layer 2 for redacted response policy
```

### 14.5 Event replay and stale state

Every workflow-impacting UI event must include:

```text
frame_id
action_id
surface_version
interaction_id
idempotency_key or derived idempotency token
signature
expiry
```

Reject when:

```text
frame expired
surface version stale
action already consumed
signature invalid
interaction_id replayed
action not present in frame
values do not match value_schema_ref
decision target no longer pending
```

---

## 15. Integration Contracts with Other Layers

### 15.1 With Layer 1

Layer 1 owns the concrete frontend.

Layer 5 sends:

```text
A2UIFrame
fallback text
catalog version metadata
action bindings
surface version
approval card data
```

Layer 1 sends back:

```text
UserInteractionEvent
client-render error
accessibility/render diagnostics
hardware MFA result metadata if relevant and policy-approved
```

Layer 5 must not:

```text
own React/Next.js implementation
own native app code
directly manipulate DOM
execute browser automation
start hardware MFA challenge
```

### 15.2 With Layer 2

Layer 5 relies on Layer 2 for:

```text
protected-disclosure policy
outbound response guard
tool/action authorization
approval validity checks
PII/secret redaction
prompt-injection handling
policy state
```

Layer 5 must send mutating or protected decisions through Layer 2 before Layer 3 or Layer 6 can act.

Layer 5 must not duplicate:

```text
intent classifier
policy server
tool allowlist engine
LLM firewall
credential policy
```

### 15.3 With Layer 3

Layer 3 is the primary producer of presentation intents.

Layer 3 sends:

```text
PresentationIntent
SurfacePatchIntent
ApprovalRequestIntent
ExternalTaskPresentation
A2UI generation hints
fallback text
```

Layer 5 returns:

```text
NormalizedSurfaceEvent
ApprovalDecisionEvent
A2AInboundNormalizedRequest
A2ARenderedResponse
render-failure notices
```

Layer 5 must not:

```text
choose workflows
assign agents
construct DAGs
choose skills
create tool requests
decide next task node
```

### 15.4 With Layer 4

Layer 5 may use UI schemas and output contracts that align with Layer 4 assets, but it must not read skill files directly.

Allowed:

```text
accept output contract names from Layer 3
render output contracts as surfaces
use shared schema names
```

Forbidden:

```text
read skill-registry.yaml
load SKILL.md bodies
load skill assets directly
inspect skill scripts
decide which skill triggered
```

### 15.5 With Layer 6

Layer 5 does not execute anything.

Layer 5 may render:

```text
sandbox status summary
execution progress
approval card for a proposed high-risk action
safe artifact download panel
```

Layer 5 must not:

```text
spawn runtime
execute tool
execute script
mount filesystem
route egress
mint JIT token
revoke tool access
quarantine runtime
```

### 15.6 With Layer 7

Layer 5 may render public connection status or normalized source data. It must not own data access.

Layer 5 must not:

```text
connect to MCP
inspect MCP tool registry
call connectors
query vector store
read memory docs
write memory docs
manage mTLS
manage tenant partitioning
```

A2A transport security belongs to Layer 7. Layer 5 owns A2A application schema validation.

### 15.7 With Layer 8

Layer 5 emits sanitized events. Layer 8 stores and evaluates them.

Layer 5 sends:

```text
ui_frame_rendered
ui_frame_blocked
ui_action_received
ui_action_rejected
approval_card_rendered
approval_decision_received
a2a_envelope_received
a2a_envelope_rejected
a2a_response_sent
```

Layer 5 must not:

```text
store telemetry
query traces for user-facing display unless given a sanitized presentation intent
score evals
compute drift
trigger quarantine
run red/blue/green loops
```

---

## 16. Event Emission

### 16.1 UI event

```go
type Layer5Event struct {
    TraceID       string            `json:"trace_id,omitempty"`
    EventType     string            `json:"event_type"`
    Surface       string            `json:"surface,omitempty"`
    FrameID       string            `json:"frame_id,omitempty"`
    ActionKind    string            `json:"action_kind,omitempty"`
    Decision      string            `json:"decision,omitempty"`
    ReasonCode    string            `json:"reason_code,omitempty"`
    CatalogID     string            `json:"catalog_id,omitempty"`
    CatalogVersion string           `json:"catalog_version,omitempty"`
    Metadata      map[string]string `json:"metadata,omitempty"`
    CreatedAt     string            `json:"created_at"`
}
```

Allowed event types:

```text
presentation_intent_received
presentation_intent_rejected
a2ui_frame_built
a2ui_frame_blocked
a2ui_validation_failed
a2ui_fallback_used
ui_action_received
ui_action_rejected
ui_action_normalized
surface_patch_rendered
surface_patch_rejected
approval_card_rendered
approval_decision_received
approval_decision_rejected
a2a_agent_card_loaded
a2a_envelope_received
a2a_envelope_rejected
a2a_request_normalized
a2a_response_sent
protected_metadata_blocked
```

Events must not include:

```text
raw user prompt
raw secret
raw token
raw cookie
raw PII
raw hidden chain-of-thought
raw SKILL.md
raw tool payload
raw connector payload
raw telemetry trace
full canvas content unless explicitly redacted and necessary
```

---

## 17. Presentation Token Budget Policy

Layer 5 is not an LLM context owner, but A2UI payloads can still bloat sessions.

```yaml
presentation_budget:
  max_a2ui_frame_bytes: 262144
  max_nodes_per_frame: 250
  max_depth: 12
  max_text_chars_per_node: 12000
  max_markdown_chars_per_frame: 50000
  max_datagrid_rows_default: 100
  max_datagrid_rows_hard: 500
  max_actions_per_frame: 50
  max_approval_cards_pending: 20
```

Degradation behavior:

```text
summarize long tables
paginate rows
collapse secondary details
render source links as expandable details
prefer fallback text if payload too large
never drop risk warnings
never drop approval diff
never drop action labels
never silently truncate user-critical fields
```

---

## 18. Production Test Matrix

### 18.1 Schema tests

```text
PresentationIntent validates
A2UIFrame validates
A2UINode validates
UIActionBinding validates
UserInteractionEvent validates
ApprovalRequestIntent validates
VibeDiff validates
A2A Agent Card validates
A2A envelope validates
invalid enum rejected
oversized payload rejected
unknown surface rejected
unknown component rejected
unknown action rejected
```

### 18.2 Security tests

```text
raw HTML rejected
script tag rejected
javascript URL rejected
data HTML URL rejected
inline event handler rejected
CSS injection rejected
SVG script rejected
hidden zero-width payload rejected
prompt-injection text in props rejected or neutralized
raw internal tool name blocked
workflow ID blocked from user-visible props
MCP endpoint blocked
SKILL.md excerpt blocked
trace ID blocked where user-facing
token/cookie/secret blocked
credential-bearing URL blocked
```

### 18.3 A2UI tests

```text
deterministic converter produces valid frame
LLM-generated A2UI invalid component rejected
LLM-generated A2UI invalid props rejected
cyclic child graph rejected
orphan child ref rejected
duplicate node ID rejected
catalog version mismatch rejected
fallback used after validation failure
basic catalog allowed only in configured environments
app catalog required in production for configured surfaces
```

### 18.4 Surface tests

```text
chat fallback renders safely
canvas patch displays before/after diff
brief update renders field-level diff
SEO settings overwrite shows explicit approval card
audit dashboard paginates issue table
source evidence table blocks unsafe links
memory update proposal does not persist memory
quality checklist does not score evals
approval card includes expiry and risk summary
```

### 18.5 Interaction tests

```text
valid view-only sort handled locally
valid answer_question event normalized
valid approve event normalized
reject-with-comment normalized
unknown action rejected
stale surface version rejected
expired frame rejected
replayed interaction rejected
invalid signature rejected
values exceeding schema rejected
approval decision without pending request rejected
MFA-required approval without MFA metadata rejected before policy handoff
```

### 18.6 A2A tests

```text
Agent Card loads
Agent Card contains public capabilities only
Agent Card excludes internal tools
A2A envelope validates
unknown capability rejected
oversized request rejected
malformed extension rejected
external request normalized but not executed
natural-language A2A message sent to Layer 2
A2A response excludes internal workflow IDs
A2A transport details not exposed by Layer 5
```

### 18.7 Boundary tests

Layer 5 must prove it does not:

```text
normalize raw user intent
classify intent
choose workflow
assign agent
load skill
execute tool
execute script
connect MCP
query connector
query memory
write memory
mint token
enforce sandbox
store telemetry
score evals
detect drift
quarantine runtime
```

---

## 19. Deployment Requirements

Layer 5 v2 is production-ready only when:

```text
all schemas validate
all component catalogs validate
catalog-lock.yaml exists
all component mappings are tested
all production surfaces use app catalog
basic catalog disabled or restricted in production
all action bindings are registered
all mutating action bindings require Layer 2 policy path
all approval cards include Vibe Diff or explicit no-diff reason
all user interaction events require signature and idempotency
stale event rejection enabled
protected metadata scanner enabled
markdown sanitizer enabled
URL sanitizer enabled
A2UI fallback renderer enabled
A2A Agent Card public disclosure scan passes
A2A envelope validation passes
Layer 2 boundary tests pass
Layer 3 presentation intent tests pass
Layer 4 non-read tests pass
Layer 6 no-execution tests pass
Layer 7 no-transport-ownership tests pass
Layer 8 event redaction tests pass
accessibility tests pass
snapshot golden surface tests pass
```

### 19.1 Executable presentation-security baseline

The repository implementation must retain:

```text
all presentation, A2UI, UI-event, approval, surface-patch, Agent Card, and A2A envelope schemas compile at startup
catalog-lock.yaml validates against the A2UI catalog schema at startup
a2ui_schema_manager enforces unique nodes, valid references, one root, acyclic reachability, and depth bounds
a2ui_part_converter allows only catalogued components, properties, required properties, child rules, and non-executable content
hidden_payload_scanner rejects hidden keys, active HTML, event handlers, CSS hiding, instruction overrides, invisible controls, bidi controls, and encoded hidden payloads
markdown_sanitizer rejects raw HTML, comments, controls, and unsafe links
ui_payload_signer signs canonical frames with a bounded expiry
a2a_envelope_validator verifies schema, Ed25519 sender identity, timestamps, payload safety, and bounded nonce replay protection
URL sanitization rejects userinfo, localhost, private literals, and unsupported schemes
cmd/api compiles and validates embedded Layer 5 contracts before opening a listener
```

Required tests include graph cycles and unreachable nodes, unknown component
and property rejection, active-content and hidden-payload attacks, signature
tampering and expiry, A2A replay and sender mismatch, and unsafe URL cases.

### 19.2 Repository readiness evidence

The repository-level production check is:

```text
go run ./cmd/readiness -root .
```

For Layer 5, the check must fail when required presentation, A2UI, A2A,
approval, event, surface-patch, or catalog artifacts are missing, unreadable,
or placeholder-only. It must also report absent executable source controls
explicitly required by this spec and explicit prototype markers in Layer 5
production sources.

Source-file presence is not proof that the control works. Schema tests,
security tests, A2UI tests, interaction tests, A2A tests, accessibility tests,
snapshot tests, and boundary tests remain mandatory.

`cmd/readiness` and `internal/releasegate` are read-only platform CI tooling.
They must not render UI, generate presentation payloads, normalize live user
events, authorize actions, execute tools, or perform any Layer 5 runtime
behavior.

---

## 20. Acceptance Criteria

Layer 5 v2 is accepted when:

1. It converts structured presentation intents into safe A2UI frames.
2. It supports chat, canvas, brief, SEO settings, audit dashboard, source/evidence, approval stream, and A2A external surfaces.
3. It supports deterministic UI conversion for core surfaces.
4. It permits LLM-generated A2UI only under strict catalog and schema validation.
5. It uses hybrid output with fallback text.
6. It never emits executable UI code.
7. It never executes tools or scripts.
8. It never connects to MCP or external APIs.
9. It never classifies raw user intent.
10. It never chooses workflows or agents.
11. It never loads SKILL.md bodies.
12. It never persists memory or business state directly.
13. It renders approval and Vibe Diff cards without authorizing actions.
14. It normalizes user events with idempotency, signature, expiry, and surface-version checks.
15. It exposes a safe public A2A Agent Card without internal metadata.
16. It validates A2A envelopes but routes execution decisions to Layer 2/3.
17. It emits only sanitized events to Layer 8.
18. It blocks protected metadata, secrets, unsafe URLs, raw HTML, scripts, and hidden payloads.
19. It passes all zero-overlap boundary tests.
20. It can degrade gracefully to safe text when UI validation fails.

---

## 21. Final Non-Goals

```text
Layer 5 must not own raw frontend implementation.
Layer 5 must not manipulate DOM.
Layer 5 must not execute browser automation.
Layer 5 must not classify user intent.
Layer 5 must not run an LLM firewall.
Layer 5 must not choose workflows.
Layer 5 must not assign agents.
Layer 5 must not construct DAGs.
Layer 5 must not load skill registry.
Layer 5 must not load SKILL.md.
Layer 5 must not inspect skill scripts.
Layer 5 must not execute tools.
Layer 5 must not execute scripts.
Layer 5 must not call APIs.
Layer 5 must not connect MCP servers.
Layer 5 must not query RAG.
Layer 5 must not read memory documents.
Layer 5 must not write memory documents.
Layer 5 must not persist canvas content directly.
Layer 5 must not persist brief state directly.
Layer 5 must not generate final article content.
Layer 5 must not generate SEO metadata independently.
Layer 5 must not analyze websites independently.
Layer 5 must not crawl pages.
Layer 5 must not score audit issues independently.
Layer 5 must not authorize high-risk actions.
Layer 5 must not initiate MFA challenge.
Layer 5 must not mint or revoke credentials.
Layer 5 must not own network transport.
Layer 5 must not own mTLS.
Layer 5 must not enforce sandboxing.
Layer 5 must not store telemetry.
Layer 5 must not score evals.
Layer 5 must not track AgBOM.
Layer 5 must not compute intent drift.
Layer 5 must not run red/blue/green SecOps loops.
Layer 5 must not quarantine runtime.
Layer 5 must not expose exact internal tool inventory, workflow IDs, DAG details, SKILL.md bodies, skill file paths, MCP endpoints, connector internals, trace internals, policy file contents, source code, secrets, tokens, cookies, or raw PII to end users or external agents.
```

---

## 22. One-Line Architecture Summary

Layer 5 v2 is the production-grade presentation and extension boundary that transforms trusted workflow results into safe A2UI frames, approval cards, user interaction events, and A2A application envelopes, while leaving intake, routing, skills, execution, data access, persistence, telemetry, evaluation, and recovery to their owning layers.

---

## 23. Conversation Presentation Contracts Addendum

Layer 5 exposes presentation-safe conversation objects only:

```go
type ConversationSummary struct {
    ID        string    `json:"id"`
    Agent     ChatAgent `json:"agent"`
    Title     string    `json:"title"`
    Starred   bool      `json:"starred"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type ChatAttachment struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    ContentType string `json:"content_type"`
    Size        int64  `json:"size"`
}

type ChatMessage struct {
    ID          string           `json:"id"`
    Role        string           `json:"role"`
    Content     string           `json:"content"`
    Attachments []ChatAttachment `json:"attachments,omitempty"`
    CreatedAt   time.Time        `json:"created_at"`
}
```

Allowed roles are `user` and `assistant`. Internal ADK author names, invocation IDs, branch names, attachment hashes, storage references, tool calls, thought parts, signatures, and model metadata must not enter this contract.

The frontend must consume these contracts through a typed client, use `credentials: include`, keep Audit and Content history state isolated, and send files as multipart form data. Connecting these contracts must not introduce a second styling system or duplicate visual implementations for the two agents.

---

## 24. Knowledge Base and Content-Type Presentation Addendum

Layer 5 exposes one presentation-safe knowledge document at a time:

```go
type KnowledgeDocument struct {
    Section     KnowledgeSection       `json:"section"`
    Version     uint64                 `json:"version"`
    Profile     *KnowledgeProfile      `json:"profile,omitempty"`
    EEAT        *KnowledgeEEAT         `json:"eeat,omitempty"`
    Competitors *KnowledgeCompetitors  `json:"competitors,omitempty"`
    Topics      *KnowledgeTopics       `json:"topics,omitempty"`
    Tone        *KnowledgeTone         `json:"tone,omitempty"`
    Memory      *KnowledgeMemory       `json:"memory,omitempty"`
    UpdatedAt   *time.Time             `json:"updated_at,omitempty"`
}
```

Exactly one section payload is present. The frontend reads and writes through:

```text
GET /v1/knowledge/{section}
PUT /v1/knowledge/{section}
```

The `PUT` body contains the last-read `version`, `approved=true`, and exactly one matching section payload. Existing Save/Add/Remove controls are the explicit user decision surface. Layer 5 renders loading, empty, saving, conflict, and safe error states; it does not perform persistence or silently retry version conflicts.

`ConversationSummary` adds:

```go
ContentType ContentType `json:"content_type,omitempty"`
```

Allowed presentation values are `article`, `blog_post`, `linkedin_post`, `youtube_description`, and `product_description`. The Content composer sends the selected value in both conversation creation JSON and message multipart data. Audit requests omit it.

Attachment selection remains inside the existing Tailwind composer. Selected files render as bounded horizontal preview cards: browser-local image thumbnails for images, type/name/size cards for other files, and an accessible remove control. Preview rendering must stay inside the browser and must not send preview bytes through a remote image optimizer. The composer retains all existing actions and colors. Its resting container and textarea have no visible border; keyboard focus remains visible through the owning focus state.
