# Agentic Engineering Service Specification

Status: Draft v1 (language-agnostic, domain-neutral)

Purpose: Define an executable, fail-closed platform for building, operating, securing, evaluating,
and governing domain agents.

## 1. Problem Statement

An agent can appear complete while containing only directories, schemas, interfaces, prompts, or
tests that never exercise real behavior. Compilation proves syntax. File presence proves structure.
Mocked tests prove a test harness. None of those facts prove that an agent can safely select a
workflow, load a skill, authorize a tool, execute it in isolation, retrieve tenant-scoped evidence,
produce a safe result, or pass a live evaluation.

This specification solves six engineering problems:

- It separates intake, orchestration, skills, presentation, runtime execution, integrations, and
  evaluation into explicit control layers with no silent ownership overlap.
- It turns security requirements into executable gates instead of comments or empty policy files.
- It makes tools, connectors, workflows, and skills independently versioned and testable.
- It prevents a skill from being called production-ready when its required tools or workflows do
  not yet exist.
- It provides domain-neutral evaluation contracts that work for software, biology, finance,
  research, operations, content, or another bounded domain.
- It distinguishes repository evidence from live-provider, infrastructure, canary, and human
  approval evidence.

The platform is designed for implementation by any capable coding agent or engineering team. It
does not depend on a particular model vendor, agent application, programming language, cloud, or
product domain.

Important boundary:

- This specification defines an agent platform and its executable gates.
- A domain implementation supplies its own intents, tools, connectors, workflows, skills, evidence
  rules, risk classes, output contracts, and domain-specific evaluations.
- A domain pack is incomplete until those artifacts exist as working implementations.
- Trigger routing can be evaluated before tools exist, but end-to-end skill execution cannot.
- No release gate may infer an execution pass from a schema, filename, mocked response, or empty
  implementation.

Normative language:

- `MUST` and `MUST NOT` are mandatory for conformance.
- `SHOULD` and `SHOULD NOT` are recommended unless a documented risk decision says otherwise.
- `MAY` is optional.
- `fail closed` means deny, block, quarantine, or require review when required evidence is missing,
  malformed, stale, ambiguous, unverifiable, or unavailable.

## 2. Goals and Non-Goals

### 2.1 Goals

- Accept untrusted user and system input through a bounded, normalized intake contract.
- Detect prompt injection and protected metadata before unsafe content reaches execution.
- Classify only into a closed, versioned intent enum and request clarification when ambiguous.
- Require an explicit authorization decision before every executable action.
- Route work through validated workflows represented as bounded directed acyclic graphs.
- Load versioned skills through checksummed manifests and progressive disclosure.
- Generate only schema-valid, signed, non-executable presentation contracts.
- Execute tools in an attested sandbox with bounded filesystem, network, dependency, credential,
  time, concurrency, and output budgets.
- Connect external systems through registered, tenant-isolated, provenance-preserving adapters.
- Record redacted traces, runtime bills of materials, drift, trust, cost, and policy outcomes.
- Evaluate routing, trajectories, outputs, safety, reliability, and live behavior with explicit
  pass criteria.
- Keep quarantine decisions separate from quarantine execution.
- Maintain append-only, signed governance evidence for promotion, rollback, and incident response.
- Make startup fail when required schemas, policies, registries, or embedded artifacts are invalid.
- Provide a reusable implementation order and conformance matrix for any domain agent.

### 2.2 Non-Goals

- Defining the domain knowledge of a specific agent.
- Inventing tools, connectors, or workflows from skill descriptions.
- Treating a model response as authorization.
- Treating a passing unit test as production certification.
- Mandating one programming language, framework, model provider, vector database, sandbox product,
  telemetry backend, or UI renderer.
- Storing hidden chain-of-thought or raw protected payloads for debugging.
- Allowing observability or evaluation code to execute tools or connectors directly.
- Allowing presentation code to authorize actions or mutate domain state.
- Allowing runtime code to choose workflows or classify raw intent.
- Automatically promoting skills, workflows, tools, connectors, or policies.
- Claiming that repository source code proves cloud, network, KMS, mTLS, sandbox, or canary
  deployment posture.

## 3. System Overview

### 3.1 Main Components

1. `Intake and Safety Gateway` (Layer 2)
   - Enforces request size and schema limits.
   - Normalizes text and identifiers.
   - Detects prompt injection and encoded instruction smuggling.
   - Classifies a closed intent or requests clarification.
   - Validates approval metadata.
   - Authorizes proposed actions.
   - Redacts protected outbound metadata.

2. `Orchestrator` (Layer 3)
   - Selects a registered workflow and execution profile.
   - Builds and validates a bounded DAG.
   - Owns short-lived run state, retries, cancellation, and idempotency.
   - Proposes tool requests and skill activation requests.
   - Produces presentation intents.
   - Never authorizes or executes an action.

3. `Skill Registry and Loader` (Layer 4)
   - Stores versioned skill metadata and `SKILL.md` bodies.
   - Stores references, assets, scripts-as-artifacts, resource manifests, and eval definitions.
   - Verifies schema, checksum, ownership, status, compatibility, and token budgets.
   - Loads only requested, compatible resources.
   - Never chooses workflows, scores evals, or executes scripts.

4. `Presentation and Agent Interchange` (Layer 5)
   - Builds schema-valid UI, approval, canvas, brief, chat, dashboard, and agent-to-agent contracts.
   - Enforces a locked component catalog and non-executable content.
   - Sanitizes Markdown, URLs, hidden payloads, and event data.
   - Signs time-bounded UI frames and verifies signed agent envelopes.
   - Never authorizes or executes actions.

5. `Sandboxed Runtime` (Layer 6)
   - Executes only Layer 2-authorized requests.
   - Creates attested, ephemeral sandboxes.
   - Enforces filesystem, process, dependency, network, credential, time, and output policy.
   - Strips ambient credentials and resolves short-lived scoped credentials.
   - Executes authenticated quarantine commands.
   - Never chooses workflows or owns connector policy.

6. `Interop and Data Plane` (Layer 7)
   - Owns MCP transport, connector adapters, external APIs, retrieval, provenance, and memory
     data-plane operations.
   - Enforces mTLS identity, token audience, tenant partitioning, field projection, schema hashes,
     taint labels, and evidence contracts.
   - Provides model-provider adapters for evaluation without giving credentials to Layer 8.
   - Never authorizes actions, loads skills, or executes arbitrary code.

7. `Observability, Evaluation, SecOps, and Governance` (Layer 8)
   - Ingests sanitized events and traces.
   - Tracks the runtime agent bill of materials.
   - Calculates drift, trust decay, loop, anomaly, cost, and reliability signals.
   - Scores trajectories, outputs, skill routing, safety, and pass^k reliability.
   - Runs bounded Red, Blue, and Green Team processes.
   - Writes signed, append-only governance evidence.
   - Decides or recommends quarantine but never executes it.

8. `Repository Readiness Aggregator`
   - Reads repository evidence from all owning layers.
   - Reports missing, empty, placeholder, malformed, or non-production artifacts.
   - Never performs runtime behavior and is not an additional architecture layer.

### 3.2 Control Layers

The platform is portable when implemented as seven isolated control layers:

1. `Safety and Authorization`
   - Input normalization, injection defense, closed intent, approval, action authorization,
     outbound redaction.

2. `Coordination`
   - Workflow selection, profile selection, DAG planning, run state, retries, proposals.

3. `Procedural Memory`
   - Skills, resource manifests, progressive disclosure, compatibility, static eval definitions.

4. `Presentation`
   - UI and agent interchange schemas, catalogs, signatures, sanitization, event normalization.

5. `Execution`
   - Sandbox, filesystem, egress, credentials, package policy, budgets, quarantine.

6. `Data and Integration`
   - MCP, APIs, retrieval, evidence, provenance, tenant isolation, memory operations.

7. `Glass-Box Assurance`
   - Traces, AgBOM, drift, trust, SecOps, evaluation, governance, retention, improvement proposals.

### 3.3 Required Cross-Layer Paths

Tool execution:

```text
Layer 3 proposes
  -> Layer 2 authorizes
  -> Layer 6 executes
  -> Layer 7 performs registered connector or data access when required
  -> Layer 8 observes and evaluates
```

Skill loading:

```text
Layer 3 requests activation
  -> Layer 4 validates and loads
  -> Layer 2 authorizes any resulting executable action
  -> Layer 6 executes authorized scripts or tools
  -> Layer 8 evaluates
```

Presentation:

```text
Layer 3 creates presentation intent
  -> Layer 5 builds and signs a presentation contract
  -> an external client renderer displays it
  -> Layer 8 observes sanitized events
```

Memory or preference update:

```text
Layer 2 validates approval
  -> Layer 3 plans
  -> Layer 5 presents the decision
  -> Layer 7 performs the tenant-scoped data operation
  -> Layer 8 audits
```

Quarantine:

```text
Layer 8 decides or recommends
  -> Layer 6 authenticates and executes
```

These paths MUST NOT be bypassed, including by administrative, retry, evaluation, or recovery code.

### 3.4 External Dependencies

A concrete implementation may require:

- One or more model providers.
- A sandbox runtime such as gVisor, a microVM, a hardened container, or an equivalent attested
  isolation boundary.
- A credential issuer, secret manager, KMS, or token exchange service.
- Connector endpoints, MCP servers, databases, vector stores, object stores, or domain APIs.
- An OpenTelemetry-compatible trace provider.
- A durable append-only governance store.
- A client renderer for presentation contracts.
- A CI environment for deterministic, live-provider, adversarial, race, vulnerability, chaos, and
  canary tests.

Repository code MUST treat each external dependency as unavailable until its identity, policy,
timeouts, limits, and failure behavior are explicitly configured and validated.

## 4. Core Domain Model

### 4.1 Entities

#### 4.1.1 Tenant Context

Security context attached to every request and data-plane operation.

Fields:

- `tenant_id` (opaque string)
  - Required for all tenant data.
- `subject_id` (opaque string)
  - User, service, or agent identity.
- `session_id` (opaque string)
- `trace_id` (opaque string)
- `request_id` (opaque string)
- `roles` (list of closed enum values)
- `data_region` (string or null)
- `policy_version` (string)
- `issued_at` (timestamp)
- `expires_at` (timestamp)

Rules:

- Empty tenant identifiers are invalid when data is involved.
- Tenant identifiers MUST be compared exactly and SHOULD use constant-time comparison where
  observable timing could matter.
- Raw user identity or PII SHOULD NOT be used as an internal identifier.

#### 4.1.2 Request Envelope

Bounded request received by Layer 2.

Fields:

- `schema_version`
- `request_id`
- `trace_id`
- `tenant_context`
- `input_kind` (closed enum)
- `input_text` (string or null)
- `input_payload` (schema-bound object or null)
- `selected_context_refs` (list of opaque references)
- `requested_mode` (closed enum or null)
- `approval_token` (opaque string or null)
- `client_capabilities` (allowlisted object)
- `received_at`

The envelope MUST NOT contain ambient credentials, arbitrary executable code, filesystem paths,
provider endpoints, or unbounded binary data unless a separately authorized input contract permits
them.

#### 4.1.3 Intake Decision

Layer 2 result after normalization, safety inspection, and intent classification.

Fields:

- `request_id`
- `trace_id`
- `normalized_input_ref`
- `intent` (closed enum)
- `mode` (closed enum)
- `confidence`
- `decision` (`allow`, `clarify`, or `block`)
- `reason_code` (closed enum)
- `risk_class`
- `sanitized_context_refs`
- `policy_version`
- `decided_at`

Raw prompts and hidden reasoning MUST NOT appear in this object.

#### 4.1.4 Authorization Decision

Layer 2 decision for one proposed action.

Fields:

- `decision_id`
- `trace_id`
- `tenant_id`
- `action_class`
- `resource_scope`
- `connector_id` (or null)
- `tool_id` (or null)
- `workflow_version`
- `skill_versions` (map)
- `approval_requirement`
- `approval_evidence_ref` (or null)
- `allowed` (boolean)
- `denial_code` (or null)
- `constraints`
  - allowed domains
  - allowed methods
  - allowed fields
  - maximum bytes
  - timeout
  - retry limit
  - credential scopes
- `issued_at`
- `expires_at`
- `signature`

Authorization MUST bind exactly to the tenant, trace, action class, resource, connector, tool, and
time window. A decision for one action MUST NOT authorize another.

#### 4.1.5 Tool Definition

Executable capability registered by the domain implementation.

Fields:

- `tool_id`
- `name`
- `version`
- `owner`
- `status`
- `risk_class`
- `action_class`
- `description`
- `input_schema`
- `output_schema`
- `error_schema`
- `side_effect_class`
- `idempotency_contract`
- `timeout_ms`
- `retry_policy`
- `required_credential_scopes`
- `allowed_runtime_class`
- `allowed_connectors`
- `data_classification`
- `redaction_contract`
- `implementation_ref`
- `test_manifest_ref`
- `checksum`

`implementation_ref` MUST resolve to executable code or a registered remote adapter. A schema-only
tool is not implemented.

#### 4.1.6 Workflow Definition

Versioned Layer 3 coordination contract.

Fields:

- `workflow_id`
- `version`
- `owner`
- `status`
- `compatible_intents`
- `compatible_modes`
- `compatible_profiles`
- `entry_node`
- `nodes`
- `edges`
- `maximum_nodes`
- `maximum_replans`
- `timeout_ms`
- `failure_policy`
- `required_tools`
- `optional_tools`
- `required_skills`
- `presentation_outcomes`
- `eval_manifest_ref`
- `checksum`

Each node declares a logical action class, dependencies, input references, output references, and
failure behavior. Workflow definitions MUST NOT embed credentials or bypass authorization.

#### 4.1.7 DAG Plan

One bounded workflow instance.

Fields:

- `plan_id`
- `workflow_id`
- `workflow_version`
- `trace_id`
- `entry_node`
- `nodes`
  - `node_id`
  - `action_class`
  - `depends_on`
  - `input_refs`
  - `expected_output_contract`
  - `timeout_ms`
  - `idempotency_key`
- `created_at`
- `expires_at`

The plan MUST be acyclic, reachable from its entry node, bounded, and free of unknown references.

#### 4.1.8 Skill Definition

Versioned procedural-memory artifact.

Fields:

- `skill_id`
- `name`
- `description`
- `version`
- `owner`
- `status`
- `tier` (`read`, `draft`, or `act`)
- `risk_class`
- `compatible_profiles`
- `compatible_intents`
- `allowed_modes`
- `capability_tags`
- `declared_action_classes`
- `required_tools`
- `required_workflows`
- `output_contracts`
- `anti_triggers`
- `token_budget`
- `resource_manifest_ref`
- `eval_manifest_ref`
- `body_checksum`

The description MUST state when to use the skill and when not to use it. A skill definition MUST
NOT claim a tool or workflow that is absent from its registered dependencies.

#### 4.1.9 Skill Activation Request and Bundle

Layer 3 requests activation; Layer 4 validates and returns a bounded bundle.

Request fields:

- `activation_id`
- `trace_id`
- `skill_id`
- `required_version`
- `intent`
- `mode`
- `profile`
- `requested_resources`
- `token_budget`

Bundle fields:

- `activation_id`
- `skill_id`
- `version`
- `body`
- `resources`
- `declared_action_classes`
- `required_tools`
- `required_workflows`
- `total_tokens`
- `checksum`
- `loaded_at`
- `expires_at`

No skill path, resource path, or exact internal inventory is user-facing.

#### 4.1.10 Presentation Intent and Frame

Layer 3 emits a semantic intent; Layer 5 returns a signed data contract.

Presentation intent fields:

- `presentation_id`
- `trace_id`
- `surface`
- `purpose`
- `data_refs`
- `allowed_actions`
- `approval_state`
- `expires_at`

Presentation frame fields:

- `schema_version`
- `frame_id`
- `surface`
- `root_node`
- `nodes`
- `events`
- `approval_contract` (or null)
- `issued_at`
- `expires_at`
- `signature`

Presentation content MUST be non-executable and valid against a locked catalog.

#### 4.1.11 Runtime Execution Request and Result

Layer 6 execution contract.

Request fields:

- `execution_id`
- `trace_id`
- `tenant_id`
- `authorization_decision`
- `action_class`
- `tool_id`
- `tool_version`
- `input_ref`
- `runtime_class`
- `filesystem_policy`
- `egress_policy`
- `dependency_policy`
- `credential_request`
- `budgets`
- `idempotency_key`
- `created_at`
- `expires_at`
- `signature`

Result fields:

- `execution_id`
- `status`
- `output_ref` (or null)
- `structured_error` (or null)
- `sandbox_attestation_ref`
- `credential_lease_ref` (or null)
- `egress_summary`
- `resource_usage`
- `started_at`
- `completed_at`
- `result_signature`

Raw credentials, arbitrary filesystem contents, and unredacted provider bodies MUST NOT be returned.

#### 4.1.12 Connector Manifest and Evidence Packet

Connector manifest fields:

- `connector_id`
- `version`
- `owner`
- `status`
- `transport`
- `authentication`
- `tenant_isolation`
- `read_only`
- `allowed_action_classes`
- `required_runtime_class`
- `source_contracts`
- `maximum_response_bytes`
- `timeout_ms`
- `redirect_policy`
- `provenance_required`
- `taint_scan_required`
- `schema_hashes`
- `checksum`

Evidence packet fields:

- `evidence_id`
- `tenant_id`
- `connector_id`
- `source_contract_id`
- `retrieved_at`
- `content_hash`
- `projected_fields`
- `provenance`
- `taint_labels`
- `freshness`
- `signature`

Evidence packets contain allowlisted fields only.

#### 4.1.13 Trace and Runtime Agent Bill of Materials

Trace fields:

- `trace_id`
- `tenant_ref_hash`
- `intent`
- `mode`
- `profile_ref_hash`
- `root_span`
- `child_spans`
- `outcome`
- `policy_outcomes`
- `cost`
- `latency`
- `token_usage`
- `trust_score`
- `drift_score`
- `created_at`

Runtime Agent Bill of Materials (`AgBOM`) fields:

- model and prompt versions
- workflow ID and version
- skill IDs and versions
- tool IDs and versions
- connector IDs and versions
- policy versions
- schema hashes
- sandbox image digest
- dependency digests
- credential lease references
- presentation schema version
- evaluation versions

Raw prompts, hidden reasoning, credentials, raw selected text, raw document bodies, and raw PII are
forbidden.

#### 4.1.14 Evaluation Manifest and Evidence

Evaluation manifest fields:

- `eval_id`
- `version`
- `owner`
- `target_kind`
- `target_id`
- `target_version`
- `corpus_checksum`
- `required_fixture_types`
- `trajectory_mode`
- `pass_k`
- `minimum_confidence`
- `rubric_id`
- `judge_model`
- `judge_prompt_version`
- `safety_cases`
- `promotion_thresholds`

Evaluation evidence fields:

- `evidence_id`
- `eval_id`
- `target_id`
- `target_version`
- `corpus_checksum`
- `model`
- `prompt_version`
- `runs`
- `passed_cases`
- `failed_cases`
- `safety_failures`
- `flake_rate`
- `minimum_confidence`
- `tool_versions`
- `workflow_versions`
- `environment_attestation_refs`
- `started_at`
- `completed_at`
- `signature`

Evidence MUST identify the exact implementation and corpus under test.

#### 4.1.15 Quarantine Decision and Command

Layer 8 produces a decision; Layer 6 accepts only an authenticated command.

Decision fields:

- `decision_id`
- `trace_id`
- `tenant_id`
- `reason_code`
- `severity`
- `evidence_refs`
- `policy_version`
- `recommended_actions`
- `requires_human_review`
- `decided_at`

Command fields:

- `command_id`
- `decision_id`
- `trace_id`
- `tenant_id`
- `actions`
- `issued_at`
- `expires_at`
- `nonce`
- `signature`

Layer 8 MUST NOT execute the command.

#### 4.1.16 Governance Record

Append-only release and incident evidence.

Fields:

- `record_id`
- `record_type`
- `target_kind`
- `target_id`
- `target_version`
- `tenant_scope_hash` (when applicable)
- `evidence_ref_hashes`
- `decision`
- `actor_ref_hash`
- `approver_ref_hash` (or null)
- `previous_record_hash`
- `record_hash`
- `signature`
- `created_at`

### 4.2 Stable Identifiers and Normalization Rules

- Identifiers are opaque, non-empty, length-bounded, and schema-validated.
- Wire fields MAY use `snake_case`; implementation symbols follow the host language conventions.
- Unicode user text is normalized with NFKC before security matching.
- Security checks preserve the original input only through a protected ephemeral reference.
- URLs are parsed with a standards-compliant parser; string-prefix URL checks are insufficient.
- Domains are converted to a canonical lower-case form and matched exactly or by explicit
  registered suffix policy.
- Schema, artifact, prompt, image, dependency, and evidence digests use a collision-resistant hash
  such as SHA-256.
- Tenant vector namespaces use an HMAC or equivalent keyed derivation, not a raw tenant ID.
- Idempotency keys bind tenant, action, resource, tool version, and intended side effect.
- Timestamps use UTC and are checked against bounded clock skew.
- Nonces are single-use within an explicit replay window.
- Error messages are stable, lower-case machine categories plus a safe operator summary.

## 5. Repository Specification (Domain Contract)

### 5.1 Artifact Layout

A conforming domain pack SHOULD use an equivalent version-controlled layout:

```text
agentic/
  schemas/
  policies/
    intake/
    runtime/
    integration/
    observability/
    evaluation/
    governance/
  tools/
    tool-registry.yaml
    <tool-id>/
      manifest.yaml
      input.schema.json
      output.schema.json
      errors.schema.json
      implementation/
      tests/
  workflows/
    workflow-registry.yaml
    <workflow-id>/
      workflow.yaml
      eval-manifest.yaml
      evals/
  skills/
    skill-registry.yaml
    <skill-id>/
      SKILL.md
      OWNERS
      CHANGELOG.md
      resource-manifest.yaml
      eval-manifest.yaml
      references/
      assets/
      scripts/
      evals/
  connectors/
    connector-registry.yaml
    source-contracts.yaml
    manifests/
  presentation/
    schemas/
    catalog-lock.yaml
  evals/
    integration/
    end-to-end/
    adversarial/
    canary/
```

Equivalent layouts are allowed, but ownership and dependency direction MUST remain explicit.

### 5.2 Artifact Status Model

Allowed status values:

1. `blocked`
   - Artifact is known but cannot be used.

2. `experimental`
   - Schema-valid and available only in development or evaluation.

3. `canary`
   - Approved for bounded production traffic with rollback.

4. `active`
   - Approved for its declared production scope.

5. `quarantined`
   - Disabled because of a security, safety, quality, or integrity signal.

6. `deprecated`
   - Available only for controlled migration.

7. `retired`
   - Unavailable for new use.

Status transitions require signed governance evidence. A registry edit alone MUST NOT promote an
artifact.

### 5.3 Tool Manifest Contract

Every tool manifest MUST:

- Validate against a versioned strict schema.
- Reference an implementation and tests.
- Declare input, output, and structured error schemas.
- Declare side effects and idempotency.
- Declare exact action class, runtime class, connector dependencies, and credential scopes.
- Declare time, size, concurrency, and retry limits.
- Declare data classification and redaction behavior.
- Declare whether human approval is required.
- Include owner, version, status, checksum, and changelog.

Tool readiness levels:

- `schema_ready`
  - Schemas validate; no execution claim.
- `implementation_ready`
  - Real implementation exists and unit tests pass.
- `integration_ready`
  - Real connector/runtime dependencies pass integration tests.
- `production_ready`
  - Security, reliability, canary, governance, and human gates pass.

### 5.4 Workflow Manifest Contract

Every workflow MUST:

- Reference only registered action classes, tools, skills, and output contracts.
- Declare a single entry node.
- Contain no duplicate, self-referential, unknown, cyclic, or unreachable nodes.
- Declare maximum nodes, replans, retries, and elapsed time.
- Declare compensation or safe failure behavior for side effects.
- Carry an eval manifest and checksum.
- Refuse activation when a required tool or skill is not at the required status.

Workflow definitions coordinate capabilities; they do not authorize them.

### 5.5 Skill File Format

Each `SKILL.md` is Markdown with strict YAML front matter.

Required front matter:

- `name`
- `description`
- `version`
- `owner`
- `status`
- `tier`
- `risk_class`
- `compatible_profiles`
- `compatible_intents`
- `allowed_modes`
- `capability_tags`
- `declared_action_classes`
- `required_tools`
- `required_workflows`
- `output_contracts`
- `token_budget`
- `resource_manifest`
- `eval_manifest`

Required body sections:

- `Purpose`
- `When to use`
- `When NOT to use`
- `Inputs expected`
- `Procedure`
- `Output contract`
- `Quality gates`
- `Boundary rules`
- `Resources`
- `Failure behavior`

The body is procedural guidance, not executable authority.

### 5.6 Resource Manifest Contract

The resource manifest enumerates every optional reference, asset, and script.

Each entry declares:

- logical resource ID
- relative artifact path
- media type
- purpose
- checksum
- size
- token estimate where applicable
- compatible skill version
- sensitivity class
- load condition
- executable boolean

Rules:

- Undeclared resources MUST NOT load.
- Path traversal and symlink escape MUST be rejected.
- Executable resources remain inert artifacts until separately authorized and run by Layer 6.
- Resource contents are loaded only when requested and within the active token budget.

### 5.7 Evaluation Manifest Contract

Every production-targeted skill and workflow MUST define:

- positive trigger cases
- negative trigger cases
- rephrasing-stability cases
- adjacent-skill collision cases
- out-of-scope and unsafe-side-effect cases
- golden output cases where output is produced
- trajectory cases using action classes
- a weighted rubric
- regression cases
- applicable tool, connector, and workflow dependencies
- required pass^k
- live-provider requirements
- promotion thresholds

Minimum trigger corpus per skill:

- at least three positive cases
- at least three negative cases
- at least one rephrase for each positive case
- at least one adjacent collision case
- at least one out-of-scope case

These are minimum quantities, not proof of domain quality.

### 5.8 Dependency Truth and Readiness Graph

The dependency graph is:

```text
schemas and policies
  -> tool implementations
  -> connector and runtime integrations
  -> workflow implementations
  -> skill activation and routing
  -> end-to-end skill execution
  -> live evaluation
  -> canary
  -> production promotion
```

Mandatory rules:

- A workflow requiring a missing tool is `blocked`.
- A skill requiring a missing workflow or tool is `blocked`.
- A trigger-routing eval MAY run against descriptions before tools exist.
- A golden-output eval MUST NOT pass until the real workflow and tools produce the output.
- A trajectory eval MUST NOT pass against a fabricated tool trace.
- A live model judge MUST NOT convert missing execution evidence into a pass.
- A readiness report MUST distinguish `not_run`, `blocked`, `failed`, and `passed`.
- `not_run` and `blocked` are not passes.

### 5.9 Checksums, Ownership, and Change Control

- Every production artifact has a checksum over canonical content.
- Registries carry the expected checksum.
- Startup and CI reject checksum mismatch.
- Every artifact has an owner and changelog.
- Version changes are mandatory for behavior or contract changes.
- Evaluation evidence binds the exact artifact versions and corpus checksum.
- Policy, skill, workflow, tool, connector, and schema changes require human review according to
  risk class.

### 5.10 Repository Validation and Error Surface

Normalized error categories:

- `schema_missing`
- `schema_invalid`
- `policy_missing`
- `policy_invalid`
- `registry_missing`
- `registry_invalid`
- `artifact_unregistered`
- `artifact_missing`
- `artifact_empty`
- `artifact_placeholder`
- `checksum_mismatch`
- `owner_missing`
- `changelog_missing`
- `dependency_missing`
- `dependency_not_ready`
- `eval_manifest_missing`
- `eval_evidence_missing`
- `protected_metadata_detected`
- `token_budget_exceeded`

Repository validation is read-only. It MUST NOT activate artifacts, call providers, execute tools,
connect external systems, score live evals, or mutate status.

## 6. Configuration Specification

### 6.1 Source Precedence and Resolution Semantics

Configuration precedence:

1. Explicit process or deployment configuration.
2. Versioned domain-pack configuration.
3. Environment-variable indirection for secrets and deployment-specific values.
4. Built-in secure defaults.

Rules:

- Secrets are supplied only through an approved secret source.
- Secret values MUST NOT appear in configuration files, command arguments, URLs, reports, or logs.
- An unresolved required environment variable is a startup error.
- An empty value is not equivalent to a configured value.
- Paths are canonicalized before policy checks.
- URIs and shell strings are not modified by filesystem path expansion.
- Unknown security-sensitive fields are rejected.

### 6.2 Startup Validation

Before opening a listener, accepting a job, or starting a worker, the process MUST:

1. Strictly decode root routing and authorization policy.
2. Compile every JSON Schema used by intake, skills, presentation, runtime, connectors, events, and
   evaluation.
3. Validate cross-field semantic security invariants.
4. Validate embedded or deployed registries.
5. Validate checksums, manifests, owners, changelogs, and required eval fixtures.
6. Validate connector and source contracts.
7. Validate telemetry, redaction, eval, SecOps, drift, trust, retention, and governance policies.
8. Verify that configured active artifacts have valid promotion evidence.
9. Fail closed on any required validation error.

Startup validation proves configuration integrity, not live dependency health.

### 6.3 Policy Families

Required policy families:

- `routing`
  - closed intents, modes, confidence thresholds, clarification behavior
- `authorization`
  - action classes, risks, approvals, scope, semantic denials
- `runtime`
  - sandbox, filesystem, process, dependency, budgets
- `egress`
  - domains, methods, headers, DNS/IP rules, redirects, byte limits
- `credentials`
  - issuers, audiences, scopes, TTLs, revocation
- `connectors`
  - status, transport, identity, tenancy, source contracts
- `presentation`
  - schemas, component catalog, signing, replay window
- `telemetry`
  - spans, sampling, cost, latency, token accounting
- `redaction`
  - never-store, hash-only, summarized, protected metadata
- `evaluation`
  - trajectory modes, pass^k, rubric, confidence, judge behavior
- `SecOps`
  - Red, Blue, Green Team boundaries and thresholds
- `drift`
  - warning, block, and quarantine-decision thresholds
- `trust`
  - score range, decay, recovery, critical signals
- `retention`
  - data classes, deletion, legal hold, encryption
- `governance`
  - append-only ledger and release requirements

Every security policy uses a version and `fail_closed` mode.

### 6.4 Safe Reload Semantics

Implementations MAY support dynamic reload.

If supported:

- Parse and validate the complete candidate configuration before activation.
- Apply changes atomically.
- Keep the last known good configuration when validation fails.
- Do not silently change in-flight authorization decisions.
- Re-evaluate future dispatch and activation against the new version.
- Emit a sanitized configuration-change event and governance record.
- Require restart when a listener, trust root, sandbox driver, or cryptographic key cannot be safely
  rebound.

### 6.5 Configuration Summary (Cheat Sheet)

This section is intentionally redundant so an implementation agent can build the configuration
layer quickly.

- `policy_mode`: required, must be `fail_closed`
- `routing.allowed_intents`: non-empty closed enum
- `routing.allowed_modes`: non-empty closed enum
- `routing.minimum_confidence`: number in `[0,1]`
- `routing.low_confidence_behavior`: must be `clarify` or `block`
- `intake.max_input_bytes`: positive bounded integer
- `intake.max_decode_depth`: positive bounded integer
- `authorization.max_ttl_ms`: positive bounded integer
- `workflow.max_nodes`: positive bounded integer
- `workflow.max_replans`: non-negative bounded integer
- `skills.max_active_tokens`: positive bounded integer
- `presentation.max_nodes`: positive bounded integer
- `presentation.max_depth`: positive bounded integer
- `runtime.default_timeout_ms`: positive bounded integer
- `runtime.max_output_bytes`: positive bounded integer
- `runtime.allowed_classes`: non-empty enum
- `egress.default`: must be `deny`
- `egress.redirects`: must be `deny` unless explicitly required and re-authorized
- `credentials.max_ttl_ms`: positive bounded integer
- `connectors.default_status`: must be `blocked`
- `retrieval.max_top_k`: positive bounded integer
- `telemetry.store_hidden_chain_of_thought`: must be `false`
- `redaction.raw_protected_retention_days`: must be `0`
- `evaluation.judge.temperature`: must be `0`
- `evaluation.judge.require_json`: must be `true`
- `evaluation.judge.position_swap`: required for pairwise scoring
- `evaluation.minimum_confidence`: number in `[0,1]`
- `evaluation.pass_k`: positive values by risk class
- `secops.red_team.production`: must be `false`
- `secops.green_team.execute_recovery`: must be `false`
- `trust.score_range`: must be bounded
- `retention.encryption_at_rest`: must be `true`
- `governance.append_only`: must be `true`
- `governance.signed_records`: must be `true`

## 7. Agent Orchestration State Machine

The orchestrator is the only component that mutates short-lived plan state. Authorization,
execution, connector state, durable memory, presentation rendering, and eval scoring remain with
their owning layers.

### 7.1 Request Lifecycle States

1. `Received`
   - Request envelope exists but has not passed Layer 2.

2. `SafetyBlocked`
   - Input violated size, schema, prompt-injection, protected-context, or policy rules.

3. `ClarificationRequired`
   - Intent, mode, scope, target, or approval is ambiguous.

4. `IntakeAccepted`
   - Normalized request has a closed intent and sanitized context.

5. `Planning`
   - Layer 3 is selecting a registered profile, workflow, and bounded DAG.

6. `PlanRejected`
   - Workflow, DAG, dependency, compatibility, or budget validation failed.

7. `SkillLoading`
   - Layer 4 is validating requested skill bundles and resources.

8. `ReadyForAuthorization`
   - Plan and skill dependencies are structurally valid.

9. `AuthorizationBlocked`
   - Layer 2 denied a proposed action or required approval is missing.

10. `ReadyForExecution`
    - A signed, unexpired authorization decision exists.

11. `Executing`
    - Layer 6 is running an authorized tool in a bounded runtime.

12. `WaitingForData`
    - Layer 7 is performing a registered connector, retrieval, or memory operation on behalf of
      Layer 6.

13. `Presenting`
    - Layer 5 is building a signed presentation contract.

14. `Completed`
    - The requested bounded outcome completed.

15. `Failed`
    - A structured non-policy failure occurred.

16. `Canceled`
    - Cancellation propagated through active work.

17. `QuarantineDecided`
    - Layer 8 emitted a quarantine decision.

18. `Quarantined`
    - Layer 6 authenticated and executed the corresponding command.

Every transition emits a sanitized event. Raw input and hidden reasoning are never state-machine
fields.

### 7.2 Run Attempt Lifecycle

One plan attempt moves through:

1. `ValidateEnvelope`
2. `NormalizeInput`
3. `InspectSafety`
4. `ClassifyIntent`
5. `SelectWorkflow`
6. `BuildDAG`
7. `ValidateDependencies`
8. `LoadSkills`
9. `ProposeAction`
10. `AuthorizeAction`
11. `PrepareRuntime`
12. `ResolveCredential`
13. `ExecuteTool`
14. `AccessConnector` (optional)
15. `ValidateResult`
16. `BuildPresentation`
17. `RecordOutcome`
18. `Evaluate`
19. `Finish`

Each phase has an explicit timeout, cancellation path, structured error, and owning layer.

### 7.3 DAG Validation

A DAG validator MUST reject:

- empty version
- empty plan ID
- missing entry node
- unknown entry node
- duplicate node IDs
- empty action class
- missing dependency
- self-dependency
- cycle
- unreachable node
- node count over policy
- dependency fan-out over policy
- timeout over policy
- action class not declared by the selected workflow
- output contract not accepted by a downstream consumer

Validation algorithm:

1. Validate top-level schema.
2. Index nodes by ID and reject duplicates.
3. Verify entry node exists.
4. Verify all dependency references exist and are not self-references.
5. Run a color-marked depth-first search or Kahn topological sort to detect cycles.
6. Traverse from the entry node and reject unreachable nodes.
7. Validate each action class against workflow and profile policy.
8. Validate budgets and output contracts.
9. Produce an immutable topological plan.

### 7.4 Action Proposal and Authorization

For every executable node:

1. Layer 3 creates a proposed action without credentials.
2. Layer 2 validates intent, mode, tenant, resource, current source safety, approval metadata, tool
   registration, action class, and policy.
3. Layer 2 issues a signed, expiring authorization decision or a denial.
4. Layer 6 revalidates the decision before execution.
5. Layer 6 refuses a mismatched tenant, trace, tool, action class, resource, scope, or time window.

The model, workflow, skill, tool, connector, and runtime MUST NOT self-authorize.

### 7.5 Retry, Cancellation, and Replanning

Retry rules:

- Retry only errors classified as transient.
- Side-effecting retries require an idempotency key.
- Retry count and elapsed time are bounded.
- Backoff is bounded exponential with jitter or an equivalent policy.
- Authorization and credentials are revalidated before every retry.
- A schema, safety, authorization, tenancy, signature, checksum, provenance, or policy failure is
  not retried as a transient error.

Cancellation rules:

- A parent context or cancellation token propagates to planning, skill loading, runtime, connectors,
  model providers, telemetry export, and eval execution.
- Cancellation stops new child work.
- Runtime cleanup and credential revocation still occur.
- Hidden background work is forbidden.

Replanning rules:

- Replanning remains within the selected workflow and profile unless Layer 2 accepts a new request.
- Maximum replans are explicit.
- Replanning cannot broaden tool, connector, resource, credential, or side-effect scope.

### 7.6 Idempotency and Recovery

- Every side effect uses a tenant-scoped idempotency key.
- Runtime stores only the minimum state required to reject a duplicate.
- Connector adapters propagate supported idempotency keys.
- A repeated quarantine command is idempotent.
- A replayed UI event or A2A envelope is rejected.
- Recovery does not infer success from process exit; it verifies the declared output or side effect.
- Partial failure produces a structured recovery recommendation, not an unapproved compensating
  action.

### 7.7 Zero-Overlap Hard Stops

Layer 2 MUST NOT:

- choose workflows
- build DAGs
- load skill bodies
- execute tools
- call connectors
- render UI
- score evals

Layer 3 MUST NOT:

- classify raw intent
- authorize actions
- read `SKILL.md` directly
- execute tools
- connect external systems
- render UI
- score evals

Layer 4 MUST NOT:

- choose workflows
- authorize actions
- execute tools or scripts
- call connectors
- render UI
- score evals

Layer 5 MUST NOT:

- choose workflows
- authorize or execute actions
- call connectors
- persist memory
- score evals

Layer 6 MUST NOT:

- classify intent
- choose workflows
- authorize actions
- own connector registration
- render UI
- score evals

Layer 7 MUST NOT:

- authorize actions
- choose workflows
- load skills
- execute arbitrary code
- render UI
- score evals

Layer 8 MUST NOT:

- authorize actions
- choose workflows
- load runtime skill bodies
- execute tools or connectors
- mutate domain memory or presentation state
- execute quarantine
- auto-apply policy, workflow, skill, or tool changes

## 8. Intake, Routing, Workflows, and Skills

### 8.1 Input Normalization

Layer 2 normalization sequence:

1. Reject input larger than policy.
2. Validate the request envelope schema.
3. Reject invalid UTF encoding.
4. Normalize security-inspection text with Unicode NFKC.
5. Normalize line endings.
6. Bound nested object depth, key count, array length, and string length.
7. Normalize only fields whose contracts permit normalization.
8. Store the original input only through a protected, ephemeral reference when operationally
   necessary.

Normalization MUST NOT turn invalid input into valid executable syntax.

### 8.2 Prompt-Injection Firewall

The firewall inspects:

- direct instruction override attempts
- system, developer, assistant, tool, or role-marker impersonation
- requests to reveal prompts, hidden reasoning, credentials, policies, tools, paths, or traces
- split-token and punctuation-separated variants
- full-width and Unicode-confusable variants
- bounded base64 or equivalent encoded text
- selected text, retrieved text, connector content, metadata, attachments, and tool output
- instruction-like content inside otherwise valid domain data

Required behavior:

- Decode only allowlisted encodings.
- Bound decode depth and decoded bytes.
- Re-run inspection after each decode.
- Distinguish quoted or analyzed hostile content from instructions where the product supports that
  use case.
- Fail closed when confidence or context is insufficient for a high-risk action.
- Return a stable reason code, not the matched secret or full attack payload.

The firewall is one control. It does not replace action authorization, sandboxing, egress, or
output validation.

### 8.3 Intent Classification

Intent classification uses a closed enum defined by the domain pack.

Rules:

- Unknown intent becomes `clarification_required`, not an invented route.
- Low confidence becomes clarification or block according to policy.
- Multiple materially different intents require disambiguation or an explicitly supported
  multi-intent workflow.
- The classifier does not choose a workflow or tool.
- Deterministic rules MAY handle obvious cases before a model classifier.
- Model output is schema-validated and mapped to the closed enum.
- Classification evidence is categorical or hashed; raw prompts are not logged.

### 8.4 Approval Metadata

Approval evidence MUST bind:

- tenant
- subject
- trace
- action class
- target resource
- proposed change summary
- before and after digests when applicable
- scope
- expiry
- policy version
- nonce

Approval is invalid when stale, replayed, ambiguous, broader than the proposed action, or signed by
an unauthorized approver.

### 8.5 Outbound Protected-Metadata Redaction

Before user-facing, provider-facing, connector-facing, or telemetry output, scan for:

- system and developer prompts
- hidden chain-of-thought
- raw user prompts when prohibited
- workflow, profile, route, and trace identifiers
- skill and internal filesystem paths
- exact internal tool inventory
- MCP and private connector endpoints
- bearer tokens
- API keys
- OAuth tokens
- cookies
- passwords and client secrets
- private keys
- raw PII
- raw document or memory bodies

Redaction replaces the complete protected value, not only its label. Output is rescanned after
redaction. A redaction failure blocks the output.

### 8.6 Workflow Selection

Layer 3 selects only from registered workflows whose:

- status permits the current environment
- compatible intent contains the accepted intent
- compatible mode contains the accepted mode
- profile is compatible
- risk class is allowed
- required tools, connectors, and skills meet readiness requirements
- input contract matches available sanitized context
- budget fits current limits

When no workflow qualifies, return clarification or a structured unsupported-capability result.
Never silently choose a close but incompatible workflow.

### 8.7 Profile Selection

A profile is a versioned set of:

- allowed intents and modes
- allowed action classes
- allowed workflow IDs
- allowed skill IDs
- maximum risk
- runtime classes
- connector classes
- presentation surfaces
- budget limits
- approval requirements

Profiles narrow authority. They never expand Layer 2 policy.

### 8.8 Skill Registry Validation

At startup and in CI, Layer 4 MUST:

- compile registry, front matter, resource, and eval schemas
- reject empty registries
- reject duplicate IDs
- validate status and semantic version
- verify every `SKILL.md` checksum
- validate all required body sections
- enforce body and resource token limits
- verify non-empty owner and changelog
- verify resource and eval manifests
- verify every declared dependency exists
- reject an `active` skill without promotion evidence
- reject undeclared files that would load or execute

### 8.9 Progressive Disclosure

Skill loading occurs in three stages:

1. Registry metadata
   - ID, description, anti-triggers, tier, compatibility, dependencies, output contracts.

2. Skill body
   - Loaded only after a Layer 3 activation request passes Layer 4 validation.

3. Resources
   - Loaded individually when the skill procedure requires them.

The loader enforces per-resource and total token budgets. It returns a checksummed bundle and never
executes resources.

### 8.10 Skill Trigger Corpus

Layer 4 provides validated definitions; Layer 8 performs scoring.

Each case includes:

- globally unique case ID
- target skill
- group
- untrusted input
- expected skill or `none`
- forbidden skill when applicable
- expected intent where applicable
- safety-critical boolean

Provider-safe skill metadata contains:

- opaque run-specific candidate alias
- description
- anti-triggers

It excludes:

- file paths
- workflow internals
- tool endpoints
- credentials
- exact user-facing internal inventory

### 8.11 Tool and Workflow Dependency Gate

Before a skill execution eval:

1. Resolve `required_tools`.
2. Verify each tool implementation exists.
3. Verify unit, negative, boundary, and contract tests pass.
4. Resolve required connectors and runtime classes.
5. Verify integration evidence exists.
6. Resolve `required_workflows`.
7. Validate each workflow DAG and referenced dependency.
8. Verify a real execution path can be invoked in the evaluation environment.
9. Only then run golden-output, trajectory, pass^k, and end-to-end skill evals.

If any step is missing:

```text
eval status = blocked
reason = dependency_not_ready
```

The result is not `failed` and not `passed`. This distinction prevents false production claims.

### 8.12 Skill Promotion

An `experimental` skill may move to `canary` only when:

- schemas and checksums pass
- trigger, negative, rephrasing, collision, and out-of-scope cases pass
- required tools and workflows are integration-ready
- golden output and trajectory cases pass
- safety regressions pass with zero safety failures
- pass^k reliability meets policy
- live-provider evidence binds model and prompt versions
- no open critical finding exists
- a rollback plan exists
- governance evidence is signed
- a named human approver accepts the scope

An `active` skill requires successful canary evidence.

## 9. Presentation and Agent Interchange Safety

### 9.1 Presentation Schemas

Layer 5 compiles and validates versioned schemas for:

- presentation intent
- UI frame
- component catalog
- user event
- approval card
- surface patch
- agent card
- agent envelope

Unknown security-sensitive fields are rejected. All schemas compile before the process accepts work.

### 9.2 Locked Component Catalog

The catalog declares:

- component type
- allowed properties
- required properties
- property schemas
- allowed child types
- maximum children
- whether user input is accepted
- allowed event types
- accessibility requirements

The converter rejects:

- unknown component types
- unknown or missing properties
- invalid child relationships
- executable HTML, script, style, URL, command, or code content
- event types not declared by the catalog

### 9.3 UI Graph Validation

Frame validation requires:

- one root
- unique node IDs
- valid references
- no cycles
- bounded depth
- bounded node count
- no unreachable nodes
- schema-valid properties
- event references to declared events only

The graph validator is deterministic and does not render.

### 9.4 Hidden-Payload Scanner

The scanner recursively rejects:

- hidden or protected object keys
- script, iframe, object, embed, or active HTML
- inline event handlers
- CSS hiding or off-screen concealment
- HTML comments used to hide instructions
- instruction override content
- invisible and bidirectional control characters
- encoded payloads that decode into prohibited content
- excessive recursion, key count, or text size

### 9.5 Markdown and URL Safety

Markdown rules:

- Raw HTML is rejected unless an allowlisted sanitizer proves a required use case.
- HTML comments are rejected.
- Control characters are rejected.
- Unsafe link schemes are rejected.
- Embedded executable content is rejected.

URL rules:

- Parse with a standards-compliant URL parser.
- Allow only explicitly configured schemes.
- Reject userinfo.
- Reject empty host.
- Reject localhost and local suffixes.
- Reject private, loopback, link-local, multicast, unspecified, metadata, carrier-grade NAT, and
  documentation IP literals.
- Do not treat URL sanitization as egress authorization.

### 9.6 Signed Presentation Frames

Frames are signed over canonical content.

Validation binds:

- frame ID
- trace
- tenant
- surface
- schema version
- issue time
- expiry
- payload digest

Expired, altered, unsigned, or wrong-tenant frames are rejected.

### 9.7 User Event Normalization and Replay

User events include:

- event ID
- frame ID
- component ID
- event type
- normalized value or opaque value reference
- tenant
- trace
- issued time
- nonce
- signature

Layer 5 validates shape and normalizes the event. Layer 2 validates any approval or executable
meaning. Replay windows are bounded and nonces are single-use.

### 9.8 Agent Cards and A2A Envelopes

An Agent Card declares public, product-approved capabilities. It MUST NOT expose internal tools,
workflows, profiles, skill paths, endpoints, or traces.

A2A envelope validation requires:

- schema validation
- registered sender
- exact receiver
- Ed25519 or equivalent signature verification
- issue and expiry checks
- bounded clock skew
- nonce replay protection
- hidden-payload scan
- payload size and nesting limits

Agent messages never bypass intake, authorization, runtime, or connector policy.

### 9.9 Approval and Vibe-Diff Contract

For a high-stakes mutation, the approval surface SHOULD present:

- product-level action summary
- affected resource
- relevant before state
- proposed after state
- concise semantic difference
- risk and reversibility
- expiry
- approve and reject actions

It MUST NOT expose raw ASTs, hidden reasoning, credentials, internal identifiers, or unredacted
content beyond what the user is authorized to review.

### 9.10 Presentation Boundary

Layer 5 output is data. A separate client renders it. Layer 5 MUST NOT:

- execute JavaScript or shell content
- perform tool or connector calls
- authorize a user event
- persist memory
- choose the next workflow
- score its own correctness

## 10. Sandboxed Runtime and Execution Safety

### 10.1 Runtime Preflight

Before execution, Layer 6 validates:

- request schema and signature
- authorization signature and expiry
- exact tenant, trace, action, tool, resource, and scope binding
- tool registration and version
- runtime class
- sandbox image digest
- filesystem policy
- dependency policy
- egress policy
- credential request
- time, CPU, memory, process, file, network, and output budgets
- idempotency state
- current quarantine state

Any missing or mismatched requirement blocks execution.

### 10.2 Sandbox Driver Contract

The platform defines an interface; a deployment supplies a concrete driver.

The driver returns an attestation containing:

- sandbox ID
- runtime implementation and version
- kernel-isolation enabled
- rootless execution enabled
- seccomp or equivalent syscall policy enabled
- isolated network namespace
- cgroup or equivalent resource controls
- no ambient credentials
- read-only base filesystem
- writable ephemeral mount list
- image digest
- policy digest
- creation time
- expiry

The runtime destroys and rejects a sandbox when attestation is incomplete or does not match policy.

A source file named after a sandbox product is not an isolation implementation. Production requires
a real driver and independent escape-test evidence.

### 10.3 Ephemeral State Manager

Rules:

- Canonicalize the configured root.
- Create state only below that root.
- Create a cryptographically random workspace name.
- Mark the workspace with an ownership marker containing execution ID and version.
- Reject symlink or reparse-point escape.
- Mount only declared paths.
- Refuse cleanup when the target is the root, outside the root, unmarked, or owned by another
  execution.
- Clean state after completion, cancellation, timeout, or quarantine according to forensics policy.

### 10.4 Filesystem Policy

Default policy is deny.

Allowed operations are explicit:

- read-only paths
- writable ephemeral paths
- maximum files
- maximum total bytes
- maximum individual file bytes
- allowed extensions or media types
- symlink policy
- executable-bit policy

Absolute paths, traversal, device files, named pipes, sockets, and undeclared mounts are rejected
unless an explicit runtime class requires and authorizes them.

### 10.5 Ambient Credential Stripping

Before launching a process:

- Start from an empty environment or a minimal allowlist.
- Remove common cloud, VCS, package, database, model-provider, shell-history, proxy, and application
  credential variables.
- Reject malformed environment names.
- Reject protected names even when case or separator variants are used.
- Inject only a short-lived scoped credential reference required by the authorized action.
- Never mount a general host credential directory.

### 10.6 JIT Credential Broker

The broker receives an authorized credential request containing:

- tenant
- trace
- execution
- connector
- audience
- exact scopes
- maximum TTL
- purpose
- authorization decision

The external issuer returns an opaque lease. The broker:

- rejects scope broadening
- bounds TTL
- stores only the minimum lookup state
- resolves only an unexpired exact-match lease
- supports single lease, trace-wide, tenant-wide, and emergency revocation according to policy
- excludes token values from serialization, logging, telemetry, and errors

No runtime component may use ambient authority as a fallback.

### 10.7 Dependency and Supply-Chain Policy

Default production policy:

- no dynamic install
- exact versions
- immutable package source
- checksum or signature verification
- approved license and vulnerability state
- prebuilt dependency image where possible

The blocker rejects:

- package-manager invocation not explicitly authorized
- shell indirection around package managers
- download-and-execute patterns
- piping network output to a shell
- unpinned versions
- mutable tags
- missing or mismatched digest
- undeclared binary execution

### 10.8 Network Egress Controller

Default policy is deny.

For an allowed request:

1. Parse the URL.
2. Require an allowed scheme.
3. Match the exact registered domain policy.
4. Validate method and request headers.
5. Resolve every address.
6. Reject the request if any resolved address is non-public or prohibited.
7. Pin the checked address for the connection.
8. Preserve the original hostname for TLS identity.
9. Disable redirects or re-authorize every redirect as a new request.
10. Bound request bytes, response bytes, headers, timeout, and connection count.
11. Validate response content type and schema where declared.
12. Emit a sanitized egress event.

Prohibited addresses include:

- private
- loopback
- link-local
- multicast
- unspecified
- metadata-service ranges
- carrier-grade NAT
- benchmarking and documentation ranges
- policy-defined internal networks

This sequence prevents common SSRF and DNS-rebinding bypasses.

### 10.9 Process and Budget Controls

Each execution has explicit:

- wall-clock timeout
- CPU budget
- memory limit
- process count
- thread count
- open-file limit
- filesystem byte limit
- request and response byte limits
- network connection limit
- tool-call limit
- retry limit
- output token or text limit where applicable

Goroutines, threads, child processes, queues, and streams are bounded and cancellable.

### 10.10 Result Validation

Before returning a result:

- Verify process status and timeout state.
- Validate output against the registered schema.
- Enforce output size.
- Scan for secrets, protected metadata, unsafe URLs, and hidden payloads.
- Verify declared side effect where possible.
- Attach sandbox, tool, connector, and policy versions.
- Sign the result contract.
- Emit only a reference to large protected output.

### 10.11 Stateful Quarantine

Layer 6 accepts only a signed, unexpired, non-replayed quarantine command linked to a Layer 8
decision.

Allowed actions may include:

- block new executions for a trace, tenant, tool, workflow, connector, or artifact version
- freeze active sandboxes
- revoke credential leases
- disable egress routes
- disable tool or connector registrations through the owning control path
- preserve bounded forensic evidence

Execution is idempotent. Layer 6 returns a signed status. It does not decide whether quarantine was
warranted.

### 10.12 Runtime Error Mapping

Normalized categories:

- `authorization_invalid`
- `authorization_expired`
- `authorization_scope_mismatch`
- `artifact_quarantined`
- `sandbox_unavailable`
- `sandbox_attestation_invalid`
- `filesystem_denied`
- `dependency_denied`
- `dependency_integrity_failure`
- `credential_denied`
- `credential_expired`
- `credential_revoked`
- `egress_denied`
- `dns_policy_failure`
- `tls_identity_failure`
- `timeout`
- `budget_exceeded`
- `output_schema_invalid`
- `output_redaction_failure`
- `canceled`

Errors never include credentials, raw protected payloads, or hidden provider responses.

## 11. Integration and Data-Plane Contract

### 11.1 Connector Registration

Every connector is blocked by default. Activation requires:

- strict manifest validation
- owner and version
- transport identity
- authentication mode
- exact action classes
- source contracts
- tenant isolation mode
- runtime class
- timeout and byte limits
- redirect policy
- provenance policy
- taint policy
- schema hashes
- contract tests
- security review
- governance evidence

An endpoint string and API key do not constitute a connector implementation.

### 11.2 Source Contracts

A source contract declares:

- contract ID and version
- connector ID
- allowed request fields
- allowed response fields
- forbidden response fields
- required provenance fields
- maximum records
- maximum string length
- freshness requirements
- data classification
- redaction rules

The connector projects the provider response into this contract before returning it. Unknown fields
are dropped or cause rejection according to policy; they are never silently forwarded.

### 11.3 Tenant Boundary Enforcement

Every data-plane operation binds:

- authenticated tenant
- requested tenant
- credential tenant
- storage namespace
- trace
- connector
- source contract

All values must agree. Empty values fail closed. Tenant checks occur before and after any backend
query so a compromised or misconfigured backend cannot return cross-tenant data unnoticed.

### 11.4 mTLS Identity

For mTLS connectors:

1. Validate certificate time bounds.
2. Validate the full chain to a configured trust root.
3. Validate key usage.
4. Validate revocation according to deployment policy.
5. Require the exact expected SPIFFE ID or equivalent service identity.
6. Reject a valid certificate for the wrong connector.
7. Bind the verified identity to the connector manifest.

TLS encryption without exact peer identity is insufficient.

### 11.5 Token Audience and Scope

Connector credentials are accepted only when:

- issuer is trusted
- signature is valid
- audience exactly matches the connector
- tenant exactly matches
- connector claim exactly matches
- required scopes are present
- no undeclared scopes are used
- issue and expiry are within bounded skew
- token is not revoked

Credential contents are never logged or returned.

### 11.6 Field Projection

Projection rules:

- Input is treated as untrusted.
- Only explicitly allowlisted fields are copied.
- Nested maps and arrays are deep-copied.
- Forbidden and unknown fields are absent from the result.
- Size and depth limits apply before and after projection.
- Projection never returns references to mutable provider objects.

### 11.7 Provenance Attestation

An evidence payload is canonicalized and hashed. The provenance attestor signs or verifies:

- content hash
- connector ID and version
- source contract ID and version
- tenant
- retrieval time
- source identity
- freshness
- transformation chain

Unsigned, mismatched, stale, or unverifiable evidence is rejected for decisions requiring
provenance.

### 11.8 Taint Tracking

Taint labels may include:

- `external_untrusted`
- `user_supplied`
- `model_generated`
- `retrieved_unverified`
- `prompt_injection_suspected`
- `stale`
- `cross_tenant_suspected`
- `schema_drift`

Taint propagates through derived evidence. Sanitization may add a new label but does not erase the
source label. High-risk actions require policy-approved evidence classes.

### 11.9 Vector and Retrieval Isolation

Tenant vector namespaces are derived with a keyed function:

```text
namespace = HMAC(namespace_key, tenant_id || corpus_id || schema_version)
```

Retrieval validation:

- credential tenant matches request tenant
- namespace is derived, never supplied raw by a model
- embedding dimension is exact and bounded
- values are finite; NaN and infinity are rejected
- `top_k` is bounded
- filters are schema-valid and tenant-safe
- backend query uses the derived namespace
- every returned item is rechecked for namespace and tenant
- every item carries provenance and taint
- cancellation and timeout propagate

Cross-tenant results trigger a critical signal and quarantine decision.

### 11.10 MCP Transport

MCP or an equivalent tool protocol requires:

- strict JSON-RPC framing
- bounded message size
- unique request IDs
- closed method allowlist
- exact protocol version
- registered server identity
- advertised capability validation
- tool and resource allowlist
- canonical schema hash verification
- timeout and cancellation
- bounded circuit breaker
- sanitized errors

Unknown methods, capabilities, tools, resources, or schema hashes fail closed.

### 11.11 MCP Handshake

Handshake validation binds:

- connector ID
- server identity
- protocol version
- transport
- capabilities
- tool IDs
- resource IDs
- prompt IDs when allowed
- schema hashes
- policy version

The server MUST NOT add capabilities after authorization without a new handshake and policy
decision.

### 11.12 Stdio Transport

Layer 7 MUST NOT spawn an arbitrary local process.

For stdio MCP:

- Layer 6 starts the registered process in an authorized sandbox.
- Layer 6 returns a bounded pipe abstraction.
- Layer 7 performs protocol framing over that pipe.
- Process, filesystem, dependency, credential, and egress controls remain Layer 6 responsibilities.

### 11.13 Streamable HTTP Transport

HTTP transport requires:

- HTTPS
- mTLS where the manifest requires it
- exact host identity
- no automatic redirects
- bounded request and response bodies
- bounded headers
- explicit content type
- context cancellation
- timeout
- connection-pool limits
- certificate and schema-drift tests

### 11.14 Circuit Breaker

The breaker has explicit states:

- `closed`
- `open`
- `half_open`

It tracks sanitized failure categories, not raw provider bodies.

Rules:

- Open after a bounded threshold.
- Reject while open.
- Allow a bounded probe count after cooldown.
- Close only after configured successes.
- Reopen on probe failure.
- Keep breaker state tenant-safe where failures can differ by tenant.

### 11.15 Model Provider Adapter

Evaluation and optional agent model access use a Layer 7 adapter. The adapter receives a
provider-neutral request from the owning layer and returns bounded structured output.

Mandatory controls:

- credential loaded only from an approved environment or secret source
- credential sent only in the provider's required authorization header
- fixed or strictly allowlisted endpoint
- redirects disabled
- bounded timeout
- bounded response bytes
- no request or response logging
- no persistent provider conversation unless explicitly required
- provider-side storage disabled where supported
- background execution disabled unless explicitly governed
- temperature zero for deterministic evaluation
- fixed seed where supported
- reasoning summaries disabled where supported
- closed structured-output schema
- provider errors reduced to status and safe category
- output independently revalidated by Layer 8

Provider-specific example:

- A Gemini Interactions API adapter reads `GEMINI_API_KEY` from the environment.
- It sends the key only in the `x-goog-api-key` header.
- It uses the configured Google Interactions endpoint.
- It sets `store=false`, `background=false`, `temperature=0`, a fixed seed, and no thinking
  summaries.
- It requests `application/json` with a closed response schema.

A key pasted into chat, source, command history, URL, report, or log is compromised and MUST be
rotated before use.

### 11.16 Memory Data Plane

Memory and preference storage are Layer 7 operations.

Requirements:

- Layer 2 validates user intent and approval.
- Layer 3 plans the update.
- Layer 5 presents the proposed semantic change.
- Layer 7 writes only the approved tenant-scoped fields.
- Layer 8 records hashes and safe summaries.
- Raw memory bodies are not logged.
- Reads return projected, bounded context.
- Deletion, export, retention, and legal hold are supported where required.

### 11.17 Integration Error Mapping

Normalized categories:

- `connector_unregistered`
- `connector_inactive`
- `source_contract_invalid`
- `tenant_mismatch`
- `identity_invalid`
- `audience_invalid`
- `scope_invalid`
- `schema_hash_mismatch`
- `protocol_invalid`
- `provenance_invalid`
- `taint_blocked`
- `retrieval_invalid`
- `cross_tenant_result`
- `circuit_open`
- `provider_unavailable`
- `provider_response_invalid`
- `response_limit_exceeded`
- `canceled`

## 12. Context Assembly and Agent Execution

### 12.1 Context Inputs

Allowed context sources:

- sanitized intake decision
- tenant and trace references
- accepted intent and mode
- selected profile
- validated workflow plan
- validated skill bundle
- registered tool input schemas
- projected evidence packets
- approved memory or preference summaries
- previous structured results within the same authorized trace
- presentation requirements

Context sources are ordered, versioned, and checksummed.

### 12.2 Forbidden Context

The active model context MUST NOT include:

- ambient credentials
- raw system or developer prompts from another component
- hidden chain-of-thought
- another tenant's data
- raw unapproved memory
- raw connector responses outside source contracts
- exact internal tool inventory beyond the current bounded candidates
- arbitrary skill or workflow paths
- MCP endpoints
- unbounded logs or traces
- unscanned retrieved instructions
- stale approval evidence

### 12.3 Context Budgeting

The assembler allocates explicit budgets for:

- platform instruction
- domain policy
- workflow
- active skill bodies
- optional skill resources
- tool schemas
- evidence
- memory
- conversation state
- output reserve

Rules:

- Budget overflow fails or triggers a declared summarization path.
- Security policy and authorization are never truncated.
- Summaries preserve provenance and taint.
- Progressive disclosure is preferred over loading all skills.
- A bounded candidate set SHOULD contain 5 to 15 skills for model-based routing.

### 12.4 Tool Discovery

The model sees only tools that are:

- registered
- compatible with the selected profile and workflow
- implemented
- allowed in the current environment
- representable within the context budget

Tool descriptions expose product-level purpose and schemas. They do not expose credentials,
endpoints, implementation paths, or unauthorized tools.

### 12.5 Tool Call Lifecycle

1. Model or deterministic workflow proposes a structured tool input.
2. Layer 3 binds it to a DAG node and action class.
3. Layer 2 validates and authorizes.
4. Layer 6 executes under runtime policy.
5. Layer 7 performs external data access when required.
6. Layer 6 validates the result.
7. Layer 3 consumes the structured result.
8. Layer 5 presents user-facing data when needed.
9. Layer 8 observes and evaluates.

Malformed, unknown, unauthorized, or unavailable tools return a structured failure. They do not
fall back to arbitrary code or a nearby tool.

### 12.6 Workflow Execution

For each topologically ready node:

- verify dependencies completed successfully
- revalidate cancellation and plan expiry
- resolve required skill bundle
- construct the proposed action
- authorize
- execute
- validate the output contract
- persist only owning-layer state
- emit a sanitized node outcome

Parallel nodes MAY run only when:

- dependencies permit it
- concurrency is bounded
- tenant and resource locks are safe
- cancellation is shared
- aggregate budgets are enforced

### 12.7 Agent Turn Contract

An implementation may use a local model, remote API, app server, or another agent protocol.

The logical contract is:

1. Initialize a bounded session.
2. Provide sanitized context.
3. Advertise the bounded tool surface.
4. Receive structured model output.
5. Reject unsupported tool or user-input requests according to policy.
6. Execute tool calls only through the required path.
7. Continue until an explicit completion, clarification, failure, cancellation, timeout, or turn
   limit.
8. Close provider and runtime resources.

Exact wire fields are provider-specific; ownership and safety semantics are not.

### 12.8 Multi-Agent and A2A Handoffs

An agent handoff requires:

- registered sender and receiver
- signed A2A envelope
- declared capability intersection
- tenant and trace binding
- bounded context references
- expiry and nonce
- explicit expected output contract

The receiving agent re-enters through Layer 2. Sender authority is not inherited.

### 12.9 Human Input and Approval

When input is required:

- Suspend or terminate the bounded attempt according to policy.
- Present a safe request through Layer 5.
- Do not leave an execution or credential lease waiting indefinitely.
- Resume only with a new validated event and, where required, new authorization.

Approval and clarification are different:

- Clarification resolves meaning or scope.
- Approval authorizes a specific high-stakes proposal.

### 12.10 Result Construction

The final product response is built from:

- validated structured outputs
- projected evidence
- safe error summaries
- presentation contracts

It MUST NOT expose:

- raw traces
- raw AgBOM
- hidden reasoning
- internal routes, profiles, workflow IDs, tool IDs, skill paths, or connector endpoints unless the
  product explicitly permits a safe public name
- credentials
- raw protected data

### 12.11 Missing Capability Semantics

When a requested tool, workflow, connector, or skill implementation is absent:

- Return `capability_unavailable` or `dependency_not_ready`.
- Identify the missing capability at an operator-safe logical level.
- Do not fabricate an answer from a model.
- Do not silently downgrade to a side-effecting alternative.
- Do not record a successful execution eval.
- Keep dependent skills and workflows blocked.

## 13. Logging, Observability, Evaluation, and Governance

### 13.1 Logging Conventions

Allowed common fields:

- hashed tenant reference
- trace ID or hashed trace reference according to exposure policy
- request ID
- action class
- artifact version hashes
- status
- reason code
- latency
- bounded resource usage

Logging requirements:

- Structured output.
- Stable field names and reason codes.
- No secrets or protected payloads.
- No hidden chain-of-thought.
- No raw user, selected, canvas, brief, memory, or connector bodies.
- No exact internal inventory in user-facing logs.
- Errors are sanitized before logging.

### 13.2 Event Redaction

The event redactor:

- recursively traverses bounded objects
- drops prohibited keys
- rejects excessive depth or key count
- detects credential patterns in values
- HMACs configured correlation fields
- stores only approved categorical values and safe summaries
- does not mutate the source object
- fails closed when safe redaction cannot be guaranteed

Redaction classes:

- `never_store`
- `store_as_hmac`
- `store_as_redacted_summary`
- `store_as_structured_metric`

### 13.3 Trace Contract

Every production trace has:

- required root span
- required child spans for executed phases
- tenant, intent, mode, and profile references
- parent-child links
- action, policy, runtime, connector, presentation, and evaluation outcomes
- cost, token, and latency metrics
- final status

OpenTelemetry implementation rules:

- Use an injected provider.
- Do not mutate global provider state from a library.
- Allowlist span names and attributes.
- Reject protected attributes.
- Link child spans to the correct parent.
- Honor cancellation.
- Use tail-based retention for errors, policy blocks, quarantine decisions, eval failures, and
  high-cost sessions.

### 13.4 Runtime Agent Bill of Materials

AgBOM is generated for every production trace and includes exact versions and hashes of everything
that influenced the run.

Rules:

- Missing required AgBOM components make the trace incomplete.
- AgBOM is immutable after trace finalization.
- Raw AgBOM is operator-restricted.
- User-facing views use product-safe summaries.
- Evaluation evidence references the AgBOM hash.

### 13.5 Drift, Trust, and Loop Detection

Intent drift compares categorical or hashed representations of:

- accepted intent
- selected workflow
- proposed action classes
- executed action classes
- connector classes
- resulting presentation

Trust starts at a bounded initial score and decays on:

- authorization block
- tool or connector expansion
- schema drift
- provenance failure
- prompt-injection signal
- approval misuse
- cross-tenant signal
- repeated failed action
- unexplained cost spike

Recovery is bounded per trace. A critical signal can require a quarantine decision regardless of
the numeric score.

Loop detection tracks:

- repeated DAG node
- repeated tool with equivalent parameters
- repeated provider failure
- replan count
- unchanged-result retries

Threshold breaches stop or block further work according to policy.

### 13.6 SecOps Triad

Red Team:

- Runs synthetic adversarial fixtures.
- Covers injection, secret extraction, scope expansion, SSRF, path traversal, dependency abuse,
  tenant crossing, replay, approval bypass, schema drift, and hidden payloads.
- Runs in CI, shadow, or staging.
- MUST NOT launch attacks against production user traffic.

Blue Team:

- Evaluates sanitized live traces.
- Detects action, tool, connector, cost, loop, approval, provenance, and tenant anomalies.
- Produces bounded findings and evidence references.

Green Team:

- Produces recovery, test, policy, workflow, skill, or tool recommendations.
- Never edits or deploys them automatically.
- Routes quarantine decisions to Layer 6.
- Requires human review for control changes.

### 13.7 Trajectory Evaluation

Supported modes:

- `EXACT`
  - Actual action-class sequence and multiplicity must exactly match.

- `IN_ORDER`
  - Required action classes occur in order; allowed intermediate actions may exist.

- `ANY_ORDER`
  - Required action classes occur in any order.

All modes enforce:

- forbidden actions absent
- required multiplicity
- maximum action count
- no unknown action class
- authorization before execution
- connector access through Layer 7
- cancellation and failure-path expectations

Trajectory evals use action classes, not provider-specific tool-call IDs.

### 13.8 Pass^k Reliability

For a target and exact version:

1. Execute the same versioned case `k` times.
2. Bound concurrency and total elapsed time.
3. Honor cancellation.
4. Validate every run independently.
5. Compute pass count and flake rate.
6. Treat any safety-critical failure as blocking.
7. Require all runs when policy defines pass^k as all-pass.

Recommended minimums:

- smoke: `1`
- read workflow: `3`
- draft workflow: `3`
- guarded write workflow: `5`
- high-impact workflow: `10`

Domains MAY require higher values.

### 13.9 LLM-as-Judge

A judge is supplemental, never sole security evidence.

Judge input:

- redacted eval ID
- versioned instruction
- candidate output or pairwise candidates
- versioned rubric
- exact model
- prompt version
- temperature zero

Judge output is strict JSON:

- rubric ID
- one integer score per declared dimension
- pairwise winner when applicable
- confidence
- short evidence-based summary

Validation:

- Reject unknown fields and trailing data.
- Reject hidden reasoning fields.
- Reject missing or extra dimensions.
- Enforce dimension bounds.
- Recompute weighted score locally.
- Require minimum confidence.
- Require human review below the review threshold.
- Use position swap for pairwise scoring.
- Mark inconsistent winners or materially different scores unstable.
- Bound retries.

The provider does not determine the final weighted score or release status.

### 13.10 Skill Trigger Evaluation

For each target skill:

1. Load the exact validated corpus.
2. Select 5 to 15 bounded candidates containing the target, expected adjacent skills, and
   deterministic distractors.
3. Map real skill IDs to run-specific opaque aliases.
4. Rotate candidate order on every run.
5. Keep expected answers local.
6. Ask the Layer 7 provider adapter for strict structured selection.
7. Require every case exactly once.
8. Validate alias, confidence, summary size, and closed schema.
9. Map aliases back locally.
10. Score positive, negative, rephrase, collision, and safety cases.
11. Detect selection instability across pass^k runs.
12. Emit a sanitized report containing case hashes, not prompts or provider summaries.

The report includes:

- corpus checksum
- target skill and version
- model and prompt version
- runs
- total, passed, and failed cases
- safety failures
- flaky cases and flake rate
- minimum confidence

Trigger evaluation proves routing behavior only.

### 13.11 Tool, Workflow, Skill, and End-to-End Evaluation

Tool eval:

- real implementation invoked
- input/output/error schemas
- success, failure, and boundary cases
- authorization mismatch
- cancellation and timeout
- idempotency
- credential scope
- filesystem and egress denial
- redaction

Connector eval:

- real adapter or certified sandbox endpoint
- identity and mTLS
- audience and scope
- tenant isolation
- source projection
- provenance
- taint
- schema drift
- circuit breaker
- chaos and timeout

Workflow eval:

- real DAG
- real dependencies
- expected action-class trajectory
- error and compensation paths
- retry and cancellation
- presentation outcome

Skill execution eval:

- real skill bundle
- real workflow
- real tools and connectors
- golden output
- forbidden output qualities
- trajectory
- safety regression
- pass^k
- live judge where policy requires it

End-to-end eval:

```text
untrusted request
  -> intake
  -> routing
  -> skill loading
  -> authorization
  -> runtime
  -> connector or data plane
  -> output validation
  -> presentation
  -> observability
  -> evaluation
```

If any dependency is mocked, the evidence is labeled `mocked` and cannot satisfy a live
end-to-end production gate.

### 13.12 Evaluation Evidence Integrity

Evidence records:

- exact target versions
- exact corpus checksum
- exact policy versions
- exact model and prompt version
- exact tool, workflow, skill, connector, schema, runtime, and image versions
- environment class
- whether each dependency was real, simulated, or mocked
- run counts and outcomes
- safety failures
- flake rate
- start and end times
- signer

Rules:

- Raw prompts and protected outputs are stored only if an explicitly approved encrypted evaluation
  vault permits them; the default is no raw retention.
- Release evidence contains hashes and safe summaries.
- Evidence is signed.
- Changed artifacts invalidate stale evidence.
- Missing evidence is `not_run`.

### 13.13 Immutable Governance Ledger

The ledger is:

- append-only
- hash-chained
- signed
- timestamped
- tenant-scoped where applicable
- defensively copied on read
- independently verifiable

It records:

- artifact promotion and demotion
- policy changes
- named approvals
- eval evidence acceptance
- canary start and completion
- rollback
- quarantine decisions and results
- incident findings
- retention and deletion actions

Tampering with any record or link invalidates the chain.

### 13.14 Correction Mining and Improvement

Correction mining consumes sanitized labels and hashes, not raw conversations.

A recommendation requires:

- minimum occurrence count
- minimum distinct tenant count where multi-tenant privacy applies
- stable failure classification
- evidence references
- owning layer
- proposed test
- proposed control change
- risk and rollback summary

Outputs are proposals only. No recommendation auto-modifies production.

### 13.15 Retention and Privacy

Required controls:

- encryption at rest
- retention classes
- raw protected retention of zero by default
- tenant deletion
- legal hold where required
- immutable governance records
- deletion audit
- restricted trace and AgBOM access
- safe export

Deletion never breaks governance-chain verification; it removes protected payloads and retains only
permitted hashes and records.

### 13.16 Monitoring Interface

A synchronous operator snapshot SHOULD include:

- active traces
- policy blocks
- authorization denials
- tool and connector failures
- quarantine state
- cost and latency
- token usage
- drift and trust distributions
- loop signals
- eval status by artifact
- artifact status
- missing dependencies
- canary state

The snapshot is data only and contains no raw protected payloads.

### 13.17 Test and Conformance Matrix

Every implementation runs:

1. `Formatting and build`
   - formatter
   - compile
   - dependency lock verification

2. `Unit tests`
   - success
   - failure
   - boundary
   - table-driven or equivalent systematic cases

3. `Schema and configuration tests`
   - valid documents
   - unknown fields
   - empty documents
   - placeholder documents
   - semantic cross-field failures

4. `Security tests`
   - injection variants
   - secret patterns
   - hidden payloads
   - replay
   - signature tampering
   - tenant mismatch
   - SSRF and DNS rebinding
   - path traversal and symlink escape
   - dependency abuse

5. `Concurrency tests`
   - race detector or equivalent
   - cancellation
   - bounded queues
   - idempotency
   - duplicate command handling

6. `Static analysis`
   - language vet or compiler lints
   - unreachable and unused production code
   - security analyzers
   - dependency vulnerability scan

7. `Integration tests`
   - sandbox driver
   - credential issuer
   - connectors
   - MCP transports
   - retrieval backend
   - telemetry exporter
   - governance store

8. `Evaluation`
   - tool
   - workflow
   - skill trigger
   - skill execution
   - trajectory
   - pass^k
   - judge
   - adversarial
   - end-to-end

9. `Deployment tests`
   - chaos
   - load
   - canary
   - rollback
   - incident and quarantine drill

The test report distinguishes unit, mock integration, certified integration, live provider, canary,
and production evidence.

### 13.18 Repository Readiness Gate

The read-only gate MUST:

- enumerate required artifacts
- validate schemas and policies
- reject empty and placeholder content
- verify checksums, owners, changelogs, manifests, and fixtures
- verify executable control sources or packages required by the implementation profile
- report missing dependencies
- report non-production statuses
- return machine-readable findings
- exit non-zero while a blocker exists

It MUST NOT:

- execute tools
- call connectors
- start sandboxes
- read secrets
- run live model evals
- infer eval success from files
- mutate artifacts
- promote status

Repository gate pass is necessary and insufficient for production.

### 13.19 Production Readiness State

The release state is calculated from explicit evidence:

```text
repository_ready
AND tools_ready
AND connectors_ready
AND workflows_ready
AND skills_ready
AND runtime_attested
AND security_tests_passed
AND live_evals_passed
AND canary_passed
AND governance_approved
```

If any term is absent or false, production readiness is false.

Required release evidence:

- all required schemas and policies valid
- no placeholder control artifacts
- all required tool implementations tested
- all required connectors certified
- all required workflows executed
- all production skills evaluated against real dependencies
- zero open critical security findings
- race and concurrency tests passed
- static and vulnerability analysis passed
- sandbox and credential infrastructure attested
- tenant isolation tested
- pass^k and adversarial suites passed
- canary completed
- rollback plan tested
- named approver recorded

### 13.20 Implementation Order for Coding Agents

An implementation agent SHOULD work in this order:

1. Read repository instructions and language style guides.
2. Define domain intents, modes, risk classes, action classes, and public output contracts.
3. Create strict schemas and fail-closed policy parsers.
4. Implement Layer 2 normalization, firewall, classification, authorization, and redaction.
5. Implement Layer 3 DAG validation and orchestration without execution.
6. Implement the tool registry and real tool implementations.
7. Implement Layer 6 runtime boundaries required by those tools.
8. Implement Layer 7 connectors and tenant-safe evidence contracts.
9. Write and integration-test real workflows using those tools and connectors.
10. Create Layer 4 skills that reference only existing workflows and tools.
11. Implement Layer 5 presentation and A2A contracts.
12. Implement Layer 8 telemetry, AgBOM, drift, trust, SecOps, eval, and governance.
13. Add startup validation for all embedded or deployed contracts.
14. Add the read-only readiness aggregator.
15. Run static trigger evals.
16. Run real tool, connector, workflow, and skill execution evals.
17. Run live-provider, adversarial, pass^k, and end-to-end suites.
18. Run canary and rollback.
19. Record human approval.
20. Promote only the exact evaluated versions.

Stop rules:

- Stop when ownership is unclear.
- Stop when a required spec or style guide is missing.
- Stop when a dependency is not implemented.
- Stop when a secret is exposed.
- Stop when a live test would exceed approved budget or egress.
- Stop when an action requires new authority.
- Report the exact blocker; never replace it with a stub pass.

### 13.21 Acceptance Criteria

This specification is implemented when:

1. Layer ownership is enforced in code and tests.
2. Every executable action follows proposal, authorization, runtime, data-plane, and observation
   paths.
3. Input, output, events, and artifacts are schema-bound and size-bounded.
4. Prompt injection and protected metadata controls are executable.
5. Workflows are bounded, acyclic, reachable, and dependency-valid.
6. Skills are checksummed, progressively disclosed, dependency-truthful, and status-gated.
7. Presentation and A2A payloads are non-executable, sanitized, signed, and replay-safe.
8. Runtime execution is sandboxed, ephemeral, credential-scoped, egress-controlled, and
   supply-chain constrained.
9. Connectors enforce identity, audience, tenant isolation, projection, provenance, and taint.
10. Traces are redacted, linked, and accompanied by an AgBOM.
11. Drift, trust, loops, anomalies, and quarantine decisions are calculated from safe evidence.
12. Tool, workflow, skill, trajectory, pass^k, judge, adversarial, and end-to-end evals exist.
13. Missing real tools or workflows block dependent execution evals.
14. Governance evidence is signed and append-only.
15. Repository, infrastructure, live evaluation, canary, and human evidence are reported
    separately.
16. No component claims production readiness from compilation, file presence, mocks, or empty
    artifacts.

### 13.22 Final Non-Goals and Boundaries

```text
The platform does not invent domain tools.
The platform does not infer workflows from skill prose.
The platform does not let a model authorize itself.
The platform does not let skills execute directly.
The platform does not let presentation mutate state.
The platform does not let runtime choose policy.
The platform does not let connectors bypass tenancy.
The platform does not let observability execute recovery.
The platform does not store hidden chain-of-thought.
The platform does not expose secrets or protected internals.
The platform does not count blocked or not-run evals as passes.
The platform does not equate repository readiness with production readiness.
```

One-line architecture summary:

Agentic Engineering is a contract-first, fail-closed, domain-neutral platform in which intake
authorizes nothing implicitly, orchestration coordinates without executing, skills provide
checksummed procedural memory, presentation remains non-executable, runtime enforces isolation,
integrations preserve identity and tenancy, and glass-box evaluation governs promotion using real
dependency-backed evidence.
