---
name: anti-redundancy
description: Use when writing any multi-section content to ensure every paragraph earns its place. Covers rules for eliminating repeated points across sections, a pre-section redundancy check, and common redundancy traps to avoid. Trigger for long-form articles, SEO content, reports, or any structured writing where section-by-section discipline is needed. Do NOT use for single-section copy, metadata-only tasks, or unrelated editing.
version: 1.0.0
owner_team: content_platform
tier: draft
risk_class: medium
compatible_profiles:
    - content_execution_guard
compatible_intents:
    - optimize_content
allowed_modes:
    - write
    - edit
    - optimize
capability_tags:
    - anti_redundancy
declared_action_classes:
    - read_brand_context
    - canvas_write
output_contracts:
    - anti_redundancy_draft
    - quality_summary
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Anti-Redundancy

## ANTI-REDUNDANCY (CRITICAL)

- NEVER make the same point in two different sections. Each section must advance the argument.
- If you established a distinction (e.g. "Protocol A is search-driven, Protocol B is conversational") in section 2, do NOT re-explain it in sections 4, 5, and 6. Reference it briefly ("as discussed above") but add NEW analysis, not the same observation in different words.
- Before writing each section, mentally review what you've already said. If a point was covered, skip it.
- Common redundancy traps to avoid:
  - Restating the core thesis in every section introduction
  - Re-explaining what two things are after the explainer sections
  - Repeating the same shared characteristic of two things in 3+ sections
  - Summarizing the entire article in the conclusion instead of advancing it

## REDUNDANCY CHECK BEFORE EACH SECTION

Before writing any section, ask:

1. Have I made this point already?
2. If yes — can I advance the argument instead of restating it?
3. If I can't advance it — should this section exist at all?

The goal: every paragraph earns its place by adding something the previous paragraph didn't say.

## Purpose

Provide procedural guidance to remove repeated ideas while preserving distinct analysis.

## When to use

- Use when the authorized intent is `optimize_content` and the request is to remove repeated ideas while preserving distinct analysis.

## When NOT to use

- Do not use when the request belongs to `formatting_rules`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `anti_redundancy_draft`
- `quality_summary`

## Quality gates

- Keep claims tied to supplied evidence.
- Separate facts, inferences, and recommendations.
- Reject protected metadata and unsupported certainty.
- Confirm the output matches the declared contract.

## Boundary rules

This skill provides procedural guidance only.

It must not:
- classify raw user intent
- choose workflows or agents
- authorize or execute tools or scripts
- connect to MCP servers or external APIs
- read or write memory documents directly
- mutate canvas, brief, chat, dashboard, or UI state
- store telemetry or score evaluations
- expose internal identifiers, endpoints, traces, credentials, or protected metadata

## Resources

No runtime references, assets, or scripts are declared for this version.

## Failure behavior

Fail closed and return a safe request for the missing context, evidence, mode, or approval. Never fabricate data or silently broaden scope.
