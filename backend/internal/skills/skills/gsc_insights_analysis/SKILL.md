---
name: gsc-insights-analysis
description: Analyzes Search Console queries, pages, clicks, impressions, CTR, positions, and trends to identify search opportunities and quick wins. Use when the user wants ranking, CTR, or performance insights from GSC data. Do NOT use for GA4 behavior analysis, unsupported SEO advice, or content drafting.
version: 1.0.0
owner_team: audit_platform
tier: read
risk_class: low
compatible_profiles:
    - seo_aeo_auditor
compatible_intents:
    - traffic_analysis
allowed_modes:
    - audit
    - read
capability_tags:
    - gsc_insights_analysis
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - gsc_insights_analysis_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# GSC Insights Expert

## Trigger Conditions

- User asks about Google Search Console data
- User wants to understand rankings, clicks, impressions, CTR, or position
- User wants to find keyword opportunities
- User wants to improve CTR from search results
- User asks which queries are close to page 1
- User wants to know what pages or queries are underperforming
- User asks about content gaps in search performance
- User wants to find quick wins from high-impression, low-click queries
- User wants to understand ranking changes over time

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Identify the search goal
- Determine whether the user wants:
    * query performance analysis
    * page performance analysis
    * CTR optimization
    * ranking improvement
    * opportunity discovery
    * trend comparison
    * content gap analysis
- If unclear, ask one question:
    "Are you looking at queries, pages, CTR, or ranking opportunities?"

### STEP 2 - Use the right data window
- Default to the last 28 or 30 days if no range is specified
- If the user asks for movement, compare to the previous equivalent period
- If the user asks for trends, include daily or period-over-period context

### STEP 3 - Analyze query performance
For each query, evaluate:
  - clicks
  - impressions
  - CTR
  - average position
  - trend direction
  - landing page association
  - whether the query is a winner, loser, or opportunity

### STEP 4 - Analyze page performance
For each page, evaluate:
  - total clicks
  - total impressions
  - CTR
  - average position
  - query mix
  - whether the page is underperforming relative to its potential
  - whether it needs content, title, or internal link improvements

### STEP 5 - Identify high-value opportunity types
Flag:
  - high impressions, low CTR queries
  - positions 4-10 queries that can be pushed into top 3
  - positions 11-20 queries that can move to page 1
  - pages with strong impressions but weak clicks
  - queries that are declining but still viable
  - queries missing dedicated landing pages
  - content gaps where search demand exists but coverage is weak

### STEP 6 - Analyze trends
Compare:
  - current period vs previous period
  - winners vs losers
  - rising vs declining queries
  - page-level changes over time
  - whether movements are sitewide or isolated

### STEP 7 - Prioritize by impact
Rank recommendations in this order:

  HIGH PRIORITY:
    - high-impression CTR fixes
    - queries sitting just outside page 1
    - pages with ranking upside and high business value
    - declining queries tied to important content
    - pages with strong search demand but weak click yield

  MEDIUM PRIORITY:
    - content expansion for partial matches
    - internal link support
    - snippet optimization
    - secondary query targeting

  LOWER PRIORITY:
    - low-volume query tweaks
    - tiny position changes with little traffic value
    - cosmetic reporting issues

### STEP 8 - Build specific recommendations
For each recommendation, include:
  - the query or page
  - the observed issue or opportunity
  - why it matters
  - what to change
  - where to change it
  - expected effect on traffic or rankings

### STEP 9 - Validate before returning
Check that:
  - recommendations are tied to actual GSC data
  - clicks, impressions, CTR, and position are treated distinctly
  - opportunity thresholds are clear
  - movement is interpreted in context
  - no generic SEO advice appears without data backing

### STEP 10 - Build output
Present as:
  1. Query performance
     (best, worst, and opportunity queries)
  2. Page performance
     (pages with ranking or CTR issues)
  3. Trend changes
     (what is rising or falling)
  4. Quick wins
     (high-impression, low-CTR, page 2 targets)
  5. Content gaps
     (search demand without strong coverage)
  6. ONE highest-impact fix
     (best move this week)
  7. ONE follow-up question

### GUARDRAILS:
- Never confuse clicks with impressions
- Never confuse CTR with ranking position
- Never treat a page 2 query like a failure when it is a quick win
- Never recommend changes without tying them to observed query/page behavior
- Always look for quick wins before long-term projects
- Always separate query-level and page-level insights

## Purpose

Provide procedural guidance to interpret Search Console query, page, and indexing evidence.

## When to use

- Use when the authorized intent is `traffic_analysis` and the request is to interpret Search Console query, page, and indexing evidence.

## When NOT to use

- Do not use when the request belongs to `ga4_analysis`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `gsc_insights_analysis_report`

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
