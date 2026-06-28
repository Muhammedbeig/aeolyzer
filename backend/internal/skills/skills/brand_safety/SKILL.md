---
name: brand-safety
description: 'Rules for brand-safe content creation: avoiding competitor references, positioning the user''s brand as authority, and using outbound links responsibly. Use when writing any article or content piece to ensure it serves the user''s brand rather than elevating competitors. Do NOT use for neutral competitor research, technical audits, or private working notes.'
version: 1.0.0
owner_team: content_platform
tier: read
risk_class: low
compatible_profiles:
    - content_collaborator
compatible_intents:
    - draft_article
allowed_modes:
    - plan
    - read
capability_tags:
    - brand_safety
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - brand_safety_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Brand Safety

## BRAND AWARENESS: THE USER'S IMPLICIT INTENT

The user is publishing content to build authority for THEIR brand. Every piece should serve that goal,
even when the tone is educational rather than promotional. This doesn't mean being salesy; it means
being strategically aware of who benefits from the narrative you construct.

**IMPORTANT:** You MUST call getSourcesInsights early in the session to get the exact competitor domain
blocklist. This gives you a definitive list (not guesses) of which domains are competitors.

### Rules

- NEVER cite, link to, or reference any domain from the competitor blocklist (returned by
  getSourcesInsights). This is an absolute rule. Not as sources, not as references, not as examples.
- NEVER use a direct competitor as the hero case study. If a competitor has useful data, frame it as
  industry context or cautionary evidence, not as a success story the reader should admire or emulate.
  BAD: "[Competitor] built a framework that tripled their market share. Here's how they did it."
  GOOD: "Industry data shows market share can shift dramatically: one brand went from 13% to 32%
  in a quarter by focusing on bottom-funnel signals."
- If the user's brand offers a solution in the space being discussed, the article should naturally
  position the user as a credible authority. This happens through: the depth of analysis, the quality
  of insight, and the implicit message that "the people writing this understand this problem deeply."
- You have access to the project context (brand name, domain, working memory). Use it to understand
  what space the user operates in. If the article topic overlaps with their product/service category,
  be especially careful about elevating competitors.
- When competitor data IS useful (and it often is), use it one of these ways:
  1. Anonymize: "One major platform found that..." (when the specific name isn't essential)
  2. Frame as industry trend: "Across the industry, brands investing in this approach saw..."
  3. Use as contrast/cautionary: "Even [competitor] discovered their results weren't translating
     to revenue, proving that surface metrics alone aren't enough."
- For outbound links, ONLY use domains from the authoritySources list (from getSourcesInsights)
  or well-known institutional sources (.gov, .edu, major publications). Aim for 3-4 authority
  site links per article in similar topic neighborhoods but NOT direct competitors.
- The user's brand doesn't need to be mentioned in every article. But the article should never
  make a reader think "I should go buy [competitor's product]"; that's a strategic failure.
- When in doubt, ask yourself: "Would the user be proud to publish this under their brand?"

## Purpose

Provide procedural guidance to review content for competitor promotion and brand-safety risks.

## When to use

- Use when the authorized intent is `draft_article` and the request is to review content for competitor promotion and brand-safety risks.

## When NOT to use

- Do not use when the request belongs to `competitor_intelligence`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `brand_safety_report`

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
