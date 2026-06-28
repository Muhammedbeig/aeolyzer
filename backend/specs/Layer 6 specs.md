# Layer 6 Production-Ready Specs v2
## Security Control Plane, Isolated Runtime, JIT Identity, Egress Governance, Supply-Chain Defense, and Stateful Recovery

**Status:** Production-ready upgrade  
**Supersedes:** Layer 6 v1 where sandboxing, network egress, JIT credentials, file enforcement, and Green Team recovery were sketched but not operationally specified  
**Primary rule:** Layer 6 is the execution boundary. It runs already-authorized tool and script actions inside a zero-trust runtime. It never classifies user intent, chooses workflows, loads skills, decides policy, renders UI, owns MCP/data connectors, stores telemetry, evaluates quality, or decides whether quarantine is needed.

---

## 1. Upgrade Decision

Upgrade is required.

Layer 6 v1 correctly identified the core runtime responsibilities:

```text
gVisor / kernel-level isolation
ephemeral sandbox reset
non-interactive egress proxy
supply-chain defense
JIT token broker
file-tree enforcement
stateful quarantine
auto-refactoring loop
```

But the rest of the platform has evolved into a production SEO/AEO auditor and content-writing agent with multiple execution classes:

```text
read-only crawls
page scraping and rendering
Search Console / analytics connector calls
schema validation
robots.txt and sitemap validation
deterministic skill scripts
content QA parsers
metadata length calculators
orphan-page graph scans
source/evidence extractors
safe file patching
human-approved content writes
external A2A task execution
```

Layer 6 must become a production-grade execution control plane that can safely run high-velocity agent actions without turning the agent into an over-privileged automation account.

The upgrade preserves the architecture:

```text
Layer 1: user surface, frontend/client shell, device UI, local CLI, hardware MFA challenge UX
Layer 2: intake, policy, semantic safety, protected disclosure guard, tool authorization
Layer 3: DAG orchestration, workflow routing, agent assignment, execution sequencing
Layer 4: progressive-disclosure skills, instructions, bundled scripts, references, assets
Layer 5: A2UI, presentation surfaces, approval cards, A2A application envelopes
Layer 6: isolated runtime, sandbox, file controls, egress controls, JIT credentials, runtime recovery execution
Layer 7: MCP transport, API connectors, data mesh, vector/RAG stores, durable memory retrieval
Layer 8: observability, evals, drift detection, AgBOM analytics, Red/Blue/Green decisioning
```

Layer 6 must be treated as a hardened execution kernel, not a policy brain.

---

## 2. Layer Objective

Layer 6 executes approved actions inside a constrained environment with the minimum possible authority.

It must:

1. Accept only signed, policy-approved `RuntimeExecutionRequest` envelopes.
2. Reject raw prompts, unapproved tool calls, unstructured commands, and direct shell requests.
3. Create a per-task sandbox lease using deny-by-default runtime profiles.
4. Execute tools, scripts, headless browsers, validators, and controlled file operations only through approved runtime adapters.
5. Strip ambient credentials before execution.
6. Issue or request only task-scoped, short-lived JIT credentials when explicitly required.
7. Confine filesystem reads and writes to declared mounts.
8. Force all network access through controlled egress paths.
9. Block dynamic dependency installation except through pinned, vetted registries.
10. Preserve deterministic artifacts required by the workflow and discard all other ephemeral state.
11. Enforce resource limits: CPU, memory, wall-clock time, file size, process count, network budget, and output token budget.
12. Capture stdout, stderr, exit codes, resource usage, and sanitized execution summaries.
13. Redact secrets, tokens, cookies, raw PII, and protected metadata before any runtime result leaves the layer.
14. Emit sanitized runtime events to Layer 8.
15. Execute stateful quarantine commands received from Layer 8.
16. Support safe auto-refactoring only inside quarantine or repair scopes, never as uncontrolled self-modification.
17. Fail closed on unknown tools, malformed envelopes, expired approvals, missing scopes, path traversal, symlink escape, egress denial, dependency drift, sandbox profile mismatch, or output-policy violation.

Layer 6 exists to answer this question:

```text
Given an already-authorized action, how can it be executed with the smallest possible blast radius?
```

Layer 6 must not answer:

```text
What did the user ask?
Is the request safe?
Which workflow should run?
Which agent should execute?
Which skill should load?
Which MCP server should be called?
Which UI should be shown?
Did the final answer satisfy the user?
Is the agent drifting?
Should quarantine be triggered?
```

---

## 3. Strict Layer Boundary

### 3.1 Layer 6 owns

```text
runtime execution gateway
sandbox lease creation
sandbox profile selection from approved request metadata
kernel-level isolation profile enforcement
ephemeral state reset
process execution limits
headless browser sandboxing
script runner hardening
controlled file read/write enforcement
path canonicalization
symlink and hardlink escape prevention
workspace mount enforcement
output artifact staging
network egress routing
DNS and outbound socket denial
dynamic dependency blocking
pinned dependency registry enforcement
runtime package cache policy
runtime SBOM capture
runtime secret stripping
runtime environment variable filtering
JIT token issuance/retrieval enforcement
JIT token attachment to approved tool contexts
JIT token revocation on completion
credential TTL and scope enforcement
stdout/stderr capture
runtime result redaction
runtime event emission
stateful quarantine execution
runtime freeze/thaw mechanics
tool access revocation execution
ephemeral sandbox teardown
repair-only auto-refactoring execution
runtime health probes
runtime acceptance tests
```

### 3.2 Layer 6 does not own

| Responsibility | Owning layer |
|---|---|
| User input normalization | Layer 2 |
| Prompt-injection detection | Layer 2 |
| Protected disclosure detection | Layer 2 |
| Tool allow/deny policy | Layer 2 |
| Semantic pre-execution safety judgment | Layer 2 |
| User intent classification | Layer 2 |
| Workflow selection | Layer 3 |
| DAG sequencing | Layer 3 |
| Agent assignment | Layer 3 |
| Tool-choice reasoning | Layer 3 |
| Skill metadata or skill body ownership | Layer 4 |
| Skill trigger logic | Layer 4 / Layer 3 contract |
| A2UI / dashboards / approval UX | Layer 5 |
| A2A application envelope formatting | Layer 5 |
| MCP transport and API connector implementations | Layer 7 |
| Durable memory retrieval | Layer 7 |
| Data-at-rest encryption ownership | Layer 7 data mesh / platform infra |
| Telemetry storage | Layer 8 |
| Eval scoring | Layer 8 |
| Intent drift detection | Layer 8 |
| Quarantine decisioning | Layer 8 |
| Red Team and Blue Team simulation/analysis | Layer 8 |
| Final response composition | Layer 3 / assistant response layer, guarded by Layer 2 and rendered by Layer 5 |

### 3.3 Non-overlap rule for Layer 6

Allowed:

```text
Execute an already-authorized crawl action through a controlled web-fetch adapter.
Run a deterministic skill script in an ephemeral sandbox.
Open only declared read-only and write-only mounts.
Attach a task-scoped token to one approved connector invocation.
Deny direct network access from an untrusted script.
Kill a process that exceeds memory or wall-clock limits.
Reset sandbox state between action attempts.
Freeze tool access after a Layer 8 quarantine command.
Preserve short-term state for forensics during quarantine.
Emit sanitized execution metadata to Layer 8.
```

Forbidden:

```text
Classify the user's request.
Re-authorize a blocked tool call.
Let an agent open a raw shell because it asks nicely.
Load SKILL.md bodies.
Choose which skill should run.
Open MCP sockets directly.
Call external APIs outside Layer 7 connectors.
Render approval cards.
Display diffs.
Store telemetry.
Score task quality.
Decide that intent drift exists.
Start quarantine without a trusted command.
Expose filesystem paths, source code, secrets, raw logs, tokens, cookies, or tool implementations to the user.
```

---

## 4. Whitepaper-Aligned Design Principles

### 4.1 Secure the harness, not just the model

A raw model becomes an enterprise agent only when wrapped in a harness that gives it tools, state, feedback loops, and constraints. Layer 6 is one of the hardest boundaries in that harness. It assumes model output may be wrong, manipulated, stale, or adversarial.

Production rule:

```text
Never trust the model to self-limit execution authority.
The runtime must enforce limits externally.
```

### 4.2 Context-as-a-perimeter

Static credentials and static RBAC are insufficient for agentic systems because a valid token does not prove aligned intent. Layer 6 enforces runtime context constraints:

```text
who requested the action
which intent was classified
which workflow node is active
which tool was approved
which data scope is allowed
which environment is active
which mount is allowed
which egress path is allowed
which token TTL is active
```

Runtime context must be checked at execution time, not only during login.

### 4.3 Zero ambient authority

No agent runtime may inherit the broad permissions of the human user, service account, host process, developer laptop, IDE, CI worker, or orchestrator.

Required behavior:

```text
start with no credentials
attach only scoped JIT credentials
expire credentials at task completion
revoke credentials on failure, timeout, or quarantine
redact credential material before event emission
```

### 4.4 Ephemeral by default, persistent by exception

Every sandbox starts clean and ends disposable. Persistent state must be explicitly declared as an output artifact, not accidentally inherited from a previous run.

Layer 6 may preserve:

```text
approved output files
sanitized execution summaries
forensic runtime snapshots under quarantine
approved cache entries from vetted registries
```

Layer 6 must discard:

```text
working directories
temporary files
shell history
process-local secrets
package manager caches not explicitly trusted
browser cookies
session storage
unapproved downloads
transient generated scripts
```

### 4.5 Non-interactive egress only

Agent code must not browse the internet freely. Any external fetch must go through controlled tools, crawlers, connector adapters, or Layer 7 transport.

Production rule:

```text
No raw outbound sockets from untrusted scripts.
No direct curl/wget/browser network escape from generated code.
No arbitrary package registry access.
No DNS exfiltration.
No internal metadata service access.
```

### 4.6 Dependency installation is a security event

Dynamic dependency installation is not a convenience feature. It is a supply-chain risk.

Default rule:

```text
Generated scripts may not install dependencies at runtime.
```

Exception rule:

```text
A dependency may be used only if it is pre-approved, pinned, hash-verified, fetched from a vetted registry, and recorded in the runtime SBOM.
```

### 4.7 Small batch sizes reduce blast radius

Large autonomous changes increase review burden and failure cost. Layer 6 enforces batch limits on runtime actions:

```text
max files changed
max output bytes
max URLs crawled
max pages rendered
max records processed
max write operations
max retry attempts
max wall-clock time
```

The agent may ask for a larger batch, but Layer 6 does not grant it unless the authorized request explicitly contains the larger budget.

### 4.8 Quarantine is stateful, not destructive

When Layer 8 detects intent drift, unsafe behavior, or trust decay, Layer 6 must be able to freeze the runtime without destroying forensic context or user-visible work.

Stateful quarantine means:

```text
freeze external actions
revoke tokens
block network egress
preserve short-term runtime memory for investigation
preserve approved user canvas state
prevent additional tool execution
return a safe stopped-state summary
```

It does not mean blindly killing the container and losing the evidence.

### 4.9 Layer 6 is enforcement, not explanation

Layer 6 can report that it denied a path, token, network route, package, or runtime profile. It must not reveal protected implementation details. User-facing explanations must be safe summaries produced through Layer 2 and Layer 5.

---

## 5. Required Directory Upgrade

### 5.1 Final Layer 6 tree

```text
/layer_06_runtime
  ├── README.md
  ├── layer6-boundary.md
  ├── runtime-execution.schema.json
  ├── sandbox-lease.schema.json
  ├── runtime-result.schema.json
  ├── jit-token.schema.json
  ├── quarantine-command.schema.json
  ├── dependency-policy.schema.json
  ├── filesystem-policy.schema.json
  ├── egress-policy.schema.json
  ├── runtime-changelog.md
  │
  ├── /config
  │   ├── runtime_profiles.yaml
  │   ├── sandbox_profiles.yaml
  │   ├── tool_runtime_map.yaml
  │   ├── filesystem_mount_policy.yaml
  │   ├── network_egress_policy.yaml
  │   ├── dependency_allowlist.yaml
  │   ├── registry_trust_policy.yaml
  │   ├── jit_scope_policy.yaml
  │   ├── resource_budget_policy.yaml
  │   ├── output_redaction_policy.yaml
  │   ├── quarantine_policy.yaml
  │   ├── repair_policy.yaml
  │   └── runtime_event_policy.yaml
  │
  ├── /execution_gateway
  │   ├── runtime_request_validator.go
  │   ├── approved_envelope_verifier.go
  │   ├── runtime_profile_resolver.go
  │   ├── execution_dispatcher.go
  │   ├── tool_adapter_registry.go
  │   ├── script_runner.go
  │   ├── browser_runner.go
  │   ├── validator_runner.go
  │   ├── file_operation_runner.go
  │   ├── retry_guard.go
  │   ├── timeout_guard.go
  │   ├── resource_budget_guard.go
  │   └── runtime_result_builder.go
  │
  ├── /sandbox_environment
  │   ├── gvisor_isolation.go
  │   ├── sandbox_lease_manager.go
  │   ├── namespace_isolator.go
  │   ├── seccomp_profile_manager.go
  │   ├── rootless_user_manager.go
  │   ├── cgroup_limiter.go
  │   ├── process_tree_guard.go
  │   ├── tempdir_manager.go
  │   ├── ephemeral_state_manager.go
  │   ├── cleanroom_image_resolver.go
  │   ├── sandbox_health_checker.go
  │   └── sandbox_teardown.go
  │
  ├── /filesystem_control
  │   ├── file_tree_enforcer.go
  │   ├── mount_policy_loader.go
  │   ├── path_canonicalizer.go
  │   ├── symlink_escape_guard.go
  │   ├── hardlink_escape_guard.go
  │   ├── read_scope_guard.go
  │   ├── write_scope_guard.go
  │   ├── artifact_stager.go
  │   ├── file_quota_guard.go
  │   ├── diff_limit_guard.go
  │   ├── file_redactor.go
  │   └── forensic_snapshotter.go
  │
  ├── /network_egress
  │   ├── egress_proxy_controller.go
  │   ├── outbound_socket_blocker.go
  │   ├── dns_policy_enforcer.go
  │   ├── url_scope_guard.go
  │   ├── metadata_service_blocker.go
  │   ├── crawler_proxy_adapter.go
  │   ├── connector_proxy_adapter.go
  │   ├── bandwidth_budget_guard.go
  │   ├── request_header_redactor.go
  │   └── egress_event_builder.go
  │
  ├── /supply_chain_defense
  │   ├── supply_chain_defender.go
  │   ├── dynamic_install_blocker.go
  │   ├── package_name_confusion_guard.go
  │   ├── lockfile_enforcer.go
  │   ├── hash_verifier.go
  │   ├── registry_pin_enforcer.go
  │   ├── sbom_generator.go
  │   ├── vulnerability_policy_checker.go
  │   ├── license_policy_checker.go
  │   └── dependency_cache_manager.go
  │
  ├── /iam_context
  │   ├── ambient_credential_stripper.go
  │   ├── jit_token_broker.go
  │   ├── token_scope_resolver.go
  │   ├── token_ttl_guard.go
  │   ├── token_audience_guard.go
  │   ├── token_attach_guard.go
  │   ├── token_revoker.go
  │   ├── secret_env_filter.go
  │   └── credential_redactor.go
  │
  ├── /runtime_observation
  │   ├── stdout_stderr_capture.go
  │   ├── exit_code_mapper.go
  │   ├── resource_usage_meter.go
  │   ├── runtime_event_emitter.go
  │   ├── execution_summary_redactor.go
  │   ├── runtime_agbom_reporter.go
  │   ├── policy_outcome_reporter.go
  │   └── anomaly_signal_forwarder.go
  │
  ├── /secops_green_team_ops
  │   ├── quarantine_command_validator.go
  │   ├── stateful_quarantine.go
  │   ├── runtime_freezer.go
  │   ├── tool_access_revoker.go
  │   ├── network_freezer.go
  │   ├── token_emergency_revoker.go
  │   ├── forensic_bundle_builder.go
  │   ├── quarantine_status_reporter.go
  │   ├── auto_refactoring_loop.go
  │   ├── repair_patch_runner.go
  │   └── recovery_release_guard.go
  │
  └── /tests
      ├── runtime_request_schema_test.go
      ├── approved_envelope_verification_test.go
      ├── sandbox_profile_test.go
      ├── gvisor_isolation_test.go
      ├── ephemeral_reset_test.go
      ├── path_traversal_test.go
      ├── symlink_escape_test.go
      ├── write_scope_test.go
      ├── egress_denial_test.go
      ├── metadata_service_block_test.go
      ├── dynamic_install_block_test.go
      ├── package_confusion_test.go
      ├── hash_pinning_test.go
      ├── ambient_credential_strip_test.go
      ├── jit_scope_test.go
      ├── jit_expiry_test.go
      ├── token_revocation_test.go
      ├── output_redaction_test.go
      ├── stdout_stderr_redaction_test.go
      ├── resource_budget_test.go
      ├── quarantine_freeze_test.go
      ├── forensic_preservation_test.go
      ├── repair_scope_test.go
      ├── no_intent_classification_boundary_test.go
      ├── no_policy_decision_boundary_test.go
      ├── no_workflow_selection_boundary_test.go
      ├── no_skill_loading_boundary_test.go
      ├── no_mcp_transport_boundary_test.go
      └── golden_runtime_flows_test.go
```

---

## 6. Canonical Runtime Classes

Layer 6 must support a finite set of runtime classes. Unknown classes are rejected.

```yaml
runtime_classes:
  read_only_fetch:
    purpose: "Fetch or render external web content through approved crawl/fetch adapters."
    filesystem: "read_temp_write_artifact_only"
    network: "egress_proxy_only"
    credentials: "none_or_fetch_scoped"
    mutation: false

  analytics_connector_call:
    purpose: "Invoke pre-approved analytics/search data connector through Layer 7 adapter boundary."
    filesystem: "temp_only"
    network: "connector_proxy_only"
    credentials: "jit_connector_scoped"
    mutation: false

  deterministic_skill_script:
    purpose: "Run bundled deterministic scripts from approved skills, such as parsers, counters, validators, graph analyzers, and metadata calculators."
    filesystem: "declared_input_mounts_plus_output_mount"
    network: "none_by_default"
    credentials: "none_by_default"
    mutation: "artifact_only"

  validation_runner:
    purpose: "Validate JSON-LD, XML, robots rules, metadata length, internal links, and structured outputs."
    filesystem: "read_input_write_result"
    network: "none_by_default"
    credentials: "none"
    mutation: "result_only"

  headless_browser_runner:
    purpose: "Render pages, inspect layout, capture screenshots, and test basic UI interactions inside a browser sandbox."
    filesystem: "temp_plus_artifact"
    network: "egress_proxy_only"
    credentials: "none_or_site_scoped"
    mutation: false

  safe_file_patch:
    purpose: "Apply already-approved file patches to allowed workspace artifacts."
    filesystem: "declared_rw_mounts_only"
    network: "none"
    credentials: "none"
    mutation: true

  approval_bound_action:
    purpose: "Perform an action that Layer 2/3 marked as requiring valid approval metadata."
    filesystem: "declared_by_request"
    network: "declared_by_request"
    credentials: "jit_scoped_if_required"
    mutation: "declared_by_request"

  quarantine_repair:
    purpose: "Run repair-only scripts after Layer 8/Green Team initiated quarantine."
    filesystem: "forensic_snapshot_plus_repair_workspace"
    network: "none_unless_repair_policy_allows"
    credentials: "none"
    mutation: "repair_scope_only"
```

### Required runtime class behavior

```text
Every runtime request must declare exactly one runtime_class.
Runtime profile must match tool/action type.
Runtime class cannot be upgraded by the agent.
Runtime class cannot be changed after approval.
Runtime class determines default filesystem, network, credential, and resource budgets.
```

---

## 7. Core Input and Output Contracts

### 7.1 Input from Layer 2 / Layer 3: `RuntimeExecutionRequest`

Layer 6 accepts only an execution request that has passed Layer 2 tool-policy authorization and Layer 3 orchestration sequencing.

```go
type RuntimeExecutionRequest struct {
    TraceID             string                 `json:"trace_id"`
    SessionID           string                 `json:"session_id"`
    TaskID              string                 `json:"task_id"`
    WorkflowID          string                 `json:"workflow_id,omitempty"`
    NodeID              string                 `json:"node_id"`
    AgentID             string                 `json:"agent_id"`
    Intent              string                 `json:"intent"`
    RuntimeClass        string                 `json:"runtime_class"`
    ActionType          string                 `json:"action_type"`
    ToolContractID      string                 `json:"tool_contract_id,omitempty"`
    ScriptRef           string                 `json:"script_ref,omitempty"`
    InputRefs           []ArtifactRef          `json:"input_refs,omitempty"`
    Params              map[string]interface{} `json:"params"`
    FilesystemPolicyRef string                 `json:"filesystem_policy_ref"`
    EgressPolicyRef     string                 `json:"egress_policy_ref"`
    CredentialPolicyRef string                 `json:"credential_policy_ref"`
    ResourceBudgetRef   string                 `json:"resource_budget_ref"`
    ApprovalRef         string                 `json:"approval_ref,omitempty"`
    PolicyDecisionID    string                 `json:"policy_decision_id"`
    RequestSignature    string                 `json:"request_signature"`
    ExpiresAt           string                 `json:"expires_at"`
    Metadata            map[string]string      `json:"metadata,omitempty"`
}

type ArtifactRef struct {
    RefID       string `json:"ref_id"`
    RefType     string `json:"ref_type"`
    AccessMode  string `json:"access_mode"`
    ContentHash string `json:"content_hash,omitempty"`
}
```

Layer 6 must reject requests with:

```text
missing PolicyDecisionID
missing RequestSignature
expired ExpiresAt
unknown RuntimeClass
unknown ActionType
unknown ToolContractID for tool actions
unregistered ScriptRef for script actions
missing filesystem policy
missing egress policy
missing credential policy
missing resource budget
raw shell command strings outside registered adapters
inline executable code unless explicitly allowed by a registered repair or validation adapter
network=unrestricted
filesystem=root_rw
credential_policy=ambient
approval-required action without ApprovalRef
params containing unresolved high-risk placeholders
```

### 7.2 Internal contract: `SandboxLease`

```go
type SandboxLease struct {
    LeaseID          string            `json:"lease_id"`
    TraceID          string            `json:"trace_id"`
    TaskID           string            `json:"task_id"`
    RuntimeClass     string            `json:"runtime_class"`
    SandboxProfile   string            `json:"sandbox_profile"`
    RootlessUser      string            `json:"rootless_user"`
    Mounts           []MountSpec       `json:"mounts"`
    EgressProfile    string            `json:"egress_profile"`
    CredentialRefs   []string          `json:"credential_refs,omitempty"`
    ResourceBudgets  ResourceBudgets   `json:"resource_budgets"`
    CreatedAt        string            `json:"created_at"`
    ExpiresAt        string            `json:"expires_at"`
    Labels           map[string]string `json:"labels,omitempty"`
}

type MountSpec struct {
    LogicalName string `json:"logical_name"`
    Mode        string `json:"mode"`
    Scope       string `json:"scope"`
    MaxBytes    int64  `json:"max_bytes,omitempty"`
}

type ResourceBudgets struct {
    MaxWallClockMs int64 `json:"max_wall_clock_ms"`
    MaxCPUPercent  int64 `json:"max_cpu_percent"`
    MaxMemoryMB    int64 `json:"max_memory_mb"`
    MaxProcesses   int64 `json:"max_processes"`
    MaxOutputBytes int64 `json:"max_output_bytes"`
    MaxNetworkMB   int64 `json:"max_network_mb"`
    MaxRetries     int64 `json:"max_retries"`
}
```

The lease is internal to Layer 6 and must not be exposed to the user.

### 7.3 Output to Layer 3 / Layer 5: `RuntimeExecutionResult`

```go
type RuntimeExecutionResult struct {
    TraceID           string                 `json:"trace_id"`
    SessionID         string                 `json:"session_id"`
    TaskID            string                 `json:"task_id"`
    NodeID            string                 `json:"node_id"`
    RuntimeClass      string                 `json:"runtime_class"`
    Status            string                 `json:"status"`
    ExitCode          int                    `json:"exit_code,omitempty"`
    Summary           string                 `json:"summary"`
    SanitizedStdoutRef string                `json:"sanitized_stdout_ref,omitempty"`
    SanitizedStderrRef string                `json:"sanitized_stderr_ref,omitempty"`
    OutputArtifacts   []ArtifactRef          `json:"output_artifacts,omitempty"`
    ResourceUsage     map[string]interface{} `json:"resource_usage,omitempty"`
    PolicyOutcomes    []string               `json:"policy_outcomes,omitempty"`
    Retryable         bool                   `json:"retryable"`
    ErrorClass        string                 `json:"error_class,omitempty"`
    RedactionApplied  bool                   `json:"redaction_applied"`
    Metadata          map[string]string      `json:"metadata,omitempty"`
}
```

Allowed statuses:

```text
success
blocked_by_runtime_policy
blocked_by_filesystem_policy
blocked_by_egress_policy
blocked_by_supply_chain_policy
blocked_by_credential_policy
blocked_by_resource_budget
failed_execution
failed_validation
timeout
quarantined
redacted_output
```

### 7.4 Input from Layer 8: `QuarantineCommand`

```go
type QuarantineCommand struct {
    TraceID          string            `json:"trace_id"`
    SessionID        string            `json:"session_id"`
    TargetScope      string            `json:"target_scope"`
    TriggerReason    string            `json:"trigger_reason"`
    Severity         string            `json:"severity"`
    RequestedActions []string          `json:"requested_actions"`
    PreserveState    bool              `json:"preserve_state"`
    DecisionRef      string            `json:"decision_ref"`
    ExpiresAt        string            `json:"expires_at"`
    Signature        string            `json:"signature"`
    Metadata         map[string]string `json:"metadata,omitempty"`
}
```

Layer 6 validates the command signature and executes only allowed actions:

```text
freeze_runtime
revoke_jit_tokens
block_egress
revoke_tool_access
preserve_forensic_snapshot
stop_new_executions
allow_read_only_status
start_repair_scope
```

Layer 6 does not decide whether the command is justified. Layer 8 owns detection and decisioning.

---

## 8. Runtime Request Flow

### 8.1 Happy path

```text
01 Layer 3 proposes action
02 Layer 2 validates/sanitizes/authorizes the action
03 Layer 2 emits authorized runtime envelope
04 Layer 6 validates envelope signature and expiry
05 Layer 6 resolves runtime class and sandbox profile
06 Layer 6 strips ambient credentials
07 Layer 6 creates sandbox lease
08 Layer 6 mounts declared input/output scopes
09 Layer 6 configures egress proxy or network deny profile
10 Layer 6 attaches scoped JIT token only if required
11 Layer 6 executes adapter, script, validator, browser, or file patch
12 Layer 6 captures stdout/stderr/resources/artifacts
13 Layer 6 revokes JIT tokens
14 Layer 6 redacts execution result
15 Layer 6 tears down ephemeral sandbox
16 Layer 6 emits sanitized runtime event to Layer 8
17 Layer 6 returns RuntimeExecutionResult to Layer 3
18 Layer 3 decides next DAG step
```

### 8.2 Fail-closed path

```text
request invalid
-> reject before sandbox creation
-> emit blocked_by_runtime_policy
-> return safe blocked result
```

```text
filesystem escape attempt
-> stop process
-> revoke credentials
-> preserve forensic metadata if configured
-> emit blocked_by_filesystem_policy
-> return safe blocked result
```

```text
dynamic dependency install attempt
-> block package manager call
-> stop process unless policy allows graceful continuation
-> emit blocked_by_supply_chain_policy
-> return safe blocked result
```

```text
egress violation
-> deny request
-> stop or continue according to egress severity
-> emit blocked_by_egress_policy
-> return safe blocked result
```

```text
Layer 8 quarantine command
-> freeze runtime
-> revoke tokens
-> block network
-> preserve state
-> return quarantined status
```

---

## 9. Sandbox Environment Specs

### 9.1 Kernel-level isolation

Layer 6 must run untrusted or model-influenced execution inside kernel-level isolation.

Required controls:

```text
gVisor or equivalent sandbox runtime
rootless execution
user namespace isolation
PID namespace isolation
network namespace isolation
mount namespace isolation
seccomp profile enforcement
read-only base image
no privileged containers
no host PID namespace
no host network namespace
no hostPath mounts
no Docker socket mounts
no Kubernetes service account token auto-mount
no metadata service access
cgroup CPU/memory/process limits
```

### 9.2 Sandbox profiles

```yaml
sandbox_profiles:
  no_network_script:
    runtime: "gvisor"
    rootless: true
    network: "none"
    filesystem: "temp_plus_declared_mounts"
    seccomp: "strict"
    max_processes: 16

  egress_proxy_browser:
    runtime: "gvisor"
    rootless: true
    network: "proxy_only"
    filesystem: "temp_plus_artifacts"
    seccomp: "browser_constrained"
    max_processes: 64

  connector_invocation:
    runtime: "gvisor"
    rootless: true
    network: "connector_proxy_only"
    filesystem: "temp_only"
    seccomp: "strict"
    max_processes: 16

  file_patch:
    runtime: "gvisor"
    rootless: true
    network: "none"
    filesystem: "declared_rw_only"
    seccomp: "strict"
    max_processes: 8

  quarantine_repair:
    runtime: "gvisor"
    rootless: true
    network: "none_by_default"
    filesystem: "forensic_snapshot_plus_repair_workspace"
    seccomp: "strict"
    max_processes: 8
```

### 9.3 Sandbox lifecycle

```text
Create clean temp root.
Attach declared mounts.
Apply network profile.
Apply seccomp/cgroup limits.
Run single action.
Capture outputs.
Redact outputs.
Stage approved artifacts.
Revoke credentials.
Delete temp root.
Destroy process namespace.
Emit runtime event.
```

### 9.4 Ephemeral reset requirements

A sandbox reset must occur:

```text
before every new tool execution
after every successful execution
after every failed execution
after every retry attempt
before switching runtime class
before switching agent ID
before and after quarantine repair
```

Sandbox reset must clear:

```text
temporary filesystem
process tree
environment variables
browser profile
cookies
local/session storage
package caches unless approved
credential files
SSH agent sockets
cloud metadata tokens
shell history
```

---

## 10. Filesystem Control Specs

### 10.1 Deny-by-default filesystem policy

Layer 6 must deny all filesystem access unless explicitly declared.

Default:

```yaml
filesystem_default:
  read: false
  write: false
  execute: false
  list: false
```

Allowed access is granted only by logical mount, not by arbitrary host path.

```yaml
mounts:
  input_artifacts:
    mode: read_only
    scope: workflow_declared_inputs

  output_artifacts:
    mode: write_only
    scope: task_declared_outputs

  canvas_output:
    mode: read_write
    scope: approved_canvas_artifact_only

  temp:
    mode: read_write
    scope: sandbox_ephemeral_only

  skill_script_bundle:
    mode: read_execute
    scope: layer4_declared_script_ref_only
```

### 10.2 Path canonicalization

Every filesystem operation must pass:

```text
normalize path
resolve relative segments
reject null bytes
reject control characters
resolve symlinks
resolve hardlinks
compare canonical path to mount scope
check file type
check file size quota
check read/write mode
execute operation
```

Reject patterns:

```text
../ traversal
absolute host paths
hidden root paths
symlink to outside mount
hardlink to outside mount
named pipes unless explicitly allowed
Unix sockets unless explicitly allowed
block devices
character devices
procfs/sysfs writes
credential file paths
source-code implementation paths not declared as input
production manifest paths
```

### 10.3 Write operation policy

Writes require all of the following:

```text
approved action envelope
write-enabled logical mount
declared artifact target
quota available
file extension allowed for runtime class
no overwrite unless approved
atomic write path
post-write validation
artifact hash recorded
```

### 10.4 File patch policy

Patch actions must be:

```text
small-batch
line-scoped when possible
reviewable
approval-bound for protected destinations
non-recursive by default
reversible with generated diff
validated after apply
```

Required patch limits:

```yaml
patch_limits:
  max_files_changed_default: 3
  max_lines_changed_default: 200
  max_single_file_bytes_default: 262144
  max_total_patch_bytes_default: 1048576
  overwrite_requires_approval: true
  protected_file_write_default: deny
```

### 10.5 Artifact staging

Layer 6 returns references to staged artifacts, not raw file trees.

Artifact metadata:

```go
type StagedArtifact struct {
    ArtifactID  string `json:"artifact_id"`
    TaskID      string `json:"task_id"`
    Kind        string `json:"kind"`
    ContentHash string `json:"content_hash"`
    ByteSize    int64  `json:"byte_size"`
    Redacted    bool   `json:"redacted"`
    CreatedAt   string `json:"created_at"`
}
```

Layer 6 must not expose absolute filesystem paths to the user.

---

## 11. Network Egress Specs

### 11.1 Default network posture

```yaml
network_default:
  outbound_tcp: deny
  outbound_udp: deny
  dns: deny
  metadata_service: deny
  internal_networks: deny
  public_internet: deny
```

### 11.2 Allowed egress profiles

```yaml
egress_profiles:
  none:
    outbound: deny_all

  web_fetch_proxy_only:
    outbound: proxy_only
    allowed_methods: [GET, HEAD]
    raw_sockets: deny
    dns: proxy_resolved_only
    private_ip_ranges: deny
    metadata_service: deny
    credential_forwarding: deny_by_default

  connector_proxy_only:
    outbound: layer7_connector_proxy_only
    raw_sockets: deny
    dns: connector_controlled_only
    private_ip_ranges: deny_unless_connector_declares
    metadata_service: deny

  vetted_registry_only:
    outbound: registry_proxy_only
    package_install: hash_pinned_only
    unpinned_versions: deny
    typosquat_guard: enforce
```

### 11.3 URL scope enforcement

For web-facing SEO/AEO work, Layer 6 may receive a domain or URL scope from Layer 3 through an already-authorized request. Layer 6 enforces the scope but does not choose it.

Checks:

```text
scheme must be http or https unless explicitly allowed
host must match approved domain scope
redirect chain must remain in allowed scope unless approved
private IP resolution is denied
localhost and link-local are denied
cloud metadata endpoints are denied
credential-bearing headers are stripped unless connector policy requires them
cookies are not persisted between runs
```

### 11.4 Crawl/render budget

```yaml
crawl_budget_default:
  max_urls_per_task: 25
  max_redirects_per_url: 5
  max_response_bytes: 5242880
  max_total_network_mb: 100
  max_wall_clock_ms: 120000
  respect_robots_policy_ref: layer3_or_tool_contract
```

Layer 6 enforces budgets. It does not decide crawl strategy.

### 11.5 Egress denial events

Layer 6 emits sanitized denial metadata:

```json
{
  "event_type": "runtime_egress_blocked",
  "trace_id": "...",
  "task_id": "...",
  "reason": "private_ip_range_denied",
  "runtime_class": "read_only_fetch",
  "severity": "medium"
}
```

It must not emit full URLs if URLs contain tokens, secrets, query PII, or protected internal hosts.

---

## 12. Supply-Chain Defense Specs

### 12.1 Dynamic install blocker

Layer 6 must detect and block common package manager invocations from generated code unless an explicit dependency policy allows them.

Blocked by default:

```text
npm install
pnpm add
yarn add
pip install
pipenv install
poetry add
apt install
apk add
brew install
go get
cargo add
composer require
curl | sh
wget | sh
remote script execution
```

### 12.2 Vetted dependency exception

A dependency may be loaded only if all are true:

```text
package name is in dependency_allowlist.yaml
version is pinned
hash is known
registry is trusted
license is allowed
vulnerability policy passes
package name confusion guard passes
SBOM entry is recorded
request envelope permits dependency access
runtime class permits dependency access
```

### 12.3 Hallucinated package defense

Package confusion checks:

```text
exact package allowlist match
namespace/organization validation
Levenshtein similarity to high-value packages
newly-created package age gate
maintainer trust check
registry provenance check
download spike anomaly check
hash mismatch block
unexpected transitive dependency block
```

Layer 6 must treat unknown package names as hostile by default.

### 12.4 Runtime SBOM

For every execution that loads packages or binaries, Layer 6 emits a sanitized runtime SBOM summary to Layer 8.

```go
type RuntimeSBOMSummary struct {
    TraceID      string   `json:"trace_id"`
    TaskID       string   `json:"task_id"`
    RuntimeImage string   `json:"runtime_image"`
    Packages     []string `json:"packages"`
    Hashes       []string `json:"hashes,omitempty"`
    PolicyState  string   `json:"policy_state"`
}
```

No raw secrets, private registry tokens, or internal registry credentials may be emitted.

---

## 13. IAM Context and JIT Credential Specs

### 13.1 Ambient credential stripping

Before sandbox execution, Layer 6 removes:

```text
cloud default credentials
service account tokens
user OAuth tokens
SSH keys
Git credentials
browser cookies
API tokens
package registry tokens
database credentials
Kubernetes service account tokens
metadata service access
inherited environment secrets
```

The sandbox starts with no useful credentials.

### 13.2 JIT token broker

Layer 6 may attach a credential only if:

```text
Layer 2 policy approved the tool/action
runtime request references a credential policy
credential scope matches task, tool, tenant, environment, and data source
TTL is short and bounded to task completion
audience is restricted
token is non-exportable when possible
token is unavailable to child processes unless required
token is revoked after execution
```

### 13.3 JIT scope examples

```yaml
jit_scopes:
  read_search_performance:
    access: read_only
    dataset: search_performance
    ttl_seconds: 300
    write: false

  read_analytics_report:
    access: read_only
    dataset: analytics_report
    ttl_seconds: 300
    write: false

  crawl_public_site:
    access: fetch_only
    domains: declared_target_domain_only
    ttl_seconds: 180
    write: false

  write_approved_canvas_artifact:
    access: write_only
    target: declared_canvas_artifact
    ttl_seconds: 120
    requires_approval_ref: true
```

### 13.4 Credential lifecycle

```text
resolve scope
request token
attach token to approved adapter only
execute action
remove token from environment
revoke token
verify revocation
redact token references from output
emit token lifecycle metadata without token value
```

### 13.5 Confused Deputy prevention

Layer 6 must not allow an agent to use the human user's broad credentials for a narrower task.

Reject if:

```text
credential scope is broader than requested action
credential audience does not match connector/tool
credential TTL exceeds task maximum
credential allows write when action is read-only
credential tenant does not match task tenant
credential environment does not match environment
token appears in params, stdout, stderr, artifact, URL, or logs
```

---

## 14. Tool and Script Execution Gateway

### 14.1 Registered adapters only

Layer 6 does not execute arbitrary commands. It dispatches through registered runtime adapters.

```yaml
adapters:
  web_fetch_adapter:
    runtime_class: read_only_fetch
    network: web_fetch_proxy_only
    shell: false

  page_render_adapter:
    runtime_class: headless_browser_runner
    network: web_fetch_proxy_only
    shell: false

  skill_script_adapter:
    runtime_class: deterministic_skill_script
    network: none_by_default
    shell: constrained

  schema_validator_adapter:
    runtime_class: validation_runner
    network: none_by_default
    shell: false

  file_patch_adapter:
    runtime_class: safe_file_patch
    network: none
    shell: false

  connector_invocation_adapter:
    runtime_class: analytics_connector_call
    network: connector_proxy_only
    shell: false
```

### 14.2 Script runner constraints

Scripts may run only if:

```text
script reference came from Layer 4 as a declared bundled resource
script hash matches expected value
script interpreter is allowed
script has no dynamic dependency install
script has no direct network access unless request permits it
script runs in clean temp working directory
script receives inputs via files or sanitized JSON only
script output is bounded and redacted
script exits within budget
```

### 14.3 Inline code policy

Inline code from the model is denied by default.

Exceptions require:

```text
runtime_class=quarantine_repair or validation_runner
Layer 2 explicit policy approval
small code size
no network
no credentials
read/write limits
post-run redaction
forensic capture
```

### 14.4 Headless browser constraints

Headless browsers must run with:

```text
sandboxed browser profile
no persistent cookies
no extension installation
no file download except approved artifacts
no clipboard access
no microphone/camera
no local network access
no metadata service access
request interception through proxy
navigation budget
screenshot size limit
DOM output sanitization
```

### 14.5 Retry policy

Retries are allowed only when the failure class is retryable.

Retryable:

```text
transient network timeout
connector 429 with retry-after
temporary browser render failure
validator temporary file lock
```

Not retryable:

```text
policy denial
path traversal
credential violation
egress violation
dynamic install attempt
secret leak detected
quarantine active
```

Retries must create a fresh sandbox.

---

## 15. Runtime Output and Redaction Specs

### 15.1 Captured output

Layer 6 captures:

```text
exit code
sanitized stdout
sanitized stderr
resource usage
artifact hashes
blocked policy outcomes
runtime class
adapter name
high-level execution summary
```

Layer 6 must not expose:

```text
raw tokens
raw cookies
raw credentials
raw PII
absolute host paths
internal filesystem layout
source implementation paths
MCP server internals
policy file contents
hidden prompts
chain-of-thought
raw trace internals
```

### 15.2 Redaction sequence

```text
capture raw stream in sandbox-local restricted buffer
scan for secrets/tokens/cookies/PII/protected metadata
redact stream
truncate to output budget
stage sanitized stream
emit sanitized reference
secure-delete raw buffer during teardown unless quarantine preservation required
```

### 15.3 Output size limits

```yaml
output_budget_default:
  stdout_max_bytes: 65536
  stderr_max_bytes: 65536
  artifact_max_bytes: 10485760
  artifact_count_max: 10
  summary_max_chars: 2000
```

Oversized outputs are truncated and marked:

```json
{
  "status": "success",
  "redaction_applied": true,
  "policy_outcomes": ["stdout_truncated"]
}
```

---

## 16. Runtime AgBOM Contribution

Layer 8 owns the Runtime Agent Bill of Materials, but Layer 6 contributes execution facts.

Layer 6 emits:

```text
runtime_class
adapter used
sandbox profile
credential scope class, not value
network profile
filesystem mount classes
package SBOM summary
artifact count and hashes
blocked runtime policy events
quarantine execution status
```

Layer 6 does not store the full AgBOM or analyze drift. Layer 8 owns that.

---

## 17. Stateful Quarantine Specs

### 17.1 Quarantine trigger source

Layer 6 accepts quarantine only from trusted signed commands.

Valid sources:

```text
Layer 8 Blue Team ABA / SecOps triad decision
Layer 8 observability circuit breaker
human security operator command routed through approved control channel
```

Invalid sources:

```text
agent self-command
model text instruction
raw user prompt
untrusted UI event
unverified webhook
unsigned service call
```

### 17.2 Quarantine actions

```yaml
quarantine_actions:
  freeze_runtime:
    effect: "stop starting new executions"

  revoke_jit_tokens:
    effect: "immediately revoke active scoped credentials"

  block_egress:
    effect: "deny all outbound network except status reporting"

  revoke_tool_access:
    effect: "mark session tool leases inactive"

  preserve_forensic_snapshot:
    effect: "capture sanitized runtime state, process metadata, artifact hashes, and policy events"

  allow_read_only_status:
    effect: "permit status-only reporting to orchestrator and UI"

  start_repair_scope:
    effect: "create isolated repair workspace with no credentials and no default network"
```

### 17.3 State preservation

Preserve:

```text
runtime request envelope metadata
policy decision refs
sandbox profile
resource usage
artifact hashes
sanitized stdout/stderr
process tree metadata
package SBOM summary
network denial summaries
file operation summaries
```

Do not preserve in general telemetry:

```text
raw tokens
raw credentials
cookies
raw PII
full hidden prompts
private source code not needed for forensics
```

Forensic preservation may retain restricted raw materials only in a separate access-controlled incident store, not in normal user-visible outputs.

### 17.4 Quarantine status result

```go
type QuarantineStatus struct {
    TraceID       string   `json:"trace_id"`
    SessionID     string   `json:"session_id"`
    Status        string   `json:"status"`
    ActionsTaken  []string `json:"actions_taken"`
    TokensRevoked bool     `json:"tokens_revoked"`
    EgressBlocked bool     `json:"egress_blocked"`
    StatePreserved bool    `json:"state_preserved"`
    SafeSummary    string   `json:"safe_summary"`
}
```

---

## 18. Auto-Refactoring and Repair Scope

### 18.1 Repair is constrained execution

Auto-refactoring is not general agent self-modification. It is a constrained repair workflow executed after a trusted quarantine or repair command.

Allowed repair cases:

```text
remove unsafe dependency invocation
replace dynamic install with pinned dependency reference
patch generated script to remove network call
patch generated script to reduce filesystem scope
patch generated script to validate inputs
patch generated script to redact outputs
add deterministic guard around known unsafe pattern
```

Forbidden repair cases:

```text
expand tool permissions
weaken policy checks
remove sandboxing
add credentials
increase network access
modify Layer 2 policies
modify Layer 3 workflows
modify Layer 4 skill bodies without skill-review workflow
modify Layer 7 connectors
modify Layer 8 evals
ship directly to production without approval
```

### 18.2 Repair output

Repair output must be a proposed patch artifact, not an automatic production change, unless the request is explicitly scoped to a disposable generated script.

```text
repair patch
reason summary
risk summary
tests run
remaining blockers
requires human review flag
```

---

## 19. Security Policy Matrix

| Risk | Runtime control | Default behavior |
|---|---|---|
| Container escape | gVisor/rootless/seccomp/namespaces | Block |
| Host filesystem access | deny-by-default mounts | Block |
| Path traversal | canonicalization + symlink/hardlink guard | Block |
| Credential leakage | ambient stripping + output redaction | Block and redact |
| Confused deputy | task-scoped JIT credentials | Block broad scope |
| Dynamic dependency install | package manager blocker | Block |
| Slopsquatting / hallucinated package | allowlist + hash pinning + registry trust | Block |
| Raw internet access | egress proxy only | Block |
| Metadata service exfiltration | metadata service blocker | Block |
| DNS exfiltration | DNS proxy/deny | Block |
| Large autonomous blast radius | resource and batch budgets | Block over-budget |
| Agent drift execution | Layer 8 quarantine command execution | Freeze |
| Unsafe repair | repair scope guard | Block |
| Raw protected output | redaction and truncation | Redact or block |

---

## 20. Integration Contracts with Other Layers

### 20.1 Layer 2 -> Layer 6

Layer 2 provides:

```text
policy-authorized execution envelope
sanitized parameters
policy decision ID
tool/action authorization state
approval metadata verification state
protected output constraints
```

Layer 6 returns:

```text
runtime policy outcome
execution result
runtime denial reason
credential lifecycle summary
sanitized event facts
```

Layer 6 must not reinterpret policy. If Layer 2 blocks, Layer 6 does not execute.

### 20.2 Layer 3 -> Layer 6

Layer 3 provides:

```text
workflow context
task ID
node ID
runtime class
tool/action contract reference
input artifact refs
execution sequencing metadata
```

Layer 6 returns:

```text
runtime result
artifact refs
retryability class
error class
sanitized execution summary
```

Layer 6 must not choose the next workflow step.

### 20.3 Layer 4 -> Layer 6

Layer 4 may provide declared script references and resource bundles through Layer 3/2-approved requests.

Layer 6 verifies:

```text
script ref exists in request
script hash matches expected hash
script access is read_execute only
script execution uses approved interpreter
```

Layer 6 must not inspect or load skill bodies for routing purposes.

### 20.4 Layer 5 -> Layer 6

Layer 5 does not directly execute runtime actions. It sends user decisions/events to Layer 2/3. Layer 6 receives only authorized execution requests.

Layer 6 may return status data that Layer 5 presents, but it must be sanitized.

### 20.5 Layer 6 -> Layer 7

Layer 7 owns MCP transport, connectors, data mesh, and memory retrieval. Layer 6 may run a constrained connector invocation adapter only through a Layer 7-approved connector proxy or local transport boundary.

Layer 6 must not:

```text
open arbitrary MCP connections
inspect connector internals
choose data sources
read vector stores directly
persist durable memory
bypass mTLS/data mesh controls
```

### 20.6 Layer 6 -> Layer 8

Layer 6 emits:

```text
execution started
execution completed
execution failed
runtime policy blocked
filesystem policy blocked
egress policy blocked
credential issued/revoked metadata
resource budget exceeded
supply-chain block
sandbox reset
quarantine action executed
repair scope executed
```

Layer 8 stores, analyzes, evaluates, detects drift, and decides future quarantine.

---

## 21. Event Emission

### 21.1 Runtime event schema

```go
type RuntimeEvent struct {
    EventID      string                 `json:"event_id"`
    TraceID      string                 `json:"trace_id"`
    SessionID    string                 `json:"session_id"`
    TaskID       string                 `json:"task_id"`
    EventType    string                 `json:"event_type"`
    RuntimeClass string                 `json:"runtime_class,omitempty"`
    Severity     string                 `json:"severity"`
    Summary      string                 `json:"summary"`
    Metadata     map[string]interface{} `json:"metadata,omitempty"`
    CreatedAt    string                 `json:"created_at"`
}
```

### 21.2 Required event types

```text
runtime_request_received
runtime_request_rejected
sandbox_lease_created
sandbox_reset_completed
execution_started
execution_completed
execution_failed
filesystem_access_denied
egress_access_denied
metadata_service_access_denied
dynamic_dependency_install_blocked
package_hash_mismatch_blocked
jit_token_issued_metadata
jit_token_revoked_metadata
ambient_credentials_stripped
resource_budget_exceeded
output_redaction_applied
artifact_staged
quarantine_command_received
quarantine_applied
repair_scope_started
repair_patch_generated
sandbox_teardown_completed
```

### 21.3 Event redaction

Events must not include:

```text
raw prompts
hidden chain-of-thought
raw user PII
raw secrets
token values
cookies
full filesystem paths
source code
private connector internals
raw stdout/stderr unless redacted and bounded
```

---

## 22. Resource Budget Policy

### 22.1 Default budgets

```yaml
resource_budgets:
  deterministic_skill_script_default:
    max_wall_clock_ms: 30000
    max_memory_mb: 512
    max_processes: 16
    max_output_bytes: 131072
    max_file_write_bytes: 1048576
    max_retries: 1

  web_fetch_default:
    max_wall_clock_ms: 120000
    max_memory_mb: 1024
    max_processes: 32
    max_output_bytes: 262144
    max_network_mb: 100
    max_retries: 2

  headless_browser_default:
    max_wall_clock_ms: 180000
    max_memory_mb: 2048
    max_processes: 64
    max_output_bytes: 524288
    max_network_mb: 150
    max_retries: 1

  file_patch_default:
    max_wall_clock_ms: 30000
    max_memory_mb: 256
    max_processes: 8
    max_output_bytes: 65536
    max_files_changed: 3
    max_lines_changed: 200
    max_retries: 0

  quarantine_repair_default:
    max_wall_clock_ms: 60000
    max_memory_mb: 512
    max_processes: 8
    max_output_bytes: 131072
    max_files_changed: 3
    max_retries: 0
```

### 22.2 Over-budget behavior

```text
soft budget warning -> emit runtime_budget_warning
hard budget exceeded -> terminate process
credential active -> revoke immediately
network active -> close connections
artifacts partial -> mark incomplete
result -> blocked_by_resource_budget
```

---

## 23. SEO/AEO and Content-Agent Runtime Policies

### 23.1 Read-only site audit actions

```yaml
site_audit_runtime:
  runtime_class: read_only_fetch
  network: web_fetch_proxy_only
  credentials: none
  filesystem: temp_plus_artifacts
  mutation: false
  max_urls_default: 25
  output: sanitized_findings_artifact
```

Layer 6 enforces read-only behavior. It must not write to the target site.

### 23.2 Search and analytics data actions

```yaml
search_analytics_runtime:
  runtime_class: analytics_connector_call
  network: connector_proxy_only
  credentials: jit_read_only_data_scope
  filesystem: temp_only
  mutation: false
  output: sanitized_report_artifact
```

Layer 6 attaches only read-scoped JIT tokens and revokes them after use.

### 23.3 Schema and XML validation

```yaml
schema_validation_runtime:
  runtime_class: validation_runner
  network: none_by_default
  credentials: none
  filesystem: read_input_write_result
  mutation: result_only
```

### 23.4 Internal-link graph analysis

```yaml
internal_link_graph_runtime:
  runtime_class: deterministic_skill_script
  network: none_by_default
  credentials: none
  filesystem: read_input_write_artifact
  mutation: artifact_only
```

### 23.5 Content artifact writes

```yaml
content_write_runtime:
  runtime_class: safe_file_patch
  network: none
  credentials: none_or_write_artifact_scope
  filesystem: declared_rw_artifact_only
  approval_required_for_overwrite: true
  mutation: true
```

Layer 6 writes only approved artifact targets, never arbitrary CMS or production files unless a separate approved connector/action exists.

---

## 24. Production Test Matrix

### 24.1 Boundary tests

| Test | Expected result |
|---|---|
| Raw prompt sent to Layer 6 | Reject |
| Unapproved tool request | Reject |
| Missing policy decision ID | Reject |
| Expired execution envelope | Reject |
| Runtime class mismatch | Reject |
| Layer 6 asked to classify intent | Reject |
| Layer 6 asked to choose workflow | Reject |
| Layer 6 asked to load SKILL.md | Reject |
| Layer 6 asked to open MCP directly | Reject |
| Layer 6 asked to render UI | Reject |
| Layer 6 asked to store telemetry | Reject |

### 24.2 Sandbox tests

| Test | Expected result |
|---|---|
| Container escape attempt | Block |
| Host filesystem probe | Block |
| Root user request | Block |
| Privileged container request | Block |
| Docker socket mount | Block |
| Kubernetes service token access | Block |
| Process fork bomb | Kill by budget |
| Memory exhaustion | Kill by cgroup |
| Long-running process | Timeout |

### 24.3 Filesystem tests

| Test | Expected result |
|---|---|
| `../` traversal | Block |
| Symlink outside mount | Block |
| Hardlink outside mount | Block |
| Absolute host path read | Block |
| Secret file read | Block |
| Production manifest write | Block |
| Approved artifact write | Allow |
| Overwrite without approval | Block |
| Oversized artifact | Block or truncate per policy |

### 24.4 Egress tests

| Test | Expected result |
|---|---|
| Direct socket from script | Block |
| curl to public internet | Block unless proxy adapter |
| DNS exfiltration attempt | Block |
| Metadata service call | Block |
| Private IP fetch | Block |
| Redirect to private IP | Block |
| Approved domain fetch through proxy | Allow |
| Connector call through Layer 7 proxy | Allow if approved |

### 24.5 Supply-chain tests

| Test | Expected result |
|---|---|
| `npm install` from generated script | Block |
| `pip install` from generated script | Block |
| Typosquatted package | Block |
| Unpinned package version | Block |
| Hash mismatch | Block |
| Untrusted registry | Block |
| Approved pinned package | Allow only through registry policy |
| SBOM generation | Required |

### 24.6 JIT credential tests

| Test | Expected result |
|---|---|
| Ambient credential present | Strip |
| Broad user token requested | Block |
| Read-only action asks write token | Block |
| Token TTL exceeds max | Block |
| Token appears in stdout | Redact and flag |
| Token not revoked after task | Fail test |
| Quarantine active with token | Revoke immediately |

### 24.7 Quarantine tests

| Test | Expected result |
|---|---|
| Unsigned quarantine command | Reject |
| Signed quarantine command | Freeze runtime |
| Active token during quarantine | Revoke |
| Network active during quarantine | Block |
| Approved user canvas state | Preserve |
| Forensic snapshot | Preserve sanitized state |
| New execution during quarantine | Reject |
| Repair scope tries to expand permission | Block |

---

## 25. Deployment Requirements

### 25.1 Infrastructure

```text
gVisor or equivalent sandbox runtime enabled
rootless runtime enforced
seccomp profiles installed
cgroup limits enabled
network namespaces enabled
egress proxy deployed
metadata service blocked from sandbox
container image signing enabled
runtime images pinned
policy/config mounted read-only
runtime event bus configured
secret manager integration supports JIT token issuance
emergency revocation path tested
quarantine control channel signed and authenticated
```

### 25.2 CI/CD gates

Required gates:

```text
runtime schema tests pass
sandbox escape tests pass
filesystem boundary tests pass
egress policy tests pass
supply-chain blocking tests pass
JIT token lifecycle tests pass
output redaction tests pass
quarantine tests pass
boundary non-overlap tests pass
load tests within resource budget pass
```

### 25.3 Rollout stages

```text
local sandbox dry run
staging with no credentials
staging with read-only JIT credentials
shadow mode for runtime events
limited canary for read-only actions
limited canary for artifact writes
human-approved production rollout
continuous drift and runtime policy monitoring
```

---

## 26. Acceptance Criteria

Layer 6 is production-ready only when all are true:

```text
No runtime action executes without a signed approved envelope.
No sandbox starts with ambient credentials.
No untrusted process runs outside kernel-level isolation.
No undeclared filesystem access is possible.
No raw network access is possible from untrusted scripts.
No dynamic dependency installation succeeds by default.
No JIT token outlives its task.
No JIT token value appears in returned output or events.
No approved artifact exposes host filesystem paths.
No quarantine command can be spoofed.
No quarantine destroys required forensic state.
No repair scope can expand privileges.
All runtime events are sanitized before Layer 8 ingestion.
All boundary tests prove Layer 6 does not classify, orchestrate, render, load skills, own MCP, store telemetry, or evaluate quality.
```

---

## 27. Final Non-Goals

Layer 6 must not become:

```text
an intent classifier
a policy engine
a workflow orchestrator
a skill router
a skill registry
a UI renderer
a connector implementation layer
a memory/RAG layer
a telemetry database
an eval system
a drift detector
a general shell service
a developer workstation proxy
a generic web browser for agents
a package manager gateway for arbitrary generated code
a production deployment system
```

---

## 28. One-Line Architecture Summary

```text
Layer 6 is the zero-trust execution cell: it accepts only authorized actions, strips ambient authority, runs them in ephemeral constrained sandboxes, enforces filesystem/network/credential limits, emits sanitized runtime facts, and executes stateful recovery when Layer 8 tells it to freeze.
```
