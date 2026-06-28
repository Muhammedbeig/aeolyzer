---
name: strategic-intelligence
description: Use when suggesting topics, choosing angles, or evaluating what to write next. Covers the gap map mindset, how to find gaps with available tools, four types of content opportunities (competitive, timing, depth, audience), and how to present opportunities strategically. Trigger for topic selection, content planning, angle choice, and any decision about what to publish or how to differentiate it. Do NOT use for full drafting, technical audits, or metadata-only tasks.
version: 1.0.0
owner_team: content_platform
tier: read
risk_class: low
compatible_profiles:
    - content_collaborator
compatible_intents:
    - seo_planning
allowed_modes:
    - plan
    - read
capability_tags:
    - strategic_intelligence
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - strategic_intelligence_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Strategic Intelligence

Every content suggestion should be informed by competitive intelligence thinking. You don't have a full competitive gap analysis product yet, but you can approximate the strategic framework with the tools you have.

## THE GAP MAP MINDSET

When suggesting topics, angles, or evaluating what to write, think in terms of GAPS:

- What are competitors publishing that the user ISN'T covering?
- What topics are getting traction in the industry with NO good coverage yet?
- What angle on a hot topic has everyone missed?
- What does the user's audience need to know that nobody is explaining well?

## HOW TO FIND GAPS WITH YOUR CURRENT TOOLS

1. quickContext gives you: competitors, industry, brand positioning
2. webSearch "[competitor] blog [topic]" shows what they've published
3. webSearch "[industry] [topic] 2026" shows the current landscape
4. The DELTA between what exists and what's missing = content opportunity

## FOUR TYPES OF CONTENT OPPORTUNITIES (prioritized)

1. **COMPETITIVE GAPS:** Competitor X published about [topic] but missed [angle]. User can own that angle. Highest strategic value.
2. **TIMING GAPS:** Something just happened (launch, regulation, trend) and nobody in the user's space has published a thoughtful take yet. First-mover advantage.
3. **DEPTH GAPS:** Existing coverage is shallow (press release summaries, generic overviews). User can go deeper with original analysis, data, or expert perspective.
4. **AUDIENCE GAPS:** Content exists but it's written for the wrong audience. E.g. technical content exists but nothing for the decision-maker, or vice versa.

## WHEN PRESENTING OPPORTUNITIES

- Name the gap type: "This is a competitive gap - [Competitor] covered X but missed Y"
- Quantify when possible: "There are 3 articles on this topic but none address [angle]"
- Connect to strategy: "Publishing this positions you as [authority on X] and captures [audience segment] who are currently reading [competitor]'s take"
- Be honest about difficulty: "This is a crowded topic but the [angle] is untouched"

This framework applies to EVERY content decision:

- Topic selection (what gap does this fill?)
- Angle choice (what perspective is missing from existing coverage?)
- CTA design (what action captures the reader into the user's ecosystem?)
- Differentiation (what makes this piece worth reading over what already exists?)

## Purpose

Provide procedural guidance to synthesize market, audience, competitor, and source evidence into strategy.

## When to use

- Use when the authorized intent is `seo_planning` and the request is to synthesize market, audience, competitor, and source evidence into strategy.

## When NOT to use

- Do not use when the request belongs to `content_strategy`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `strategic_intelligence_report`

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
