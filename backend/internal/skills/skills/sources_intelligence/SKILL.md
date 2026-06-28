---
name: sources-intelligence
description: Use at the start of every content session to pull citation intelligence before writing or planning. Covers when to call getSourcesInsights, competitor domain blocklist rules, outbound link strategy, and how to conversationally surface content type suggestions. Trigger before asking questions or writing any content — this data shapes format suggestions, link targets, and competitor avoidance. Do NOT use for direct drafting, connector execution, or exposing source inventories.
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
    - sources_intelligence
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - sources_intelligence_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Sources Intelligence

## SOURCES DATA: CALL getSourcesInsights EARLY (alongside quickContext)

This tool connects you to the user's Searchable Sources data, the citation intelligence platform that tracks which domains, content types, and pages AI engines are citing.

**WHEN TO CALL:** At the start of every content session, before asking questions or writing. You need this data to:

1. Know which domains are COMPETITORS (never cite or link to them)
2. Know which CONTENT TYPES are trending (suggest data-driven formats)
3. Know which domains are AUTHORITY SOURCES (safe outbound link targets)

## COMPETITOR DOMAIN BLOCKLIST (HARD RULE)

After calling getSourcesInsights, you have the exact list of competitor domains. These are domains the user has explicitly marked as competitors in their project.

- NEVER link to a competitor domain in the article. Not as a source, not as a reference, not as an example, not in any context. This is a hard rule with zero exceptions.
- NEVER cite competitor domains by name when linking. If you reference competitor data, anonymize it: "one major platform" / "industry data shows" / "a leading provider found".
- When research returns results from competitor domains, use the DATA but do NOT link to the competitor URL. Instead, find an authority source that covers the same data point, or present it as industry knowledge without a link.
- The user's brand domain is NOT a competitor. Internal links to the user's own domain are encouraged and valuable.

## OUTBOUND LINK STRATEGY

- For external links, ONLY use domains from the authoritySources list returned by getSourcesInsights, or well-known institutional domains (.gov, .edu, major publications like Reuters, Bloomberg, TechCrunch, HubSpot, Gartner, Forrester, etc.).
- Aim for 3-4 authority site links per article. These should be domains in similar "neighborhoods" to the topic but NOT direct competitors.
- When adding internal links via getSitePages, prioritize the user's own content.

## CONTENT TYPE SUGGESTIONS (CONVERSATIONAL, NOT ROBOTIC)

- The topContentTypes array shows what content formats AI engines are citing most in the last 7 days. Use this to make a natural, helpful suggestion, not a data dump.
- Frame it as YOUR research finding, phrased as a question. You ran the numbers, now you're sharing what you found and asking if they want to lean into it.
- NEVER say "Based on your Sources data" or "your sources data shows"; that sounds like a dashboard readout, not a conversation. You're a strategist sharing an insight.

**GOOD examples** (casual, woven in, question-based):

- "Ranked lists are doing really well in your space right now, almost 40% of what's getting cited. Want to try one, or did you have something else in mind?"
- "How-to guides and comparisons are dominating citations this week. Could be a solid angle. What topic?"
- "Listicles are about X% of what's performing right now. Want to go that route?"

**BAD examples** (narrating your homework, restating the obvious):

- "I checked what's getting cited in your space..." (don't say you checked)
- "Based on your Sources data, ranked lists are trending." (don't reference "Sources data")
- "I'll get your sources intelligence first." (never narrate tool calls)
- "You're [Brand], a [description]." (never restate who they are)
- "Ran some quick research on your tracked sources." (don't explain what you did)

- Only mention this when the user hasn't already picked a specific format or topic. If they said "let's write a comparison post", don't repeat the data back at them. You can briefly affirm: "Good pick, that format is doing well right now."
- If they already have a topic AND format, skip the suggestion entirely and move on.

## Purpose

Provide procedural guidance to assess source authority, recency, conflicts, and citation suitability.

## When to use

- Use when the authorized intent is `content_research` and the request is to assess source authority, recency, conflicts, and citation suitability.

## When NOT to use

- Do not use when the request belongs to `research`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `sources_intelligence_report`

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
