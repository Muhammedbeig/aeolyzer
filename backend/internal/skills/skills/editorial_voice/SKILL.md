---
name: editorial-voice
description: Use when writing or editing any article, blog post, or long-form content to ensure it has a strong, consistent analytical voice. Covers taking positions, "so what" analysis, direct language patterns, blog post voice specifics, and voice consistency checks. Trigger for SEO articles, blog posts, op-eds, or any content that should analyze rather than just inform. Do NOT use for data extraction, metadata-only tasks, or neutral factual summaries.
version: 1.0.0
owner_team: content_platform
tier: draft
risk_class: medium
compatible_profiles:
    - content_execution_guard
compatible_intents:
    - draft_article
allowed_modes:
    - write
    - edit
    - optimize
capability_tags:
    - editorial_voice
declared_action_classes:
    - read_brand_context
    - canvas_write
output_contracts:
    - editorial_voice_draft
    - quality_summary
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Editorial Voice

## EDITORIAL VOICE

- Take positions. Don't just present information - interpret it. Say what it MEANS.
- After presenting a fact, add "so what" analysis: what does this mean for the reader's business/decisions?
- Use direct analytical language:
  - "Here's what this actually means for [audience]:"
  - "The part most coverage misses:"
  - "The real risk is..."
  - "This matters because..."
  - "What no one is talking about:"
- NEVER hedge more than once per paragraph. Pick a position and defend it.
  - BAD: "This could potentially maybe lead to some changes"
  - GOOD: "This will force decision-makers to rethink their entire approach"
- Distinguish your piece from a press release. Press releases inform. Your article should ANALYZE. For every claim or feature described, ask: "So what? Why should the reader care?"

## BLOG POST VOICE SPECIFICS

- Write from a first-person perspective (I/we). Conversational, personality-forward, relatable.
- The reader should feel like they're hearing from a real person, not a brand.
- Use personal anecdotes, specific experiences, and candid opinions.
- Replace corporate language with direct human language:
  - BAD: "Leverage synergistic cross-functional alignment"
  - GOOD: "Get your teams talking to each other"
- Shorter sentences. Shorter paragraphs. More white space.
- Bold the key takeaway in each section so skimmers get the point.

## VOICE CONSISTENCY

- The voice established in the introduction must carry through to the conclusion.
- If the intro is confident and direct, the body cannot be hedging and passive.
- If the intro uses first person, the body cannot switch to third person.
- Read the first paragraph and the last paragraph back to back. Do they sound like the same writer with the same conviction? If not, revise.

## Purpose

Provide procedural guidance to apply an approved editorial voice without inventing brand preferences.

## When to use

- Use when the authorized intent is `draft_article` and the request is to apply an approved editorial voice without inventing brand preferences.

## When NOT to use

- Do not use when the request belongs to `writing`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `editorial_voice_draft`
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
