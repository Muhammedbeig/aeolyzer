# Layer 3 Production-Ready Specs v2
## Distributed Orchestration and Agent Fabric for SEO/AEO Auditor + Content Agent

**Status:** Production-ready upgrade  
**Supersedes:** Layer 3 production specs v1 where content collaboration was clarification-only  
**Primary rule:** Layer 3 plans and coordinates only. It never sanitizes, authorizes, executes, retrieves, renders, connects, observes, quarantines, silently updates memory, or persists long-term memory.

---

## 1. Upgrade Decision

Upgrade is needed.

The previous Layer 3 spec was strong for SEO/AEO auditing, site health, AI visibility, and clarification flows. It was not complete for the added content-agent workflow set:

- topic discovery
- content brief
- research
- SEO planning
- page analysis
- article planning
- drafting
- optimization
- repurposing
- memory and tone management

This v2 keeps the old rules and adds content-agent routing, mode gates, output-surface gates, and approval gates.

---

## 2. Layer Boundary

Layer 3 owns:

```text
validated-intent consumption
workflow blueprint selection
DAG construction and sequencing
capability profile selection
profile/workflow compatibility validation
short-lived orchestration state
context-window budgeting from sanitized context
proposed tool-call request creation
skill activation request creation
content-generation task contract creation
presentation intent creation
abstract data-need declaration
A2A/file-message-bus handoff envelopes
orchestration events
Plan/Write/Edit/Optimize/Audit mode gating
brief/canvas/chat surface routing by contract
memory-update approval gating
deep-research approval gating
```

Layer 3 does not own:

```text
raw input normalization
prompt-injection detection
protected-disclosure classification
intent classification
tool authorization
tool execution
skill loading
SKILL.md reading
skill script execution
memory document reading
memory writing
canvas rendering
brief rendering
chat rendering
MCP transport
connector execution
sandbox execution
credential minting
telemetry storage
eval scoring
quarantine
```

Ownership remains:

| Responsibility | Owner |
|---|---|
| Raw intake, safety, classification, tool authorization, redaction | Layer 2 |
| Skill registry, progressive disclosure, SKILL.md, references, assets, scripts | Layer 4 |
| Canvas, brief, chat, UI cards, dashboards, A2UI | Layer 5 |
| Runtime, filesystem, network egress, sandboxing, package controls, credentials | Layer 6 |
| MCP, APIs, connectors, RAG, retrieval | Layer 7 |
| Telemetry, evals, drift, AgBOM, SecOps, quarantine, audit storage | Layer 8 |

---

## 3. Required Directory Upgrade

Add the following to the v1 directory tree:

```text
/config/capability_profiles
  ├── content_collaborator.yaml
  └── content_execution_guard.yaml

/workflows
  ├── topic-discovery.bp.yaml
  ├── content-brief.bp.yaml
  ├── content-research.bp.yaml
  ├── seo-planning.bp.yaml
  ├── page-analysis.bp.yaml
  ├── article-planning.bp.yaml
  ├── article-drafting.bp.yaml
  ├── content-optimization.bp.yaml
  ├── content-repurposing.bp.yaml
  └── memory-tone-management.bp.yaml

/internal/orchestrator
  ├── content_workflow_test.go
  ├── mode_gate_test.go
  ├── memory_update_test.go
  ├── output_surface_test.go
  ├── /state
  │   ├── brief_state_guard.go
  │   └── canvas_state_guard.go
  ├── /handoff
  │   └── content_task_contracts.go
  └── /requests
      └── approval_request_builder.go
```

Do not add `/skills` to Layer 3. Skills remain Layer 4.

Do not add MCP, connector, or tool implementation files to Layer 3. They remain Layers 6 and 7.

---

## 4. Updated Orchestration Config

### 4.1 Runtime contract additions

```yaml
version: 2
policy_mode: fail_closed

runtime_contract:
  layer: 3
  owns_raw_intake: false
  owns_tool_authorization: false
  owns_tool_execution: false
  owns_skill_loading: false
  owns_mcp_transport: false
  owns_ui_rendering: false
  owns_canvas_rendering: false
  owns_brief_rendering: false
  owns_chat_rendering: false
  owns_memory_persistence: false
  owns_long_term_telemetry: false

mode_gates:
  default_mode: plan
  allowed_modes:
    - plan
    - write
    - edit
    - optimize
    - audit
  write_requires_layer2_mode_flag: true
  edit_requires_selected_text: true
  memory_write_requires_approval: true
  deep_research_requires_approval: true

output_surfaces:
  allowed:
    - canvas
    - brief
    - chat
    - report
    - table
    - dashboard
  canvas_writes_must_use_contract: true
  brief_updates_must_use_contract: true
  chat_must_be_coordination_only: true

state:
  store_raw_prompts: false
  store_raw_tool_payloads: false
  store_secrets: false
  store_raw_content: false
  store_memory_doc_body: false
  store_tone_doc_body: false
```

### 4.2 Profiles allowed

```yaml
profiles:
  default: seo_aeo_auditor
  allowed:
    - seo_aeo_auditor
    - content_collaborator
    - content_execution_guard
```

### 4.3 Workflows allowed

```yaml
workflows:
  allowed:
    - website-audit.bp
    - site-health.bp
    - ai-visibility.bp
    - content-strategy.bp
    - clarification.bp
    - topic-discovery.bp
    - content-brief.bp
    - content-research.bp
    - seo-planning.bp
    - page-analysis.bp
    - article-planning.bp
    - article-drafting.bp
    - content-optimization.bp
    - content-repurposing.bp
    - memory-tone-management.bp
```

### 4.4 New intent routes

```yaml
intent_routes:
  topic_discovery:
    profile: content_collaborator
    workflow: topic-discovery.bp
    required_context: []
    allowed_modes: [plan]
    enabled: true

  content_brief:
    profile: content_collaborator
    workflow: content-brief.bp
    required_context: []
    allowed_modes: [plan]
    enabled: true

  content_research:
    profile: content_collaborator
    workflow: content-research.bp
    required_context:
      - topic
    allowed_modes: [plan, write]
    enabled: true

  seo_planning:
    profile: content_collaborator
    workflow: seo-planning.bp
    required_context: []
    allowed_modes: [plan, optimize]
    enabled: true

  page_analysis:
    profile: content_collaborator
    workflow: page-analysis.bp
    required_context:
      - target_url
    allowed_modes: [plan, audit, optimize]
    enabled: true

  article_planning:
    profile: content_collaborator
    workflow: article-planning.bp
    required_context:
      - topic
    allowed_modes: [plan]
    enabled: true

  draft_article:
    profile: content_execution_guard
    workflow: article-drafting.bp
    required_context:
      - topic
      - audience
      - intent
      - target_word_count
    allowed_modes: [write]
    enabled: true

  optimize_content:
    profile: content_execution_guard
    workflow: content-optimization.bp
    required_context: []
    allowed_modes: [optimize, edit]
    enabled: true

  repurpose_content:
    profile: content_execution_guard
    workflow: content-repurposing.bp
    required_context:
      - content_type
    allowed_modes: [edit, write]
    enabled: true

  switch_content_type:
    profile: content_execution_guard
    workflow: content-repurposing.bp
    required_context:
      - content_type
    allowed_modes: [edit, write]
    enabled: true

  edit_existing:
    profile: content_execution_guard
    workflow: content-optimization.bp
    required_context:
      - selected_text
    allowed_modes: [edit]
    enabled: true

  memory_tone_management:
    profile: content_collaborator
    workflow: memory-tone-management.bp
    required_context: []
    allowed_modes: [plan, write, edit, optimize]
    enabled: true

  update_memory:
    profile: content_collaborator
    workflow: memory-tone-management.bp
    terminal_behavior: request_approval_then_delegate
    enabled: true
```

### 4.5 Route behavior rules

```text
unknown intent -> reject and emit orchestration_route_blocked
route.enabled=false -> reject and emit orchestration_route_disabled
protected_disclosure_request -> no workflow planned
out_of_bounds -> no workflow planned
fallback_clarification -> clarification.bp only; no tools; no skills; no remote agents
missing required context -> clarification.bp only; ask for missing fields
draft_article while mode != write -> blocked or converted to article-planning.bp
edit_existing without selected_text -> clarification.bp only
deepResearch without approval -> approval request only
memory update without approval -> approval request only
user-supplied workflow_id/agent_id/tool_call/skill_id/output_path/memory_write -> reject
```

---

## 5. Capability Profiles

### 5.1 `content_collaborator.yaml`

```yaml
profile_id: content_collaborator
profile_version: 2
mode: plan
role: content_strategy_planner

allowed_intents:
  - topic_discovery
  - content_brief
  - content_research
  - seo_planning
  - page_analysis
  - article_planning
  - memory_tone_management
  - update_memory
  - fallback_clarification
  - capability_explanation
  - documentation_lookup

workflow_allowlist:
  - topic-discovery.bp
  - content-brief.bp
  - content-research.bp
  - seo-planning.bp
  - page-analysis.bp
  - article-planning.bp
  - memory-tone-management.bp
  - clarification.bp

capability_tags:
  - content_strategy
  - topic_discovery
  - brief_building
  - content_research
  - source_intelligence
  - seo_planning
  - page_analysis
  - hidden_intent
  - audience_strategy
  - tone_adaptation
  - memory_update_proposal
  - safe_summary
  - clarification

public_capability_summary:
  - topic discovery
  - content brief planning
  - source-backed research planning
  - SEO planning
  - page analysis planning
  - article planning
  - tone and brand preference handling
  - clarification questions

protected_metadata_policy:
  expose_internal_tool_names: false
  expose_skill_inventory: false
  expose_workflow_ids_to_user: false
  expose_model_provider: false
  expose_mcp_servers: false
  expose_trace_internals: false
  expose_memory_file_paths: false
  expose_canvas_internal_state: false
```

### 5.2 `content_execution_guard.yaml`

```yaml
profile_id: content_execution_guard
profile_version: 1
mode: write_or_edit_guarded
role: guarded_content_execution_coordinator

allowed_intents:
  - draft_article
  - optimize_content
  - repurpose_content
  - switch_content_type
  - edit_existing

workflow_allowlist:
  - article-drafting.bp
  - content-optimization.bp
  - content-repurposing.bp
  - clarification.bp

capability_tags:
  - guarded_drafting
  - section_by_section_writing
  - targeted_editing
  - content_repurposing
  - seo_support_outputs
  - post_write_quality
  - output_surface_routing
  - no_silent_memory_write

public_capability_summary:
  - guarded article drafting
  - section-by-section writing coordination
  - targeted editing coordination
  - content repurposing coordination
  - SEO support output planning
  - post-write quality coordination

protected_metadata_policy:
  expose_internal_tool_names: false
  expose_skill_inventory: false
  expose_workflow_ids_to_user: false
  expose_model_provider: false
  expose_mcp_servers: false
  expose_trace_internals: false
  expose_memory_file_paths: false
  expose_canvas_internal_state: false
```

### 5.3 Profile validation rules

Reject a profile if:

```text
profile_id missing
profile_version missing
allowed_intents empty
workflow_allowlist contains unknown workflows
capability_tags empty
public_capability_summary contains internal tool names
protected_metadata_policy permits internal disclosure
profile contains SKILL.md body text
profile contains skill file names as hard dependencies
profile contains tool implementation code
profile contains MCP server names or URLs
profile contains credentials or secret-like values
profile contains model/provider names
profile contains filesystem implementation paths
profile allows draft_article without write-mode guard
profile allows update_memory without approval guard
```

---

## 6. Workflow Blueprint Schema Additions

Every content workflow blueprint must use this schema extension:

```yaml
workflow_id: string
workflow_version: integer
owner_layer: layer_3_orchestration
default_profile: string
allowed_intents: [string]
required_context: [string]
allowed_modes: [string]

nodes:
  - node_id: string
    type: orchestration | proposed_tool_request | skill_activation_request | content_generation_task | presentation_intent | data_need | approval_request | terminal
    purpose: string
    required_capabilities: [string]
    requires_layer2_tool_authorization: boolean
    requires_user_approval: boolean
    optional: boolean
    timeout_seconds: integer
    output_contract: string
    surface_hint: string

edges:
  - from: string
    to: string
    condition: string

failure_policy:
  missing_context: string
  mode_not_allowed: string
  approval_missing: string
  layer2_denied_tool: string
  layer4_skill_unavailable: string
  layer5_surface_unavailable: string
  layer6_runtime_failed: string
  layer7_data_unavailable: string
  protected_disclosure: string
```

Validation must reject workflow blueprints containing:

```text
raw prompts
hidden prompt text
credentials
secrets
bearer tokens
cookies
API URLs
MCP server names
model/provider names
tool implementation code
shell snippets
direct network instructions
SKILL.md bodies
memory document bodies
canvas document bodies
A2UI component JSON
```

---

## 7. New Workflow Blueprints

### 7.1 `topic-discovery.bp.yaml`

```yaml
workflow_id: topic-discovery.bp
workflow_version: 1
owner_layer: layer_3_orchestration
default_profile: content_collaborator
allowed_intents: [topic_discovery]
required_context: []
allowed_modes: [plan]

nodes:
  - node_id: validate_topic_scope
    type: orchestration
    purpose: confirm available brand, domain, audience, and seed topic context
    required_capabilities: [topic_discovery]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 5
    output_contract: topic_scope_validated
    surface_hint: chat

  - node_id: collect_brand_context
    type: proposed_tool_request
    purpose: propose brand, positioning, and competitor context collection
    required_capabilities: [content_strategy]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 30
    output_contract: brand_context_summary
    surface_hint: chat

  - node_id: collect_source_intelligence
    type: proposed_tool_request
    purpose: propose competitor blocklist, authority source, and format trend signal collection
    required_capabilities: [source_intelligence]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 30
    output_contract: source_intelligence_summary
    surface_hint: chat

  - node_id: collect_topic_landscape
    type: proposed_tool_request
    purpose: propose competitor coverage, industry news, audience questions, and site search
    required_capabilities: [topic_discovery, content_research]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: false
    timeout_seconds: 60
    output_contract: topic_landscape_summary
    surface_hint: chat

  - node_id: inspect_relevant_pages
    type: proposed_tool_request
    purpose: propose reading specific first-party or competitor pages only when a result materially affects topic choice
    required_capabilities: [page_analysis]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 60
    output_contract: inspected_page_summary
    surface_hint: chat

  - node_id: check_existing_overlap
    type: proposed_tool_request
    purpose: propose site page discovery and cannibalization check if topic is close to existing content
    required_capabilities: [seo_planning]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 60
    output_contract: topic_overlap_summary
    surface_hint: chat

  - node_id: request_topic_decision
    type: approval_request
    purpose: ask for topic, audience, angle, or length only after research narrows options
    required_capabilities: [clarification]
    requires_layer2_tool_authorization: false
    requires_user_approval: true
    optional: true
    timeout_seconds: 10
    output_contract: topic_decision_request
    surface_hint: chat

  - node_id: save_topic_angle
    type: proposed_tool_request
    purpose: propose saving chosen topic and angle to brief
    required_capabilities: [brief_building]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: false
    timeout_seconds: 30
    output_contract: brief_topic_saved
    surface_hint: brief

  - node_id: present_topic_options
    type: presentation_intent
    purpose: hand off topic options and evidence summary to Layer 5
    required_capabilities: [topic_discovery]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 10
    output_contract: presentation_intent
    surface_hint: chat

edges:
  - from: validate_topic_scope
    to: collect_brand_context
    condition: always
  - from: collect_brand_context
    to: collect_source_intelligence
    condition: always
  - from: collect_source_intelligence
    to: collect_topic_landscape
    condition: always
  - from: collect_topic_landscape
    to: inspect_relevant_pages
    condition: if_relevant_pages_found
  - from: inspect_relevant_pages
    to: check_existing_overlap
    condition: always
  - from: check_existing_overlap
    to: request_topic_decision
    condition: if_user_decision_required
  - from: request_topic_decision
    to: save_topic_angle
    condition: if_user_selected_topic
  - from: save_topic_angle
    to: present_topic_options
    condition: always

failure_policy:
  missing_context: continue_with_available_context_or_ask
  mode_not_allowed: stop_and_request_plan_mode
  approval_missing: wait_for_user_decision
  layer2_denied_tool: skip_node_and_record_blocked_dependency
  layer4_skill_unavailable: continue_with_reduced_capability_if_optional
  layer5_surface_unavailable: stop_and_return_safe_response
  layer6_runtime_failed: mark_node_failed_and_continue_if_optional
  layer7_data_unavailable: continue_with_data_gap_if_optional
  protected_disclosure: stop_and_return_safe_response
```

### 7.2 `content-brief.bp.yaml`

```yaml
workflow_id: content-brief.bp
workflow_version: 1
owner_layer: layer_3_orchestration
default_profile: content_collaborator
allowed_intents: [content_brief]
required_context: []
allowed_modes: [plan]

nodes:
  - node_id: collect_brief_defaults
    type: proposed_tool_request
    purpose: propose brand context and stored tone preference lookup through approved contracts
    required_capabilities: [brief_building, tone_adaptation]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 30
    output_contract: brief_defaults_summary
    surface_hint: brief

  - node_id: request_missing_brief_inputs
    type: approval_request
    purpose: ask for missing topic, angle, audience, intent, CTA, length, or notes in a batched decision request
    required_capabilities: [clarification]
    requires_layer2_tool_authorization: false
    requires_user_approval: true
    optional: false
    timeout_seconds: 10
    output_contract: brief_input_request
    surface_hint: chat

  - node_id: save_brief
    type: proposed_tool_request
    purpose: propose saving topic, angle, audience, intent, CTA, and notes to brief
    required_capabilities: [brief_building]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: false
    timeout_seconds: 30
    output_contract: brief_saved
    surface_hint: brief

  - node_id: propose_reusable_preference
    type: approval_request
    purpose: request approval before proposing reusable brand or tone preference update
    required_capabilities: [memory_update_proposal]
    requires_layer2_tool_authorization: false
    requires_user_approval: true
    optional: true
    timeout_seconds: 10
    output_contract: memory_update_approval_request
    surface_hint: chat

  - node_id: present_brief_status
    type: presentation_intent
    purpose: hand off brief summary and next available actions to Layer 5
    required_capabilities: [brief_building]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 10
    output_contract: presentation_intent
    surface_hint: brief

edges:
  - from: collect_brief_defaults
    to: request_missing_brief_inputs
    condition: if_brief_incomplete
  - from: request_missing_brief_inputs
    to: save_brief
    condition: if_inputs_available
  - from: save_brief
    to: propose_reusable_preference
    condition: if_reusable_preference_detected
  - from: propose_reusable_preference
    to: present_brief_status
    condition: always

failure_policy:
  missing_context: ask_only_for_required_missing_fields
  mode_not_allowed: stop_and_request_plan_mode
  approval_missing: wait_for_user_decision
  layer2_denied_tool: terminal_safe_response
  layer4_skill_unavailable: terminal_safe_response
  layer5_surface_unavailable: stop_and_return_safe_response
  layer6_runtime_failed: terminal_safe_response
  layer7_data_unavailable: terminal_safe_response
  protected_disclosure: stop_and_return_safe_response
```

### 7.3 `content-research.bp.yaml`

```yaml
workflow_id: content-research.bp
workflow_version: 1
owner_layer: layer_3_orchestration
default_profile: content_collaborator
allowed_intents: [content_research]
required_context: [topic]
allowed_modes: [plan, write]

nodes:
  - node_id: validate_research_scope
    type: orchestration
    purpose: confirm topic, evidence depth, source restrictions, and competitor blocklist availability
    required_capabilities: [content_research]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 5
    output_contract: research_scope_validated
    surface_hint: chat

  - node_id: collect_source_intelligence
    type: proposed_tool_request
    purpose: propose authority-source and competitor-blocklist collection before research
    required_capabilities: [source_intelligence]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 30
    output_contract: source_intelligence_summary
    surface_hint: chat

  - node_id: run_fast_research
    type: proposed_tool_request
    purpose: propose fast web research for stats, quotes, current developments, competitor coverage, and search intent
    required_capabilities: [content_research]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: false
    timeout_seconds: 60
    output_contract: fast_research_summary
    surface_hint: chat

  - node_id: request_deep_research_approval
    type: approval_request
    purpose: ask for user approval before slower evidence-heavy research
    required_capabilities: [content_research]
    requires_layer2_tool_authorization: false
    requires_user_approval: true
    optional: true
    timeout_seconds: 10
    output_contract: deep_research_approval_request
    surface_hint: chat

  - node_id: run_deep_research
    type: proposed_tool_request
    purpose: propose deeper research only after explicit approval
    required_capabilities: [content_research]
    requires_layer2_tool_authorization: true
    requires_user_approval: true
    optional: true
    timeout_seconds: 180
    output_contract: deep_research_summary
    surface_hint: chat

  - node_id: scrape_important_sources
    type: proposed_tool_request
    purpose: propose clean extraction from important URLs only when required by evidence plan
    required_capabilities: [page_analysis]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 60
    output_contract: extracted_source_summary
    surface_hint: chat

  - node_id: present_research_summary
    type: presentation_intent
    purpose: hand off safe research summary and usable evidence set to Layer 5
    required_capabilities: [content_research]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 10
    output_contract: presentation_intent
    surface_hint: chat

edges:
  - from: validate_research_scope
    to: collect_source_intelligence
    condition: always
  - from: collect_source_intelligence
    to: run_fast_research
    condition: always
  - from: run_fast_research
    to: request_deep_research_approval
    condition: if_evidence_heavy_piece
  - from: request_deep_research_approval
    to: run_deep_research
    condition: if_approved
  - from: run_deep_research
    to: scrape_important_sources
    condition: if_important_urls_found
  - from: scrape_important_sources
    to: present_research_summary
    condition: always

failure_policy:
  missing_context: stop_and_request_clarification
  mode_not_allowed: stop_and_request_supported_mode
  approval_missing: skip_deep_research
  layer2_denied_tool: skip_node_and_record_blocked_dependency
  layer4_skill_unavailable: continue_with_reduced_capability_if_optional
  layer5_surface_unavailable: stop_and_return_safe_response
  layer6_runtime_failed: mark_node_failed_and_continue_if_optional
  layer7_data_unavailable: continue_with_data_gap_if_optional
  protected_disclosure: stop_and_return_safe_response
```

### 7.4 `seo-planning.bp.yaml`

```yaml
workflow_id: seo-planning.bp
workflow_version: 1
owner_layer: layer_3_orchestration
default_profile: content_collaborator
allowed_intents: [seo_planning]
required_context: []
allowed_modes: [plan, optimize]

nodes:
  - node_id: validate_seo_scope
    type: orchestration
    purpose: confirm topic, target domain, URL, or existing article context from sanitized state
    required_capabilities: [seo_planning]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 5
    output_contract: seo_scope_validated
    surface_hint: brief

  - node_id: collect_serp_landscape
    type: proposed_tool_request
    purpose: propose search intent and SERP landscape checks
    required_capabilities: [seo_planning, content_research]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 60
    output_contract: serp_landscape_summary
    surface_hint: brief

  - node_id: collect_internal_link_candidates
    type: proposed_tool_request
    purpose: propose first-party page discovery for internal linking opportunities
    required_capabilities: [seo_planning]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 60
    output_contract: internal_link_candidates_summary
    surface_hint: brief

  - node_id: verify_link_candidate_pages
    type: proposed_tool_request
    purpose: propose page reads for candidate internal links before recommending them
    required_capabilities: [page_analysis]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 60
    output_contract: verified_internal_links_summary
    surface_hint: brief

  - node_id: check_cannibalization
    type: proposed_tool_request
    purpose: propose overlap and cannibalization risk analysis
    required_capabilities: [seo_planning]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 60
    output_contract: cannibalization_summary
    surface_hint: brief

  - node_id: save_keyword_direction
    type: proposed_tool_request
    purpose: propose storing keyword, intent, and SEO direction in brief notes
    required_capabilities: [brief_building]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 30
    output_contract: seo_direction_saved
    surface_hint: brief

  - node_id: present_seo_plan
    type: presentation_intent
    purpose: hand off SEO plan, search intent, internal link opportunities, and overlap risks to Layer 5
    required_capabilities: [seo_planning]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 10
    output_contract: presentation_intent
    surface_hint: brief

edges:
  - from: validate_seo_scope
    to: collect_serp_landscape
    condition: always
  - from: collect_serp_landscape
    to: collect_internal_link_candidates
    condition: always
  - from: collect_internal_link_candidates
    to: verify_link_candidate_pages
    condition: if_candidates_found
  - from: verify_link_candidate_pages
    to: check_cannibalization
    condition: always
  - from: check_cannibalization
    to: save_keyword_direction
    condition: if_direction_available
  - from: save_keyword_direction
    to: present_seo_plan
    condition: always

failure_policy:
  missing_context: ask_only_for_required_missing_fields
  mode_not_allowed: stop_and_request_supported_mode
  approval_missing: terminal_safe_response
  layer2_denied_tool: skip_node_and_record_blocked_dependency
  layer4_skill_unavailable: continue_with_reduced_capability_if_optional
  layer5_surface_unavailable: stop_and_return_safe_response
  layer6_runtime_failed: mark_node_failed_and_continue_if_optional
  layer7_data_unavailable: continue_with_data_gap_if_optional
  protected_disclosure: stop_and_return_safe_response
```

### 7.5 `page-analysis.bp.yaml`

```yaml
workflow_id: page-analysis.bp
workflow_version: 1
owner_layer: layer_3_orchestration
default_profile: content_collaborator
allowed_intents: [page_analysis]
required_context: [target_url]
allowed_modes: [plan, audit, optimize]

nodes:
  - node_id: validate_page_scope
    type: orchestration
    purpose: confirm URL, brand context need, and comparison purpose
    required_capabilities: [page_analysis]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 5
    output_contract: page_scope_validated
    surface_hint: report

  - node_id: collect_primary_page
    type: proposed_tool_request
    purpose: propose clean page extraction as primary page analysis input
    required_capabilities: [page_analysis]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: false
    timeout_seconds: 60
    output_contract: page_extract_summary
    surface_hint: report

  - node_id: fallback_search_if_needed
    type: proposed_tool_request
    purpose: propose backup search only when scraping is thin or unavailable
    required_capabilities: [content_research]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 60
    output_contract: fallback_page_context
    surface_hint: report

  - node_id: collect_brand_frame
    type: proposed_tool_request
    purpose: propose brand-positioning context collection to frame the page
    required_capabilities: [content_strategy]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 30
    output_contract: brand_frame_summary
    surface_hint: report

  - node_id: check_page_overlap
    type: proposed_tool_request
    purpose: propose cannibalization or overlap checks when page relates to planned content
    required_capabilities: [seo_planning]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 60
    output_contract: page_overlap_summary
    surface_hint: report

  - node_id: present_page_analysis
    type: presentation_intent
    purpose: hand off page findings, evidence, and recommended next steps to Layer 5
    required_capabilities: [page_analysis]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 10
    output_contract: presentation_intent
    surface_hint: report

edges:
  - from: validate_page_scope
    to: collect_primary_page
    condition: always
  - from: collect_primary_page
    to: fallback_search_if_needed
    condition: if_scrape_thin_or_failed
  - from: fallback_search_if_needed
    to: collect_brand_frame
    condition: always
  - from: collect_brand_frame
    to: check_page_overlap
    condition: if_overlap_check_needed
  - from: check_page_overlap
    to: present_page_analysis
    condition: always

failure_policy:
  missing_context: stop_and_request_clarification
  mode_not_allowed: stop_and_request_supported_mode
  approval_missing: terminal_safe_response
  layer2_denied_tool: skip_node_and_record_blocked_dependency
  layer4_skill_unavailable: continue_with_reduced_capability_if_optional
  layer5_surface_unavailable: stop_and_return_safe_response
  layer6_runtime_failed: mark_node_failed_and_continue_if_optional
  layer7_data_unavailable: continue_with_data_gap_if_optional
  protected_disclosure: stop_and_return_safe_response
```

### 7.6 `article-planning.bp.yaml`

```yaml
workflow_id: article-planning.bp
workflow_version: 1
owner_layer: layer_3_orchestration
default_profile: content_collaborator
allowed_intents: [article_planning]
required_context: [topic]
allowed_modes: [plan]

nodes:
  - node_id: validate_article_plan_scope
    type: orchestration
    purpose: confirm topic, audience, angle, length, CTA, screenshot needs, and target surface
    required_capabilities: [content_strategy]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 5
    output_contract: article_plan_scope
    surface_hint: canvas

  - node_id: collect_strategy_context
    type: proposed_tool_request
    purpose: propose brand context and source intelligence collection before planning
    required_capabilities: [content_strategy, source_intelligence]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 30
    output_contract: strategy_context_summary
    surface_hint: canvas

  - node_id: collect_plan_evidence
    type: proposed_tool_request
    purpose: propose research evidence before planning
    required_capabilities: [content_research]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: false
    timeout_seconds: 60
    output_contract: planning_evidence_summary
    surface_hint: canvas

  - node_id: request_missing_plan_inputs
    type: approval_request
    purpose: ask for audience, angle, length, CTA, screenshots, or format choices when missing
    required_capabilities: [clarification]
    requires_layer2_tool_authorization: false
    requires_user_approval: true
    optional: true
    timeout_seconds: 10
    output_contract: article_plan_input_request
    surface_hint: chat

  - node_id: save_planning_brief
    type: proposed_tool_request
    purpose: propose saving finalized brief before plan finalization
    required_capabilities: [brief_building]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: false
    timeout_seconds: 30
    output_contract: planning_brief_saved
    surface_hint: brief

  - node_id: build_article_plan_contract
    type: content_generation_task
    purpose: create section-by-section blueprint task contract; Layer 3 does not write plan prose directly
    required_capabilities: [content_strategy, hidden_intent, seo_planning]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 60
    output_contract: article_plan_contract
    surface_hint: canvas

  - node_id: present_article_plan
    type: presentation_intent
    purpose: hand off article plan contract to Layer 5 for canvas or brief display
    required_capabilities: [content_strategy]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 10
    output_contract: presentation_intent
    surface_hint: canvas

edges:
  - from: validate_article_plan_scope
    to: collect_strategy_context
    condition: always
  - from: collect_strategy_context
    to: collect_plan_evidence
    condition: always
  - from: collect_plan_evidence
    to: request_missing_plan_inputs
    condition: if_inputs_missing
  - from: request_missing_plan_inputs
    to: save_planning_brief
    condition: if_inputs_available
  - from: save_planning_brief
    to: build_article_plan_contract
    condition: always
  - from: build_article_plan_contract
    to: present_article_plan
    condition: always

failure_policy:
  missing_context: stop_and_request_clarification
  mode_not_allowed: stop_and_request_plan_mode
  approval_missing: wait_for_user_decision
  layer2_denied_tool: skip_node_and_record_blocked_dependency
  layer4_skill_unavailable: continue_with_reduced_capability_if_optional
  layer5_surface_unavailable: stop_and_return_safe_response
  layer6_runtime_failed: mark_node_failed_and_continue_if_optional
  layer7_data_unavailable: continue_with_data_gap_if_optional
  protected_disclosure: stop_and_return_safe_response
```

### 7.7 `article-drafting.bp.yaml`

```yaml
workflow_id: article-drafting.bp
workflow_version: 1
owner_layer: layer_3_orchestration
default_profile: content_execution_guard
allowed_intents: [draft_article]
required_context: [topic, audience, intent, target_word_count]
allowed_modes: [write]

nodes:
  - node_id: enforce_write_mode
    type: orchestration
    purpose: fail closed unless Layer 2 supplied authorized write-mode flag
    required_capabilities: [guarded_drafting]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 5
    output_contract: write_mode_validated
    surface_hint: canvas

  - node_id: collect_tone_rules
    type: proposed_tool_request
    purpose: propose reading stored tone and brand rules through approved contracts
    required_capabilities: [tone_adaptation]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 30
    output_contract: tone_rule_summary
    surface_hint: canvas

  - node_id: refresh_positioning_context
    type: proposed_tool_request
    purpose: propose brand positioning refresh if needed for drafting
    required_capabilities: [content_strategy]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 30
    output_contract: positioning_summary
    surface_hint: canvas

  - node_id: collect_supporting_page_context
    type: proposed_tool_request
    purpose: propose page extraction only when draft must reflect or quote a page
    required_capabilities: [page_analysis]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 60
    output_contract: supporting_page_summary
    surface_hint: canvas

  - node_id: collect_internal_link_targets
    type: proposed_tool_request
    purpose: propose first-party page discovery for internal linking targets
    required_capabilities: [seo_planning]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 60
    output_contract: internal_link_target_summary
    surface_hint: canvas

  - node_id: build_section_generation_contract
    type: content_generation_task
    purpose: create section-by-section drafting contract with word-count tracking; Layer 3 does not write section prose
    required_capabilities: [section_by_section_writing, guarded_drafting]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 60
    output_contract: section_generation_contract
    surface_hint: canvas

  - node_id: propose_canvas_write
    type: proposed_tool_request
    purpose: propose writing approved section content to canvas through canvas write contract
    required_capabilities: [output_surface_routing]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: false
    timeout_seconds: 60
    output_contract: canvas_write_result
    surface_hint: canvas

  - node_id: request_post_write_quality
    type: skill_activation_request
    purpose: request post-write quality checks without loading skill bodies in Layer 3
    required_capabilities: [post_write_quality]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 30
    output_contract: post_write_quality_request
    surface_hint: canvas

  - node_id: present_draft_status
    type: presentation_intent
    purpose: hand off drafting status, word-count summary, quality-check status, and optional SEO next step to Layer 5
    required_capabilities: [guarded_drafting]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 10
    output_contract: presentation_intent
    surface_hint: canvas

edges:
  - from: enforce_write_mode
    to: collect_tone_rules
    condition: always
  - from: collect_tone_rules
    to: refresh_positioning_context
    condition: always
  - from: refresh_positioning_context
    to: collect_supporting_page_context
    condition: if_page_context_required
  - from: collect_supporting_page_context
    to: collect_internal_link_targets
    condition: always
  - from: collect_internal_link_targets
    to: build_section_generation_contract
    condition: always
  - from: build_section_generation_contract
    to: propose_canvas_write
    condition: if_content_contract_ready
  - from: propose_canvas_write
    to: request_post_write_quality
    condition: after_each_section_or_final_section
  - from: request_post_write_quality
    to: present_draft_status
    condition: when_target_word_count_reached_or_user_stops

failure_policy:
  missing_context: stop_and_request_clarification
  mode_not_allowed: stop_and_request_write_mode
  approval_missing: wait_for_user_decision
  layer2_denied_tool: block_required_canvas_write_or_skip_optional_dependency
  layer4_skill_unavailable: block_if_required
  layer5_surface_unavailable: stop_and_return_safe_response
  layer6_runtime_failed: mark_node_failed_and_stop_if_canvas_write_failed
  layer7_data_unavailable: continue_with_data_gap_if_optional
  protected_disclosure: stop_and_return_safe_response
```

### 7.8 `content-optimization.bp.yaml`

```yaml
workflow_id: content-optimization.bp
workflow_version: 1
owner_layer: layer_3_orchestration
default_profile: content_execution_guard
allowed_intents: [optimize_content, edit_existing]
required_context: []
allowed_modes: [optimize, edit]

nodes:
  - node_id: validate_optimization_scope
    type: orchestration
    purpose: confirm live page, canvas content, selected text, SEO settings, or optimization goal
    required_capabilities: [seo_support_outputs, targeted_editing]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 5
    output_contract: optimization_scope
    surface_hint: canvas

  - node_id: analyze_current_page
    type: proposed_tool_request
    purpose: propose live page analysis when URL is present
    required_capabilities: [page_analysis]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 60
    output_contract: current_page_summary
    surface_hint: canvas

  - node_id: compare_current_serp
    type: proposed_tool_request
    purpose: propose current SERP or newer-source comparison when optimization needs external context
    required_capabilities: [seo_planning, content_research]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 60
    output_contract: serp_comparison_summary
    surface_hint: canvas

  - node_id: find_internal_link_improvements
    type: proposed_tool_request
    purpose: propose internal link candidates and overlap reduction checks
    required_capabilities: [seo_planning]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 60
    output_contract: optimization_link_summary
    surface_hint: canvas

  - node_id: build_targeted_edit_contract
    type: content_generation_task
    purpose: create targeted edit or SEO-output contract; Layer 3 does not rewrite content directly
    required_capabilities: [targeted_editing, seo_support_outputs]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 60
    output_contract: optimization_edit_contract
    surface_hint: canvas

  - node_id: propose_canvas_edit
    type: proposed_tool_request
    purpose: propose exact-match surgical edit or SEO-field update through approved contract; never overwrite existing SEO values unless authorized
    required_capabilities: [targeted_editing, output_surface_routing]
    requires_layer2_tool_authorization: true
    requires_user_approval: true
    optional: false
    timeout_seconds: 60
    output_contract: canvas_edit_result
    surface_hint: canvas

  - node_id: request_style_memory_approval
    type: approval_request
    purpose: request approval before proposing recurring tone or style memory updates
    required_capabilities: [memory_update_proposal]
    requires_layer2_tool_authorization: false
    requires_user_approval: true
    optional: true
    timeout_seconds: 10
    output_contract: memory_update_approval_request
    surface_hint: chat

  - node_id: present_optimization_status
    type: presentation_intent
    purpose: hand off optimization plan, edit status, SEO outputs, and unresolved risks to Layer 5
    required_capabilities: [seo_support_outputs]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 10
    output_contract: presentation_intent
    surface_hint: canvas

edges:
  - from: validate_optimization_scope
    to: analyze_current_page
    condition: if_url_available
  - from: analyze_current_page
    to: compare_current_serp
    condition: if_external_context_needed
  - from: compare_current_serp
    to: find_internal_link_improvements
    condition: always
  - from: find_internal_link_improvements
    to: build_targeted_edit_contract
    condition: always
  - from: build_targeted_edit_contract
    to: propose_canvas_edit
    condition: if_exact_edit_available
  - from: propose_canvas_edit
    to: request_style_memory_approval
    condition: if_recurring_style_rule_detected
  - from: request_style_memory_approval
    to: present_optimization_status
    condition: always

failure_policy:
  missing_context: ask_only_for_required_missing_fields
  mode_not_allowed: stop_and_request_supported_mode
  approval_missing: wait_for_user_decision
  layer2_denied_tool: skip_node_and_record_blocked_dependency
  layer4_skill_unavailable: continue_with_reduced_capability_if_optional
  layer5_surface_unavailable: stop_and_return_safe_response
  layer6_runtime_failed: mark_node_failed_and_stop_if_edit_failed
  layer7_data_unavailable: continue_with_data_gap_if_optional
  protected_disclosure: stop_and_return_safe_response
```

### 7.9 `content-repurposing.bp.yaml`

```yaml
workflow_id: content-repurposing.bp
workflow_version: 1
owner_layer: layer_3_orchestration
default_profile: content_execution_guard
allowed_intents: [repurpose_content, switch_content_type]
required_context: [content_type]
allowed_modes: [edit, write]

nodes:
  - node_id: validate_repurposing_scope
    type: orchestration
    purpose: confirm target content type, source content availability, and tone requirements
    required_capabilities: [content_repurposing]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 5
    output_contract: repurposing_scope
    surface_hint: canvas

  - node_id: set_content_type_contract
    type: proposed_tool_request
    purpose: propose content type switch through approved content-type contract
    required_capabilities: [content_repurposing]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: false
    timeout_seconds: 30
    output_contract: content_type_set
    surface_hint: canvas

  - node_id: collect_tone_rules
    type: proposed_tool_request
    purpose: propose reading stored tone rules through approved contract
    required_capabilities: [tone_adaptation]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 30
    output_contract: tone_rule_summary
    surface_hint: canvas

  - node_id: save_new_format_intent
    type: proposed_tool_request
    purpose: propose saving new format intent to brief
    required_capabilities: [brief_building]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 30
    output_contract: repurposing_brief_saved
    surface_hint: brief

  - node_id: build_repurposed_content_contract
    type: content_generation_task
    purpose: create repurposed-content contract; Layer 3 does not write final content directly
    required_capabilities: [content_repurposing, guarded_drafting]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 60
    output_contract: repurposed_content_contract
    surface_hint: canvas

  - node_id: propose_canvas_write_or_edit
    type: proposed_tool_request
    purpose: propose canvas write or exact-match edit through approved contract
    required_capabilities: [output_surface_routing]
    requires_layer2_tool_authorization: true
    requires_user_approval: true
    optional: false
    timeout_seconds: 60
    output_contract: repurposed_canvas_result
    surface_hint: canvas

  - node_id: present_repurposing_status
    type: presentation_intent
    purpose: hand off content-type switch, repurposing status, and unresolved decisions to Layer 5
    required_capabilities: [content_repurposing]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 10
    output_contract: presentation_intent
    surface_hint: canvas

edges:
  - from: validate_repurposing_scope
    to: set_content_type_contract
    condition: always
  - from: set_content_type_contract
    to: collect_tone_rules
    condition: always
  - from: collect_tone_rules
    to: save_new_format_intent
    condition: always
  - from: save_new_format_intent
    to: build_repurposed_content_contract
    condition: always
  - from: build_repurposed_content_contract
    to: propose_canvas_write_or_edit
    condition: if_contract_ready
  - from: propose_canvas_write_or_edit
    to: present_repurposing_status
    condition: always

failure_policy:
  missing_context: stop_and_request_clarification
  mode_not_allowed: stop_and_request_supported_mode
  approval_missing: wait_for_user_decision
  layer2_denied_tool: block_required_canvas_change
  layer4_skill_unavailable: continue_with_reduced_capability_if_optional
  layer5_surface_unavailable: stop_and_return_safe_response
  layer6_runtime_failed: mark_node_failed_and_stop_if_canvas_change_failed
  layer7_data_unavailable: continue_with_data_gap_if_optional
  protected_disclosure: stop_and_return_safe_response
```

### 7.10 `memory-tone-management.bp.yaml`

```yaml
workflow_id: memory-tone-management.bp
workflow_version: 1
owner_layer: layer_3_orchestration
default_profile: content_collaborator
allowed_intents: [memory_tone_management, update_memory]
required_context: []
allowed_modes: [plan, write, edit, optimize]

nodes:
  - node_id: validate_memory_scope
    type: orchestration
    purpose: determine whether request is read-only inspection, reusable preference proposal, or explicit update request
    required_capabilities: [memory_update_proposal, tone_adaptation]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 5
    output_contract: memory_scope_validated
    surface_hint: chat

  - node_id: inspect_memory_or_tone
    type: proposed_tool_request
    purpose: propose memory or tone inspection through approved read contract; Layer 3 never reads document bodies directly
    required_capabilities: [tone_adaptation]
    requires_layer2_tool_authorization: true
    requires_user_approval: false
    optional: true
    timeout_seconds: 30
    output_contract: memory_tone_summary
    surface_hint: chat

  - node_id: request_memory_update_approval
    type: approval_request
    purpose: request explicit approval before proposing any persistent memory or tone update
    required_capabilities: [memory_update_proposal]
    requires_layer2_tool_authorization: false
    requires_user_approval: true
    optional: false
    timeout_seconds: 10
    output_contract: memory_update_approval_request
    surface_hint: chat

  - node_id: propose_memory_update
    type: proposed_tool_request
    purpose: propose persistent memory or tone update only after approval; never apply silently
    required_capabilities: [memory_update_proposal]
    requires_layer2_tool_authorization: true
    requires_user_approval: true
    optional: false
    timeout_seconds: 30
    output_contract: memory_update_proposed
    surface_hint: chat

  - node_id: present_memory_status
    type: presentation_intent
    purpose: hand off memory/tone summary, approval state, or proposal status to Layer 5
    required_capabilities: [memory_update_proposal]
    requires_layer2_tool_authorization: false
    requires_user_approval: false
    optional: false
    timeout_seconds: 10
    output_contract: presentation_intent
    surface_hint: chat

edges:
  - from: validate_memory_scope
    to: inspect_memory_or_tone
    condition: if_read_needed
  - from: inspect_memory_or_tone
    to: request_memory_update_approval
    condition: if_update_needed
  - from: request_memory_update_approval
    to: propose_memory_update
    condition: if_approved
  - from: propose_memory_update
    to: present_memory_status
    condition: always

failure_policy:
  missing_context: ask_only_for_required_missing_fields
  mode_not_allowed: stop_and_request_supported_mode
  approval_missing: wait_for_user_decision
  layer2_denied_tool: terminal_safe_response
  layer4_skill_unavailable: terminal_safe_response
  layer5_surface_unavailable: stop_and_return_safe_response
  layer6_runtime_failed: terminal_safe_response
  layer7_data_unavailable: terminal_safe_response
  protected_disclosure: stop_and_return_safe_response
```

---

## 8. New Public Contracts

### 8.1 Updated intake

```go
type IntakeDecision struct {
    TraceID          string                 `json:"trace_id"`
    Intent           string                 `json:"intent"`
    Confidence       float64                `json:"confidence"`
    SanitizedContext map[string]string      `json:"sanitized_context"`
    DisclosureStatus string                 `json:"disclosure_status,omitempty"`
    PolicyState      string                 `json:"policy_state,omitempty"`
    Mode             string                 `json:"mode,omitempty"`
    ApprovedActions  []string               `json:"approved_actions,omitempty"`
    Metadata         map[string]interface{} `json:"metadata,omitempty"`
}
```

### 8.2 Content generation task

```go
type ContentGenerationTask struct {
    TraceID              string            `json:"trace_id"`
    TaskID               TaskID            `json:"task_id"`
    Intent               string            `json:"intent"`
    Mode                 string            `json:"mode"`
    SurfaceHint          string            `json:"surface_hint"`
    RequiredCapabilities []string          `json:"required_capabilities"`
    Inputs               map[string]string `json:"inputs"`
    OutputContract       string            `json:"output_contract"`
    Constraints          []string          `json:"constraints"`
}
```

Layer 3 may build this contract. It must not generate final article prose, final edits, or final visible output directly.

### 8.3 Approval request

```go
type ApprovalRequest struct {
    TraceID     string                 `json:"trace_id"`
    TaskID      TaskID                 `json:"task_id"`
    ApprovalFor string                 `json:"approval_for"`
    Reason      string                 `json:"reason"`
    Options     []string               `json:"options,omitempty"`
    Payload     map[string]interface{} `json:"payload,omitempty"`
}
```

Layer 3 builds approval requests. Layer 5 presents them. Layer 2 validates returned approval metadata.

---

## 9. New Go Type Additions

```go
type OrchestrationMode string

const (
    ModePlan     OrchestrationMode = "plan"
    ModeWrite    OrchestrationMode = "write"
    ModeEdit     OrchestrationMode = "edit"
    ModeOptimize OrchestrationMode = "optimize"
    ModeAudit    OrchestrationMode = "audit"
)

const (
    NodeTypeContentGeneration NodeType = "content_generation_task"
    NodeTypeApprovalRequest   NodeType = "approval_request"
)

const (
    WorkflowTopicDiscovery      WorkflowID = "topic-discovery.bp"
    WorkflowContentBrief        WorkflowID = "content-brief.bp"
    WorkflowContentResearch     WorkflowID = "content-research.bp"
    WorkflowSEOPlanning         WorkflowID = "seo-planning.bp"
    WorkflowPageAnalysis        WorkflowID = "page-analysis.bp"
    WorkflowArticlePlanning     WorkflowID = "article-planning.bp"
    WorkflowArticleDrafting     WorkflowID = "article-drafting.bp"
    WorkflowContentOptimization WorkflowID = "content-optimization.bp"
    WorkflowContentRepurposing  WorkflowID = "content-repurposing.bp"
    WorkflowMemoryTone          WorkflowID = "memory-tone-management.bp"

    AgentContentCollaborator    AgentID = "content_collaborator"
    AgentContentExecutionGuard  AgentID = "content_execution_guard"
)
```

---

## 10. New File Specs

### 10.1 `brief_state_guard.go`

Required functions:

```go
func ValidateBriefRequiredFields(ctx map[string]string, required []string) error
func BuildBriefUpdateProposal(plan DAGPlan, task TaskNode, fields map[string]string) (ProposedToolRequest, error)
func ValidateNoBriefOverwrite(fields map[string]string, existingKeys []string) error
```

Rules:

```text
accept sanitized brief state only
never read raw brief files directly
never render brief UI
never overwrite existing brief fields unless route explicitly allows
send all brief updates through ProposedToolRequest
```

### 10.2 `canvas_state_guard.go`

Required functions:

```go
func ValidateCanvasWriteMode(decision IntakeDecision) error
func ValidateSelectedTextForEdit(decision IntakeDecision) error
func BuildCanvasChangeProposal(plan DAGPlan, task TaskNode, change ContentGenerationTask) (ProposedToolRequest, error)
```

Rules:

```text
reject draft_article without write mode
reject edit_existing without selected_text
never read raw canvas file directly
never render canvas UI
never write article text directly
never edit canvas text directly
send all canvas changes through ProposedToolRequest
```

### 10.3 `content_task_contracts.go`

Required functions:

```go
func BuildContentGenerationTask(plan DAGPlan, task TaskNode, ctx PlanningContext) (ContentGenerationTask, error)
func ValidateContentGenerationTask(task ContentGenerationTask) error
func ValidateSurfaceHint(surface string) error
```

Rules:

```text
build task contract only
include mode, surface_hint, capability tags, constraints, and sanitized summaries
do not include raw prompt
do not include secrets
do not include internal skill filenames
do not include internal tool inventory
do not include memory document bodies
do not include raw canvas content unless Layer 2 provided sanitized selected_text
```

### 10.4 `approval_request_builder.go`

Required functions:

```go
func BuildApprovalRequest(plan DAGPlan, task TaskNode, approvalFor string, reason string) (ApprovalRequest, error)
func ValidateApprovalRequirement(intent string, mode string, task TaskNode) error
func ValidateApprovalResult(decision IntakeDecision, approvalFor string) error
```

Rules:

```text
deepResearch requires explicit approval
memory update requires explicit approval
edit proposal requires explicit approval when it changes existing canvas content
write mode must be explicitly supplied by Layer 2
never infer approval from free text inside sanitized context
never treat tool success as approval
```

---

## 11. Core Flows

### 11.1 Topic discovery

```text
01 receive topic_discovery
02 validate plan mode
03 collect brand context if available
04 collect source intelligence
05 collect fast research and site discovery
06 inspect material pages only
07 check cannibalization if topic overlaps existing pages
08 request topic/angle decision after research
09 propose brief update
10 hand off topic options to Layer 5
```

### 11.2 Brief building

```text
01 receive content_brief
02 collect default brand/tone summaries through approved requests
03 ask missing brief fields in a batch
04 propose brief save
05 if reusable preference detected, request approval before memory proposal
06 hand off brief state to Layer 5
```

### 11.3 Research

```text
01 receive content_research
02 require topic
03 collect source intelligence
04 propose fast research
05 if evidence-heavy, request deepResearch approval
06 only after approval, propose deepResearch
07 scrape important sources only when needed
08 hand off safe research summary
```

### 11.4 Article planning

```text
01 receive article_planning
02 require topic and plan mode
03 collect strategy context
04 collect evidence before planning
05 ask for audience, angle, length, CTA, screenshots if missing
06 propose brief save
07 build article plan content-generation task contract
08 hand off plan contract to Layer 5
09 do not write article sections
```

### 11.5 Drafting

```text
01 receive draft_article
02 require write mode from Layer 2
03 require topic, audience, intent, target_word_count
04 propose tone and positioning context reads
05 propose supporting page/context reads if needed
06 build section-generation task contract
07 propose canvas write via Layer 2
08 repeat section-by-section
09 request post-write quality skill activation
10 hand off status to Layer 5
11 never write directly to canvas
```

### 11.6 Optimization and editing

```text
01 receive optimize_content or edit_existing
02 validate optimize/edit mode
03 require selected_text for edit_existing
04 propose page, SERP, link, and overlap checks as needed
05 build targeted edit or SEO-output contract
06 propose exact-match canvas edit or null-field SEO update
07 request approval before existing-content changes
08 hand off optimization result
09 never overwrite existing SEO fields silently
```

### 11.7 Repurposing

```text
01 receive repurpose_content or switch_content_type
02 require target content_type
03 propose content-type switch through contract
04 propose tone rule read
05 propose brief update for new format intent
06 build repurposed-content task contract
07 propose canvas write or edit
08 hand off status
```

### 11.8 Memory and tone

```text
01 receive memory_tone_management or update_memory
02 determine read-only vs update
03 propose memory/tone inspection through approved contract if needed
04 if update needed, request approval
05 after approval, propose memory update through Layer 2
06 hand off status
07 never update silently
08 never store memory or tone document bodies
```

---

## 12. New Error Codes

```go
const (
    ErrModeNotAllowed              = "MODE_NOT_ALLOWED"
    ErrWriteModeRequired           = "WRITE_MODE_REQUIRED"
    ErrEditSelectionRequired       = "EDIT_SELECTION_REQUIRED"
    ErrApprovalRequired            = "APPROVAL_REQUIRED"
    ErrDeepResearchApprovalMissing = "DEEP_RESEARCH_APPROVAL_MISSING"
    ErrMemoryApprovalMissing       = "MEMORY_APPROVAL_MISSING"
    ErrSilentMemoryWriteRejected   = "SILENT_MEMORY_WRITE_REJECTED"
    ErrCanvasWriteRejected         = "CANVAS_WRITE_REJECTED"
    ErrBriefOverwriteRejected      = "BRIEF_OVERWRITE_REJECTED"
    ErrOutputSurfaceInvalid        = "OUTPUT_SURFACE_INVALID"
    ErrContentTaskInvalid          = "CONTENT_TASK_INVALID"
)
```

Safe user-facing examples:

```text
I need the topic before I can plan the article.
I need a target length before I can start drafting.
I can plan this, but writing requires Write mode.
I need approval before using deeper research.
I need approval before proposing a saved tone or memory update.
I need selected text before I can prepare a surgical edit.
```

Unsafe user-facing examples:

```text
article-drafting.bp failed at propose_canvas_write
content_execution_guard rejected writeSection
Layer 2 denied writeSection
memory-tone-management.bp blocked proposeMemoryUpdate
```

---

## 13. Production Test Matrix Additions

### 13.1 Content routing tests

```text
topic_discovery + plan mode -> topic-discovery.bp with content_collaborator
content_brief + plan mode -> content-brief.bp with content_collaborator
content_research + topic -> content-research.bp
article_planning + topic + plan mode -> article-planning.bp
draft_article + write mode + required fields -> article-drafting.bp
draft_article + plan mode -> blocked or converted to article-planning.bp
optimize_content + optimize mode -> content-optimization.bp
edit_existing without selected_text -> clarification-only plan
repurpose_content + content_type -> content-repurposing.bp
memory_tone_management -> memory-tone-management.bp
update_memory without approval -> approval request only
unknown content intent -> blocked
```

### 13.2 Mode gate tests

```text
draft_article without mode -> blocked
draft_article with mode=plan -> blocked
draft_article with mode=write but no target_word_count -> clarification-only plan
edit_existing with no selected_text -> blocked
optimize_content with mode=plan -> blocked unless route allows planning-only optimization
deepResearch proposed without approval -> blocked
memory update proposed without approval -> blocked
```

### 13.3 Output surface tests

```text
canvas output through content_generation_task only -> allowed
direct canvas file write by Layer 3 -> blocked
brief update through ProposedToolRequest only -> allowed
direct brief rendering by Layer 3 -> blocked
chat presentation intent for coordination only -> allowed
presentation intent containing A2UI component JSON -> blocked
handoff to arbitrary output path -> blocked
```

### 13.4 Memory tests

```text
readMemoryDoc proposed through Layer 2 -> allowed
Layer 3 reading memory document body -> blocked
proposeMemoryUpdate after approval -> allowed
proposeMemoryUpdate without approval -> blocked
memory doc body stored in state -> blocked
tone doc body stored in state -> blocked
silent memory update -> blocked
```

### 13.5 Content workflow validation tests

```text
workflow with writeSection node not requiring Layer 2 authorization -> blocked
workflow with proposeEdit node not requiring selected_text -> blocked
workflow with memory update node not requiring approval -> blocked
workflow with deepResearch node not requiring approval when evidence-heavy -> blocked
workflow with direct tool execution instructions -> blocked
workflow with raw prompt text -> blocked
workflow with skill filenames as hard dependency -> blocked
workflow with hidden model/provider names -> blocked
```

### 13.6 Boundary tests

Layer 3 must prove it does not:

```text
classify raw intent
inspect prompt injection
authorize tools
execute tools
call quickContext
call getSourcesInsights
call webSearch
call deepResearch
call scrapePage
call getSitePages
call getCannibalizationReport
call askQuestion
call updateBrief
call proposePlan
call readMemoryDoc
call proposeMemoryUpdate
call setContentType
call writeSection
call proposeEdit
load SKILL.md
execute skill scripts
read memory docs
write memory docs
read canvas raw text
write canvas raw text
render UI
connect to MCP
query vector stores
store OpenTelemetry directly
calculate evals
trigger quarantine
persist long-term memory
```

---

## 14. Deployment Requirements

Layer 3 v2 is production-ready only if:

```text
all v1 SEO/AEO workflows still validate
all new content workflows validate at startup
content profiles validate at startup
content routes are internally configured only
unknown content intents fail closed
drafting requires write-mode flag from Layer 2
editing requires selected_text from Layer 2
memory update requires approval
deep research requires approval
canvas writes are proposed, not executed
brief updates are proposed, not executed
chat remains coordination-only by presentation intent
content generation is a task contract, not direct Layer 3 prose generation
Layer 3 never reads memory or tone document bodies
Layer 3 never stores raw canvas content
Layer 3 never stores raw article drafts
Layer 3 never exposes exact tool names, workflow IDs, profile IDs, skill files, or route internals to users
content workflow red-team tests pass
boundary tests still prove no overlap with Layers 1, 2, 4, 5, 6, 7, or 8
```

---

## 15. Acceptance Criteria

Layer 3 v2 is accepted when:

1. It consumes only Layer 2 `IntakeDecision` objects.
2. It rejects raw prompts and user-supplied routes.
3. It maps SEO/AEO and content intents to safe terminal behaviors or allowed workflows.
4. It supports topic discovery, content brief, research, SEO planning, page analysis, article planning, drafting, optimization, repurposing, and memory/tone management.
5. It routes planning to `content_collaborator`.
6. It routes guarded writing/editing to `content_execution_guard`.
7. It enforces Plan mode vs Write mode.
8. It blocks drafting in Plan mode.
9. It blocks silent memory updates.
10. It requires approval before deep research.
11. It requires selected text before surgical edits.
12. It creates content-generation task contracts but does not write final content directly.
13. It creates proposed canvas and brief requests but does not render or mutate surfaces directly.
14. It creates skill activation requests but never loads skill bodies.
15. It creates proposed tool requests but never authorizes or executes them.
16. It emits sanitized orchestration events only.
17. It stores only short-lived orchestration state and sanitized summaries.
18. It fails closed on unknown, malformed, stale, unsafe, unauthorized, or boundary-violating state.
19. It proves through tests that it does not overlap with Layers 1, 2, 4, 5, 6, 7, or 8.

---

## 16. Final Non-Goals

```text
Layer 3 must not classify raw user intent.
Layer 3 must not inspect prompt injection.
Layer 3 must not redact outbound text.
Layer 3 must not authorize tool calls.
Layer 3 must not execute tool calls.
Layer 3 must not run quickContext, getSourcesInsights, webSearch, deepResearch, scrapePage, getSitePages, getCannibalizationReport, askQuestion, updateBrief, proposePlan, readMemoryDoc, proposeMemoryUpdate, setContentType, writeSection, or proposeEdit.
Layer 3 must not load skill-registry.yaml.
Layer 3 must not read SKILL.md.
Layer 3 must not execute scripts from skills.
Layer 3 must not write article sections directly.
Layer 3 must not edit canvas text directly.
Layer 3 must not render dashboards.
Layer 3 must not generate A2UI components.
Layer 3 must not serve A2A Agent Cards.
Layer 3 must not open MCP connections.
Layer 3 must not call MCP servers.
Layer 3 must not enforce sandbox mounts.
Layer 3 must not manage filesystem write controls.
Layer 3 must not manage network egress.
Layer 3 must not mint JIT credentials.
Layer 3 must not revoke credentials.
Layer 3 must not trigger quarantine.
Layer 3 must not store long-term telemetry.
Layer 3 must not calculate evaluations.
Layer 3 must not persist long-term memory.
Layer 3 must not silently update memory.
Layer 3 must not expose workflow IDs, profile IDs, route IDs, exact tool names, skill files, memory paths, or internal handoff details to end users.
```

---

## 17. One-Line Architecture Summary

Layer 3 v2 converts Layer 2 validated SEO/AEO and content-agent intents into safe workflow plans, profile routes, DAG nodes, approval gates, and handoff contracts while leaving sanitization, authorization, execution, skill loading, memory persistence, rendering, connection, observation, and recovery to their owning layers.
