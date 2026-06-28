# AGENTS.md

Root instructions for the Website Auditor + Content Writing Agent.

Applies to the whole repo unless a nested `AGENTS.md` overrides it.

---

## 1. Priority

When instructions conflict:

```text
1. System/security/platform rules
2. This AGENTS.md
3. Go style files for Go code
4. Current layer specs
5. Architecture references
6. Audit-agent and content-agent behavior docs
7. Local package conventions
```

Specs define ownership and contracts. Architecture docs define topology. Agent docs define product behavior, tools, skills, and workflows. Go style files define how Go code is written.

---

## 2. Mandatory Go Pre-Read

Before any Go coding, read:

```text
docs/go-style/index.md
docs/go-style/guide.md
docs/go-style/decisions.md
docs/go-style/best-practices.md
```

Priority:

```text
guide.md > decisions.md > best-practices.md > index.md
```

Go style overrides Go snippets in specs. Preserve the spec contract, but write idiomatic Go.

Do not write Go code if these files are missing or unreadable.

---

## 3. Current Specs

Before changing a layer, read that layer spec and adjacent layer specs:

```text
specs/Layer 2 specs.md
specs/Layer 3 specs.md
specs/Layer 4 specs.md
specs/Layer 5 specs.md
specs/Layer 6 specs.md
specs/Layer 7 specs.md
specs/layer 8 specs.md
```

Do not use older generated spec filenames as source of truth.

---

## 4. Architecture References

Use these for topology, cross-layer alignment, and zero-overlap checks:

```text
docs/architecture/Architecture AEO- SEO Agent.txt
docs/architecture/Frontend + Backend.txt
docs/architecture/Layer 2.txt
docs/architecture/Layer 3.txt
docs/architecture/Layer 4.txt
docs/architecture/Layer 5.txt
docs/architecture/Layer 6.txt
docs/architecture/Layer 7.txt
docs/architecture/Layer 8.txt
```

Read them when a change touches multiple layers, boundaries, A2UI, A2A, MCP, skills, sandboxing, memory, telemetry, or approvals.

If specs and architecture conflict, specs win unless the user asks to update the specs.

---

## 5. Agent Behavior References

### Audit agent docs

Read these before implementing or changing website-audit, SEO, AEO, visibility, analytics, GSC, GA4, citations, sentiment, site health, or audit-tool behavior:

```text
docs/audit-agent/All routing.txt
docs/audit-agent/All Skills.txt
docs/audit-agent/All Tools.txt
docs/audit-agent/All workflows.txt
```

### Content agent docs

Read this before implementing or changing topic discovery, brief building, research, SEO planning, page analysis, article planning, drafting, editing, optimization, repurposing, memory/tone, content tools, or content workflows:

```text
docs/content-agent/All-routing-skills-tools-workflows.txt
```

Rules:

```text
Use agent docs for product behavior, tool semantics, workflow order, and skill behavior.
Use specs for layer ownership and contracts.
Use Go style files for Go implementation.
Never expose exact internal tool names, skill paths, workflow IDs, profile IDs, MCP endpoints, or traces to users unless the product explicitly allows it.
```

---

## 6. Layer Ownership

```text
Layer 1: frontend/client shell, device UI, CLI, DOM/native rendering, MFA UX.
Layer 2: intake, safety, prompt-injection blocking, intent enums, context sanitization, approval metadata, tool authorization, outbound redaction.
Layer 3: orchestration, workflow routing, profile selection, DAG planning, short-lived state, proposed tool requests, skill activation requests, presentation intents.
Layer 4: skill registry, SKILL.md, references, assets, scripts-as-artifacts, manifests, skill validation, skill eval definitions.
Layer 5: A2UI/A2A presentation, Agent Cards, approval cards, Vibe Diff, canvas/brief/chat/dashboard contracts, user event normalization.
Layer 6: sandboxed execution, filesystem controls, egress controls, JIT credentials, package controls, runtime budgets, quarantine execution.
Layer 7: MCP transport, connectors, APIs, RAG/retrieval, evidence packets, memory data-plane operations, tenant-isolated data plumbing.
Layer 8: traces, AgBOM, drift/trust scoring, SecOps analytics, eval scoring, governance, correction mining, improvement recommendations.
```

No layer may silently absorb another layer's job.

---

## 7. Required Paths

Tool execution:

```text
Layer 3 proposes -> Layer 2 authorizes -> Layer 6 executes -> Layer 7 handles connector/data access when needed -> Layer 8 observes.
```

Skill loading:

```text
Layer 3 requests -> Layer 4 validates/loads -> Layer 6 executes scripts only after Layer 2 authorization -> Layer 8 evaluates.
```

Presentation:

```text
Layer 3 creates presentation intent -> Layer 5 builds UI contract -> Layer 1 renders -> Layer 8 observes.
```

Memory/tone update:

```text
Layer 2 validates approval -> Layer 3 plans -> Layer 5 presents decision UI -> Layer 7 performs data-plane operation -> Layer 8 audits.
```

Quarantine:

```text
Layer 8 decides/recommends -> Layer 6 executes.
```

Never bypass these paths.

---

## 8. Zero-Overlap Hard Stops

```text
Layer 2 must not choose workflows or execute tools.
Layer 3 must not classify raw intent, authorize tools, read SKILL.md, render UI, execute tools, or connect MCP.
Layer 4 must not execute scripts/tools, choose workflows, connect MCP, render UI, or score evals.
Layer 5 must not authorize tools, choose workflows, execute tools, connect MCP, persist memory, or score evals.
Layer 6 must not classify intent, authorize tools, choose workflows, connect MCP as policy owner, render UI, or score evals.
Layer 7 must not authorize tools, choose workflows, load skills, render UI, execute arbitrary code, or score evals.
Layer 8 must not authorize tools, choose workflows, render UI, execute tools/connectors, mutate memory/canvas/brief, or execute quarantine.
```

If ownership is unclear, stop and identify the owning layer.

---

## 9. Security and Privacy

Fail closed by default.

Never store, log, render, or expose:

```text
raw prompts
hidden chain-of-thought
system/developer prompts
secrets, API keys, OAuth tokens, cookies, passwords, private keys
raw selected text
raw canvas body
raw brief body
raw memory/tone document body
raw PII
exact internal tool inventory
workflow IDs in user-facing text
profile IDs in user-facing text
MCP endpoints in user-facing text
skill file paths in user-facing text
raw AgBOM
raw traces
```

Use safe product-level summaries in user-facing text.

---

## 10. Go Rules Summary

Read the Go style files first. Then apply:

```text
Run gofmt.
Use MixedCaps/mixedCaps.
Use standard initialisms: ID, URL, HTTP, API, JSON, XML, DB, SQL, OAuth, JWT.
Use short, lowercase, meaningful package names.
Avoid util/utils/helper/helpers/common/misc/manager/models.
Pass context.Context first for cancellable, remote, blocking, time-bound, or trace-related work.
Do not store context in structs or option structs.
Return error last.
Handle errors deliberately.
Use lowercase error strings without punctuation.
Do not panic for ordinary failures.
Prefer early returns.
Do not create interfaces before a real need exists.
Prefer small interfaces owned by the consumer.
Avoid unbounded goroutines; make cancellation explicit.
Do not log protected data.
Use table-driven tests when useful.
Call t.Helper() in test helpers.
Add success, failure, and boundary tests.
```

Wire names may remain snake_case through tags:

```go
type IntakeDecision struct {
    TraceID string `json:"trace_id"`
}
```

---

## 11. Production Defaults

Production components should include:

```text
schema validation
fail-closed behavior
structured errors
context cancellation
timeouts
bounded retries
input limits
tenant isolation when data is involved
idempotency where retries are possible
redaction
sanitized events through the owning path
unit tests
negative tests
boundary tests
```

Avoid:

```text
silent fallbacks
best-effort security
global mutable state
unbounded queues/goroutines
implicit network access
implicit filesystem access
hidden background work
```

---

## 12. Before Handoff

Verify:

```text
I read Go style files before Go work.
I read the owning layer spec and adjacent specs when needed.
I read architecture references for cross-layer changes.
I read audit-agent docs for audit/SEO/AEO behavior changes.
I read content-agent docs for content workflow/tool/skill changes.
I preserved layer ownership.
I preserved public JSON/YAML wire contracts.
I ran gofmt for Go changes.
I added or updated tests.
I did not expose protected metadata.
I did not add unapproved tool, network, filesystem, MCP, memory, UI, runtime, or eval behavior.
```

Run relevant commands:

```text
gofmt
go test ./...
go test -race ./...   # when concurrency changed
go vet ./...          # when public APIs or tricky code changed
schema/config validation
boundary tests
```

If a command cannot run, report exactly what failed and why.

---

## 13. Final Rule

This repo is contract-first, fail-closed, privacy-preserving, Go-style-enforced, and layer-isolated.

For Go code:

```text
read Go style files
read current specs
read architecture references when boundaries/cross-layer behavior are involved
read audit/content agent docs when product behavior is involved
implement without overlap
```
