# Layer 7 Production-Ready Specs v2
## Grounding Interop, MCP Transport, Connector Contracts, Tenant-Isolated Retrieval, and Data Security Mesh

**Status:** Production-ready upgrade  
**Supersedes:** Layer 7 v1 where MCP transport, data connectors, encryption, and vector retrieval were sketched but not operationally specified  
**Primary rule:** Layer 7 is the grounding and interoperability data plane. It connects already-authorized agent actions to verified external/internal resources through typed MCP/A2A-compatible transport, schema-locked connector contracts, tenant-isolated retrieval, provenance-preserving evidence packets, and encrypted data plumbing. It never classifies user intent, chooses workflows, loads skills, renders UI, executes arbitrary code, grants policy approval, mints credentials, owns quarantine, stores telemetry as the source of truth, or evaluates agent quality.

---

## 1. Upgrade Decision

Upgrade is required.

Layer 7 v1 correctly identified the right macro responsibilities:

```text
MCP transport socket plane
stdio JSON-RPC piping
remote Streamable HTTP / SSE transport
enterprise MCP server connectors
mTLS for A2A and MCP traffic
CMEK-backed encryption
vector retrieval
tenant partitioning
```

But v1 is not yet production-ready because it does not define:

```text
connector trust lifecycle
MCP server identity verification
schema hash pinning
tool spoofing defense
transport-level fail-closed behavior
tenant-scoped retrieval contracts
source provenance packets
data classification propagation
connector result tainting
credential consumption boundaries
cross-layer event contracts
connection health SLOs
registry governance
CI/CD gates for connector manifests
safe degradation modes
```

Layer 7 v2 upgrades this layer from "resource plumbing" into a hardened enterprise data plane that can safely ground an SEO/AEO auditor and content-writing agent without turning the agent into a web of bespoke, ungoverned API wrappers.

---

## 2. Layer 7 Mission

Layer 7 exists to answer one question:

```text
How does a validated agent action safely reach grounded data, tools, and memory without bespoke wrappers, cross-tenant leakage, connector spoofing, or context poisoning?
```

It does this by owning:

```text
1. MCP transport clients
2. MCP server registry and manifest validation
3. Connector-local schema contracts
4. Source-specific adapters
5. Tenant-isolated retrieval
6. Encrypted data plumbing
7. Provenance and taint metadata on every returned payload
8. Connector health and degradation state
9. Sanitized interop events emitted to Layer 8
```

Layer 7 does not decide whether the agent is allowed to call a tool. That decision belongs to Layer 2.  
Layer 7 does not execute untrusted code or own process isolation. That belongs to Layer 6.  
Layer 7 does not decide which workflow or agent should use the connector. That belongs to Layer 3.  
Layer 7 does not teach the agent how to use the connector procedurally. That belongs to Layer 4.  
Layer 7 does not render connector results to the user. That belongs to Layer 5.  
Layer 7 does not evaluate result quality or trigger recovery. That belongs to Layer 8.

Layer 7 receives a narrow, already-authorized, already-sandboxed, credential-scoped request envelope and returns a typed, provenance-rich, tenant-safe result envelope.

---

## 3. Non-Negotiable Production Rules

### 3.1 Fail Closed

Layer 7 must fail closed for:

```text
unknown connector_id
unknown tool_id
unverified MCP server identity
schema hash mismatch
unsupported MCP protocol version
missing tenant_id
missing request_id
missing policy_decision_id
missing sandbox_execution_id for runtime-originating calls
missing jit_credential_ref for authenticated calls
wrong token audience
wrong tenant namespace
cross-tenant vector namespace access
unbounded query
unbounded pagination
oversized response
malformed JSON-RPC
unexpected result schema
tainted result returned without taint metadata
remote server returns tool list different from pinned manifest
remote server changes method signature without registry approval
remote connector asks for interactive credential flow
connector attempts direct egress outside Layer 6 egress constraints
```

Layer 7 must never "try anyway" because connector failure is safer than silent exfiltration, schema drift, or hallucinated grounding.

### 3.2 Data Must Remain Data

Every external or retrieved payload must be wrapped as data, never injected as instruction.

Layer 7 must label payloads with:

```text
source_type
source_id
tenant_id
data_classification
provenance_hash
retrieval_timestamp
freshness_ttl
taint_flags
license_or_usage_constraints
confidence
source_url_or_resource_ref when safe
field_projection_applied
```

If a web page, CMS note, repository file, review, search result, or vector memory contains prompt-like text, Layer 7 must not interpret it. It must return it as quoted evidence with taint metadata for downstream policy and evaluation layers.

### 3.3 No Ambient Data Access

Layer 7 may only use scoped credentials passed by Layer 6. It must not:

```text
read environment secrets
read developer machine tokens
use human user OAuth sessions directly
store long-lived API keys
refresh credentials
mint credentials
broaden credentials
borrow another tenant credential
```

Layer 6 owns JIT token issuance and revocation. Layer 7 validates token audience, scope, tenant, and expiry, then consumes the token for the specific connector call.

### 3.4 No Bespoke Connector Drift

Every connector must be described by a versioned manifest. No connector may be added by "just writing code" and wiring it directly into the agent.

A production connector requires:

```text
manifest
schema
transport contract
authentication mode
tenant boundary
field projection rules
rate limits
pagination bounds
result size limits
provenance rules
taint behavior
golden contract tests
security tests
observability event mapping
rollback plan
```

### 3.5 No Cross-Layer Shortcuts

Layer 7 must not call Layer 3 to ask which workflow is active.  
Layer 7 must not call Layer 4 to inspect SKILL.md.  
Layer 7 must not call Layer 5 to create UI.  
Layer 7 must not call Layer 8 to decide whether a connector result is "good."  
Layer 7 must not call Layer 6 to bypass sandbox or egress policy.  
Layer 7 must not call Layer 2 to reinterpret intent.

Cross-layer communication happens through typed request envelopes, sanitized events, and explicit contracts only.

---

## 4. Zero-Overlap Boundary Contract

| Capability | Owning Layer | Layer 7 Responsibility | Layer 7 Must Not Do |
|---|---:|---|---|
| User input, "vibe diff", MFA UI | Layer 1 | None | Render approvals or collect MFA |
| Prompt injection blocking, intent enum, policy allow/deny | Layer 2 | Verify `policy_decision_id` exists and matches requested connector/tool | Classify intent, approve tools, redact final user output |
| Workflow, DAG routing, agent assignment | Layer 3 | Accept structured call envelopes from DAG nodes | Choose workflow, assign agent, mutate orchestration state |
| Procedural skills and SKILL.md loading | Layer 4 | Expose typed resources that skills may call through the orchestrator | Load skills, evaluate skill triggers, modify SKILL.md |
| A2UI rendering and A2A endpoint exposure | Layer 5 | Provide safe resource payloads for UI conversion and support mTLS data plane | Render cards, define UI schema, expose public agent card |
| Sandbox, filesystem, egress proxy, JIT token broker | Layer 6 | Consume scoped runtime references, attach to sanctioned stdio pipes, use approved egress route | Spawn unsandboxed processes, mint credentials, open arbitrary sockets |
| Grounding interop and resource plumbing | Layer 7 | Own MCP transport, connector registry, data mesh, vector retrieval, provenance packets | N/A |
| Observability, SecOps, quality loops | Layer 8 | Emit sanitized interop events and health facts | Store authoritative traces, evaluate trajectory, trigger quarantine |

---

## 5. Production Directory Layout

```text
/layer_07_interop
│
├── /config
│   ├── connector-registry.yaml
│   ├── mcp-server-manifest.schema.json
│   ├── source-contracts.yaml
│   ├── data-classification-map.yaml
│   ├── endpoint-bindings.yaml
│   ├── vector-namespaces.yaml
│   ├── field-projection-rules.yaml
│   ├── response-budget.yaml
│   ├── provenance-policy.yaml
│   └── connector-slo.yaml
│
├── /mcp_transport_plane
│   ├── jsonrpc_codec.go
│   ├── request_envelope_validator.go
│   ├── stdio_client.go
│   ├── streamable_http_client.go
│   ├── sse_stream_reader.go
│   ├── protocol_version_negotiator.go
│   ├── handshake_validator.go
│   ├── schema_hash_verifier.go
│   ├── mcp_server_identity_verifier.go
│   ├── connection_pool.go
│   ├── retry_backoff_controller.go
│   ├── circuit_breaker.go
│   ├── transport_timeout_guard.go
│   ├── payload_size_guard.go
│   └── transport_event_emitter.go
│
├── /connector_registry
│   ├── registry_loader.go
│   ├── manifest_validator.go
│   ├── connector_capability_index.go
│   ├── connector_status_store.go
│   ├── version_pin_resolver.go
│   ├── schema_migration_guard.go
│   ├── registry_signature_verifier.go
│   └── public_registry_quarantine.go
│
├── /enterprise_mcp_servers
│   ├── search_console_mcp_server.go
│   ├── analytics_mcp_server.go
│   ├── repository_mcp_server.go
│   ├── cms_mcp_server.go
│   ├── review_intelligence_mcp_server.go
│   ├── web_index_mcp_server.go
│   ├── page_crawl_mcp_server.go
│   ├── schema_validation_mcp_server.go
│   ├── rank_tracking_mcp_server.go
│   ├── local_listing_mcp_server.go
│   └── connector_shared.go
│
├── /data_security_mesh
│   ├── mtls_identity_verifier.go
│   ├── token_audience_validator.go
│   ├── credential_ref_resolver.go
│   ├── tenant_context_enforcer.go
│   ├── field_projection_enforcer.go
│   ├── row_scope_enforcer.go
│   ├── cmek_envelope_store.go
│   ├── storage_encryption.go
│   ├── transit_encryption.go
│   ├── result_data_classifier.go
│   ├── provenance_attestor.go
│   ├── taint_marker.go
│   └── data_mesh_event_emitter.go
│
├── /vector_rag_store
│   ├── vector_namespace_resolver.go
│   ├── vector_retrieval_engine.go
│   ├── hybrid_retrieval_engine.go
│   ├── tenant_partitioning.go
│   ├── metadata_filter_enforcer.go
│   ├── retrieval_budget_guard.go
│   ├── prompt_injection_taint_detector.go
│   ├── poisoning_guard.go
│   ├── evidence_packet_builder.go
│   └── memory_provenance_index.go
│
├── /resource_packaging
│   ├── evidence_packet.go
│   ├── connector_result_envelope.go
│   ├── citation_ref_builder.go
│   ├── freshness_ttl_resolver.go
│   ├── confidence_calibrator.go
│   ├── result_truncator.go
│   ├── deterministic_sorter.go
│   └── user_safe_summary_adapter.go
│
├── /interop_events
│   ├── event_schema.go
│   ├── sanitized_event_emitter.go
│   ├── health_event_emitter.go
│   └── audit_hash_builder.go
│
├── /tests
│   ├── /contract
│   ├── /transport
│   ├── /tenant_isolation
│   ├── /security
│   ├── /schema_drift
│   ├── /chaos
│   └── /golden_payloads
│
└── /docs
    ├── connector-onboarding.md
    ├── mcp-hardening.md
    ├── vector-retrieval-runbook.md
    ├── schema-migration-runbook.md
    ├── incident-disable-connector.md
    └── zero-overlap-boundaries.md
```

---

## 6. Configuration Contracts

### 6.1 `connector-registry.yaml`

This is the authoritative Layer 7 registry. It describes connectors, but does not authorize their use. Layer 2 authorizes use. Layer 7 verifies identity, schema, endpoint, tenant boundary, and protocol conformance.

```yaml
version: 2
registry_owner: layer_07_interop
default_mode: fail_closed

connectors:
  - connector_id: search_console
    display_category: search_performance_data
    status: enabled
    environment: production
    source_type: first_party_account_data
    protocol: mcp
    transport:
      allowed_modes:
        - streamable_http
      endpoint_binding: gsc_remote_prod
      requires_mtls: true
      requires_sse: false
      max_response_bytes: 1048576
      request_timeout_ms: 10000
      idle_timeout_ms: 30000
    auth:
      credential_source: layer_06_jit_ref
      required_audience: search_console_mcp
      required_scopes:
        - readonly_search_performance
      no_refresh_allowed: true
    tenant:
      required: true
      isolation_key: tenant_id
      account_binding: verified_site_property
      cross_tenant_join_allowed: false
    schema:
      manifest_version: 1.4.2
      manifest_hash: sha256:PINNED_MANIFEST_HASH
      input_schema_hash: sha256:PINNED_INPUT_SCHEMA_HASH
      output_schema_hash: sha256:PINNED_OUTPUT_SCHEMA_HASH
    result_controls:
      default_page_size: 100
      max_page_size: 500
      max_pages: 10
      field_projection_rule: gsc_standard_projection
      provenance_required: true
      taint_detection_required: true
      freshness_ttl_seconds: 86400
    observability:
      emit_tool_span: true
      emit_health_events: true
      redact_raw_payloads: true
```

### 6.2 `mcp-server-manifest.schema.json`

Every MCP server must publish a manifest that Layer 7 can verify.

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "title": "Layer7MCPServerManifest",
  "type": "object",
  "required": [
    "server_id",
    "server_version",
    "protocol_version",
    "transport_modes",
    "tool_contracts",
    "auth",
    "tenant_boundary",
    "schema_hashes",
    "provenance"
  ],
  "properties": {
    "server_id": { "type": "string", "pattern": "^[a-z0-9_\\-]+$" },
    "server_version": { "type": "string" },
    "protocol_version": { "type": "string" },
    "transport_modes": {
      "type": "array",
      "items": { "enum": ["stdio", "streamable_http"] }
    },
    "tool_contracts": {
      "type": "array",
      "items": {
        "type": "object",
        "required": [
          "tool_id",
          "input_schema_hash",
          "output_schema_hash",
          "risk_tier",
          "data_classes_returned"
        ],
        "properties": {
          "tool_id": { "type": "string" },
          "input_schema_hash": { "type": "string" },
          "output_schema_hash": { "type": "string" },
          "risk_tier": { "enum": ["read_only", "draft", "state_changing"] },
          "data_classes_returned": {
            "type": "array",
            "items": { "type": "string" }
          }
        }
      }
    },
    "auth": {
      "type": "object",
      "required": ["credential_source", "audience", "scopes"],
      "properties": {
        "credential_source": { "enum": ["layer_06_jit_ref", "none"] },
        "audience": { "type": "string" },
        "scopes": { "type": "array", "items": { "type": "string" } }
      }
    },
    "tenant_boundary": {
      "type": "object",
      "required": ["tenant_key_required", "cross_tenant_access"],
      "properties": {
        "tenant_key_required": { "type": "boolean" },
        "cross_tenant_access": { "const": false }
      }
    },
    "schema_hashes": {
      "type": "object",
      "required": ["manifest_hash"],
      "properties": {
        "manifest_hash": { "type": "string" }
      }
    },
    "provenance": {
      "type": "object",
      "required": ["source_identity_required", "result_hash_required"],
      "properties": {
        "source_identity_required": { "const": true },
        "result_hash_required": { "const": true }
      }
    }
  }
}
```

### 6.3 `source-contracts.yaml`

Layer 7 uses source contracts to normalize heterogenous sources into consistent grounding payloads.

```yaml
version: 2

source_contracts:
  search_console_query:
    source_type: first_party_search_performance
    freshness_class: recent
    required_fields:
      - query
      - page
      - country
      - device
      - clicks
      - impressions
      - ctr
      - position
      - date_range
    optional_fields:
      - search_appearance
      - branded_query_label
      - intent_bucket
    forbidden_fields:
      - raw_oauth_token
      - account_email
      - unmasked_user_identifier
    default_projection:
      - query
      - page
      - clicks
      - impressions
      - ctr
      - position
      - date_range
    evidence_usage:
      allowed_in_answer: true
      citation_required: true
      may_trigger_action: false

  analytics_page_report:
    source_type: first_party_behavioral_analytics
    freshness_class: recent
    required_fields:
      - page_path
      - sessions
      - users
      - engagement_rate
      - conversions
      - date_range
    forbidden_fields:
      - client_id
      - user_id
      - advertising_id
      - raw_ip
    default_projection:
      - page_path
      - sessions
      - engagement_rate
      - conversions
      - date_range
    evidence_usage:
      allowed_in_answer: true
      citation_required: false
      may_trigger_action: false

  brand_guideline_memory:
    source_type: tenant_memory
    freshness_class: durable
    required_fields:
      - memory_id
      - excerpt
      - source_document
      - updated_at
    forbidden_fields:
      - raw_secret
      - raw_personal_data
    default_projection:
      - memory_id
      - excerpt
      - source_document
      - updated_at
    evidence_usage:
      allowed_in_answer: true
      citation_required: false
      may_trigger_action: false
```

### 6.4 `vector-namespaces.yaml`

```yaml
version: 2

namespaces:
  - namespace_id: tenant_brand_memory
    tenant_key_required: true
    namespace_template: "tenant/{{tenant_id}}/brand_memory"
    allowed_query_sources:
      - layer_03_dag_node
      - layer_04_skill_script_via_layer_06
    allowed_document_classes:
      - brand_guideline
      - tone_profile
      - editorial_preference
      - product_context
    prohibited_document_classes:
      - raw_credentials
      - raw_user_pii
      - hidden_prompts
      - system_messages
      - connector_secrets
    retrieval_limits:
      top_k_default: 8
      top_k_max: 20
      max_total_chars: 12000
      max_doc_age_days: 730
    poisoning_defense:
      prompt_like_text_taint: true
      cross_tenant_similarity_reject: true
      source_document_required: true
```

---

## 7. Request and Response Envelopes

Layer 7 accepts only typed envelopes.

### 7.1 Inbound Connector Request

```go
type InteropRequest struct {
    RequestID           string            `json:"request_id"`
    SessionID           string            `json:"session_id"`
    TenantID            string            `json:"tenant_id"`
    ActorAgentID        string            `json:"actor_agent_id"`
    WorkflowNodeID      string            `json:"workflow_node_id"`
    ConnectorID         string            `json:"connector_id"`
    ToolID              string            `json:"tool_id"`
    IntentEnum          string            `json:"intent_enum"`
    PolicyDecisionID    string            `json:"policy_decision_id"`
    SandboxExecutionID  string            `json:"sandbox_execution_id,omitempty"`
    JITCredentialRef    string            `json:"jit_credential_ref,omitempty"`
    DataPurpose         string            `json:"data_purpose"`
    Input               map[string]any    `json:"input"`
    FieldProjection     []string          `json:"field_projection,omitempty"`
    PageSize            int               `json:"page_size,omitempty"`
    PageToken           string            `json:"page_token,omitempty"`
    DeadlineUnixMs      int64             `json:"deadline_unix_ms"`
    CallerTraceContext  map[string]string `json:"caller_trace_context,omitempty"`
}
```

Required validation:

```text
request_id must be globally unique
tenant_id must match credential tenant
connector_id must exist in connector-registry.yaml
tool_id must exist in pinned manifest for connector_id
policy_decision_id must be present
sandbox_execution_id must be present for runtime-originating calls
jit_credential_ref must be present for authenticated connectors
deadline must be within connector timeout budget
field projection must be a subset of source contract
input must validate against pinned schema hash
page_size must be <= connector max_page_size
```

### 7.2 Outbound Connector Result

```go
type InteropResult struct {
    RequestID          string                 `json:"request_id"`
    TenantID           string                 `json:"tenant_id"`
    ConnectorID        string                 `json:"connector_id"`
    ToolID             string                 `json:"tool_id"`
    Status             string                 `json:"status"`
    ResultSchemaHash   string                 `json:"result_schema_hash"`
    DataClassification string                 `json:"data_classification"`
    EvidencePackets    []EvidencePacket       `json:"evidence_packets"`
    Aggregates         map[string]any         `json:"aggregates,omitempty"`
    NextPageToken      string                 `json:"next_page_token,omitempty"`
    ProvenanceHash     string                 `json:"provenance_hash"`
    TaintFlags         []string               `json:"taint_flags,omitempty"`
    FreshnessTTL       int64                  `json:"freshness_ttl_seconds"`
    ConnectorHealth    string                 `json:"connector_health"`
    SanitizedEvents    []InteropEventRef      `json:"sanitized_events,omitempty"`
    Error              *InteropError          `json:"error,omitempty"`
}
```

### 7.3 Evidence Packet

```go
type EvidencePacket struct {
    EvidenceID          string            `json:"evidence_id"`
    SourceType          string            `json:"source_type"`
    SourceRef           string            `json:"source_ref"`
    SourceTitle         string            `json:"source_title,omitempty"`
    RetrievedAtUnixMs   int64             `json:"retrieved_at_unix_ms"`
    FreshnessTTLSeconds int64             `json:"freshness_ttl_seconds"`
    Confidence          float64           `json:"confidence"`
    DataClassification  string            `json:"data_classification"`
    TaintFlags          []string          `json:"taint_flags,omitempty"`
    Fields              map[string]any    `json:"fields"`
    ProvenanceHash      string            `json:"provenance_hash"`
    CitationRef         string            `json:"citation_ref,omitempty"`
}
```

Evidence packets are the only safe way for Layer 7 data to enter model context. Raw connector payloads must never be sent directly to the model.

---

## 8. MCP Transport Plane

### 8.1 `jsonrpc_codec.go`

Responsibilities:

```text
encode MCP JSON-RPC requests
decode MCP JSON-RPC responses
reject duplicate ids
reject missing ids
reject unknown method responses
reject batch requests unless connector manifest explicitly allows them
enforce UTF-8 validity
enforce max message size
canonicalize params before schema validation
```

Boundary:

```text
Does not authorize a method.
Does not retry failed business logic.
Does not redact user-visible output.
```

### 8.2 `stdio_client.go`

Production rule:

```text
Layer 7 owns the stdio JSON-RPC pipe, not host process isolation.
Layer 6 owns sandboxed process creation.
```

Flow:

```text
1. Layer 3 emits an already-authorized DAG tool envelope.
2. Layer 6 creates an isolated runtime process if the connector is local-stdio.
3. Layer 6 returns a sanctioned pipe handle reference to Layer 7.
4. Layer 7 attaches stdio_client.go to the pipe.
5. Layer 7 performs MCP handshake.
6. Layer 7 validates the server manifest and schema hashes.
7. Layer 7 executes the JSON-RPC call.
8. Layer 7 returns evidence packets.
9. Layer 6 tears down or resets the sandbox according to runtime policy.
```

Hard requirements:

```text
no unsandboxed local server launch
no shell expansion
no inherited environment secrets
no long-lived pipe reuse across tenants
no stderr leakage into model context
stderr may be summarized into sanitized health events only
```

### 8.3 `streamable_http_client.go`

Responsibilities:

```text
send MCP messages over remote Streamable HTTP
support SSE only when connector manifest allows it
validate mTLS identity
validate server certificate pin or workload identity
attach JIT credential only to the approved endpoint
enforce request timeout
enforce idle timeout
enforce response size limits
reject redirect unless endpoint binding explicitly allows redirect target
reject downgrade from HTTPS
reject interactive auth prompts
```

Hard requirements:

```text
all remote calls route through Layer 6 egress path
all remote calls use mTLS unless a connector is explicitly configured as public read-only
all authenticated remote calls require Layer 6 JIT credential reference
no cookies persisted by Layer 7
no browser session reuse
no OAuth refresh flow inside Layer 7
```

### 8.4 `handshake_validator.go`

The handshake validator must verify:

```text
server_id matches registry
server_version allowed by registry
protocol_version compatible
tool list exactly matches pinned manifest for production
input/output schemas match pinned hashes
server declared auth matches registry auth
server declared tenant boundary requires tenant key
server capabilities do not exceed registry capabilities
```

If a server returns extra tools, Layer 7 must fail closed. Extra tools are a spoofing signal, not a harmless addition.

### 8.5 `schema_hash_verifier.go`

All schemas are pinned by cryptographic hash. The verifier must run before a tool call and after a response.

```text
before call:
  validate input against pinned input schema

after call:
  validate result against pinned output schema
  reject additional fields unless schema permits them
  reject missing provenance fields
  reject unexpected sensitive fields
```

### 8.6 `circuit_breaker.go`

Layer 7 owns connector health state, not agent quarantine.

Circuit breaker states:

```text
closed: connector healthy
open: connector disabled for calls
half_open: test calls allowed using synthetic non-sensitive payloads only
degraded: connector allowed for read-only low-risk calls with stricter budgets
```

Circuit breaker triggers:

```text
schema mismatch
mTLS failure
auth audience mismatch
5xx burst
timeout burst
latency SLO breach
tool list drift
malformed JSON-RPC burst
cross-tenant access attempt
tainted payload spike
```

When the breaker opens, Layer 7 emits `connector.health.opened` to Layer 8. It does not quarantine the agent.

---

## 9. Connector Registry and Manifest Lifecycle

### 9.1 Connector Lifecycle States

```text
draft
dev_verified
staging_verified
shadow
canary
production
degraded
disabled
retired
```

Promotion gates:

```text
manifest validates
schema hashes pinned
mTLS identity tested
JIT credential audience tested
tenant isolation tests pass
field projection tests pass
provenance tests pass
taint detection tests pass
golden payload tests pass
chaos transport tests pass
Layer 8 receives sanitized event mapping
rollback manifest exists
```

### 9.2 Registry Governance

No production connector may be sourced directly from a public registry.

Allowed production sources:

```text
internal registry
official managed connector
vendor-reviewed connector
signed private connector artifact
```

Public community MCP servers may be used only in isolated development, with:

```text
no production credentials
no customer data
no shared tenant memory
no privileged filesystem mounts
no write-capable scopes
explicit "untrusted_public_registry" label
```

### 9.3 Manifest Drift Handling

When the live MCP server manifest differs from the pinned registry manifest:

```text
1. Stop the call.
2. Mark connector instance as schema_drift_detected.
3. Emit sanitized event to Layer 8.
4. Open circuit breaker for that connector/tool pair.
5. Return a structured failure to Layer 3.
6. Do not attempt automatic compatibility mapping.
```

---

## 10. Enterprise MCP Server Connectors

Layer 7 connector servers expose typed enterprise resources as MCP-compatible endpoints. Each connector must be narrow and source-specific.

### 10.1 Search Performance Connector

Owned capabilities:

```text
query performance retrieval
page performance retrieval
query/page comparison
date range trend extraction
country/device segmentation
branded/non-branded labeling when configured
```

Must not do:

```text
choose SEO strategy
write recommendations
decide ranking priority
render audit dashboard
change Search Console settings
store raw OAuth tokens
```

### 10.2 Analytics Connector

Owned capabilities:

```text
traffic report retrieval
page engagement metrics
conversion metrics
channel/source grouping
AI referral detection if configured
date range comparison
```

Must not do:

```text
deanonymize users
retrieve raw user identifiers
export unrestricted event logs
create marketing strategy
```

### 10.3 Repository Connector

Owned capabilities:

```text
read whitelisted repository files
search code within authorized scopes
list branches/tags
retrieve package manifests
retrieve structured metadata needed by audits
```

Must not do:

```text
write source code
open pull requests
merge branches
read secrets
read CI secrets
read hidden prompts
read unapproved directories
```

Write actions, if ever supported, must route through Layer 1 approval, Layer 2 policy, Layer 6 runtime, and dedicated write-capable connector contracts. They are not a default Layer 7 capability.

### 10.4 CMS Connector

Owned capabilities:

```text
read page metadata
read content drafts
read published URLs
read content type schemas
retrieve authoring fields according to projection
optionally draft non-published edits only when policy and workflow permit
```

Must not do:

```text
publish content
delete content
overwrite existing fields without approval
bypass CMS workflow
retrieve unrestricted editor PII
```

### 10.5 Review Intelligence Connector

Owned capabilities:

```text
retrieve product/service reviews
retrieve competitor review summaries when allowed
aggregate rating distributions
surface review themes with provenance
```

Must not do:

```text
scrape private customer records
fabricate review summaries
post review responses
store raw reviewer identities unless explicitly classified and permitted
```

### 10.6 Web Index / SERP Connector

Owned capabilities:

```text
retrieve public search results from approved sources
return ranked URLs with snippets
return SERP feature observations
return citation candidates
```

Must not do:

```text
simulate a user browser session
solve CAPTCHAs
bypass robots or site restrictions
use unapproved scraping endpoints
hide source provenance
```

### 10.7 Page Crawl Connector

Owned capabilities:

```text
retrieve page HTML through approved crawl tools
return status code, canonical, title, headers, visible text sample
extract structured data
extract internal links
extract robots directives
```

Must not do:

```text
execute arbitrary JavaScript outside approved render service
open direct shell network access
perform destructive scanning
bypass Layer 6 egress
```

### 10.8 Schema Validation Connector

Owned capabilities:

```text
validate JSON-LD syntax
validate schema type requirements
return deterministic validation errors
return normalized schema graph
```

Must not do:

```text
invent schema strategy
modify site code
publish schema
```

### 10.9 Rank Tracking Connector

Owned capabilities:

```text
retrieve rank observations
retrieve keyword position trends
retrieve competitor visibility where contracted
```

Must not do:

```text
infer business strategy without Layer 4 skill guidance
join rank data across tenants
```

### 10.10 Local Listing Connector

Owned capabilities:

```text
retrieve location profile status
retrieve listing completeness signals
retrieve local ranking observations where available
retrieve review counts and rating snapshots
```

Must not do:

```text
reply to reviews
change business hours
publish profile updates
```

---

## 11. Data Security Mesh

### 11.1 `mtls_identity_verifier.go`

Layer 7 must verify machine identity on every remote connector path.

Required checks:

```text
certificate chain valid
workload identity matches endpoint binding
tenant-bound connector identity matches requested tenant
certificate not expired
certificate not revoked
protocol uses TLS 1.3 or approved equivalent
server name matches manifest binding
```

No plaintext remote connector calls are allowed in production.

### 11.2 `token_audience_validator.go`

Layer 7 validates Layer 6 JIT credentials. It does not mint or refresh them.

Checks:

```text
token audience equals connector required audience
token scope is subset of required scopes
token tenant claim equals request tenant_id
token expiry exceeds expected call duration but remains short-lived
token purpose equals data_purpose
token is not reusable across connector_id mismatch
```

If validation fails, Layer 7 returns:

```json
{
  "status": "blocked",
  "error": {
    "code": "TOKEN_AUDIENCE_REJECTED",
    "safe_message": "The connector credential was not valid for this data source."
  }
}
```

### 11.3 `tenant_context_enforcer.go`

Every connector call must include a tenant context. The tenant context enforcer ensures:

```text
tenant_id present
tenant_id bound to credential
tenant_id bound to account/property/site/resource
tenant namespace selected before retrieval
tenant filters applied before query
cross-tenant joins disabled by default
tenant leakage tests pass in CI
```

### 11.4 `field_projection_enforcer.go`

Layer 7 must minimize data before it enters agent context.

Field projection rules:

```text
default projection is source-specific
requested projection must be subset of allowed projection
forbidden fields are always removed
sensitive fields require explicit policy and contract support
result envelope records projection applied
```

Layer 7 must not rely on the model to ignore sensitive fields.

### 11.5 `cmek_envelope_store.go`

Layer 7 provides encryption primitives for stored connector artifacts, caches, and retrieval indexes.

Rules:

```text
tenant-specific key references
envelope encryption for cached payloads
key rotation supported
no plaintext cached payloads at rest
cache purge by tenant
cache purge by connector
cache purge by data classification
```

Layer 7 does not own the business artifact lifecycle. It only provides encrypted storage plumbing for the data plane.

### 11.6 `provenance_attestor.go`

Every result must have a provenance hash.

Provenance input includes:

```text
connector_id
tool_id
server_id
server_version
schema_hash
tenant_id
source_ref
retrieval_timestamp
field_projection
canonicalized result payload
```

The hash allows downstream layers to prove that an answer or UI card was grounded in a specific data payload without storing raw sensitive data in logs.

### 11.7 `taint_marker.go`

Taint flags are mandatory whenever payload text may influence the model.

Possible taint flags:

```text
external_untrusted_text
prompt_like_text_detected
html_hidden_text_present
zero_width_chars_removed
script_or_style_removed
mixed_tenant_similarity_rejected
source_freshness_expired
schema_partial_result
connector_degraded
public_registry_source
low_confidence_source
```

Taint flags are metadata, not a policy decision. Layer 2 and Layer 8 use them for gating and analysis.

---

## 12. Vector RAG Store

Layer 7 owns retrieval plumbing, not procedural knowledge.

### 12.1 What Layer 7 Retrieval Is For

Allowed retrieval purposes:

```text
brand guidelines
tone preferences
editorial constraints
known product/service facts
approved entity descriptions
historical audit findings
approved style examples
previous user-approved briefs
durable tenant memory
```

Not allowed:

```text
hidden prompts
system messages
raw logs
trace dumps
raw user PII
connector secrets
unreviewed adversarial content
cross-tenant memory
```

### 12.2 Retrieval Flow

```text
1. Receive InteropRequest with tenant_id and retrieval purpose.
2. Resolve vector namespace using tenant_id.
3. Apply metadata filters before similarity search.
4. Reject any namespace mismatch.
5. Run hybrid retrieval if configured.
6. Apply top_k and character budget.
7. Run prompt-like taint detection.
8. Build evidence packets.
9. Attach source provenance.
10. Return only packaged evidence, never raw vector rows.
```

### 12.3 Poisoning Defense

Layer 7 must protect against memory poisoning and cross-tenant vector poisoning.

Controls:

```text
tenant namespace isolation
source document required
embedding ingest provenance
document class allowlist
prompt-like text tainting
untrusted source quarantine label
similarity outlier detection
cross-tenant nearest-neighbor guard
stale memory decay
human-approved memory flag for durable preferences
```

Layer 7 may mark suspicious retrieval as tainted or rejected. Layer 8 may use these events for security analysis.

### 12.4 Evidence Packet Context Budget

Layer 7 must respect context economics.

Default retrieval budgets:

```yaml
retrieval_budget:
  brand_memory:
    top_k: 8
    max_total_chars: 12000
    max_chars_per_packet: 2000
  audit_history:
    top_k: 5
    max_total_chars: 8000
    max_chars_per_packet: 1600
  source_citation_candidates:
    top_k: 12
    max_total_chars: 10000
    max_chars_per_packet: 1200
```

Layer 7 must deterministically truncate and mark truncation. It must not ask the model to decide which raw memory rows to keep.

---

## 13. Resource Packaging

Layer 7 output must be optimized for downstream reasoning without polluting the context window.

### 13.1 Evidence Packet Rules

Each evidence packet must be:

```text
small
typed
tenant-bound
source-attributed
time-stamped
hash-attested
classification-labeled
taint-labeled
field-projected
deterministically ordered
```

### 13.2 Deterministic Sorting

For repeatability, Layer 7 must sort results deterministically unless the source has a meaningful rank order.

Sort order examples:

```text
search results: source rank, then URL
analytics pages: descending sessions, then page_path
GSC queries: descending clicks, then impressions, then query
memory retrieval: score desc, then updated_at desc, then memory_id
validation errors: severity desc, then path asc
```

### 13.3 Confidence Calibration

Layer 7 may provide source confidence metadata, but must not decide recommendation quality.

Confidence inputs:

```text
source authority
source freshness
schema completeness
connector health
retrieval score
field completeness
duplicate corroboration
```

Output:

```json
{
  "confidence": 0.86,
  "confidence_basis": [
    "fresh_first_party_source",
    "schema_complete",
    "connector_healthy"
  ]
}
```

---

## 14. Production Flows

### 14.1 Connector Registration Flow

```text
01 developer submits connector manifest
02 manifest_validator.go validates structure
03 registry_signature_verifier.go validates artifact signature
04 schema_hash_verifier.go pins input/output schema hashes
05 tenant isolation tests run
06 mTLS identity test runs
07 JIT token audience test runs
08 golden payload tests run
09 schema drift tests run
10 chaos transport tests run
11 shadow mode runs with synthetic or read-only data
12 canary mode runs on low-risk traffic
13 production state approved
14 connector_registry emits sanitized promotion event
```

### 14.2 Tool Call Flow

```text
01 Layer 3 selects workflow node
02 Layer 2 has already authorized tool/action
03 Layer 6 creates sandbox and JIT credential reference if needed
04 Layer 7 receives InteropRequest
05 request_envelope_validator.go validates required fields
06 registry_loader.go loads connector contract
07 token_audience_validator.go validates credential reference
08 tenant_context_enforcer.go binds tenant/resource/account
09 schema_hash_verifier.go validates input schema
10 transport client performs MCP handshake
11 mcp_server_identity_verifier.go verifies server identity
12 connector call executes
13 schema_hash_verifier.go validates output schema
14 field_projection_enforcer.go minimizes result fields
15 provenance_attestor.go hashes result
16 taint_marker.go labels risk metadata
17 evidence_packet_builder.go packages result
18 sanitized_event_emitter.go emits interop events to Layer 8
19 InteropResult returns to caller
```

### 14.3 Retrieval Flow

```text
01 Layer 3 requests tenant memory retrieval through an authorized node
02 Layer 7 validates retrieval purpose and tenant
03 namespace resolver selects tenant namespace
04 metadata filters apply before search
05 vector/hybrid retrieval runs
06 prompt-like text taint detector labels suspicious evidence
07 poisoning guard rejects cross-tenant or source-invalid rows
08 evidence packets are built
09 result budget guard truncates deterministically
10 evidence packets return with provenance and taint metadata
```

### 14.4 Schema Drift Flow

```text
01 Live server returns changed tool schema
02 handshake_validator detects schema hash mismatch
03 Layer 7 blocks call
04 connector_status_store marks connector/tool as schema_drift_detected
05 circuit_breaker opens connector/tool pair
06 sanitized event emitted to Layer 8
07 structured failure returned to Layer 3
08 connector requires manifest update and CI gates before re-enable
```

### 14.5 Connector Health Degradation Flow

```text
01 latency, timeout, 5xx, malformed response, or drift threshold exceeded
02 circuit_breaker changes connector state to degraded or open
03 connector_result_envelope returns safe failure or degraded flag
04 health_event_emitter emits sanitized event
05 Layer 8 observes pattern and may score trajectory impact
06 Layer 6 may quarantine runtime only if Layer 8/blue-team policy triggers it
```

---

## 15. Security Controls Specific to MCP Spoofing

MCP spoofing is a Layer 7 priority because attackers can attempt to present a malicious server or altered tool list as a trusted tool source.

Layer 7 defenses:

```text
server_id pinning
mTLS identity verification
manifest signature verification
protocol version negotiation
tool list exact-match validation
input/output schema hash pinning
endpoint binding
JIT token audience validation
field projection
extra-tool fail-closed behavior
schema drift circuit breaker
sanitized event emission
```

Spoofing examples Layer 7 must block:

```text
trusted server returns an extra write-capable tool
remote endpoint redirects to unregistered host
server advertises same tool name with changed parameters
server returns additional sensitive fields
server changes auth mode from JIT to interactive OAuth
server requests broad OAuth scope
stdio server starts with a different binary hash
public registry connector claims same display name as internal connector
```

---

## 16. Data Classes

Layer 7 must propagate data class labels. It must not guess if the registry declares the class.

Recommended classes:

```text
public_web_data
first_party_search_performance
first_party_analytics_aggregate
tenant_brand_memory
tenant_content_draft
repository_metadata
repository_source_readonly
cms_metadata
cms_content_readonly
review_aggregate
local_listing_metadata
structured_validation_output
sensitive_business_data
restricted_personal_data
secret_or_credential
hidden_internal_metadata
```

Forbidden by default in Layer 7 result packets:

```text
secret_or_credential
hidden_internal_metadata
raw_personal_data
raw_oauth_token
raw_cookie
system_prompt
hidden_chain_of_thought
raw_trace_dump
```

If a connector unexpectedly returns a forbidden class, Layer 7 must block the result and emit a sanitized event.

---

## 17. Observability Event Contract to Layer 8

Layer 7 emits facts, not raw payloads.

### 17.1 Event Types

```text
interop.request.accepted
interop.request.rejected
mcp.handshake.started
mcp.handshake.succeeded
mcp.handshake.failed
mcp.schema.verified
mcp.schema.mismatch
mcp.server_identity.verified
mcp.server_identity.failed
connector.call.started
connector.call.succeeded
connector.call.failed
connector.health.degraded
connector.health.opened
connector.health.recovered
connector.registry.promoted
connector.registry.disabled
connector.public_registry.rejected
data_mesh.tenant_bound
data_mesh.cross_tenant_blocked
data_mesh.token_audience_rejected
data_mesh.field_projection_applied
data_mesh.forbidden_field_blocked
retrieval.query.started
retrieval.query.succeeded
retrieval.query.rejected
retrieval.poisoning_suspected
retrieval.prompt_like_text_tainted
result.provenance_attested
result.payload_truncated
```

### 17.2 Event Shape

```go
type Layer7Event struct {
    EventID             string            `json:"event_id"`
    EventType           string            `json:"event_type"`
    TimestampUnixMs     int64             `json:"timestamp_unix_ms"`
    RequestID           string            `json:"request_id"`
    SessionID           string            `json:"session_id,omitempty"`
    TenantHash          string            `json:"tenant_hash"`
    ConnectorID         string            `json:"connector_id,omitempty"`
    ToolID              string            `json:"tool_id,omitempty"`
    WorkflowNodeID      string            `json:"workflow_node_id,omitempty"`
    PolicyDecisionID    string            `json:"policy_decision_id,omitempty"`
    SandboxExecutionID  string            `json:"sandbox_execution_id,omitempty"`
    Outcome             string            `json:"outcome"`
    ReasonCode          string            `json:"reason_code,omitempty"`
    LatencyMs           int64             `json:"latency_ms,omitempty"`
    PayloadBytes        int64             `json:"payload_bytes,omitempty"`
    SchemaHash          string            `json:"schema_hash,omitempty"`
    ProvenanceHash      string            `json:"provenance_hash,omitempty"`
    TaintFlags          []string          `json:"taint_flags,omitempty"`
    Attributes          map[string]string `json:"attributes,omitempty"`
}
```

Rules:

```text
no raw connector payloads
no raw prompts
no raw retrieved memory
no secrets
no cookies
no OAuth tokens
no account emails unless explicitly hashed
no hidden reasoning
tenant identifiers must be hashed or pseudonymous in events
```

---

## 18. Error Model

Layer 7 errors must be structured and safe.

```go
type InteropError struct {
    Code           string            `json:"code"`
    SafeMessage    string            `json:"safe_message"`
    Retryable      bool              `json:"retryable"`
    ConnectorState string            `json:"connector_state,omitempty"`
    Details        map[string]string `json:"details,omitempty"`
}
```

Allowed error codes:

```text
CONNECTOR_UNKNOWN
CONNECTOR_DISABLED
CONNECTOR_DEGRADED
POLICY_DECISION_MISSING
SANDBOX_REF_MISSING
JIT_CREDENTIAL_MISSING
TOKEN_AUDIENCE_REJECTED
TENANT_CONTEXT_MISSING
TENANT_BOUNDARY_VIOLATION
FIELD_PROJECTION_DENIED
INPUT_SCHEMA_INVALID
OUTPUT_SCHEMA_INVALID
MANIFEST_SIGNATURE_INVALID
MCP_SERVER_IDENTITY_FAILED
MCP_SCHEMA_MISMATCH
MCP_PROTOCOL_UNSUPPORTED
MCP_TOOL_LIST_DRIFT
TRANSPORT_TIMEOUT
TRANSPORT_TLS_FAILED
TRANSPORT_REDIRECT_DENIED
PAYLOAD_TOO_LARGE
PAGINATION_LIMIT_EXCEEDED
FORBIDDEN_FIELD_RETURNED
PROVENANCE_MISSING
TAINT_METADATA_MISSING
RETRIEVAL_NAMESPACE_DENIED
RETRIEVAL_POISONING_SUSPECTED
SOURCE_FRESHNESS_EXPIRED
```

Safe message examples:

```text
"The requested connector is not enabled for this environment."
"The connector response did not match its approved schema."
"The credential was not valid for this data source."
"The requested retrieval namespace is not available for this tenant."
"The data source is temporarily degraded; try again later."
```

Never include raw stack traces, internal filesystem paths, secret names, raw endpoint URLs, or token values in errors returned to the agent.

---

## 19. Testing Strategy

### 19.1 Contract Tests

Required for every connector:

```text
manifest validates
input schemas reject bad inputs
output schemas reject unexpected fields
field projection removes forbidden fields
pagination limits hold
result envelopes include provenance
taint flags appear for prompt-like data
safe errors never leak internals
```

### 19.2 Transport Tests

```text
stdio JSON-RPC happy path
stdio malformed JSON rejection
stdio oversized payload rejection
streamable HTTP happy path
SSE allowed only when manifest permits
remote redirect denied
TLS downgrade denied
mTLS wrong identity denied
timeout behavior
retry budget behavior
circuit breaker opens and recovers
```

### 19.3 Tenant Isolation Tests

```text
Tenant A cannot query Tenant B namespace
Tenant A token cannot access Tenant B connector account
cross-tenant vector nearest-neighbor is rejected
repository account binding enforces tenant
CMS property binding enforces tenant
analytics property binding enforces tenant
cache purge by tenant removes encrypted payloads
```

### 19.4 Security Tests

```text
MCP spoofed server id blocked
MCP extra tool blocked
MCP changed schema hash blocked
public registry connector rejected in production
connector returns forbidden sensitive field blocked
prompt-like retrieved memory tainted
hidden HTML text tainted
zero-width text tainted or normalized
JIT token with broad audience rejected
expired credential rejected
interactive OAuth request rejected
```

### 19.5 Schema Drift Tests

```text
minor compatible schema change with same hash impossible
new field without manifest approval rejected
removed required field rejected
changed enum rejected
changed data classification rejected
tool rename rejected
```

### 19.6 Chaos Tests

```text
remote 500 burst
remote slow streaming
partial SSE chunks
connection reset mid-response
malformed batch result
duplicate JSON-RPC id
out-of-order response id
large page token loop
connector returns stale cache marker
mTLS certificate rotation
KMS transient failure
```

---

## 20. CI/CD Gates

A connector cannot merge to production unless all gates pass.

```text
go test ./layer_07_interop/...
manifest signature check
schema hash pinning check
contract golden tests
tenant isolation tests
transport chaos tests
forbidden-field tests
safe-error tests
event schema tests
provenance hash determinism tests
taint detector tests
SLO budget tests
public registry production ban
Layer 2 policy compatibility check
Layer 6 credential audience compatibility check
Layer 8 event schema compatibility check
```

Required artifacts:

```text
connector manifest
input/output schemas
golden request fixtures
golden response fixtures
tenant isolation fixture
taint fixture
schema drift fixture
rollback manifest
deprecation plan
owner
runbook
```

### 20.1 Executable interop and provider baseline

The repository implementation must retain:

```text
connector registry and source contracts decode strictly and fail closed at startup
the MCP manifest schema compiles at startup
tenant context comparison rejects empty or mismatched tenant identifiers
mTLS identity validation verifies the full certificate chain and exact expected SPIFFE identity
credential validation binds audience, tenant, connector, exact scopes, and validity window
field projection returns only allowlisted fields and deep-copies returned values
provenance attestation verifies an Ed25519 signature over the canonical payload hash
taint handling marks untrusted evidence without returning raw protected source bodies
vector namespaces are tenant-HMAC-derived and are checked before and after backend retrieval
retrieval rejects mismatched credentials, unbounded dimensions or top-k, NaN values, and missing provenance
MCP transport validates strict JSON-RPC, handshake identity and capabilities, schema hashes, bounded circuit breaking, HTTPS/mTLS, no redirects, and Layer 6-owned stdio pipes
```

The Gemini evaluation adapter is a Layer 7 connector and must:

```text
read credentials only from GEMINI_API_KEY at process construction
send the credential only in the x-goog-api-key header
use a fixed Google Interactions API endpoint
reject redirects
use a bounded timeout and response size
set store=false and background=false
request application/json structured output with a closed schema
set temperature=0, a fixed seed, and thinking_summaries=none
never log requests, candidates, test prompts, provider bodies, or credentials
return raw structured output to Layer 8 for independent validation and scoring
```

Unit tests must use an injected transport and synthetic credential. A live API
test is valid only with a newly provisioned environment credential and is not
proof of connector production readiness by itself.

### 20.2 Repository readiness evidence

The repository-level production check is:

```text
go run ./cmd/readiness -root .
```

For Layer 7, the check must fail when required connector-registry,
source-contract, or MCP-manifest schema artifacts are missing, unreadable, or
placeholder-only. It must also report absent executable source controls
explicitly required by this spec, including MCP transport, handshake and
schema-hash validation, mTLS identity, credential audience validation,
projection enforcement, provenance and taint handling, and tenant-partitioned
retrieval.

Artifact or source-file presence is not connector certification. Manifest
signature checks, contract golden tests, tenant isolation tests, transport
chaos tests, schema-drift tests, provenance tests, SLO checks, and deployment
attestations remain mandatory.

`cmd/readiness` and `internal/releasegate` are read-only platform CI tooling.
They must not connect MCP, call connectors, retrieve data, consume credentials,
access tenant data, mutate connector state, or perform any Layer 7 runtime
behavior.

---

## 21. Performance and SLOs

Layer 7 must optimize transport and retrieval without sacrificing safety.

Default SLOs:

```yaml
transport_slo:
  mcp_handshake_p95_ms: 250
  local_stdio_call_p95_ms: 500
  remote_read_call_p95_ms: 2500
  remote_heavy_report_p95_ms: 10000
  vector_retrieval_p95_ms: 800
  schema_validation_p95_ms: 100
  evidence_packaging_p95_ms: 100

availability_slo:
  production_connector_monthly: 99.5
  tenant_isolation_failure_rate: 0
  forbidden_field_leakage_rate: 0
  schema_drift_unblocked_rate: 0

budget_limits:
  default_max_response_bytes: 1048576
  default_max_evidence_packets: 20
  default_max_total_evidence_chars: 12000
  default_remote_timeout_ms: 10000
  default_retries: 2
  retry_jitter_required: true
```

Layer 7 may cache only when:

```text
cache policy exists
tenant key exists
data classification allows cache
payload is encrypted at rest
freshness TTL is recorded
provenance hash is preserved
cache key includes tenant_id and field projection
```

---

## 22. Safe Degradation

Layer 7 degradation must preserve trust.

Allowed degradation behaviors:

```text
return partial evidence with partial_result taint flag
reduce page size
disable expensive aggregation
use cached result if TTL valid and source contract permits
fallback from remote connector to internal verified cache
mark connector health degraded
return structured safe error
```

Forbidden degradation behaviors:

```text
switch to unverified public connector
drop tenant filters
skip schema validation
skip mTLS verification
use broad credentials
return raw payload
ask model to "reason around" missing data
invent missing metrics
```

---

## 23. Connector Onboarding Checklist

A new connector is production-ready only when this checklist is complete.

```text
[ ] Connector has a single clear source responsibility.
[ ] Connector manifest is complete.
[ ] Owner and escalation path are defined.
[ ] Transport mode is declared.
[ ] Endpoint binding is pinned.
[ ] mTLS identity is verified.
[ ] JIT credential audience is declared.
[ ] No long-lived credentials are stored in Layer 7.
[ ] Tenant boundary is explicit.
[ ] Field projection rules exist.
[ ] Forbidden fields are listed.
[ ] Input schema hash is pinned.
[ ] Output schema hash is pinned.
[ ] Result provenance is mandatory.
[ ] Taint behavior is tested.
[ ] Pagination limits exist.
[ ] Response budget exists.
[ ] Retry and timeout budget exists.
[ ] Circuit breaker behavior is configured.
[ ] Safe error messages are tested.
[ ] Sanitized events are mapped.
[ ] Golden payloads exist.
[ ] Tenant isolation tests pass.
[ ] Security tests pass.
[ ] Chaos tests pass.
[ ] Rollback manifest exists.
[ ] Public registry usage is banned in production unless explicitly reviewed as read-only and no-credential.
```

---

## 24. Deployment Modes

### 24.1 Development

```text
may use public connectors
no production credentials
synthetic tenant only
no customer data
verbose local logs permitted if secret scanner passes
schema drift allowed only with local warning
```

### 24.2 Staging

```text
internal or managed connectors only
staging credentials only
staging tenant data only
schema drift blocks call
mTLS required
event emission required
tenant isolation tests required
```

### 24.3 Production

```text
signed connector manifests only
pinned schemas only
mTLS required
JIT credentials only
CMEK encryption required
tenant isolation required
public registry banned by default
schema drift fails closed
safe errors only
sanitized events only
```

---

## 25. Layer 7 Active Session State

The active session state is not a public API. It is internal state passed across typed envelopes and sanitized telemetry.

```text
session.state.layer_07 {
  interop_request_id        : "req_..."
  active_connector_category : "search_performance_data"
  transport_mode            : "streamable_http"
  mcp_protocol_version      : "pinned_supported"
  mcp_manifest_status       : "verified"
  schema_hash_status        : "verified"
  tenant_partition_status   : "isolated"
  credential_ref_status     : "validated_audience_scope_tenant"
  field_projection_status   : "applied"
  provenance_status         : "attested"
  taint_status              : "none_or_labeled"
  connector_health          : "healthy|degraded|open"
  evidence_packet_count     : 0
  result_budget_status      : "within_limit|truncated"
}
```

Layer 7 state must not include:

```text
raw token
raw secret
raw connector payload
raw retrieved memory
raw hidden prompt
raw trace
raw user PII
```

---

## 26. Implementation Interfaces

### 26.1 Connector Interface

```go
type Connector interface {
    ID() string
    Manifest(ctx context.Context) (ConnectorManifest, error)
    ValidateInput(ctx context.Context, req InteropRequest) error
    Execute(ctx context.Context, req InteropRequest, credential ScopedCredential) (ConnectorRawResult, error)
    ValidateOutput(ctx context.Context, raw ConnectorRawResult) error
    Package(ctx context.Context, raw ConnectorRawResult, req InteropRequest) (InteropResult, error)
}
```

### 26.2 Transport Interface

```go
type MCPTransport interface {
    Connect(ctx context.Context, binding EndpointBinding) (MCPConnection, error)
    Handshake(ctx context.Context, conn MCPConnection) (MCPServerManifest, error)
    Call(ctx context.Context, conn MCPConnection, method string, params map[string]any) (json.RawMessage, error)
    Close(ctx context.Context, conn MCPConnection) error
}
```

### 26.3 Data Mesh Interface

```go
type DataMesh interface {
    ValidateCredential(ctx context.Context, req InteropRequest) (ScopedCredential, error)
    EnforceTenant(ctx context.Context, req InteropRequest, manifest ConnectorManifest) error
    ApplyProjection(ctx context.Context, sourceContract SourceContract, raw map[string]any) (map[string]any, error)
    AttestProvenance(ctx context.Context, packet EvidencePacket) (string, error)
    Classify(ctx context.Context, sourceContract SourceContract, fields map[string]any) (string, error)
    MarkTaint(ctx context.Context, fields map[string]any) ([]string, error)
}
```

### 26.4 Retrieval Interface

```go
type RetrievalEngine interface {
    ResolveNamespace(ctx context.Context, tenantID string, purpose string) (VectorNamespace, error)
    Query(ctx context.Context, namespace VectorNamespace, q RetrievalQuery) ([]RetrievedDocument, error)
    BuildEvidence(ctx context.Context, docs []RetrievedDocument, req InteropRequest) ([]EvidencePacket, error)
}
```

---

## 27. Threat Model

### 27.1 Threats Layer 7 Must Handle

```text
malicious MCP server impersonation
MCP tool list spoofing
schema drift
remote endpoint redirect
credential audience confusion
confused-deputy connector call
cross-tenant vector retrieval
cross-tenant analytics property access
forbidden field leakage
prompt injection embedded in retrieved content
hidden HTML or zero-width text in crawled pages
public registry connector substitution
stale cache presented as fresh data
oversized result flooding context
pagination loops
malformed JSON-RPC responses
SSE stream injection
mTLS downgrade
silent source provenance loss
```

### 27.2 Threats Owned Elsewhere

```text
raw user prompt injection blocking: Layer 2
agent workflow drift: Layer 3 and Layer 8
skill trigger overlap: Layer 4 and Layer 8
UI injection in rendered cards: Layer 5
sandbox escape: Layer 6
JIT credential minting/revocation: Layer 6
quarantine and auto-refactoring: Layer 6 / Layer 8
trajectory quality evaluation: Layer 8
```

---

## 28. SEO/AEO Grounding Requirements

Because the platform is an SEO/AEO auditor and content-writing agent, Layer 7 must support the following source categories with strict source separation:

```text
search performance data
analytics aggregate data
public SERP data
public crawl data
structured data validation
repository metadata
CMS content metadata
brand memory
editorial guidelines
review intelligence
local listing metadata
rank tracking data
```

For every SEO/AEO evidence packet, Layer 7 must include:

```text
source category
source freshness
tenant binding
URL or page identifier when safe
metric date range when applicable
projection fields
provenance hash
confidence
taint flags
```

For answer-engine optimization and audit use cases, Layer 7 must preserve source traceability. If a recommendation later appears in an audit report, Layer 8 and Layer 5 must be able to refer back to a safe evidence identifier rather than raw connector payload.

---

## 29. Runbooks

### 29.1 Disable a Compromised Connector

```text
1. Set connector status to disabled in connector-registry.yaml or runtime status store.
2. Open circuit breaker for connector_id.
3. Revoke active Layer 6 JIT tokens by requesting revocation through the credential control contract.
4. Emit connector.registry.disabled event.
5. Confirm no active sessions can call connector.
6. Preserve sanitized evidence hashes for forensics.
7. Do not delete raw logs in Layer 8 retention unless retention policy requires.
8. Patch manifest or rollback to last known good version.
9. Re-run security, tenant isolation, and schema drift tests.
10. Promote only after staging and canary pass.
```

Layer 7 requests revocation through a contract; it does not revoke directly.

### 29.2 Rotate Connector Endpoint Certificate

```text
1. Add new certificate identity to endpoint binding in staging.
2. Verify mTLS handshake.
3. Run schema and golden payload tests.
4. Deploy new binding to production with overlap window.
5. Remove old identity after all connections drain.
6. Emit endpoint identity rotation event.
```

### 29.3 Schema Migration

```text
1. Add new schema version as draft.
2. Generate schema hashes.
3. Run backward compatibility check.
4. Run golden fixtures.
5. Shadow new schema with read-only traffic.
6. Canary connector under new manifest.
7. Promote manifest.
8. Retire old schema after configured TTL.
```

### 29.4 Vector Namespace Incident

```text
1. Disable namespace for reads.
2. Mark namespace status as suspect.
3. Emit retrieval.poisoning_suspected.
4. Export provenance hashes for impacted evidence only.
5. Rebuild index from approved source documents.
6. Re-run cross-tenant isolation tests.
7. Re-enable namespace after security approval.
```

---

## 30. Production Acceptance Criteria

Layer 7 is production-ready only when all conditions are true:

```text
[ ] All connectors are manifest-driven.
[ ] All production manifests are signed and pinned.
[ ] All production remote transports use mTLS.
[ ] All authenticated calls consume Layer 6 JIT credential references only.
[ ] No long-lived secrets are stored in Layer 7.
[ ] All tool inputs and outputs validate against pinned schemas.
[ ] Every result is wrapped in an InteropResult envelope.
[ ] Every evidence item has provenance hash.
[ ] Every payload has data classification.
[ ] Every external text payload has taint evaluation.
[ ] Tenant partitioning is enforced before query execution.
[ ] Vector namespaces are tenant-isolated.
[ ] Field projection removes forbidden fields before model context.
[ ] Public registry connectors are banned in production by default.
[ ] Schema drift fails closed.
[ ] Connector health circuit breakers work.
[ ] Safe degradation never bypasses trust controls.
[ ] Sanitized events reach Layer 8.
[ ] No raw secrets, raw tokens, raw user PII, raw traces, or hidden prompts are logged.
[ ] Contract, security, transport, tenant isolation, chaos, and schema drift tests pass.
```

---

## 31. Final Layer 7 Mental Model

Layer 7 is the verified socket and evidence layer.

It does not decide what the agent wants.  
It does not decide whether the action is allowed.  
It does not decide how to reason about the data.  
It does not execute unsafe code.  
It does not render the result.  
It does not evaluate success.

It does exactly this:

```text
verified request in
verified connector out
tenant-safe data back
provenance attached
taint labeled
schema validated
events emitted
```

That is the production boundary.
