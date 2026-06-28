---
name: content-strategy
description: Builds content roadmaps, editorial plans, topic clusters, and gap priorities around traffic, authority, visibility, audience, and business goals. Use when the user needs a coordinated content program. Do NOT use for writing a single draft, metadata-only work, or technical site audits.
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
    - content_strategy
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - content_strategy_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Content Strategy Expert

## Trigger Conditions

- User asks for a content strategy
- User asks what to write about
- User wants a content plan or editorial calendar
- User wants to build topical authority
- User asks how to grow organic traffic with content
- User wants a content roadmap
- User asks how to structure content pillars
- User wants to identify content gaps vs competitors

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Understand the strategy goal
- Identify what the user is trying to achieve:
    * grow organic traffic
    * build topical authority
    * improve AI visibility and citations
    * capture a specific audience segment
    * fill content gaps vs competitors
    * build a content calendar
    * create a pillar and cluster structure
- If unclear, ask one clarifying question:
    "What is the primary goal of your content strategy -
     traffic growth, topical authority, or
     competitive positioning?"

### STEP 2 - Gather context
- Use quickContext for brand/domain positioning
- Identify:
    * current niche and audience
    * existing content strengths
    * competitor content landscape
    * domain authority level
    * publishing capacity (how much can they publish?)
- Pull GSC data to identify:
    * current ranking topics
    * queries with high impressions but low clicks
    * topics already showing traction
- Pull keyword data to identify:
    * topic clusters with opportunity
    * low competition entry points
    * high volume targets for later

### STEP 3 - Map the topical landscape
- Identify the core topic universe for the niche
- Group topics into:
    PILLAR TOPICS:
      - broad, high-volume, high-authority topics
      - require comprehensive long-form coverage
      - example: "how search engines work"

    CLUSTER TOPICS:
      - specific subtopics supporting the pillar
      - lower volume, lower difficulty
      - example: "what is crawl budget"
                 "how does googlebot work"
                 "what is a canonical tag"

    QUICK WIN TOPICS:
      - low difficulty, decent volume
      - can be published fast
      - build early authority signals

    COMPETITOR GAP TOPICS:
      - topics competitors rank for
      - user does not cover yet
      - represent missed traffic opportunities

### STEP 4 - Build the content architecture
- Design a pillar and cluster model:
    * one pillar page per core topic
    * 5-10 cluster pages per pillar
    * internal links connecting cluster to pillar
    * cluster pages link back to pillar
- Assign each topic a:
    * content type (guide, tutorial, comparison, FAQ)
    * target keyword
    * estimated difficulty
    * estimated traffic potential
    * priority level

### STEP 5 - Build the editorial calendar
- Sequence content in publishing order:
    MONTH 1:
      - quick wins first to build early traffic
      - foundational explainer pages
      - high-intent low-difficulty topics

    MONTH 2-3:
      - cluster pages around first pillar
      - competitor gap topics
      - supporting FAQ and how-to content

    MONTH 4+:
      - pillar pages once cluster authority builds
      - higher difficulty targets
      - content refresh of early pieces if needed

- Assign realistic publishing cadence:
    * 1-2 posts per week for new sites
    * 3-4 posts per week for established sites
    * quality over volume always

### STEP 6 - Optimize for AI visibility
- Identify topics where AI systems frequently cite sources
- Prioritize content types that get cited:
    * definition guides
    * step-by-step tutorials
    * comparison pages
    * FAQ content
    * statistical roundups
- Build answer-ready blocks into every piece
- Structure content so AI systems can extract clear answers
- Target topics where competitors lack strong AI-friendly content

### STEP 7 - Validate mentally before returning
Check for:
  - clear pillar and cluster architecture
  - realistic publishing cadence
  - data-backed topic prioritization
  - mix of quick wins and long-term targets
  - AI visibility considerations included
  - competitor gap coverage

### STEP 8 - Build output
Present as:
  1. Strategic overview
     (goal, audience, core topic universe)
  2. Pillar and cluster map
     (pillars with supporting cluster topics)
  3. Quick wins list
     (publish first for early traction)
  4. Editorial calendar
     (month by month publishing sequence)
  5. AI visibility content priorities
     (topics most likely to get AI citations)
  6. ONE follow-up question if needed

### GUARDRAILS:
- Never build a strategy without keyword data
- Never ignore competitor content landscape
- Never recommend publishing volume over quality
- Never skip the pillar and cluster architecture
- Never build a calendar without prioritizing quick wins first
- Always connect content topics to real search demand
- Always include AI visibility as a strategic dimension
- Always sequence content to build authority progressively

## Purpose

Provide procedural guidance to plan a content portfolio aligned to audience and business goals.

## When to use

- Use when the authorized intent is `seo_planning` and the request is to plan a content portfolio aligned to audience and business goals.

## When NOT to use

- Do not use when the request belongs to `strategic_intelligence`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `content_strategy_report`

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
