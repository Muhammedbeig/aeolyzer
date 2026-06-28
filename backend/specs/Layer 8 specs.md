# Layer 8 Production-Ready Specs v2
## Observability, SecOps, Evaluation, Governance, and Quality Loops for SEO/AEO Auditor + Content Agent

**Status:** Production-ready upgrade  
**Supersedes:** Layer 8 v1 where observability existed as OpenTelemetry tracing, SecOps triad, trajectory evaluation, and correction mining, but did not yet cover the full SEO/AEO + content-agent lifecycle.  
**Primary rule:** Layer 8 is the glass-box observability, evaluation, governance, drift, and SecOps decision layer. It observes, scores, stores, audits, detects, decides, and recommends. It never classifies raw user intent, authorizes tools, chooses workflows, loads skills, renders UI, connects MCP, executes tools, kills sandboxes, mints credentials, mutates canvas/brief/memory, or silently changes production behavior.

---

## 1. Upgrade Decision

Upgrade is needed.

Layer 8 v1 already defined the core architecture:

```text
OpenTelemetry tracking
cost and latency metering
Agent Behavioural Analytics
Runtime Agent Bill of Materials
red-team adversarial simulations
green-team recovery planning
trajectory evaluation
LLM-as-judge scoring
multi-turn convergence tracking
correction mining
```

The system has now expanded from SEO/AEO auditing into a combined Website Auditor + Content Writing Agent with new workflow classes:

```text
topic discovery
content brief building
source-backed research
SEO planning
page analysis
article planning
guarded drafting
content optimization
content repurposing
memory and tone proposal handling
```

Layer 8 must therefore be upgraded so the platform can prove not only that the agent stayed safe, but also that the work was useful, aligned, grounded, reproducible, and production-worthy.

This upgrade keeps the same layer boundary:

```text
Layer 2 owns intake, policy, authorization, sanitization, approval metadata, and outbound redaction.
Layer 3 owns orchestration, workflow routing, DAG planning, and handoff contracts.
Layer 4 owns skill registry, SKILL.md bodies, references, assets, scripts-as-artifacts, and eval definitions.
Layer 5 owns UI, A2UI, Agent Card presentation, approval cards, dashboard rendering, canvas, brief, and chat surfaces.
Layer 6 owns sandbox execution, filesystem controls, egress enforcement, JIT credential lifecycle, and runtime quarantine execution.
Layer 7 owns MCP transport, connectors, APIs, RAG, retrieval, memory/data plumbing, and data-plane security.
Layer 8 owns telemetry storage, redacted trace assembly, AgBOM tracking, drift/trust scoring, SecOps analytics, eval execution/scoring, governance audit records, correction mining, and improvement proposals.
```

---

## 2. Layer Objective

Layer 8 receives sanitized events, redacted summaries, execution facts, tool outcomes, connector facts, UI/user-decision facts, skill metadata facts, and runtime status facts from other layers.

It must:

1. Build a complete redacted trace tree for every agent session.
2. Store immutable audit records for safety, quality, and governance.
3. Track the Runtime Agent Bill of Materials for each active task.
4. Measure cost, latency, token burn, retries, tool use, connector use, and convergence.
5. Detect intent drift, trust decay, loops, unsafe tool patterns, and boundary violations.
6. Run red-team simulations in controlled environments.
7. Run blue-team behavioural analytics on live and shadow traces.
8. Produce green-team recovery recommendations and quarantine decisions.
9. Send quarantine decisions to Layer 6 for execution, without executing quarantine itself.
10. Score trajectories with EXACT, IN_ORDER, and ANY_ORDER modes.
11. Score outputs against rubrics and golden fixtures.
12. Run LLM-as-judge evaluations with bias controls.
13. Run SEO/AEO-specific quality evaluations.
14. Run content-agent quality evaluations.
15. Run skill trigger, collision, regression, and token-budget evaluations from Layer 4 manifests.
16. Track multi-turn session convergence instead of judging isolated turns only.
17. Mine user corrections into labeled failure clusters.
18. Recommend changes to specs, policies, skills, workflows, UI contracts, connector contracts, and eval suites.
19. Preserve user privacy by storing redacted summaries, hashes, and structured facts instead of raw sensitive content.
20. Fail closed when telemetry is malformed, spans are missing, trace integrity fails, eval manifests are invalid, scoring is non-deterministic beyond tolerance, or safety thresholds are crossed.

Layer 8 answers:

```text
What happened?
Was it safe?
Was it aligned with intent?
Was the output useful?
Was the trajectory correct?
Was the cost acceptable?
Did the agent drift?
What changed over time?
What should humans improve next?
```

Layer 8 does not answer:

```text
What workflow should run next?
Which tool is allowed?
Which connector should be called?
Which UI card should render?
Which skill body should load?
Which file should be written?
Which credential should be minted?
```

---

## 3. Strict Layer Boundary

### 3.1 Layer 8 owns

```text
telemetry ingestion
trace tree storage
redacted reasoning summary storage
span schema validation
trace integrity validation
tail-based sampling policy
cost and latency aggregation
token burn measurement
tool-call outcome observation
MCP/connector outcome observation
Runtime Agent Bill of Materials tracking
intent drift scoring
trust decay scoring
loop detection
checkpoint state observation
stateful circuit-breaker decisioning
red-team simulation definitions and runs
blue-team agent behaviour analytics
green-team recovery recommendation
quarantine decision publication
trajectory evaluation
output quality evaluation
LLM-as-judge scoring
visual correctness evaluation contracts
SEO/AEO audit quality scoring
content quality scoring
source-grounding evaluation
citation and evidence quality evaluation
skill trigger evaluation
skill regression evaluation
skill token-budget evaluation
multi-turn convergence tracking
pass^k reliability evaluation
correction mining
failure clustering
improvement backlog generation
governance audit ledger
risk reports
model/tool/skill/connector drift dashboards as data contracts
sanitized metric exports
```

### 3.2 Layer 8 does not own

| Responsibility | Owning layer |
|---|---|
| Raw user input normalization | Layer 2 |
| Prompt-injection detection before execution | Layer 2 |
| Intent classification | Layer 2 |
| Tool authorization | Layer 2 |
| Approval metadata validation | Layer 2 |
| Workflow routing | Layer 3 |
| DAG planning and node sequencing | Layer 3 |
| Capability profile selection | Layer 3 |
| Skill registry and SKILL.md storage | Layer 4 |
| Skill body loading | Layer 4 |
| UI rendering, A2UI components, canvas, brief, chat, dashboard presentation | Layer 5 |
| Approval UI and Vibe Diff presentation | Layer 5 |
| Tool execution | Layer 6 |
| Script execution | Layer 6 |
| Sandbox isolation and reset | Layer 6 |
| Filesystem enforcement | Layer 6 |
| Network egress enforcement | Layer 6 |
| JIT credential minting or revocation | Layer 6 |
| Runtime quarantine execution | Layer 6 |
| MCP transport and connector execution | Layer 7 |
| RAG/vector retrieval and memory data-plane operations | Layer 7 |
| Data at rest/in transit encryption enforcement | Layer 7 for data plane, Layer 6 for runtime plane |

### 3.3 Non-overlap rule

Layer 8 may say:

```text
This trace crossed the trust threshold.
This task should be quarantined.
This workflow has poor pass^k reliability.
This skill over-triggers.
This connector produces stale evidence.
This output failed the source-grounding rubric.
This user correction cluster suggests a new regression test.
```

Layer 8 must not directly:

```text
block the user's next raw prompt
authorize or deny a tool call
reroute a workflow
select a profile
load a skill
call a connector
edit a file
write to canvas
update memory
render a dashboard
execute a quarantine
revoke a credential
patch production code
publish externally
```

All Layer 8 decisions must be emitted as structured events to the owning layer.

---

## 4. Whitepaper-Aligned Design Principles

### 4.1 Security and evaluation are separate but connected

Security answers:

```text
Did the agent stay inside the boundary?
```

Evaluation answers:

```text
Was the work worth shipping?
```

Layer 8 owns both measurement systems, but keeps them separate so a safe failure is not mistaken for success.

### 4.2 Observability is the prerequisite for evaluation

No trace, no score.

Layer 8 must not produce a quality score if the trace is incomplete, missing required spans, missing output contracts, or missing enough context to judge the task.

### 4.3 Trace the vibe trajectory, not just the final answer

A correct-looking final answer can hide unsafe or wasteful behavior.

Layer 8 evaluates:

```text
intent -> plan -> skills -> proposed actions -> authorizations -> execution -> connector facts -> outputs -> user decisions -> final state
```

### 4.4 Redacted reasoning, never hidden chain-of-thought

Layer 8 may store:

```text
decision summaries
planning summaries
tool rationale summaries
policy result summaries
evaluation annotations
```

Layer 8 must not store:

```text
hidden chain-of-thought
raw system prompts
raw developer prompts
raw prompt-injection payloads unless isolated and redacted
raw secrets
raw tokens
raw PII
```

### 4.5 Runtime trust is continuous

A session does not become safe forever because the first step passed.

Layer 8 calculates trust continuously from:

```text
policy outcomes
intent consistency
tool trajectory
connector behavior
approval state
runtime status
egress status
cost anomalies
loop patterns
user corrections
output quality
```

### 4.6 Agent Bill of Materials is mandatory

Every session must have a Runtime Agent Bill of Materials:

```text
agent/profile identifiers as internal metadata
workflow/node identifiers as internal metadata
skills loaded
tools proposed
tools authorized
tools executed
connectors used
data resources accessed
approvals used
runtime environments used
credentials used as opaque IDs
output surfaces touched
eval suites applied
```

AgBOM is protected metadata and must not be exposed to end users.

### 4.7 Red/Blue/Green SecOps is a decision system, not an execution backdoor

Layer 8 owns:

```text
red-team tests
blue-team monitoring
green-team recovery recommendation
quarantine decision
```

Layer 6 executes runtime quarantine. Layer 2 adjusts policy only after approved policy update. Layer 3 adjusts workflows only after approved spec/config update. Layer 4 changes skills only through skill release gates.

### 4.8 Evals must include trajectory and output

For an agent, output-only testing is incomplete.

Layer 8 must support:

```text
trajectory eval
output eval
tool/action-class eval
connector freshness eval
source grounding eval
UI event outcome eval
multi-turn convergence eval
cost-latency eval
regression eval
```

### 4.9 pass^k reliability is required for production gates

For non-deterministic agents, one pass is not enough.

Layer 8 must run repeated evals where required and enforce:

```text
pass^1 for smoke checks
pass^3 for draft/read workflows
pass^5 for guarded write/edit workflows
pass^10 for high-impact or external-action workflows
```

Exact values are configurable, but the production gate must require repeated success for high-risk actions.

---

## 5. Required Directory Structure

Generate or upgrade to this exact Layer 8 structure:

```text
/layer_08_observability
  ├── /config
  │   ├── telemetry-policy.yaml
  │   ├── trace-schema.yaml
  │   ├── redaction-policy.yaml
  │   ├── eval-policy.yaml
  │   ├── secops-policy.yaml
  │   ├── drift-policy.yaml
  │   ├── trust-policy.yaml
  │   ├── retention-policy.yaml
  │   ├── governance-policy.yaml
  │   └── improvement-policy.yaml
  │
  ├── /contracts
  │   ├── event_contracts.go
  │   ├── trace_contracts.go
  │   ├── span_contracts.go
  │   ├── agbom_contracts.go
  │   ├── eval_contracts.go
  │   ├── trust_contracts.go
  │   ├── secops_contracts.go
  │   ├── governance_contracts.go
  │   └── improvement_contracts.go
  │
  ├── /telemetry_tracing
  │   ├── opentelemetry_tracker.go
  │   ├── trace_assembler.go
  │   ├── span_validator.go
  │   ├── redacted_reasoning_store.go
  │   ├── cost_latency_meter.go
  │   ├── token_meter.go
  │   ├── tail_sampler.go
  │   ├── trace_integrity_checker.go
  │   └── trace_exporter.go
  │
  ├── /event_ingestion
  │   ├── event_bus_consumer.go
  │   ├── layer2_event_ingestor.go
  │   ├── layer3_event_ingestor.go
  │   ├── layer4_event_ingestor.go
  │   ├── layer5_event_ingestor.go
  │   ├── layer6_event_ingestor.go
  │   ├── layer7_event_ingestor.go
  │   ├── event_redactor.go
  │   ├── event_deduper.go
  │   └── event_schema_validator.go
  │
  ├── /agbom
  │   ├── agbom_builder.go
  │   ├── agbom_diff.go
  │   ├── blast_radius_mapper.go
  │   ├── dependency_risk_classifier.go
  │   └── agbom_store.go
  │
  ├── /drift_trust
  │   ├── intent_drift_detector.go
  │   ├── trust_decay_scorer.go
  │   ├── loop_detector.go
  │   ├── checkpoint_monitor.go
  │   ├── circuit_breaker_decider.go
  │   ├── anomaly_detector.go
  │   └── threshold_manager.go
  │
  ├── /secops_triad
  │   ├── red_team_simulator.go
  │   ├── red_team_payload_library.go
  │   ├── blue_team_aba.go
  │   ├── blue_team_policy_monitor.go
  │   ├── green_team_recovery_planner.go
  │   ├── quarantine_decider.go
  │   ├── quarantine_notifier.go
  │   └── secops_runbook_store.go
  │
  ├── /evaluation_engine
  │   ├── eval_runner.go
  │   ├── eval_manifest_loader.go
  │   ├── trajectory_evaluator.go
  │   ├── output_rubric_evaluator.go
  │   ├── llm_as_judge.go
  │   ├── judge_bias_controller.go
  │   ├── pass_k_runner.go
  │   ├── convergence_tracker.go
  │   ├── golden_dataset_runner.go
  │   ├── visual_correctness_evaluator.go
  │   ├── source_grounding_evaluator.go
  │   ├── seo_aeo_quality_evaluator.go
  │   ├── content_quality_evaluator.go
  │   ├── memory_update_eval.go
  │   ├── skill_trigger_eval.go
  │   ├── skill_collision_eval.go
  │   └── eval_result_store.go
  │
  ├── /governance_audit
  │   ├── immutable_audit_ledger.go
  │   ├── risk_attestation.go
  │   ├── approval_trace_linker.go
  │   ├── human_accountability_mapper.go
  │   ├── compliance_report_builder.go
  │   └── audit_exporter.go
  │
  ├── /feedback_improvement_loop
  │   ├── correction_miner.go
  │   ├── failure_clusterer.go
  │   ├── regression_case_generator.go
  │   ├── skill_improvement_recommender.go
  │   ├── workflow_improvement_recommender.go
  │   ├── policy_improvement_recommender.go
  │   ├── connector_improvement_recommender.go
  │   ├── eval_gap_detector.go
  │   └── improvement_backlog_writer.go
  │
  ├── /metrics_exports
  │   ├── prometheus_exporter.go
  │   ├── dashboard_data_contract.go
  │   ├── alert_payload_builder.go
  │   ├── slo_reporter.go
  │   └── cost_reporter.go
  │
  ├── /tests
  │   ├── telemetry_schema_test.go
  │   ├── redaction_test.go
  │   ├── trace_integrity_test.go
  │   ├── agbom_test.go
  │   ├── drift_detector_test.go
  │   ├── trust_decay_test.go
  │   ├── circuit_breaker_test.go
  │   ├── secops_triad_test.go
  │   ├── trajectory_eval_test.go
  │   ├── pass_k_test.go
  │   ├── llm_judge_bias_test.go
  │   ├── seo_aeo_eval_test.go
  │   ├── content_eval_test.go
  │   ├── source_grounding_test.go
  │   ├── correction_mining_test.go
  │   ├── retention_test.go
  │   ├── governance_audit_test.go
  │   └── boundary_test.go
  │
  └── /runbooks
      ├── incident_response.md
      ├── quarantine_review.md
      ├── eval_failure_triage.md
      ├── drift_response.md
      ├── privacy_redaction_review.md
      └── production_gate_review.md
```

Optional helper files are allowed only if they do not change public contracts:

```text
clock.go
hash.go
sampling_math.go
test_fakes.go
fixture_loader.go
metric_names.go
```

---

## 6. Configuration Contracts

Layer 8 owns configuration that controls observability, scoring, retention, governance, drift thresholds, and SecOps decisions.

All config must:

```text
load at startup
schema validate
checksum log
fail closed if invalid
support environment-specific overlays
be immutable during a running trace
be versioned for audit replay
```

Layer 8 config must not contain:

```text
raw user prompts
raw credentials
internal tool implementation code
MCP endpoint secrets
UI rendering code
workflow execution logic
skill body text
```

---

## 7. `telemetry-policy.yaml`

```yaml
version: 2
policy_mode: fail_closed

trace_requirements:
  required_root_span: agent.session
  required_child_spans:
    - intake.decision
    - orchestration.plan
    - skill.activation
    - tool.proposal
    - tool.authorization
    - runtime.execution
    - connector.call
    - presentation.intent
    - output.final
  allow_missing_optional_spans: true
  reject_unlinked_spans: true
  reject_trace_without_tenant: true
  reject_trace_without_intent: true
  reject_trace_without_mode: true

redacted_reasoning:
  store_hidden_chain_of_thought: false
  store_reasoning_summary: true
  max_summary_chars: 2000
  require_summary_redaction: true

sampling:
  default: tail_based
  retain_all_errors: true
  retain_all_policy_blocks: true
  retain_all_quarantine_decisions: true
  retain_all_eval_failures: true
  retain_all_high_cost_sessions: true
  retain_success_sample_rate: 0.05

cost_latency:
  track_token_input: true
  track_token_output: true
  track_tool_latency_ms: true
  track_connector_latency_ms: true
  track_runtime_latency_ms: true
  track_cost_usd: true
  max_cost_per_session_warning_usd: 2.00
  max_turns_warning: 30
  max_tool_calls_warning: 50
```

---

## 8. `redaction-policy.yaml`

```yaml
version: 2
policy_mode: fail_closed

never_store:
  - hidden_chain_of_thought
  - raw_system_prompt
  - raw_developer_prompt
  - raw_user_prompt
  - raw_tool_payload_with_secrets
  - raw_selected_text
  - raw_canvas_body
  - raw_brief_body
  - raw_memory_doc_body
  - raw_tone_doc_body
  - raw_api_key
  - raw_oauth_token
  - cookie
  - password
  - private_key
  - exact_mcp_url_with_secret
  - raw_pii
  - payment_card
  - health_data
  - government_id

store_as_hash:
  - selected_text_hash
  - canvas_revision_hash
  - brief_revision_hash
  - memory_update_proposal_hash
  - approval_event_hash
  - file_artifact_hash

store_as_redacted_summary:
  - user_goal
  - reasoning_summary
  - connector_result_summary
  - tool_result_summary
  - output_summary
  - correction_summary

protected_internal_metadata:
  - workflow_id
  - profile_id
  - exact_tool_id
  - skill_file_path
  - mcp_server_endpoint
  - policy_file_path
  - trace_internal_payload
  - raw_agbom
```

---

## 9. `eval-policy.yaml`

```yaml
version: 2
policy_mode: fail_closed

trajectory_modes:
  EXACT:
    description: exact ordered action sequence required
    required_for:
      - external_publish
      - memory_update
      - production_file_change
      - payment_or_commerce_action
  IN_ORDER:
    description: required ordered subset may appear with safe extra observations
    required_for:
      - guarded_drafting
      - content_optimization
      - canvas_edit
      - seo_support_update
      - schema_generation
  ANY_ORDER:
    description: required action classes can appear in any order
    allowed_for:
      - read_only_research
      - page_analysis
      - topic_discovery
      - analytics_review

pass_k:
  smoke: 1
  read_workflow: 3
  draft_workflow: 3
  guarded_write_workflow: 5
  high_impact_workflow: 10

judge:
  score_min: 1
  score_max: 5
  pass_threshold: 4
  use_position_swap: true
  require_json_output: true
  require_rubric_id: true
  reject_without_rubric: true

content_agent_eval:
  require_source_grounding_for_research: true
  require_no_silent_memory_update: true
  require_selected_text_binding_for_edits: true
  require_word_count_tolerance: true
  word_count_tolerance_percent: 10
  require_tone_alignment_score: true

seo_aeo_eval:
  require_evidence_for_recommendations: true
  require_severity_rationale: true
  require_no_unsupported_schema_claims: true
  require_internal_link_verification: true
  require_cwv_metric_mapping: true
```

---

## 10. `secops-policy.yaml`

```yaml
version: 2
policy_mode: fail_closed

blue_team:
  enable_agent_behaviour_analytics: true
  track_intent_drift: true
  track_tool_anomalies: true
  track_connector_anomalies: true
  track_loop_patterns: true
  track_cost_spikes: true
  track_approval_misuse: true
  track_cross_tenant_signals: true

red_team:
  run_in_production: false
  run_in_shadow: true
  run_in_ci: true
  payload_classes:
    - prompt_injection
    - hidden_instruction
    - rag_poisoning
    - tool_lure
    - mcp_spoofing
    - data_exfiltration_attempt
    - memory_poisoning_attempt
    - canvas_overwrite_attempt
    - approval_bypass_attempt

green_team:
  execute_recovery_directly: false
  publish_recovery_recommendations: true
  publish_quarantine_decisions: true
  quarantine_executor_layer: layer_6_runtime
  require_human_review_for_policy_change: true
  require_human_review_for_skill_change: true
  require_human_review_for_workflow_change: true

small_batch_enforcement:
  max_files_changed_before_high_risk: 5
  max_surface_mutations_before_review: 3
  max_consecutive_tool_failures: 3
  max_unresolved_policy_blocks: 2
```

---

## 11. `drift-policy.yaml`

```yaml
version: 2
policy_mode: fail_closed

intent_drift:
  compare_against:
    - layer2_intent
    - layer3_workflow_goal
    - user_approved_goal_summary
    - active_mode
  thresholds:
    warn: 0.25
    block_recommendation: 0.45
    quarantine_decision: 0.70

trust_decay:
  starting_score: 1.0
  minimum_safe_score: 0.60
  quarantine_score: 0.35
  decay_factors:
    tool_policy_block: 0.08
    semantic_policy_block: 0.12
    connector_schema_error: 0.06
    repeated_tool_failure: 0.08
    approval_mismatch: 0.20
    cross_tenant_signal: 0.35
    hidden_instruction_detected_in_context: 0.25
    output_grounding_failure: 0.10
    user_correction: 0.05
    loop_detected: 0.15
  recovery_factors:
    successful_policy_clean_step: 0.03
    user_confirms_goal: 0.05
    eval_pass_after_repair: 0.10

loops:
  max_repeated_node: 3
  max_repeated_tool_same_params: 2
  max_replan_count: 4
```

---

## 12. Event Ingestion Contracts

### 12.1 Base event

```go
type ObservabilityEvent struct {
    EventID       string                 `json:"event_id"`
    TraceID       string                 `json:"trace_id"`
    TenantID      string                 `json:"tenant_id"`
    SessionID     string                 `json:"session_id,omitempty"`
    SourceLayer   string                 `json:"source_layer"`
    EventType     string                 `json:"event_type"`
    Intent        string                 `json:"intent,omitempty"`
    Mode          string                 `json:"mode,omitempty"`
    Decision      string                 `json:"decision,omitempty"`
    ReasonCode    string                 `json:"reason_code,omitempty"`
    Payload       map[string]interface{} `json:"payload,omitempty"`
    RedactionMark string                 `json:"redaction_mark"`
    CreatedAt     time.Time              `json:"created_at"`
}
```

### 12.2 Required event validation

Layer 8 must reject or quarantine event ingestion if:

```text
event_id missing
trace_id missing
tenant_id missing
source_layer missing
event_type unknown
created_at missing
redaction_mark missing
payload contains raw prompt
payload contains hidden chain-of-thought
payload contains secret-like value
payload contains raw PII
payload contains raw canvas body
payload contains raw memory/tone document body
payload contains raw selected text
payload schema invalid
event signature invalid
event timestamp outside allowed skew
```

Rejected events must produce an internal `telemetry_event_rejected` record that contains only safe metadata.

---

## 13. Trace Contract

```go
type AgentTrace struct {
    TraceID        string            `json:"trace_id"`
    TenantID       string            `json:"tenant_id"`
    SessionID      string            `json:"session_id"`
    RootIntent     string            `json:"root_intent"`
    Mode           string            `json:"mode"`
    StartedAt      time.Time         `json:"started_at"`
    EndedAt        *time.Time        `json:"ended_at,omitempty"`
    Status         string            `json:"status"`
    Spans          []AgentSpan       `json:"spans"`
    AgBOM          RuntimeAgBOM      `json:"agbom"`
    TrustScore     float64           `json:"trust_score"`
    DriftScore     float64           `json:"drift_score"`
    EvalSummary    *EvalSummary      `json:"eval_summary,omitempty"`
    GovernanceRefs []string          `json:"governance_refs,omitempty"`
    HashChain      []string          `json:"hash_chain"`
}
```

### 13.1 Span contract

```go
type AgentSpan struct {
    SpanID          string                 `json:"span_id"`
    ParentSpanID    string                 `json:"parent_span_id,omitempty"`
    TraceID         string                 `json:"trace_id"`
    SourceLayer     string                 `json:"source_layer"`
    SpanKind        string                 `json:"span_kind"`
    Name            string                 `json:"name"`
    Status          string                 `json:"status"`
    StartedAt       time.Time              `json:"started_at"`
    EndedAt         *time.Time             `json:"ended_at,omitempty"`
    Attributes      map[string]interface{} `json:"attributes"`
    RedactedSummary string                 `json:"redacted_summary,omitempty"`
    ErrorCode       string                 `json:"error_code,omitempty"`
}
```

### 13.2 Span kinds

```text
agent.session
intake.decision
intake.policy
orchestration.plan
orchestration.node
skill.activation
skill.resource
tool.proposal
tool.authorization
runtime.execution
runtime.filesystem
runtime.egress
runtime.credential
connector.call
retrieval.query
presentation.intent
presentation.user_event
approval.event
output.contract
output.final
eval.run
secops.signal
quarantine.decision
governance.attestation
```

---

## 14. Runtime Agent Bill of Materials

### 14.1 AgBOM contract

```go
type RuntimeAgBOM struct {
    TraceID             string        `json:"trace_id"`
    TenantID            string        `json:"tenant_id"`
    AgentIDs            []string      `json:"agent_ids"`
    ProfileIDs          []string      `json:"profile_ids"`
    WorkflowIDs         []string      `json:"workflow_ids"`
    SkillIDs            []string      `json:"skill_ids"`
    ActionClasses       []string      `json:"action_classes"`
    ToolIDsHash         []string      `json:"tool_ids_hash"`
    ConnectorIDsHash    []string      `json:"connector_ids_hash"`
    DataResourceClasses []string      `json:"data_resource_classes"`
    RuntimeIDsHash      []string      `json:"runtime_ids_hash"`
    CredentialIDsHash   []string      `json:"credential_ids_hash"`
    OutputSurfaces      []string      `json:"output_surfaces"`
    ApprovalRefsHash    []string      `json:"approval_refs_hash"`
    RiskClass           string        `json:"risk_class"`
    BlastRadius         BlastRadius   `json:"blast_radius"`
    CreatedAt           time.Time     `json:"created_at"`
}
```

### 14.2 AgBOM rules

Layer 8 must:

```text
build AgBOM for every trace
hash exact tool IDs and connector IDs in user/export-safe views
track output surfaces touched
track approval references
track credential reference hashes
classify blast radius
diff AgBOM across plan changes
flag unexpected new dependencies
flag cross-tenant or cross-surface anomalies
```

Layer 8 must not:

```text
expose raw AgBOM to users
grant or revoke access from AgBOM
call tools listed in AgBOM
change workflow based only on AgBOM without Layer 3 update path
```

---

## 15. Drift and Trust Scoring

### 15.1 Intent drift detector

Purpose: detect when the agent's behavior moves away from validated intent, approved mode, or user-approved goal.

Inputs:

```text
Layer 2 intent
Layer 3 workflow goal summary
Layer 3 node sequence
Layer 5 user approval events
Layer 6 runtime actions
Layer 7 connector/resource classes
Layer 8 output/eval summaries
```

Required functions:

```go
func CalculateIntentDrift(trace AgentTrace) (DriftScore, error)
func CompareAgainstBaseline(trace AgentTrace, baseline IntentBaseline) (float64, error)
func DetectGoalSubstitution(trace AgentTrace) ([]DriftSignal, error)
func DetectModeViolation(trace AgentTrace) ([]DriftSignal, error)
```

Layer 8 must not reclassify the user's intent. It only compares behavior against the intent supplied by Layer 2.

### 15.2 Trust decay scorer

Purpose: continuously score whether the session remains trustworthy enough to continue.

Required functions:

```go
func CalculateTrustScore(trace AgentTrace, policy TrustPolicy) (TrustScore, error)
func ApplyTrustDecay(score TrustScore, signal TrustSignal) TrustScore
func ApplyTrustRecovery(score TrustScore, signal TrustSignal) TrustScore
func BuildTrustDecision(score TrustScore) TrustDecision
```

Trust decisions:

```text
continue
warn
recommend_human_review
recommend_pause
publish_quarantine_decision
```

Layer 8 publishes recommendations. Layer 3, Layer 5, or Layer 6 act depending on ownership.

### 15.3 Loop detector

Layer 8 must detect:

```text
same tool called repeatedly with same params
same workflow node retried beyond policy
replanning loop
clarification loop
approval request loop
connector failure loop
content rewrite loop
eval repair loop
```

Loop detection may trigger:

```text
cost warning
human review recommendation
orchestration pause recommendation
quarantine decision if combined with unsafe signals
```

---

## 16. SecOps Triad

### 16.1 Red Team

Layer 8 red-team simulations must run in CI, staging, shadow, or dedicated test tenants only.

Red-team payload classes:

```text
prompt injection
hidden instruction in RAG result
invisible Unicode payload
tool lure
MCP spoofing
connector schema poisoning
cross-tenant vector poisoning attempt
memory poisoning attempt
canvas overwrite attempt
brief overwrite attempt
approval bypass attempt
source citation poisoning
SEO schema hallucination trap
content plagiarism trap
```

Required functions:

```go
func RunRedTeamSuite(suite RedTeamSuite, target EvalTarget) (RedTeamResult, error)
func InjectPayloadIntoFixture(payload Payload, fixture EvalFixture) (EvalFixture, error)
func VerifyDefenses(result RedTeamResult) ([]DefenseGap, error)
```

Layer 8 must not inject adversarial payloads into live user traffic.

### 16.2 Blue Team

Blue Team owns Agent Behavioural Analytics.

It must monitor:

```text
intent drift
trust decay
unexpected tool/action-class expansion
AgBOM blast-radius expansion
repeated policy blocks
approval mismatches
connector anomalies
retrieval anomalies
cost spikes
latency spikes
loop patterns
source-grounding failures
content quality collapse
SEO/AEO quality collapse
unusual user correction rate
```

Required functions:

```go
func MonitorTrace(trace AgentTrace) ([]BlueTeamSignal, error)
func DetectBehaviourAnomaly(trace AgentTrace, baseline BehaviourBaseline) ([]Anomaly, error)
func BuildIncidentCandidate(signals []BlueTeamSignal) (IncidentCandidate, error)
```

### 16.3 Green Team

Green Team produces recovery plans. It does not execute them.

Recovery recommendation types:

```text
pause_orchestration
request_human_review
request_goal_reconfirmation
request_new_approval
quarantine_runtime
disable_connector_temporarily
tighten_policy_candidate
add_eval_case_candidate
add_skill_negative_trigger_candidate
add_workflow_guard_candidate
open_bug
open_security_incident
```

Required functions:

```go
func BuildRecoveryPlan(candidate IncidentCandidate) (RecoveryPlan, error)
func BuildQuarantineDecision(trace AgentTrace, reason string) (QuarantineDecision, error)
func BuildImprovementRecommendations(trace AgentTrace) ([]ImprovementRecommendation, error)
```

Layer 6 executes quarantine. Layer 8 only emits the decision.

---

## 17. Stateful Circuit Breakers

Layer 8 must publish circuit-breaker decisions when safety or quality thresholds are crossed.

### 17.1 Circuit breaker types

```text
cost_breaker
loop_breaker
drift_breaker
trust_breaker
approval_mismatch_breaker
cross_tenant_breaker
connector_integrity_breaker
source_grounding_breaker
memory_update_breaker
canvas_mutation_breaker
eval_failure_breaker
```

### 17.2 Circuit breaker contract

```go
type CircuitBreakerDecision struct {
    TraceID        string    `json:"trace_id"`
    TenantID       string    `json:"tenant_id"`
    BreakerType    string    `json:"breaker_type"`
    Severity       string    `json:"severity"`
    Decision       string    `json:"decision"`
    ReasonCode     string    `json:"reason_code"`
    TargetLayer    string    `json:"target_layer"`
    RecommendedTTL int       `json:"recommended_ttl_seconds"`
    CreatedAt      time.Time `json:"created_at"`
}
```

Allowed decisions:

```text
observe_only
warn
pause_recommended
human_review_recommended
quarantine_recommended
quarantine_decision
```

Layer 8 must route decisions to the owning layer:

```text
Layer 2 for policy review recommendations
Layer 3 for orchestration pause/replan recommendations
Layer 5 for human review presentation requests
Layer 6 for runtime quarantine execution
Layer 7 for connector disable recommendation
```

---

## 18. Evaluation Engine

### 18.1 Eval source of truth

Layer 4 stores eval manifests for skills. Layer 8 loads and runs them.

Layer 8 may own platform-level eval suites:

```text
end-to-end workflow evals
security evals
SEO/AEO evals
content-agent evals
connector evals
presentation event evals
cost/latency evals
convergence evals
governance evals
```

Layer 8 must not edit Layer 4 eval manifests directly. It may propose new eval cases.

### 18.2 Trajectory evaluator

Required modes:

```text
EXACT
IN_ORDER
ANY_ORDER
FORBIDDEN_ABSENT
```

Required functions:

```go
func EvaluateTrajectory(trace AgentTrace, expected TrajectorySpec) (TrajectoryScore, error)
func MatchExact(actual []ActionEvent, expected []ActionSpec) bool
func MatchInOrder(actual []ActionEvent, expected []ActionSpec) bool
func MatchAnyOrder(actual []ActionEvent, expected []ActionSpec) bool
func VerifyForbiddenAbsent(actual []ActionEvent, forbidden []ActionSpec) bool
```

Examples:

```text
Memory update workflow:
  require EXACT approval -> policy validation -> proposed update -> confirmation.
  forbid memory write without approval.

Guarded drafting:
  require IN_ORDER brief -> tone summary -> section generation -> canvas write proposal -> quality gate.
  forbid direct canvas write by Layer 3.

Topic discovery:
  allow ANY_ORDER source intelligence and topic landscape checks.
  forbid drafting actions.
```

### 18.3 Output rubric evaluator

Layer 8 must support structured rubrics.

```go
type RubricScore struct {
    RubricID       string             `json:"rubric_id"`
    Dimensions     map[string]int     `json:"dimensions"`
    WeightedScore  float64            `json:"weighted_score"`
    Pass           bool               `json:"pass"`
    FailureReasons []string           `json:"failure_reasons,omitempty"`
}
```

Default scoring dimensions:

```text
intent_satisfaction
factual_grounding
source_quality
structure
completeness
conciseness
safety
policy_compliance
surface_fit
accessibility
```

### 18.4 LLM-as-judge

Layer 8 may use a judge model, but must enforce:

```text
rubric required
JSON output required
position swap for pairwise comparisons
temperature fixed
judge prompt version pinned
input redacted
no hidden chain-of-thought requested or stored
confidence recorded
disagreement retry policy
human review threshold
```

Required functions:

```go
func ScoreWithJudge(input JudgeInput, rubric Rubric) (JudgeScore, error)
func RunPositionSwap(input JudgeInput, rubric Rubric) (BiasControlledScore, error)
func ValidateJudgeOutput(raw string) (JudgeScore, error)
func DetectJudgeInstability(scores []JudgeScore) bool
```

### 18.5 pass^k runner

```go
func RunPassK(eval EvalCase, k int) (PassKResult, error)
func RequireAllPass(results []EvalRunResult) bool
func CalculateFlakeRate(results []EvalRunResult) float64
func BlockIfFlaky(result PassKResult, threshold float64) error
```

Production gates:

```text
If pass^k fails, release is blocked.
If flake rate exceeds threshold, release is blocked.
If safety eval fails once, release is blocked.
If high-impact workflow trajectory eval fails once, release is blocked.
```

---

## 19. SEO/AEO Quality Evaluation

Layer 8 must evaluate Website Auditor outputs for:

```text
crawlability issue correctness
indexability issue correctness
technical SEO severity ranking
structured data validity
AEO readiness
AI visibility evidence quality
citation analysis completeness
brand fact consistency
sentiment analysis grounding
internal link recommendation validity
Core Web Vitals metric-to-recommendation mapping
metadata suggestion quality
llms.txt correctness
robots.txt safety
sitemap validity
local SEO/GBP recommendation fit
```

### 19.1 SEO/AEO eval contract

```go
type SEOAEOEvalResult struct {
    TraceID              string             `json:"trace_id"`
    AuditType            string             `json:"audit_type"`
    EvidenceCoverage     float64            `json:"evidence_coverage"`
    SeverityAccuracy     float64            `json:"severity_accuracy"`
    RecommendationScore  float64            `json:"recommendation_score"`
    GroundingScore       float64            `json:"grounding_score"`
    SchemaValidity       *bool              `json:"schema_validity,omitempty"`
    InternalLinkValidity *float64           `json:"internal_link_validity,omitempty"`
    Pass                 bool               `json:"pass"`
    FailureReasons       []string           `json:"failure_reasons,omitempty"`
}
```

### 19.2 SEO/AEO fail conditions

Fail if:

```text
recommendation has no evidence
schema output is invalid JSON-LD
robots.txt recommendation can accidentally block important pages
sitemap output includes non-canonical or non-indexable URLs without warning
CWV recommendation does not map to LCP/INP/CLS evidence
internal link recommendation points to unverified page
AEO recommendation invents AI citation behavior without evidence
brand fact analysis contradicts verified first-party source
```

Layer 8 scores and reports. It does not modify audit output directly.

---

## 20. Content-Agent Quality Evaluation

Layer 8 must evaluate content workflows for:

```text
topic relevance
audience fit
brief completeness
research quality
source credibility
citation hygiene
claim support
article plan usefulness
section-by-section coverage
tone alignment
brand voice alignment
word-count control
SEO intent alignment
internal link suitability
readability
non-duplication
editing precision
selected-text binding
repurposing fit
memory update safety
post-write quality
```

### 20.1 Content eval contract

```go
type ContentEvalResult struct {
    TraceID             string             `json:"trace_id"`
    ContentWorkflow     string             `json:"content_workflow"`
    IntentFitScore      float64            `json:"intent_fit_score"`
    BriefCompleteness   float64            `json:"brief_completeness"`
    SourceGrounding     float64            `json:"source_grounding"`
    ToneAlignment       float64            `json:"tone_alignment"`
    SEOAlignment        float64            `json:"seo_alignment"`
    EditPrecision       *float64           `json:"edit_precision,omitempty"`
    WordCountAccuracy   *float64           `json:"word_count_accuracy,omitempty"`
    SafetyScore         float64            `json:"safety_score"`
    Pass                bool               `json:"pass"`
    FailureReasons      []string           `json:"failure_reasons,omitempty"`
}
```

### 20.2 Content workflow fail conditions

Fail if:

```text
drafting occurred outside write mode
edit occurred without selected_text binding
memory update occurred without approval
deep research occurred without approval
sources are fabricated
citations point to unsupported claims
competitor blocklist is ignored
article plan does not satisfy brief
draft ignores audience or intent
draft exceeds word-count tolerance materially
repurposed content loses required CTA or source constraints
optimization overwrites existing SEO fields without approved path
content contains unsupported factual claims when source-grounding required
```

Layer 8 scores and emits improvement recommendations. It does not edit content.

---

## 21. Source Grounding and Citation Evaluation

Layer 8 must evaluate evidence packets from Layer 7 and output contracts from Layer 3/5.

### 21.1 Source grounding checks

```text
claim has supporting source
source is allowed by source policy summary
source freshness meets task requirement
first-party source preferred for brand facts
competitor source excluded when policy says so
quote/stat not distorted
citation attached to correct claim
unsupported claims flagged
```

### 21.2 Source grounding contract

```go
type SourceGroundingResult struct {
    TraceID             string   `json:"trace_id"`
    ClaimsChecked       int      `json:"claims_checked"`
    SupportedClaims     int      `json:"supported_claims"`
    UnsupportedClaims   int      `json:"unsupported_claims"`
    SourceFreshnessPass bool     `json:"source_freshness_pass"`
    CompetitorPolicyPass bool    `json:"competitor_policy_pass"`
    CitationAccuracy    float64  `json:"citation_accuracy"`
    Pass                bool     `json:"pass"`
    FailureReasons      []string `json:"failure_reasons,omitempty"`
}
```

Layer 8 does not fetch sources directly. It evaluates provided evidence packets and may request a connector eval through proper Layer 7/L2/L3 routes.

---

## 22. Skill Evaluation

Layer 8 runs skill evals defined by Layer 4.

### 22.1 Skill eval types

```text
trigger accuracy
negative trigger precision
adjacent-skill collision
output rubric
tool/action-class trajectory
token-budget regression
context-overflow regression
co-loaded skill regression
security boundary regression
```

### 22.2 Skill eval rules

Layer 8 must fail a skill release if:

```text
positive trigger accuracy below threshold
negative trigger precision below threshold
adjacent skill collision unresolved
skill body causes context budget failure
trajectory expects exact internal tool IDs instead of action classes
skill output rubric score below threshold
security eval fails
protected metadata appears in output
```

Layer 8 may propose skill changes but must not edit SKILL.md directly.

---

## 23. Governance and Audit Ledger

Layer 8 owns immutable governance records.

### 23.1 Governance record

```go
type GovernanceRecord struct {
    RecordID          string            `json:"record_id"`
    TraceID           string            `json:"trace_id"`
    TenantID          string            `json:"tenant_id"`
    ActorRefsHash     []string          `json:"actor_refs_hash"`
    AgentRefsHash     []string          `json:"agent_refs_hash"`
    ApprovalRefsHash  []string          `json:"approval_refs_hash"`
    RiskClass         string            `json:"risk_class"`
    ActionSummary     string            `json:"action_summary"`
    EvidenceRefsHash  []string          `json:"evidence_refs_hash"`
    DecisionSummary   string            `json:"decision_summary"`
    EvalRefs          []string          `json:"eval_refs,omitempty"`
    HashPrev          string            `json:"hash_prev,omitempty"`
    HashSelf          string            `json:"hash_self"`
    CreatedAt         time.Time         `json:"created_at"`
}
```

### 23.2 Required governance behaviors

```text
append-only records
hash-chain integrity
approval-to-action linking
human accountability mapping
agent identity mapping
risk class assignment
high-impact action attestation
export-safe summary generation
retention policy enforcement
legal hold support
```

Layer 8 must not expose raw governance internals to end users.

---

## 24. Correction Mining and Continuous Improvement

Layer 8 mines user corrections, eval failures, trace failures, and policy blocks.

### 24.1 Correction classes

```text
intent_misread
wrong_depth
wrong_tone
wrong_format
unsupported_claim
bad_source
bad_internal_link
missed_brief_requirement
poor_seo_priority
overwriting_existing_content
memory_update_wrong
edit_not_surgical
too_verbose
too_shallow
tool_loop
connector_failure
UI_confusion
approval_fatigue
```

### 24.2 Improvement recommendation contract

```go
type ImprovementRecommendation struct {
    RecommendationID string            `json:"recommendation_id"`
    SourceTraceIDs   []string          `json:"source_trace_ids"`
    TargetLayer      string            `json:"target_layer"`
    RecommendationType string          `json:"recommendation_type"`
    RiskClass        string            `json:"risk_class"`
    Summary          string            `json:"summary"`
    EvidenceSummary  string            `json:"evidence_summary"`
    ProposedOwner    string            `json:"proposed_owner"`
    RequiresHumanReview bool           `json:"requires_human_review"`
    Metadata         map[string]string `json:"metadata,omitempty"`
    CreatedAt        time.Time         `json:"created_at"`
}
```

### 24.3 Allowed recommendation targets

```text
Layer 2: policy rule candidate, redaction rule candidate, intent enum example candidate
Layer 3: workflow guard candidate, DAG failure policy candidate, mode gate candidate
Layer 4: skill trigger/anti-trigger candidate, eval case candidate, skill boundary candidate
Layer 5: UI copy or approval-card clarity candidate
Layer 6: runtime limit candidate, sandbox policy candidate, egress policy candidate
Layer 7: connector schema validation candidate, source freshness candidate
Layer 8: eval rubric candidate, drift threshold candidate, sampling rule candidate
```

Layer 8 must not auto-apply these changes.

---

## 25. Retention and Privacy

### 25.1 Retention classes

```yaml
retention_classes:
  routine_success:
    retain_days: 30
    raw_payloads_allowed: false

  eval_result:
    retain_days: 365
    raw_payloads_allowed: false

  policy_block:
    retain_days: 365
    raw_payloads_allowed: false

  security_incident:
    retain_days: 1095
    raw_payloads_allowed: false
    legal_hold_supported: true

  governance_record:
    retain_days: 2555
    raw_payloads_allowed: false
    append_only: true
```

### 25.2 Privacy rules

Layer 8 must:

```text
store summaries instead of raw prompts
hash sensitive refs
redact raw selected text
redact raw canvas content
redact raw brief content
redact raw memory/tone content
redact secrets and tokens
support tenant-level deletion where legally allowed
preserve governance records under legal hold
separate eval fixtures from live user data
```

---

## 26. Metrics and SLOs

### 26.1 Safety metrics

```text
prompt_injection_block_rate
semantic_policy_block_rate
approval_mismatch_rate
trust_decay_incidents
intent_drift_incidents
quarantine_decisions
cross_tenant_signal_count
memory_update_block_count
canvas_overwrite_block_count
```

### 26.2 Quality metrics

```text
workflow_success_rate
trajectory_pass_rate
output_rubric_score
source_grounding_score
content_quality_score
seo_aeo_quality_score
skill_trigger_accuracy
skill_collision_rate
pass_k_success_rate
flake_rate
```

### 26.3 Efficiency metrics

```text
tokens_per_successful_session
cost_to_converge
turns_to_converge
tool_calls_per_session
connector_latency_p95
runtime_latency_p95
eval_runtime_p95
abandoned_session_rate
user_correction_rate
```

### 26.4 Production SLO defaults

```yaml
slo:
  trace_ingestion_success_rate: 0.999
  required_span_completeness_rate: 0.995
  eval_pipeline_success_rate: 0.99
  dashboard_metric_freshness_seconds: 60
  security_signal_latency_p95_ms: 5000
  quarantine_decision_latency_p95_ms: 10000
  event_redaction_success_rate: 1.0
```

Layer 8 provides dashboard data contracts. Layer 5 renders dashboards.

---

## 27. Integration with Other Layers

### 27.1 From Layer 2

Layer 8 receives:

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

Layer 8 must not:

```text
perform intake
classify intent
authorize tools
redact outbound user response
```

### 27.2 From Layer 3

Layer 8 receives:

```text
workflow_selected
profile_selected
dag_plan_created
node_started
node_completed
node_failed
handoff_created
mode_gate_blocked
approval_request_created
content_task_contract_created
orchestration_failed
```

Layer 8 must not:

```text
choose workflow
assign profile
sequence nodes
create DAG
```

### 27.3 From Layer 4

Layer 8 receives:

```text
skill_registry_loaded
skill_activation_requested
skill_loaded
skill_load_blocked
skill_resource_loaded
skill_checksum_failed
eval_manifest_loaded
```

Layer 8 also loads eval manifests through read-only contracts.

Layer 8 must not:

```text
edit skill registry
load SKILL.md for runtime use
execute scripts
change skill status
```

### 27.4 From Layer 5

Layer 8 receives:

```text
presentation_intent_received
ui_render_success
ui_render_failed
approval_presented
approval_clicked
approval_rejected
vibe_diff_presented
canvas_revision_hash
brief_revision_hash
chat_delivery_status
```

Layer 8 must not:

```text
render UI
capture user clicks directly
write canvas
write brief
```

### 27.5 From Layer 6

Layer 8 receives:

```text
runtime_started
runtime_completed
sandbox_reset
filesystem_allowed
filesystem_blocked
egress_allowed
egress_blocked
jit_token_issued_metadata
jit_token_revoked_metadata
quarantine_executed
runtime_failed
```

Layer 8 may send:

```text
quarantine_decision
runtime_pause_recommendation
security_incident_record
```

Layer 8 must not:

```text
execute code
reset sandbox
block egress directly
mint/revoke credentials
execute quarantine
```

### 27.6 From Layer 7

Layer 8 receives:

```text
mcp_connection_opened
mcp_connection_failed
connector_call_started
connector_call_completed
connector_schema_error
retrieval_query_summary
evidence_packet_created
source_freshness_summary
tenant_partition_verified
data_security_event
```

Layer 8 must not:

```text
open MCP connection
call connector
query vector store
read memory docs
write memory docs
enforce encryption
```

---

## 28. Production Test Matrix

### 28.1 Telemetry tests

```text
root span missing -> trace rejected
tenant missing -> trace rejected
unlinked span -> trace rejected
raw prompt in event -> event rejected
hidden chain-of-thought in payload -> event rejected
secret-like payload -> event rejected
tail sampler retains errors
tail sampler retains policy blocks
tail sampler samples routine successes
trace hash-chain validates
trace replay uses config version from trace time
```

### 28.2 AgBOM tests

```text
AgBOM built for every completed trace
unexpected tool/action-class expansion flagged
connector ID hashed in export-safe view
credential ID hashed in export-safe view
approval ref linked to action
blast radius classified
cross-tenant resource class anomaly flagged
raw AgBOM not exposed to user-safe exports
```

### 28.3 Drift and trust tests

```text
mode violation increases drift
goal substitution detected
repeated policy blocks reduce trust
approval mismatch triggers trust breaker
loop pattern detected
cross-tenant signal triggers critical trust decay
safe recovery signals can improve trust within cap
quarantine decision emitted when threshold crossed
Layer 8 does not execute quarantine
```

### 28.4 SecOps tests

```text
red-team payload suite runs only in CI/shadow/staging
RAG hidden instruction test produces expected signal
MCP spoofing fixture produces expected signal
Blue Team detects repeated tool anomaly
Green Team produces recovery recommendation
quarantine decision targets Layer 6
policy change recommendation requires human review
skill change recommendation requires human review
```

### 28.5 Evaluation tests

```text
EXACT trajectory passes only exact sequence
IN_ORDER trajectory passes ordered subset
ANY_ORDER trajectory passes unordered required actions
forbidden action present -> fail
LLM judge requires rubric
LLM judge output must be JSON
position swap detects ordering bias
pass^k fails if any required run fails
flake rate above threshold blocks release
```

### 28.6 SEO/AEO eval tests

```text
unsupported SEO recommendation fails
invalid JSON-LD fails
unsafe robots.txt recommendation fails
non-canonical sitemap entry flagged
CWV recommendation without metric mapping fails
internal link target not verified fails
AI visibility claim without evidence fails
brand fact contradiction fails
```

### 28.7 Content eval tests

```text
drafting outside write mode fails
edit without selected_text hash fails
memory update without approval fails
deep research without approval fails
fabricated citation fails
brief missing audience/intent fails planning eval
draft ignoring tone fails
word count outside tolerance fails
repurposing loses CTA fails
SEO overwrite without approval fails
```

### 28.8 Correction mining tests

```text
user corrections clustered
correction clusters mapped to failure classes
regression case generated as proposal only
skill anti-trigger recommendation created
workflow guard recommendation created
policy recommendation created
no recommendation auto-applied
```

### 28.9 Governance tests

```text
governance record append-only
hash-chain validates
approval linked to action
actor refs hashed
agent refs hashed
legal hold prevents deletion
retention policy deletes eligible routine success traces
user-safe audit export redacts protected internals
```

### 28.10 Boundary tests

Layer 8 must prove it does not:

```text
normalize raw input
classify intent
authorize tools
choose workflow
assign agent
build DAG
load SKILL.md for runtime
execute skill scripts
render UI
create A2UI
capture approval clicks directly
execute tool
execute connector
open MCP connection
query vector store
read memory document
write memory document
write canvas
write brief
manage sandbox
enforce egress
mint JIT token
revoke JIT token
execute quarantine
change policies directly
change skills directly
change workflows directly
```

---

## 29. Deployment Requirements

Layer 8 v2 is production-ready only if:

```text
telemetry-policy validates at startup
trace-schema validates at startup
redaction-policy validates at startup
eval-policy validates at startup
secops-policy validates at startup
drift-policy validates at startup
trust-policy validates at startup
retention-policy validates at startup
governance-policy validates at startup
all event ingestion schemas validate
all required spans are enforced
all sensitive payload redaction tests pass
AgBOM is generated for every production trace
intent drift and trust decay are calculated for every production trace
quarantine decisions are routed to Layer 6 only
Layer 8 never executes runtime actions
Layer 8 never stores hidden chain-of-thought
Layer 8 never stores raw prompts, secrets, raw canvas, raw brief, raw memory, or raw selected text
SEO/AEO eval suites pass
content-agent eval suites pass
skill eval integration with Layer 4 manifests passes
pass^k gates are enforced
governance ledger hash-chain validates
retention policy is enforced
dashboard exports are data contracts only
all zero-overlap boundary tests pass
```

---

## 30. Acceptance Criteria

Layer 8 v2 is accepted when:

1. It stores complete redacted trace trees for SEO/AEO and content-agent sessions.
2. It validates and rejects malformed or unsafe telemetry events.
3. It tracks Runtime Agent Bill of Materials for every session.
4. It calculates cost, latency, token burn, retries, and convergence.
5. It detects drift, trust decay, loops, connector anomalies, approval misuse, and unsafe surface mutation patterns.
6. It runs red-team suites only in safe environments.
7. It monitors live traces through blue-team behavioural analytics.
8. It generates green-team recovery recommendations without executing them.
9. It emits quarantine decisions to Layer 6 without directly quarantining.
10. It scores trajectories with EXACT, IN_ORDER, ANY_ORDER, and forbidden-action checks.
11. It runs LLM-as-judge scoring with rubric, JSON, position-swap, and stability controls.
12. It runs pass^k reliability gates for production workflows.
13. It evaluates SEO/AEO outputs for evidence, severity, structured data, internal links, CWV mapping, AI visibility, and recommendation quality.
14. It evaluates content-agent outputs for brief completeness, source grounding, tone, SEO alignment, edit precision, word-count control, repurposing quality, and memory safety.
15. It mines user corrections into labeled failure clusters.
16. It proposes improvements to owning layers but never auto-applies them.
17. It maintains an immutable governance audit ledger.
18. It enforces privacy and retention policies.
19. It provides metrics and dashboard data contracts without rendering dashboards.
20. It passes all boundary tests against Layers 1 through 7.

---

## 31. Final Non-Goals

```text
Layer 8 must not normalize raw user input.
Layer 8 must not classify raw intent.
Layer 8 must not inspect prompt injection before intake.
Layer 8 must not authorize tools.
Layer 8 must not validate approval metadata for execution.
Layer 8 must not choose workflows.
Layer 8 must not assign profiles.
Layer 8 must not build DAGs.
Layer 8 must not load SKILL.md for runtime context.
Layer 8 must not edit skill files.
Layer 8 must not execute skill scripts.
Layer 8 must not render A2UI.
Layer 8 must not render dashboards.
Layer 8 must not render canvas, brief, or chat.
Layer 8 must not capture user approval clicks directly.
Layer 8 must not execute tools.
Layer 8 must not execute scripts.
Layer 8 must not manage sandbox.
Layer 8 must not enforce filesystem permissions.
Layer 8 must not enforce network egress.
Layer 8 must not mint credentials.
Layer 8 must not revoke credentials.
Layer 8 must not execute quarantine.
Layer 8 must not open MCP connections.
Layer 8 must not call connectors.
Layer 8 must not query vector stores.
Layer 8 must not read memory documents.
Layer 8 must not write memory documents.
Layer 8 must not write article drafts.
Layer 8 must not edit existing content.
Layer 8 must not update SEO fields.
Layer 8 must not publish externally.
Layer 8 must not auto-change policies, workflows, skills, UI contracts, runtime settings, or connector contracts.
Layer 8 must not expose raw traces, raw AgBOM, hidden chain-of-thought, exact internal tool inventory, workflow IDs, profile IDs, MCP endpoints, secrets, tokens, cookies, raw PII, raw canvas, raw brief, raw memory, or raw selected text to end users.
```

---

## 32. One-Line Architecture Summary

Layer 8 v2 is the production-grade glass-box observability, SecOps, evaluation, governance, and continuous-improvement layer that records redacted traces, tracks AgBOM, scores safety and quality, detects drift and trust decay, runs trajectory/output/pass^k evaluations, mines failures, and publishes decisions or recommendations while leaving intake, routing, skills, rendering, execution, MCP/data access, memory persistence, and runtime quarantine execution to their owning layers.
