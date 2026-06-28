---
name: content-refresh-strategy
description: Plans evidence-based updates for aging or declining content using performance trends, current coverage, and competitor changes. Use when the user wants to recover rankings, refresh pages, or prioritize decaying content. Do NOT use for net-new drafts, technical audits, or metadata-only changes.
version: 1.0.0
owner_team: content_platform
tier: read
risk_class: low
compatible_profiles:
    - content_collaborator
compatible_intents:
    - optimize_content
allowed_modes:
    - plan
    - read
capability_tags:
    - content_refresh_strategy
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - content_refresh_strategy_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Content Refresh Expert

## Trigger Conditions

- User wants to update old content
- User asks why content is losing rankings
- User wants to maintain or recover traffic
- User asks about content decay
- User wants to refresh existing articles
- User asks which pages need updating
- User wants to improve underperforming content
- User asks how to make old content relevant again

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Understand the refresh goal
- Identify whether the user wants:
    * a single page refreshed
    * a site-wide content decay audit
    * a prioritized refresh roadmap
    * specific sections updated
    * a full rewrite vs light update
- If unclear, ask one clarifying question:
    "Which page or pages are you looking to refresh?"

### STEP 2 - Gather performance data
- Use quickContext for brand/domain context
- Pull GSC data to identify:
    * pages with declining clicks over time
    * pages with high impressions but falling CTR
    * pages dropping from position 1-3 to 4-10
    * pages that ranked well historically but lost ground
- Pull GA4 data to identify:
    * pages with declining sessions
    * pages with high bounce rates
    * pages with low engagement time
- If a specific URL is provided, inspect the page content directly

### STEP 3 - Diagnose the decay cause
Map the performance drop to likely causes:

  CONTENT STALENESS:
    - outdated statistics or data
    - references to old tools or platforms
    - missing recent developments
    - outdated screenshots or examples

  COMPETITIVE DISPLACEMENT:
    - newer, better content published by competitors
    - competitor pages now have more depth
    - competitor pages have better structure or format

  INTENT DRIFT:
    - search intent for the keyword has shifted
    - Google now prefers a different content format
    - SERP features have changed (PAA, snippets, etc.)

  TECHNICAL DECAY:
    - page speed has degraded
    - internal links pointing to it have been removed
    - external links have been lost
    - page was accidentally noindexed or decanonized

  AUTHORITY DECAY:
    - backlinks have been lost or devalued
    - domain authority of competing pages has grown
    - E-E-A-T signals have weakened

### STEP 4 - Score refresh priority
Apply this prioritization logic:

REFRESH IMMEDIATELY if:
  - page was top 3 and dropped to 4-10
  - page has high impressions but declining CTR
  - page covers a topic with recent major developments
  - page has thin content vs current top-ranking competitors

REFRESH SOON if:
  - page is on page 2 with decent impressions
  - page has outdated statistics or examples
  - page lacks structured data or answer blocks
  - page has not been updated in 12+ months

MONITOR ONLY if:
  - page is stable with no significant decline
  - page covers evergreen topic with no recent changes
  - page has low traffic but consistent performance

CONSIDER CONSOLIDATING if:
  - page is thin and covers same topic as another page
  - page has near-zero traffic and impressions
  - page is cannibalizing a stronger page

### STEP 5 - Build the refresh plan
For each page to refresh, specify:
  - what to update (statistics, examples, sections)
  - what to add (new sections, FAQs, answer blocks)
  - what to remove (outdated information, thin sections)
  - what format changes to make (restructure, add tables, etc.)
  - whether to update the publish date
  - whether to update the title tag and meta description
  - whether to add or update schema markup
  - internal linking updates needed

### STEP 6 - Prioritize by effort vs impact
Classify each refresh as:
  LIGHT UPDATE:
    - update statistics and dates
    - add one or two new sections
    - update title and meta
    - estimated time: 1-2 hours

  MODERATE REFRESH:
    - restructure content format
    - add FAQ or answer blocks
    - add schema markup
    - update internal links
    - estimated time: 3-5 hours

  FULL REWRITE:
    - content is fundamentally outdated
    - intent has shifted completely
    - competitor content is significantly better
    - estimated time: 1-2 days

### STEP 7 - Validate mentally before returning
Check for:
  - data-backed decay diagnosis
  - clear prioritization logic
  - specific actionable refresh steps
  - effort vs impact balance
  - no vague advice like "update your content"

### STEP 8 - Build output
Present as:
  1. Decay diagnosis summary
     (what caused the drop and why)
  2. Pages prioritized for refresh
     (with decay signal and priority level)
  3. Refresh plan per page
     (what to update, add, remove, restructure)
  4. Effort classification
     (light update vs moderate refresh vs full rewrite)
  5. ONE follow-up question if needed

### GUARDRAILS:
- Never recommend refreshing without diagnosing the decay cause first
- Never suggest a full rewrite when a light update will do
- Never ignore intent drift as a possible cause
- Never update publish dates without making real content changes
- Never skip checking GSC and GA4 data before prioritizing
- Always connect the refresh recommendation to a specific performance signal
- Always separate quick light updates from heavy rewrites

## Purpose

Provide procedural guidance to prioritize existing pages for evidence-backed content refreshes.

## When to use

- Use when the authorized intent is `optimize_content` and the request is to prioritize existing pages for evidence-backed content refreshes.

## When NOT to use

- Do not use when the request belongs to `long_form_content_audit`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `content_refresh_strategy_report`

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
