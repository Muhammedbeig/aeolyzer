---
name: depth-evidence
description: Use when writing or researching any article that needs to go beyond surface-level coverage. Covers second-order effects, competitive context, counter-arguments, specific evidence standards, and data contextualization. Trigger for SEO articles, thought leadership pieces, B2B content, or any writing where depth and credibility are required. Do NOT use for lightweight copy, metadata-only tasks, or formatting-only edits.
version: 1.0.0
owner_team: content_platform
tier: read
risk_class: low
compatible_profiles:
    - content_collaborator
compatible_intents:
    - content_research
allowed_modes:
    - plan
    - read
capability_tags:
    - depth_evidence
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - depth_evidence_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Depth and Evidence

## DEPTH AND EVIDENCE (REQUIRED)

- Go beyond the press release. For every major topic, address:
  - What's missing from the standard narrative?
  - What second-order effects will this cause?
  - What's the competitive context? (Who else is doing this? What's their angle?)
  - What are the counter-arguments or risks?
- Reference SPECIFIC numbers, dates, companies, and quotes.
  - BAD: "Many companies are adopting this approach"
  - GOOD: "Over 1 million businesses, including [specific named examples], have adopted [specific thing] since [specific date]"
- Include counter-arguments. If you argue for X, acknowledge the strongest case against X. This builds credibility, not weakness.
- If the plan includes sources with key data, WEAVE that data into the prose. Don't list facts. Contextualize them: "[Stat] - which means [interpretation for reader]."

## Purpose

Provide procedural guidance to check whether major claims have sufficient credible evidence and analysis.

## When to use

- Use when the authorized intent is `content_research` and the request is to check whether major claims have sufficient credible evidence and analysis.

## When NOT to use

- Do not use when the request belongs to `sources_intelligence`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `depth_evidence_report`

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
