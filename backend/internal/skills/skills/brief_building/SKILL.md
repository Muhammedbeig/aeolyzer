---
name: brief-building
description: Builds a content brief covering topic, angle, audience, intent, CTA, keywords, length, subtype, and notes. Use before drafting begins or when planning article requirements. Do NOT use for writing sections, researching sources, or publishing content.
version: 1.0.0
owner_team: content_platform
tier: draft
risk_class: medium
compatible_profiles:
    - content_execution_guard
compatible_intents:
    - content_brief
allowed_modes:
    - write
    - edit
    - optimize
capability_tags:
    - brief_building
declared_action_classes:
    - read_brand_context
    - canvas_write
output_contracts:
    - brief_building_draft
    - quality_summary
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Brief Building

Before writing begins, you need a brief. The brief is the contract between you and the user - it defines what you're writing, for whom, and why. Without it, you're guessing.

## MINIMUM FIELDS (must have before writing)

- Topic: what are we writing about?
- Angle: what is the unique perspective or thesis?
- Audience: who is reading this?
- Intent: what should the reader know, feel, or do after reading?
- CTA: where should the reader go next?

## IDEAL FIELDS (collect when possible)

- Keywords: target search terms or search query
- Length: target word count based on content type and goal
- Subtype: opinion | comparison | topic-guide | news-analysis | how-to | alternatives
- Notes: anything specific the user wants included or avoided

## SUBTYPE SELECTION

- opinion: user has a strong POV, wants to take a position
- comparison: evaluating two or more options side by side
- topic-guide: comprehensive coverage of a subject
- news-analysis: breaking down a recent event or development
- how-to: step-by-step instructional content
- alternatives: positioning options against a category leader

## RULES

- Call `updateBrief()` before the first `writeSection()` call - always
- If the user skips brief questions and says "just write it", infer the brief from context and save it anyway
- If topic is clear but angle is missing, suggest 2-3 angle options before writing
- Never start writing without at least: topic, audience, and intent
- The brief is saved and visible alongside the canvas - treat it as a live document
- Update the brief if the user changes direction mid-article

## Purpose

Provide procedural guidance to build a structured content brief from an approved topic.

## When to use

- Use when the authorized intent is `content_brief` and the request is to build a structured content brief from an approved topic.

## When NOT to use

- Do not use when the request belongs to `outline_structure`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `brief_building_draft`
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
