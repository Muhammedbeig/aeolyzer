# Layer 4 Production-Ready Specs v2
## Progressive Disclosure Skill Directory for SEO/AEO Auditor + Content Agent

**Status:** Production-ready upgrade  
**Supersedes:** Layer 4 v1 where the library was SEO/AEO-first and had no content-agent governance model  
**Primary rule:** Layer 4 is the procedural-memory layer. It owns the skill registry, skill folders, SKILL.md bodies, references, assets, scripts-as-artifacts, skill packaging, skill validation, skill versioning, and skill-eval definitions. It never classifies raw user intent, chooses workflows, assigns agents, authorizes tools, executes tools or scripts, renders UI, connects MCP, reads/writes memory documents, stores telemetry, scores evals, or quarantines runtime.

---

## 1. Upgrade Decision

Upgrade is needed.

Layer 4 v1 correctly implements progressive disclosure for the original SEO/AEO skill set, but Layer 2 v2 and Layer 3 v2 now support a broader content-agent workflow surface:

```text
topic discovery
content brief
research
SEO planning
page analysis
article planning
drafting
optimization
repurposing
memory and tone management
```

Layer 4 must therefore be upgraded from a 22-skill SEO/AEO directory into a production-grade, versioned, governed skill library that supports both:

```text
SEO/AEO auditor capability
Content collaborator capability
Guarded content execution capability
```

The upgrade keeps the existing architecture intact:

```text
Level 1: metadata cues, always available as a compact registry
Level 2: SKILL.md body, loaded only on explicit skill activation
Level 3: bundled scripts, references, assets, and eval definitions, loaded strictly as needed
```

Layer 4 must not become an orchestrator. Layer 3 still selects workflows and requests skill activation. Layer 4 only validates, packages, and serves the requested procedural memory.

---

## 2. Layer Objective

Layer 4 turns domain expertise into portable, testable, versioned procedural memory.

It must:

1. Maintain the canonical skill registry.
2. Keep always-loaded skill metadata compact and safe.
3. Validate every skill folder at build time.
4. Validate every requested skill activation at runtime.
5. Serve SKILL.md bodies only when requested by Layer 3.
6. Serve references and assets only when explicitly required by an active SKILL.md.
7. Store deterministic helper scripts as signed artifacts only.
8. Never execute scripts directly.
9. Store eval fixtures and rubrics for Layer 8.
10. Track version, owner, tier, maturity, risk, token estimates, checksums, and compatibility.
11. Support SEO/AEO, content planning, guarded drafting, optimization, repurposing, and tone guidance without mixing responsibilities.
12. Preserve portability across agent runtimes.
13. Prevent context rot through strict progressive disclosure.
14. Prevent trigger collisions through registry-level tests.
15. Fail closed when skill metadata, body, resource manifest, checksum, tier, or compatibility is invalid.

Layer 4 exists to answer this question:

```text
Given a validated workflow node and capability profile from Layer 3, which exact procedural instructions and supporting resources are safe and appropriate to load?
```

It must not answer:

```text
What did the user ask?
Which workflow should run?
Which agent should be assigned?
Which tool should execute?
Which UI should render?
Which data source should be queried?
Was the trajectory successful?
Should the agent be quarantined?
```

---

## 3. Strict Layer Boundary

### 3.1 Layer 4 owns

```text
skill registry
skill metadata schema
SKILL.md schema
skill folder layout
skill resource manifests
skill package validation
skill versioning
skill checksums and signatures
skill owner metadata
skill maturity tiers
skill risk tiers
skill token estimates
skill compatibility metadata
skill activation validation
skill body loading
reference loading
asset loading
script artifact registration
script checksum validation
script interface metadata
skill deprecation policy
skill collision detection
trigger fixture definitions
golden dataset definitions
trajectory expectation definitions
LLM-judge rubric definitions
skill release gates
skill library changelog
```

### 3.2 Layer 4 does not own

| Responsibility | Owning layer |
|---|---|
| Raw user input normalization | Layer 2 |
| Prompt-injection detection | Layer 2 |
| Protected-disclosure classification | Layer 2 |
| Intent classification | Layer 2 |
| Tool authorization | Layer 2 |
| Workflow selection | Layer 3 |
| DAG construction and sequencing | Layer 3 |
| Capability profile selection | Layer 3 |
| Context-window budgeting from sanitized context | Layer 3 |
| Skill activation request creation | Layer 3 |
| Content-generation task contract creation | Layer 3 |
| Canvas, brief, chat, dashboard, A2UI rendering | Layer 5 |
| Approval card rendering and user decision capture | Layer 5 |
| Tool execution | Layer 6 |
| Script execution | Layer 6 |
| Runtime sandboxing | Layer 6 |
| Filesystem mount enforcement | Layer 6 |
| Network egress control | Layer 6 |
| JIT credential minting or revocation | Layer 6 |
| MCP transport | Layer 7 |
| Connector execution | Layer 7 |
| RAG and vector retrieval | Layer 7 |
| Memory document read/write | Layer 7 for retrieval, governed by Layer 2/3/5 approvals |
| Telemetry storage | Layer 8 |
| Eval scoring | Layer 8 |
| Drift scoring | Layer 8 |
| AgBOM tracking | Layer 8 |
| SecOps red/blue/green loops | Layer 8 |
| Quarantine | Layer 8 decides, Layer 6 executes |

### 3.3 Non-overlap rule for skills

A skill may describe a procedure and declare abstract action needs. It must not perform the action.

Allowed inside SKILL.md:

```text
Explain how to conduct keyword research.
Declare that live SERP evidence is needed.
Describe how to evaluate source credibility.
Reference an output template.
Reference a deterministic helper script.
Define expected output shape.
Define post-write checklist.
Define when not to use the skill.
```

Forbidden inside SKILL.md:

```text
Call a tool directly.
Name raw internal tool IDs as user-visible instructions.
Open network connections.
Connect MCP servers.
Mint credentials.
Render UI cards.
Write to canvas directly.
Write to brief directly.
Persist memory directly.
Store telemetry.
Score evals.
Start quarantine.
Override Layer 2 policy.
Override Layer 3 routing.
Override Layer 6 sandboxing.
Override Layer 7 tenant boundaries.
```

---

## 4. Whitepaper-Aligned Design Principles

Layer 4 implements the skill principles from the whitepapers as production rules.

### 4.1 Skills are procedural memory, not general context

A skill is a reusable runbook. It teaches the agent **how** to perform a task, not merely **what** facts exist.

Layer 4 must therefore keep skills focused on repeatable procedures:

```text
how to interpret GSC decline patterns
how to build a source-backed content brief
how to evaluate internal-link opportunities
how to draft section-by-section under a brief
how to audit long-form content quality
how to translate CWV metrics into frontend tickets
how to produce JSON-LD without hallucinating schema
```

Layer 4 must reject skills that are only large fact dumps. Facts belong in references, assets, RAG, memory, or external data systems, depending on ownership.

### 4.2 Progressive disclosure is mandatory

Every skill loads in three levels:

```text
Level 1: registry metadata
  Always available in compact form.
  Includes name, description, anti-triggers, tags, tier, and token estimate.

Level 2: SKILL.md body
  Loaded only after Layer 3 requests that exact skill and Layer 4 validates compatibility.

Level 3: resources
  references/, assets/, scripts/, and evals/ are loaded or surfaced only when required by the active skill.
```

No skill may force all resources into the context window.

### 4.3 Description is the routing interface

The description is the most important registry field. It is the only always-loaded text the orchestrator and model see when deciding whether a skill may be useful.

Every description must:

```text
state what the skill does
front-load concrete trigger verbs
include when to use
include when not to use
avoid vague verbs such as help, assist, improve, handle
avoid internal jargon
avoid vendor-specific runtime names
stay within token limits
```

### 4.4 One skill, one job

Layer 4 rejects multi-purpose skills.

Bad:

```text
content-helper
Does research, outlines, drafts, edits, optimizes, creates metadata, and updates memory.
```

Good:

```text
source-backed-research
Finds and organizes reliable supporting evidence for a content brief. Use when a workflow needs current sources, statistics, quotes, or claim validation. Do NOT use for drafting, metadata, or live page scraping decisions.
```

### 4.5 Skills are dependencies

Every skill is treated as code:

```text
versioned
owned
reviewed
pinned
checksummed
tested
release-gated
deprecatable
auditable
```

A skill without tests cannot reach production.

### 4.6 Scripts are deterministic helpers, not autonomy

Scripts in `scripts/` may contain deterministic parsing, counting, formatting, or validation logic, but Layer 4 only stores and signs them. Layer 6 executes them if Layer 2 authorizes the action and Layer 3 requests execution.

### 4.7 Skills compose with MCP; they do not replace MCP

Layer 4 owns know-how. Layer 7 owns reach.

A skill may say:

```text
This procedure needs source intelligence, page content, analytics data, or site page discovery.
```

It must not:

```text
open the connector
name a connector endpoint as user-visible text
execute a connector
own transport
manage authentication
```

### 4.8 Skills compose with AGENTS.md; they do not replace it

AGENTS.md remains always-loaded project DNA: conventions, stack, global behavior. Layer 4 skill registry is a compact catalog of specialized procedures. SKILL.md bodies are loaded on demand.

### 4.9 Evaluation precedes promotion

Layer 4 stores the eval definitions that prove a skill works:

```text
positive trigger cases
negative trigger cases
golden output fixtures
expected action-class trajectory
rubrics
regression fixtures
token-load fixtures
co-loaded collision fixtures
```

Layer 8 runs and scores them. Layer 4 only stores, validates, and version-controls the definitions.

---

## 5. Required Directory Upgrade

### 5.1 Final Layer 4 tree

```text
/layer_04_skills
  ├── skill-registry.yaml
  ├── registry.schema.json
  ├── skill.schema.json
  ├── resource-manifest.schema.json
  ├── eval-manifest.schema.json
  ├── skill-library.lock
  ├── skill-changelog.md
  ├── /registry_views
  │   ├── public_capability_view.yaml
  │   ├── layer3_activation_view.yaml
  │   └── layer8_eval_view.yaml
  ├── /policies
  │   ├── skill_package_policy.yaml
  │   ├── skill_quality_policy.yaml
  │   ├── skill_token_budget_policy.yaml
  │   ├── skill_compatibility_policy.yaml
  │   └── skill_deprecation_policy.yaml
  ├── /internal
  │   ├── skill_registry_loader.go
  │   ├── skill_package_validator.go
  │   ├── skill_activation_validator.go
  │   ├── skill_body_loader.go
  │   ├── resource_manifest_loader.go
  │   ├── token_estimator.go
  │   ├── checksum_verifier.go
  │   ├── collision_detector.go
  │   ├── skill_bundle_builder.go
  │   └── skill_event_emitter.go
  ├── /tests
  │   ├── registry_schema_test.go
  │   ├── skill_schema_test.go
  │   ├── resource_manifest_test.go
  │   ├── activation_compatibility_test.go
  │   ├── trigger_collision_test.go
  │   ├── token_budget_test.go
  │   ├── protected_metadata_test.go
  │   ├── no_execution_boundary_test.go
  │   ├── no_mcp_boundary_test.go
  │   ├── no_ui_boundary_test.go
  │   └── no_memory_persistence_boundary_test.go
  └── /skills
      ├── /keyword_research
      ├── /serp_analysis
      ├── /competitor_intelligence
      ├── /content_strategy
      ├── /content_creation
      ├── /content_seo_settings
      ├── /meta_optimization
      ├── /schema_generation
      ├── /llms_txt_generation
      ├── /robots_txt_generation
      ├── /sitemap_generation
      ├── /site_audit_interpretation
      ├── /core_web_vitals_optimization
      ├── /internal_linking_strategy
      ├── /long_form_content_audit
      ├── /content_refresh_strategy
      ├── /link_opportunity_discovery
      ├── /backlink_strategy
      ├── /gsc_insights_analysis
      ├── /ga4_analysis
      ├── /local_seo_optimization
      ├── /gbp_optimization
      ├── /topic_discovery
      ├── /content_brief_building
      ├── /source_backed_research
      ├── /seo_content_planning
      ├── /page_content_analysis
      ├── /article_planning
      ├── /guarded_drafting
      ├── /content_optimization
      ├── /content_repurposing
      ├── /tone_memory_guidance
      ├── /citation_source_safety
      └── /content_quality_gates
```

### 5.2 Required folder layout per skill

Every production skill must use this layout:

```text
/skills/<skill_dir>
  ├── SKILL.md
  ├── resource-manifest.yaml
  ├── eval-manifest.yaml
  ├── CHANGELOG.md
  ├── OWNERS
  ├── /references
  │   └── *.md
  ├── /assets
  │   └── *.yaml | *.json | *.md | *.txt
  ├── /scripts
  │   └── *.go | *.py | *.sh
  └── /evals
      ├── trigger_cases.yaml
      ├── golden_cases.yaml
      ├── trajectory_cases.yaml
      ├── rubric.yaml
      └── regression_cases.yaml
```

Only `SKILL.md`, `resource-manifest.yaml`, `eval-manifest.yaml`, `CHANGELOG.md`, and `OWNERS` are mandatory for every production skill.

`references/`, `assets/`, `scripts/`, and `evals/` may be empty for simple read-only skills, but their absence must be explicitly declared in the manifests.

---

## 6. Canonical Skill Library v2

### 6.1 Existing SEO/AEO skills retained

The following v1 skills remain canonical and must be upgraded to the v2 manifest format:

```text
keyword_research
serp_analysis
competitor_intelligence
content_strategy
content_creation
content_seo_settings
meta_optimization
schema_generation
llms_txt_generation
robots_txt_generation
sitemap_generation
site_audit_interpretation
core_web_vitals_optimization
internal_linking_strategy
long_form_content_audit
content_refresh_strategy
link_opportunity_discovery
backlink_strategy
gsc_insights_analysis
ga4_analysis
local_seo_optimization
gbp_optimization
```

### 6.2 New content-agent skills

Add these v2 skills:

```text
topic_discovery
content_brief_building
source_backed_research
seo_content_planning
page_content_analysis
article_planning
guarded_drafting
content_optimization
content_repurposing
tone_memory_guidance
citation_source_safety
content_quality_gates
```

### 6.3 Rationale for new skills

| Skill | Purpose | Why separate |
|---|---|---|
| `topic_discovery` | Finds content topics, gaps, audience questions, and angle candidates | Prevents topic work from bloating general content strategy |
| `content_brief_building` | Converts topic, audience, intent, CTA, and constraints into a structured brief | Brief creation is a contract-building skill, not drafting |
| `source_backed_research` | Defines research procedure, source triage, evidence extraction, and claim support | Research has different safety and recency rules from drafting |
| `seo_content_planning` | Plans search intent, keyword direction, internal linking, and cannibalization concerns | SEO planning is not the same as technical SEO audit |
| `page_content_analysis` | Reviews a live or provided page for quality, SEO, AEO, structure, and gaps | Page analysis is narrower than full site audit |
| `article_planning` | Builds section-by-section article plans from brief and evidence | Planning must remain separate from writing |
| `guarded_drafting` | Provides section-by-section drafting procedure under brief, tone, and source constraints | Drafting is high-output and must be guarded |
| `content_optimization` | Improves existing content, metadata, clarity, structure, and internal links | Optimization is edit/optimize mode, not draft mode |
| `content_repurposing` | Converts content across formats while preserving audience, message, and source constraints | Format switching has unique loss/fit rules |
| `tone_memory_guidance` | Guides use of brand/tone preferences and memory-update proposals | Layer 4 can provide procedure, but must not read/write memory |
| `citation_source_safety` | Defines source credibility, competitor exclusion, and citation hygiene | Source trust needs reusable specialized rules |
| `content_quality_gates` | Provides post-write and post-edit QA checklists | Quality gates should be shared across drafting and optimization |

---

## 7. Registry Schema v2

### 7.1 `skill-registry.yaml`

`skill-registry.yaml` is the only Layer 4 file designed for always-loaded metadata.

It must remain compact.

```yaml
version: 2
policy_mode: fail_closed
registry_owner_layer: layer_4_skills
default_status: active
metadata_token_budget:
  max_total_registry_tokens: 2500
  max_description_tokens_per_skill: 90
  max_antitrigger_tokens_per_skill: 60

skills:
  - skill_id: topic_discovery
    name: topic-discovery
    directory: skills/topic_discovery
    status: active
    version: 1.0.0
    owner_team: content_strategy
    tier: read
    risk_class: low
    description: >
      Finds content topic ideas, audience questions, competitor coverage gaps, and angle candidates.
      Use when planning what to write before a brief exists. Do NOT use for drafting or editing.
    anti_triggers:
      - user already selected topic and wants article sections
      - user asks to rewrite existing content
      - user asks for technical site audit
    compatible_profiles:
      - content_collaborator
    compatible_intents:
      - topic_discovery
      - content_strategy
    allowed_modes:
      - plan
    capability_tags:
      - topic_discovery
      - content_strategy
      - audience_strategy
      - source_intelligence
    declared_action_classes:
      - read_brand_context
      - read_source_intelligence
      - web_research
      - page_scrape
      - site_page_discovery
      - cannibalization_check
    output_contracts:
      - topic_options
      - evidence_summary
      - brief_update_proposal
    body_token_estimate: 1800
    resource_token_estimate: 2400
    eval_manifest: skills/topic_discovery/eval-manifest.yaml
    resource_manifest: skills/topic_discovery/resource-manifest.yaml
    checksum: sha256:<computed>
```

### 7.2 Registry validation rules

Reject registry startup if:

```text
version missing
policy_mode is not fail_closed
skill_id duplicate
skill_id contains spaces or unsafe characters
directory missing
directory does not exist
SKILL.md missing
status invalid
tier invalid
risk_class invalid
description missing
description too vague
description too long
anti_triggers missing for production skill
compatible_profiles empty
compatible_intents empty
allowed_modes empty
capability_tags empty
declared_action_classes contains unknown class
declared_action_classes names exact internal tool IDs
output_contracts contain UI component IDs
resource_manifest missing
eval_manifest missing
checksum missing
checksum mismatch
owner_team missing
skill has no OWNERS file
```

### 7.3 Status enum

```yaml
status:
  - draft
  - experimental
  - active
  - deprecated
  - blocked
  - retired
```

Runtime behavior:

```text
draft -> can be loaded only in development or explicit testing sessions
experimental -> can be loaded only in canary or shadow mode
active -> can be loaded in production
deprecated -> can be loaded only if no active replacement exists; emits warning
blocked -> must not load
retired -> must not load
```

### 7.4 Tier enum

```yaml
tier:
  - read
  - draft
  - act
```

Tier meaning:

```text
read:
  Skill guides analysis, interpretation, planning, or summarization.
  It may declare read-only data needs.
  It must not produce final publishable content or action side effects.

draft:
  Skill guides generation of text, plans, briefs, outlines, metadata, or suggested edits.
  It may produce proposed output contracts.
  It must not mutate persistent surfaces directly.

act:
  Skill guides procedures that may lead to external or persistent changes.
  It requires stronger trajectory expectations and approval bindings.
  Layer 4 still does not execute actions.
```

For this platform, most Layer 4 skills should be `read` or `draft`. `act` must be rare and requires Layer 2 authorization, Layer 3 approval gate, Layer 6 runtime isolation, and Layer 8 trajectory evaluation.

---

## 8. SKILL.md Schema v2

Every SKILL.md must begin with YAML frontmatter and continue with Markdown body.

### 8.1 Required frontmatter

```yaml
---
name: guarded-drafting
description: >
  Drafts article sections from an approved brief, audience, tone, and source constraints.
  Use in Write mode when the article plan is finalized. Do NOT use for topic discovery,
  research collection, or direct publishing.
version: 1.0.0
owner_team: content
tier: draft
risk_class: medium
compatible_profiles:
  - content_execution_guard
compatible_intents:
  - draft_article
  - repurpose_content
allowed_modes:
  - write
capability_tags:
  - guarded_drafting
  - section_by_section_writing
  - tone_adaptation
  - source_grounded_writing
declared_action_classes:
  - read_memory_or_tone
  - read_brand_context
  - page_scrape
  - site_page_discovery
  - canvas_write
output_contracts:
  - draft_section
  - internal_link_suggestions
  - post_write_quality_summary
token_budget:
  body_max_tokens: 2200
  references_max_tokens: 2500
  assets_max_tokens: 1500
  total_active_max_tokens: 5500
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---
```

### 8.2 Required Markdown sections

Every SKILL.md body must include:

```markdown
# <Skill Name>

## Purpose
## When to use
## When NOT to use
## Inputs expected
## Procedure
## Output contract
## Quality gates
## Boundary rules
## Resources
## Failure behavior
```

### 8.3 Required boundary language

Every SKILL.md must include a `Boundary rules` section with explicit non-overlap wording.

Example:

```markdown
## Boundary rules

This skill provides procedural guidance only.

It must not:
- classify user intent
- choose workflows or agents
- authorize tools
- execute tools or scripts
- connect to MCP servers
- read or write memory documents
- mutate canvas, brief, chat, dashboard, or UI state
- store telemetry or score evals
- expose internal tool IDs, workflow IDs, profile IDs, or protected metadata
```

### 8.4 Style requirements

Use:

```text
short Markdown headings
plain procedural steps
rationale where it improves generalization
YAML only for nested schemas or structured config
concrete examples
positive triggers
negative triggers
verifiable output constraints
references instead of long prose
scripts instead of deterministic instructions
assets instead of template text in body
```

Avoid:

```text
long generic policy essays
giant bodies over 5,000 words
hidden prompts
provider-specific runtime instructions
exact connector endpoints
raw tool payloads
secrets
API keys
cookies
trace IDs
workflow IDs
profile IDs
UI component JSON
MCP server URLs
raw memory document bodies
raw canvas document bodies
```

---

## 9. Resource Manifest Schema

Each skill must declare its resources.

```yaml
version: 1
skill_id: guarded_drafting
resources:
  references:
    - id: tone_rules_reference
      path: references/tone_rules.md
      load_policy: on_explicit_reference
      max_tokens: 1200
      checksum: sha256:<computed>
      contains_sensitive_data: false

  assets:
    - id: draft_section_contract
      path: assets/draft_section_contract.yaml
      load_policy: on_output_contract_match
      max_tokens: 800
      checksum: sha256:<computed>
      contains_ui_schema: false

  scripts:
    - id: heading_structure_checker
      path: scripts/heading_structure_checker.go
      language: go
      execution_owner_layer: layer_6_runtime
      execution_policy: requires_layer2_authorization
      network_access: false
      filesystem_access: stdin_stdout_only
      max_runtime_seconds: 10
      checksum: sha256:<computed>
      interface:
        input_contract: heading_structure_input
        output_contract: heading_structure_report
```

### 9.1 Resource validation rules

Reject resource manifest if:

```text
resource path escapes skill folder
resource path absolute
resource checksum missing
reference exceeds max token limit
asset contains A2UI component JSON
asset contains connector URL
asset contains secret-like values
script lacks checksum
script requests network access
script requests arbitrary filesystem access
script contains package install commands
script contains credentials
script contains direct MCP calls
script contains telemetry writes
script contains UI rendering code
script contains memory persistence logic
```

### 9.2 Resource load policies

```yaml
load_policy:
  - never_runtime
  - on_explicit_reference
  - on_output_contract_match
  - on_eval_only
  - on_layer8_eval_request
  - on_layer6_execution_request
```

Layer 4 may load `references` and `assets` into a `SkillBundle` when requested. It must never execute scripts.

---

## 10. Eval Manifest Schema

Each skill must include an eval manifest even if the eval files are not loaded at runtime.

```yaml
version: 1
skill_id: guarded_drafting
eval_owner_layer: layer_8_observability
stored_by_layer: layer_4_skills
minimum_gate:
  trigger_accuracy: 0.90
  negative_trigger_precision: 0.90
  output_rubric_min_score: 4
  trajectory_mode:
    read: ANY_ORDER
    draft: IN_ORDER
    act: EXACT
  regression_pass_required: true
  token_budget_pass_required: true

eval_files:
  trigger_cases: evals/trigger_cases.yaml
  golden_cases: evals/golden_cases.yaml
  trajectory_cases: evals/trajectory_cases.yaml
  rubric: evals/rubric.yaml
  regression_cases: evals/regression_cases.yaml
```

### 10.1 Trigger cases

```yaml
positive:
  - input: "Write the first section from this approved brief."
    expected_skill: guarded_drafting
    expected_mode: write

negative:
  - input: "Find topic ideas for our blog."
    expected_not_skill: guarded_drafting
    expected_skill: topic_discovery
```

Every production skill requires at minimum:

```text
3 positive trigger cases
3 negative trigger cases
1 rephrasing-stability case per positive trigger
1 adjacent-skill collision case
1 out-of-scope case
```

### 10.2 Trajectory cases

Layer 4 eval fixtures must use action classes, not exact tool IDs.

Allowed:

```yaml
expected_action_classes:
  - read_memory_or_tone
  - read_brand_context
  - site_page_discovery
  - canvas_write
trajectory_mode: IN_ORDER
```

Forbidden:

```yaml
expected_tool_ids:
  - exactInternalToolName
```

Layer 2 owns tool authorization and may internally map action classes to exact tools. Layer 4 must not expose tool inventories.

---

## 11. Skill Activation Contract

Layer 4 accepts only structured skill activation requests from Layer 3.

### 11.1 Request

```go
type SkillActivationRequest struct {
    TraceID             string            `json:"trace_id"`
    WorkflowID          string            `json:"workflow_id"`
    NodeID              string            `json:"node_id"`
    Intent              string            `json:"intent"`
    Mode                string            `json:"mode"`
    ProfileID           string            `json:"profile_id"`
    RequestedSkillIDs   []string          `json:"requested_skill_ids"`
    RequiredTags        []string          `json:"required_tags,omitempty"`
    OutputContracts     []string          `json:"output_contracts,omitempty"`
    MaxTokenBudget      int               `json:"max_token_budget"`
    SanitizedContextRef string            `json:"sanitized_context_ref,omitempty"`
    ResourceHints       []string          `json:"resource_hints,omitempty"`
    EvalMode            bool              `json:"eval_mode,omitempty"`
    Metadata            map[string]string `json:"metadata,omitempty"`
}
```

Layer 4 must treat `WorkflowID`, `ProfileID`, and `TraceID` as internal metadata. It may validate compatibility but must not expose those values in user-visible output.

### 11.2 Response

```go
type SkillActivationResponse struct {
    TraceID             string              `json:"trace_id"`
    LoadedSkills        []SkillBundle       `json:"loaded_skills"`
    OmittedSkills       []OmittedSkill      `json:"omitted_skills,omitempty"`
    TokenEstimate       int                 `json:"token_estimate"`
    ResourceHandles     []ResourceHandle    `json:"resource_handles,omitempty"`
    CompatibilityStatus string              `json:"compatibility_status"`
    Warnings            []string            `json:"warnings,omitempty"`
    SafeSummary         string              `json:"safe_summary,omitempty"`
}
```

### 11.3 Skill bundle

```go
type SkillBundle struct {
    SkillID          string            `json:"skill_id"`
    Name            string            `json:"name"`
    Version         string            `json:"version"`
    Tier            string            `json:"tier"`
    BodyMarkdown    string            `json:"body_markdown"`
    LoadedResources []LoadedResource  `json:"loaded_resources,omitempty"`
    TokenEstimate   int               `json:"token_estimate"`
    Checksums       map[string]string `json:"checksums"`
    OutputContracts []string          `json:"output_contracts"`
}
```

### 11.4 Activation validation flow

```text
01 receive SkillActivationRequest from Layer 3
02 validate trace_id exists and request shape is valid
03 validate requested_skill_ids are known
04 validate skill status permits load
05 validate intent compatibility
06 validate mode compatibility
07 validate profile compatibility
08 validate output contract compatibility
09 validate requested resource hints
10 verify checksums
11 estimate token cost
12 enforce token budget
13 load SKILL.md body only for approved skill IDs
14 load references/assets only if resource hints and load policies permit
15 expose script handles only, never execute scripts
16 build SkillActivationResponse
17 emit sanitized skill_loaded or skill_load_blocked event toward Layer 8
```

### 11.5 Fail-closed behavior

Reject activation when:

```text
unknown skill_id
blocked skill
retired skill
draft skill in production
checksum mismatch
manifest missing
invalid SKILL.md frontmatter
invalid resource manifest
invalid eval manifest
intent incompatible
mode incompatible
profile incompatible
output contract incompatible
token budget exceeded
resource path unsafe
script resource requests disallowed access
resource contains protected metadata
```

---

## 12. Runtime Skill Loading Policy

### 12.1 Always-loaded metadata

Layer 4 provides compact registry metadata to Layer 3. It does not inject it into the model directly unless the architecture uses Layer 3 to assemble context.

Metadata includes only:

```text
skill_id
name
description
anti_triggers
compatible_intents
compatible_profiles
allowed_modes
capability_tags
tier
risk_class
body_token_estimate
status
```

Metadata must not include:

```text
SKILL.md body
references
assets
scripts
raw eval cases
exact tool IDs
MCP endpoints
file paths shown to user
secrets
memory paths
trace IDs
workflow DAG details
```

### 12.2 Body loading

`SKILL.md` body loads only if:

```text
Layer 3 requested it
Layer 4 validates compatibility
Layer 4 validates checksum
Layer 4 validates token budget
skill status allows load
```

### 12.3 Reference loading

`references/` load only when:

```text
the active SKILL.md explicitly points to the reference
the resource manifest permits it
the resource stays within token limits
the requested workflow node requires it
```

### 12.4 Asset loading

`assets/` load only when:

```text
an output contract or resource hint requires them
the manifest permits loading
the asset is not UI rendering logic
the asset is not connector logic
the asset is not a memory document
```

### 12.5 Script handling

Layer 4 may return script metadata and checksum. It must never execute.

Script execution path:

```text
Layer 3 proposes script-related action class
Layer 2 authorizes proposed action
Layer 6 executes inside sandbox if allowed
Layer 8 observes execution
```

---

## 13. Capability Profile Compatibility

Layer 4 does not own capability profiles, but each skill must declare compatible profiles so Layer 3 can request valid combinations and Layer 4 can enforce them.

### 13.1 Profile-to-skill map

```yaml
profile_skill_map:
  seo_aeo_auditor:
    allowed_skills:
      - keyword_research
      - serp_analysis
      - competitor_intelligence
      - content_strategy
      - content_seo_settings
      - meta_optimization
      - schema_generation
      - llms_txt_generation
      - robots_txt_generation
      - sitemap_generation
      - site_audit_interpretation
      - core_web_vitals_optimization
      - internal_linking_strategy
      - long_form_content_audit
      - content_refresh_strategy
      - link_opportunity_discovery
      - backlink_strategy
      - gsc_insights_analysis
      - ga4_analysis
      - local_seo_optimization
      - gbp_optimization

  content_collaborator:
    allowed_skills:
      - topic_discovery
      - content_brief_building
      - source_backed_research
      - seo_content_planning
      - page_content_analysis
      - article_planning
      - citation_source_safety
      - tone_memory_guidance
      - content_strategy
      - keyword_research
      - serp_analysis
      - internal_linking_strategy

  content_execution_guard:
    allowed_skills:
      - guarded_drafting
      - content_optimization
      - content_repurposing
      - content_quality_gates
      - tone_memory_guidance
      - content_creation
      - long_form_content_audit
      - content_seo_settings
      - meta_optimization
      - internal_linking_strategy
```

Layer 4 validates this map but does not select profiles.

---

## 14. Intent-to-Skill Compatibility

Layer 4 uses intent compatibility only to validate skill loading, not to classify intent.

```yaml
intent_skill_compatibility:
  topic_discovery:
    primary:
      - topic_discovery
    supporting:
      - content_strategy
      - source_backed_research
      - citation_source_safety

  content_brief:
    primary:
      - content_brief_building
    supporting:
      - tone_memory_guidance
      - content_strategy
      - seo_content_planning

  content_research:
    primary:
      - source_backed_research
      - citation_source_safety
    supporting:
      - competitor_intelligence

  seo_planning:
    primary:
      - seo_content_planning
    supporting:
      - keyword_research
      - serp_analysis
      - internal_linking_strategy
      - content_strategy

  page_analysis:
    primary:
      - page_content_analysis
    supporting:
      - long_form_content_audit
      - content_seo_settings
      - meta_optimization
      - internal_linking_strategy

  article_planning:
    primary:
      - article_planning
    supporting:
      - content_brief_building
      - source_backed_research
      - citation_source_safety
      - seo_content_planning

  draft_article:
    primary:
      - guarded_drafting
    supporting:
      - content_quality_gates
      - tone_memory_guidance
      - internal_linking_strategy

  optimize_content:
    primary:
      - content_optimization
    supporting:
      - content_quality_gates
      - content_refresh_strategy
      - meta_optimization
      - content_seo_settings
      - long_form_content_audit

  repurpose_content:
    primary:
      - content_repurposing
    supporting:
      - tone_memory_guidance
      - content_quality_gates

  switch_content_type:
    primary:
      - content_repurposing
    supporting:
      - content_brief_building
      - tone_memory_guidance

  edit_existing:
    primary:
      - content_optimization
    supporting:
      - content_quality_gates
      - tone_memory_guidance

  memory_tone_management:
    primary:
      - tone_memory_guidance
    supporting: []

  update_memory:
    primary:
      - tone_memory_guidance
    supporting: []
```

Validation rule:

```text
Layer 4 may block a skill activation request if the requested skill is not compatible with the already-classified intent.
Layer 4 must not change the intent.
Layer 4 must not select a different workflow.
```

---

## 15. Skill Specs for New Content-Agent Skills

### 15.1 `topic_discovery`

```yaml
skill_id: topic_discovery
name: topic-discovery
tier: read
risk_class: low
compatible_profiles: [content_collaborator]
compatible_intents: [topic_discovery, content_strategy]
allowed_modes: [plan]
capability_tags:
  - topic_discovery
  - audience_strategy
  - content_gap_analysis
  - source_intelligence
declared_action_classes:
  - read_brand_context
  - read_source_intelligence
  - web_research
  - page_scrape
  - site_page_discovery
  - cannibalization_check
output_contracts:
  - topic_option_list
  - topic_evidence_summary
  - brief_topic_update_proposal
```

Purpose:

```text
Guide discovery of viable content topics, audience questions, competitor coverage gaps, and differentiated angles before the content brief exists.
```

Must not:

```text
draft article sections
write metadata
update memory
render topic cards
call web tools directly
```

Required evals:

```text
3 positive topic prompts
3 negative drafting/editing prompts
1 adjacent content_strategy collision case
1 topic already selected case
```

### 15.2 `content_brief_building`

```yaml
skill_id: content_brief_building
name: content-brief-building
tier: draft
risk_class: medium
compatible_profiles: [content_collaborator]
compatible_intents: [content_brief, article_planning, topic_discovery]
allowed_modes: [plan]
capability_tags:
  - brief_building
  - audience_strategy
  - hidden_intent
  - content_contracting
declared_action_classes:
  - read_brand_context
  - read_memory_or_tone
  - ask_user_question
  - update_brief
output_contracts:
  - content_brief
  - missing_brief_inputs
  - reusable_preference_candidate
```

Purpose:

```text
Guide creation or refinement of a structured content brief: topic, angle, audience, intent, CTA, constraints, source needs, and SEO direction.
```

Must not:

```text
write final article
overwrite brief without approval
persist memory
render brief UI
```

### 15.3 `source_backed_research`

```yaml
skill_id: source_backed_research
name: source-backed-research
tier: read
risk_class: medium
compatible_profiles: [content_collaborator]
compatible_intents: [content_research, article_planning, topic_discovery, optimize_content]
allowed_modes: [plan, optimize]
capability_tags:
  - content_research
  - evidence_gathering
  - source_intelligence
  - claim_support
declared_action_classes:
  - web_research
  - deep_research
  - page_scrape
  - read_source_intelligence
output_contracts:
  - evidence_table
  - source_summary
  - unsupported_claims
```

Purpose:

```text
Guide evidence gathering, source triage, source note extraction, quote/stat handling, and claim-to-source mapping.
```

Must not:

```text
perform deep research without approval
scrape pages directly
store source databases
write final article sections
```

### 15.4 `seo_content_planning`

```yaml
skill_id: seo_content_planning
name: seo-content-planning
tier: read
risk_class: medium
compatible_profiles: [content_collaborator, seo_aeo_auditor]
compatible_intents: [seo_planning, article_planning, optimize_content, page_analysis]
allowed_modes: [plan, optimize, audit]
capability_tags:
  - seo_planning
  - search_intent
  - internal_linking
  - cannibalization_awareness
declared_action_classes:
  - web_research
  - site_page_discovery
  - page_scrape
  - cannibalization_check
  - update_brief
output_contracts:
  - seo_plan
  - internal_link_targets
  - cannibalization_warning
```

Purpose:

```text
Guide search-intent planning, primary/secondary keyword direction, internal-link targets, SERP-fit decisions, and cannibalization risk notes.
```

Must not:

```text
run site crawl directly
edit CMS fields
overwrite metadata
```

### 15.5 `page_content_analysis`

```yaml
skill_id: page_content_analysis
name: page-content-analysis
tier: read
risk_class: medium
compatible_profiles: [content_collaborator, seo_aeo_auditor]
compatible_intents: [page_analysis, optimize_content]
allowed_modes: [audit, optimize, plan]
capability_tags:
  - page_analysis
  - content_quality
  - seo_review
  - aeo_review
declared_action_classes:
  - page_scrape
  - web_research
  - read_brand_context
  - cannibalization_check
output_contracts:
  - page_analysis_report
  - improvement_opportunities
  - risk_summary
```

Purpose:

```text
Guide content-level analysis of a specific URL or supplied page text against brand, intent, SEO, AEO, structure, and quality expectations.
```

Must not:

```text
run full technical site audit
write edits directly
render dashboard
```

### 15.6 `article_planning`

```yaml
skill_id: article_planning
name: article-planning
tier: draft
risk_class: medium
compatible_profiles: [content_collaborator]
compatible_intents: [article_planning]
allowed_modes: [plan]
capability_tags:
  - article_planning
  - outline_building
  - section_blueprint
  - source_grounding
declared_action_classes:
  - read_brand_context
  - read_source_intelligence
  - web_research
  - deep_research
  - ask_user_question
  - update_brief
output_contracts:
  - article_blueprint
  - section_plan
  - source_requirement_map
```

Purpose:

```text
Guide creation of a section-by-section article plan from a finalized brief, audience, SEO direction, source constraints, and CTA.
```

Must not:

```text
draft full article
bypass missing brief fields
perform deep research without approval
```

### 15.7 `guarded_drafting`

```yaml
skill_id: guarded_drafting
name: guarded-drafting
tier: draft
risk_class: high
compatible_profiles: [content_execution_guard]
compatible_intents: [draft_article]
allowed_modes: [write]
capability_tags:
  - guarded_drafting
  - section_by_section_writing
  - tone_adaptation
  - source_grounded_writing
declared_action_classes:
  - read_memory_or_tone
  - read_brand_context
  - page_scrape
  - site_page_discovery
  - canvas_write
output_contracts:
  - draft_section
  - draft_article_chunk
  - post_write_quality_summary
```

Purpose:

```text
Guide section-by-section drafting in Write mode using an approved brief, article plan, sources, tone constraints, and quality gates.
```

Must not:

```text
draft without Write mode
invent sources
write beyond approved scope
publish externally
persist memory
```

### 15.8 `content_optimization`

```yaml
skill_id: content_optimization
name: content-optimization
tier: draft
risk_class: high
compatible_profiles: [content_execution_guard, content_collaborator]
compatible_intents: [optimize_content, edit_existing, page_analysis]
allowed_modes: [optimize, edit]
capability_tags:
  - targeted_editing
  - content_optimization
  - seo_support_outputs
  - quality_improvement
declared_action_classes:
  - page_scrape
  - web_research
  - site_page_discovery
  - cannibalization_check
  - canvas_edit
  - seo_support_update
output_contracts:
  - edit_proposal
  - optimization_summary
  - seo_support_fields
```

Purpose:

```text
Guide improvement of existing content, selected text, metadata suggestions, internal links, clarity, structure, and source freshness.
```

Must not:

```text
edit without selected text for edit_existing
overwrite SEO fields without approval
mutate canvas directly
```

### 15.9 `content_repurposing`

```yaml
skill_id: content_repurposing
name: content-repurposing
tier: draft
risk_class: medium
compatible_profiles: [content_execution_guard]
compatible_intents: [repurpose_content, switch_content_type]
allowed_modes: [edit, write]
capability_tags:
  - content_repurposing
  - format_adaptation
  - audience_adaptation
  - message_preservation
declared_action_classes:
  - set_content_type
  - read_memory_or_tone
  - update_brief
  - canvas_write
  - canvas_edit
output_contracts:
  - repurposed_content
  - format_change_summary
```

Purpose:

```text
Guide conversion of existing content into another format while preserving source constraints, audience fit, CTA, and brand voice.
```

Must not:

```text
change content type without explicit requested output format
publish externally
silently update tone preferences
```

### 15.10 `tone_memory_guidance`

```yaml
skill_id: tone_memory_guidance
name: tone-memory-guidance
tier: read
risk_class: high
compatible_profiles: [content_collaborator, content_execution_guard]
compatible_intents: [memory_tone_management, update_memory, content_brief, draft_article, optimize_content, repurpose_content]
allowed_modes: [plan, write, edit, optimize]
capability_tags:
  - tone_adaptation
  - memory_update_proposal
  - reusable_preference_detection
  - brand_voice
declared_action_classes:
  - read_memory_or_tone
  - propose_memory_update
output_contracts:
  - tone_guidance_summary
  - memory_update_proposal
```

Purpose:

```text
Guide how to use stored brand/tone preferences and how to identify reusable preferences worth proposing for memory update.
```

Must not:

```text
read memory directly
write memory directly
silently persist preference
expose memory file paths
```

### 15.11 `citation_source_safety`

```yaml
skill_id: citation_source_safety
name: citation-source-safety
tier: read
risk_class: medium
compatible_profiles: [content_collaborator, content_execution_guard, seo_aeo_auditor]
compatible_intents: [content_research, article_planning, topic_discovery, optimize_content]
allowed_modes: [plan, write, optimize]
capability_tags:
  - source_safety
  - citation_hygiene
  - competitor_blocklist
  - authority_filtering
declared_action_classes:
  - read_source_intelligence
  - web_research
  - page_scrape
output_contracts:
  - source_quality_matrix
  - citation_rules
  - competitor_exclusion_summary
```

Purpose:

```text
Guide source selection, competitor exclusion, authority-source preference, citation hygiene, and unsupported-claim detection.
```

Must not:

```text
fetch live sources directly
store competitor lists as secrets
override user citation policy
```

### 15.12 `content_quality_gates`

```yaml
skill_id: content_quality_gates
name: content-quality-gates
tier: read
risk_class: medium
compatible_profiles: [content_execution_guard, content_collaborator, seo_aeo_auditor]
compatible_intents: [draft_article, optimize_content, repurpose_content, page_analysis]
allowed_modes: [write, edit, optimize, audit]
capability_tags:
  - post_write_quality
  - readability
  - structure_validation
  - seo_quality
declared_action_classes:
  - canvas_edit
  - seo_support_update
output_contracts:
  - quality_gate_report
  - revision_checklist
```

Purpose:

```text
Guide post-write and post-edit checks for structure, readability, source support, internal links, metadata readiness, and content-specific requirements.
```

Must not:

```text
score evals
run automated tests directly
edit content directly
```

---

## 16. Existing SEO/AEO Skill Upgrades

Every existing SEO/AEO skill must be upgraded with:

```text
frontmatter v2
resource-manifest.yaml
eval-manifest.yaml
OWNERS
CHANGELOG.md
positive and negative trigger cases
action-class trajectory cases where tools may be involved
tier classification
risk classification
token budget
profile compatibility
mode compatibility
```

### 16.1 Required skill ownership

```yaml
ownership:
  keyword_research: seo
  serp_analysis: seo
  competitor_intelligence: seo_strategy
  content_strategy: content_strategy
  content_creation: content
  content_seo_settings: seo_content
  meta_optimization: seo_content
  schema_generation: technical_seo
  llms_txt_generation: technical_seo
  robots_txt_generation: technical_seo
  sitemap_generation: technical_seo
  site_audit_interpretation: technical_seo
  core_web_vitals_optimization: frontend_performance
  internal_linking_strategy: seo_architecture
  long_form_content_audit: content_quality
  content_refresh_strategy: content_strategy
  link_opportunity_discovery: authority_growth
  backlink_strategy: authority_growth
  gsc_insights_analysis: analytics_seo
  ga4_analysis: analytics
  local_seo_optimization: local_seo
  gbp_optimization: local_seo
```

Domain teams own the skill content. Platform/AI teams own schemas, CI, validation, packaging, and release mechanics.

---

## 17. Skill Quality Policy

### 17.1 Production readiness requirements

A skill can be `active` only if:

```text
SKILL.md validates
resource manifest validates
eval manifest validates
OWNERS exists
CHANGELOG exists
description passes clarity check
description includes anti-triggers
one skill one job check passes
token budget check passes
positive trigger cases pass
negative trigger cases pass
regression cases pass
collision tests pass
no protected metadata detected
no raw tool IDs in user-visible sections
no workflow IDs in user-visible sections
no profile IDs in user-visible sections
no MCP endpoints
no secrets
no raw memory document body
no raw canvas document body
scripts are checksummed
scripts request no network access by default
references stay within token limits
assets are not UI rendering definitions
```

### 17.2 Description quality checks

Reject or require revision if the description:

```text
starts with "helps with"
is purely aspirational
does not include trigger phrases
does not include anti-triggers
overlaps strongly with another skill
uses internal jargon
mentions exact internal tool names
mentions model/provider names
mentions MCP server URLs
is longer than policy limit
```

### 17.3 SKILL.md body size

```yaml
body_size_policy:
  preferred_max_words: 1500
  hard_max_words: 5000
  if_exceeds_preferred_max: require_references_split
  if_exceeds_hard_max: reject
```

### 17.4 Verifiability

Every rule in a production skill must be one of:

```text
directly observable in output
mapped to an output contract
mapped to an expected action class
mapped to a post-write checklist
mapped to an eval rubric
```

Reject vague rules such as:

```text
Make it good.
Use best practices.
Always optimize.
Be comprehensive.
Think deeply.
```

Unless they are translated into measurable criteria.

---

## 18. Token Budget Policy

### 18.1 Registry budget

```yaml
registry_budget:
  max_total_registry_tokens: 2500
  max_skill_metadata_tokens: 120
  max_description_tokens: 90
  max_anti_trigger_tokens: 60
```

### 18.2 Active skill budget

```yaml
active_skill_budget:
  max_active_skills_per_node: 4
  preferred_active_skills_per_node: 1-2
  max_total_skill_bundle_tokens: 9000
  max_single_skill_body_tokens: 3000
  max_single_reference_tokens: 2500
  max_assets_tokens_per_node: 2500
```

Layer 4 estimates token use. Layer 3 owns final context-window budgeting.

### 18.3 Degradation behavior

If budget is exceeded:

```text
omit optional references
omit optional assets
load summary variant if available
return token_budget_exceeded warning
never silently load over budget
never remove boundary rules
never remove failure behavior
never remove output contract
```

---

## 19. Skill Source and Supply-Chain Policy

### 19.1 Source categories

```yaml
source_trust:
  first_party:
    default_trust: high
    requirements:
      - version_pin
      - checksum
      - owner
      - evals
  organization_curated:
    default_trust: medium_high
    requirements:
      - PR review
      - owner
      - evals
      - checksum
      - changelog
  community:
    default_trust: low
    requirements:
      - manual audit
      - version pin
      - dependency scan
      - sandbox-only scripts
      - no production until evals pass
```

### 19.2 External skill adoption flow

```text
01 import into quarantine directory
02 scan for secrets and protected metadata
03 validate SKILL.md schema
04 validate resource manifest
05 inspect scripts manually
06 pin version and checksum
07 rewrite description to platform standards
08 add OWNER
09 add eval cases
10 run Layer 8 eval suite in shadow mode
11 promote to draft
12 promote to active only after gates pass
```

Layer 4 stores the candidate. Layer 8 runs evals. Layer 6 executes any scripts in sandbox if testing requires it.

---

## 20. Meta-Skill and Self-Improvement Policy

Layer 4 may store skills that help draft or improve skills, but meta-skill operation must be gated.

### 20.1 Allowed

```text
Agent proposes a new SKILL.md draft from a successful trace.
Agent proposes edits to descriptions.
Agent proposes references split from overlong SKILL.md.
Agent proposes additional trigger cases.
Agent proposes regression cases after user correction.
```

### 20.2 Forbidden

```text
Agent commits new skill directly to active registry.
Agent edits production skill without review.
Agent changes eval thresholds.
Agent removes negative trigger cases.
Agent relaxes boundary rules.
Agent deletes regression fixtures.
Agent promotes itself to active.
```

### 20.3 Promotion rule

All agent-authored or agent-edited skills enter:

```text
status: draft
tier: read unless explicitly reviewed otherwise
risk_class: high until reviewed
```

Promotion requires:

```text
human review
owner approval
schema validation
trigger evals
negative trigger evals
regression evals
token budget evals
protected metadata scan
```

---

## 21. Protected Metadata Policy

Layer 4 contains protected internal artifacts. It must not expose them to end users.

Protected metadata includes:

```text
exact skill inventory
SKILL.md body
skill file paths
resource paths
script names
script contents
eval fixtures
profile IDs
workflow IDs
internal action class mappings if user-facing
MCP server names or URLs
model/provider names
policy file contents
trace IDs
raw AgBOM
secrets
tokens
cookies
raw PII
memory file paths
raw memory document bodies
raw canvas document bodies
```

Layer 4 may provide a safe public capability summary:

```text
This system includes procedural guidance for technical SEO audits, content strategy, source-backed research, structured briefs, article planning, guarded drafting, content optimization, repurposing, schema support, internal linking, analytics interpretation, and local SEO.
```

It must not provide the exact registry unless the request is an authorized internal engineering request.

---

## 22. Skill Event Emission

Layer 4 emits sanitized events to Layer 8. It does not store telemetry.

```go
type SkillEvent struct {
    TraceID      string            `json:"trace_id,omitempty"`
    EventType    string            `json:"event_type"`
    SkillID      string            `json:"skill_id,omitempty"`
    SkillVersion string            `json:"skill_version,omitempty"`
    Decision     string            `json:"decision"`
    ReasonCode   string            `json:"reason_code,omitempty"`
    TokenEstimate int              `json:"token_estimate,omitempty"`
    Metadata     map[string]string `json:"metadata,omitempty"`
    CreatedAt    time.Time         `json:"created_at"`
}
```

Allowed event types:

```text
skill_registry_loaded
skill_registry_invalid
skill_package_validated
skill_package_rejected
skill_activation_requested
skill_loaded
skill_load_blocked
skill_resource_loaded
skill_resource_blocked
skill_checksum_failed
skill_token_budget_exceeded
skill_collision_detected
skill_deprecated_warning
skill_eval_manifest_loaded
```

Events must not include:

```text
raw SKILL.md body
raw reference content
raw asset content
script contents
raw user prompt
raw selected text
secrets
trace internals beyond trace_id if allowed by event bus
raw tool payloads
```

---

## 23. Integration Contracts with Other Layers

### 23.1 With Layer 2

Layer 2 sends no raw user input to Layer 4.

Layer 4 relies on Layer 2 for:

```text
sanitized intent
protected-disclosure handling
tool authorization
approval validation
outbound redaction
prompt-injection blocking
```

Layer 4 must not duplicate:

```text
LLM firewall
intent classifier
policy server
tool policy engine
outbound response guard
```

### 23.2 With Layer 3

Layer 3 is Layer 4's primary runtime caller.

Layer 3 sends:

```text
SkillActivationRequest
profile
mode
intent
workflow node
requested skills
token budget
output contracts
resource hints
```

Layer 4 returns:

```text
SkillActivationResponse
loaded SKILL.md bodies
approved reference/asset content
script handles
token estimates
warnings
```

Layer 4 must not:

```text
choose workflow
choose profile
build DAG
sequence nodes
create tool-call requests
create content-generation contracts
route output surfaces
```

### 23.3 With Layer 5

Layer 4 may provide output schemas or content templates as assets, but Layer 5 owns rendering.

Layer 4 must not store or produce:

```text
A2UI component JSON
UI card definitions
dashboard layouts
approval UI layouts
interactive canvas behavior
```

### 23.4 With Layer 6

Layer 4 stores script artifacts. Layer 6 executes.

Layer 4 must provide:

```text
script path handle
script checksum
script language
stdin/stdout contract
resource limits
network_access=false by default
```

Layer 4 must not:

```text
execute scripts
spawn processes
install packages
open files outside skill package validation
manage sandbox
manage egress
mint credentials
```

### 23.5 With Layer 7

Layer 4 declares abstract data needs and action classes. Layer 7 owns MCP, connectors, APIs, RAG, and memory retrieval.

Layer 4 must not:

```text
connect to MCP
list MCP tools
call APIs
run web search
scrape pages
read memory docs
query vector stores
manage tenant partitioning
```

### 23.6 With Layer 8

Layer 4 stores eval manifests and emits sanitized events. Layer 8 runs evals and stores results.

Layer 4 must not:

```text
score evals
run pass^k tests
store OpenTelemetry traces
track AgBOM
score drift
run red team simulations
trigger quarantine
```

---

## 24. Production Test Matrix

### 24.1 Registry tests

```text
registry loads with all active skills
duplicate skill_id rejected
missing directory rejected
missing SKILL.md rejected
missing OWNERS rejected
missing eval manifest rejected
checksum mismatch rejected
blocked skill cannot load
retired skill cannot load
draft skill cannot load in production
description over limit rejected
description without anti-triggers rejected
```

### 24.2 SKILL.md tests

```text
frontmatter parses
required fields exist
required Markdown sections exist
boundary section exists
body under token limit
no secrets
no raw prompts
no raw tool IDs in user-visible body
no MCP URLs
no workflow IDs as instructions
no profile IDs as instructions
no UI JSON
no memory file path
no raw memory body
```

### 24.3 Activation tests

```text
Layer 3 can request known compatible skill
unknown skill blocked
intent-incompatible skill blocked
mode-incompatible skill blocked
profile-incompatible skill blocked
status-incompatible skill blocked
token-over-budget activation blocked or degraded
optional references omitted when over budget
script handle returned without execution
```

### 24.4 Progressive disclosure tests

```text
registry metadata loads without SKILL.md bodies
body loads only when requested
references load only when referenced
assets load only when output contract matches
eval files never load during normal runtime
scripts never execute in Layer 4
inactive skill body never loads
```

### 24.5 Trigger tests

```text
each skill has 3 positive triggers
each skill has 3 negative triggers
each skill has rephrasing stability cases
adjacent skill collisions detected
topic_discovery does not trigger guarded_drafting
guarded_drafting does not trigger topic_discovery
content_optimization does not trigger content_brief_building
tone_memory_guidance does not silently trigger update_memory
```

### 24.6 Resource tests

```text
resource paths cannot escape skill directory
absolute paths rejected
script with network access rejected by default
script with package install rejected
asset with A2UI JSON rejected
reference over token budget rejected or summarized
resource checksum mismatch rejected
```

### 24.7 Eval manifest tests

```text
trigger cases exist
golden cases exist for draft and act tiers
trajectory cases use action classes, not exact tool IDs
read tier allows ANY_ORDER
draft tier requires IN_ORDER
act tier requires EXACT
rubric exists for draft and act tiers
regression cases exist before active promotion
```

### 24.8 Boundary tests

Layer 4 must prove it does not:

```text
normalize raw input
detect prompt injection
classify intent
authorize tools
choose workflow
assign agent
build DAG
create tool request
execute tool
execute script
render A2UI
render canvas
render brief
render chat
connect MCP
call connector
query vector store
read memory document
write memory document
mint JIT token
manage sandbox
manage network egress
store telemetry
score evals
track AgBOM
run SecOps triad
trigger quarantine
```

---

## 25. Deployment Requirements

Layer 4 v2 is production-ready only when:

```text
all v1 SEO/AEO skills upgraded to v2 schema
new content-agent skills added
skill-registry.yaml validates
all SKILL.md files validate
all resource manifests validate
all eval manifests validate
all OWNERS files exist
all checksums pass
all descriptions have anti-triggers
all production skills have trigger evals
all draft/action skills have golden evals
all action-class trajectories use action classes, not internal tool IDs
registry token budget passes
active skill bundle token budget passes
protected metadata scan passes
script artifact scan passes
external skills are pinned and audited
Layer 3 activation compatibility tests pass
Layer 2 boundary tests pass
Layer 5 boundary tests pass
Layer 6 boundary tests pass
Layer 7 boundary tests pass
Layer 8 boundary tests pass
```

### 25.1 Executable skill-library and eval-definition baseline

The repository implementation must retain:

```text
registry.schema.json, skill.schema.json, resource-manifest.schema.json, and eval-manifest.schema.json compile as JSON Schema draft 2020-12
skills.ValidateLibrary validates the mutable repository tree for CI
skills.ValidateEmbeddedLibrary validates the immutable library shipped in the API binary
registry entries have unique skill_id values and exact SHA-256 SKILL.md checksums
every SKILL.md frontmatter and required body section validates
every skill has non-empty OWNERS, CHANGELOG.md, resource-manifest.yaml, and eval-manifest.yaml
every skill has positive, negative, rephrasing, collision, out-of-scope, golden, trajectory, rubric, and regression definitions
skills.LoadEmbeddedTriggerEvalCorpus returns provider-safe metadata, exact expectations, and a corpus checksum without exposing file paths or resource bodies
cmd/api rejects startup when the embedded library fails validation
```

Minimum trigger corpus:

```text
at least 3 positive cases per skill
at least 3 negative cases per skill
at least one rephrase for each positive case
at least one collision case per skill
at least one out-of-scope case per skill
globally unique runtime case identifiers after skill and group qualification
```

Layer 4 validates and supplies definitions only. Layer 8 scores the corpus, and
Layer 7 owns any external model connection. Generated or templated cases are
not automatically accepted as sufficient coverage; collision quality and
domain realism require human review.

No skill may move from `experimental` to `active` merely because schema,
checksum, or static-corpus validation passes. Promotion requires completed live
trigger runs, applicable golden and trajectory evaluation, regression and
safety success, signed governance evidence, and the named human approval
required by Layer 8.

### 25.2 Repository readiness evidence

The repository-level production check is:

```text
go run ./cmd/readiness -root .
```

For Layer 4, the check must:

```text
enumerate skill directories
verify every skill is represented in skill-registry.yaml
reject an empty or placeholder registry
verify required SKILL.md frontmatter fields are present
verify required SKILL.md sections are present
enforce the hard SKILL.md body-size limit
verify required ownership, changelog, resource-manifest, eval-manifest, and eval fixture files are present and non-empty
reject missing, unreadable, or placeholder Layer 4 schemas
report explicit prototype markers in Layer 4 production sources
```

File presence is not an eval pass. A production skill still requires schema
validation, checksum validation, positive and negative trigger execution,
rephrasing stability, collision tests, golden output checks where applicable,
trajectory checks where applicable, regression tests, token-budget tests, and
Layer 8 scoring.

`cmd/readiness` and `internal/releasegate` are read-only platform CI tooling.
They must not activate skills, load skill bodies into runtime context, execute
scripts, score evals, change skill status, or perform any Layer 4 runtime
behavior.

---

## 26. Acceptance Criteria

Layer 4 v2 is accepted when:

1. It remains a pure progressive-disclosure skill directory.
2. It supports SEO/AEO auditor, content collaborator, and guarded content execution profiles.
3. It adds content-agent procedural skills without duplicating Layer 3 workflows.
4. It stores skill metadata, bodies, references, assets, scripts, and eval definitions.
5. It never executes scripts or tools.
6. It never connects MCP or APIs.
7. It never reads or writes memory documents.
8. It never renders UI or output surfaces.
9. It validates skill activation requests from Layer 3.
10. It returns skill bundles within token budget.
11. It enforces skill versioning, checksums, ownership, and release gates.
12. It implements progressive disclosure at registry, body, and resource levels.
13. It supports read/draft/act maturity tiers.
14. It supports content-agent modes without owning mode decisions.
15. It stores eval manifests for Layer 8 but does not score evals.
16. It emits sanitized events only.
17. It rejects protected metadata, secrets, unsafe scripts, and oversized bodies.
18. It passes all zero-overlap boundary tests.

---

## 27. Final Non-Goals

```text
Layer 4 must not classify user intent.
Layer 4 must not choose workflows.
Layer 4 must not assign agents.
Layer 4 must not build DAGs.
Layer 4 must not create tool-call requests.
Layer 4 must not authorize tools.
Layer 4 must not execute tools.
Layer 4 must not execute scripts.
Layer 4 must not manage runtime sandboxing.
Layer 4 must not manage filesystem mounts.
Layer 4 must not manage network egress.
Layer 4 must not mint or revoke credentials.
Layer 4 must not connect MCP servers.
Layer 4 must not call connectors.
Layer 4 must not perform RAG retrieval.
Layer 4 must not read memory documents.
Layer 4 must not write memory documents.
Layer 4 must not persist tone preferences.
Layer 4 must not render A2UI.
Layer 4 must not render approval cards.
Layer 4 must not render canvas.
Layer 4 must not render brief.
Layer 4 must not render chat.
Layer 4 must not publish externally.
Layer 4 must not store telemetry.
Layer 4 must not score evals.
Layer 4 must not track AgBOM.
Layer 4 must not score intent drift.
Layer 4 must not run red/blue/green SecOps loops.
Layer 4 must not quarantine runtime.
Layer 4 must not expose SKILL.md bodies, skill file paths, exact skill inventory, scripts, eval fixtures, workflow IDs, profile IDs, exact internal tool IDs, MCP endpoints, memory paths, trace internals, secrets, tokens, cookies, or raw PII to end users.
```

---

## 28. One-Line Architecture Summary

Layer 4 v2 is the production-grade progressive-disclosure procedural-memory layer that stores, validates, versions, packages, and serves skill bodies and resources on demand, while leaving intake, routing, execution, rendering, data access, memory persistence, telemetry, evaluation scoring, and recovery to their owning layers.
