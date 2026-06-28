---
name: competitor-intelligence
description: Analyzes competitor keywords, content, backlinks, traffic, search visibility, and market positioning to identify actionable gaps. Use when the user wants competitive comparisons or an outranking strategy. Do NOT use for single-site audits, content drafting, or unsupported competitor claims.
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
    - competitor_intelligence
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - competitor_intelligence_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Competitor Intelligence

## Trigger Conditions

- User asks to analyze competitors
- User wants to find competitor keywords
- User wants to identify competitive gaps
- User wants to understand market positioning
- User asks what competitors are doing better
- User asks which competitors rank for terms they do not
- User wants to find content gaps vs competitors
- User asks how to outrank a specific competitor
- User wants traffic estimates for competitor domains

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Identify scope
- Confirm domain from project context (searchenginebasics.io)
- Confirm competitor list from project context:
    * moz.com
    * ahrefs.com
    * searchenginejournal.com
- Determine what the user wants to compare:
    * keyword gaps
    * content gaps
    * backlink gaps
    * traffic estimates
    * SERP positioning
    * AI visibility gaps
    * brand positioning
- If scope is unclear, ask one question:
    "Do you want to compare keywords, content, backlinks, or AI visibility?"

### STEP 2 - Pull competitor domain data
- Use analyzeCompetitors for the user's domain
- Pull traffic estimates using getTrafficData for each competitor
- Identify:
    * estimated monthly traffic per competitor
    * number of ranking keywords
    * domain authority signals
    * top traffic-driving pages
    * primary content categories

### STEP 3 - Identify keyword gaps
- Use getKeywords for the primary topic cluster
- Cross-reference competitor ranking keywords
- Identify:
    * keywords competitors rank for that user does not
    * keywords user ranks for that competitors do not
    * keywords all competitors rank for (high priority targets)
    * keywords no competitor ranks for (opportunity gaps)
- Filter by:
    * difficulty appropriate for user's current domain authority
    * search volume worth targeting
    * intent alignment with user's content model

### STEP 4 - Identify content gaps
- Use webSearch to find competitor content that ranks well
- Compare against user's existing content
- Identify:
    * topics competitors cover that user does not
    * content formats competitors use that user does not
    * content depth differences
    * freshness gaps
    * E-E-A-T signal differences
    * structured data usage differences

### STEP 5 - Analyze SERP positioning
- Use analyzeSERP for primary keywords
- Identify:
    * where competitors appear in top 10
    * what content types dominate each SERP
    * which SERP features competitors own
    * where user currently appears
    * realistic ranking targets based on current authority

### STEP 6 - Analyze AI visibility gaps
- Use getVisibilityReports to check user's current AI visibility
- Use getCitationSources to see which domains AI platforms cite
- Cross-reference with competitor domains:
    * which competitors appear most in AI responses
    * which AI platforms cite each competitor
    * what content types get cited for each competitor
    * where user is missing vs competitors in AI responses
- Identify:
    * AI citation gap by platform
    * topics where competitors dominate AI responses
    * content angles that would close the AI visibility gap

### STEP 7 - Analyze competitor content strategy
Use webSearch to find:
  - competitor top-performing articles by traffic signals
  - competitor content publishing frequency
  - competitor content formats (guides, tools, glossaries, studies)
  - competitor internal linking patterns
  - competitor backlink acquisition patterns
  - competitor featured snippet ownership

### STEP 8 - Identify competitive advantages and weaknesses
For each competitor, flag:

  WHERE COMPETITOR IS STRONGER:
    - higher domain authority
    - more content depth
    - more backlinks
    - better AI citation rate
    - more SERP features owned

  WHERE USER HAS AN ADVANTAGE:
    - specific topic areas not well covered by competitor
    - audience segments competitor ignores
    - content formats competitor does not produce
    - beginner-friendly angle competitors do not prioritize
    - faster page speed or better technical signals

  WHERE GAPS ARE CLOSEABLE:
    - topics with low difficulty where competitor ranks but user does not
    - content formats easy to replicate with higher quality
    - AI citation sources that link to competitors but not user

### STEP 9 - Build competitive action plan
Prioritize actions using this logic:

  TIER 1 - Immediate wins
    - keywords competitor ranks for with KD under user's threshold
    - content gaps where user has partial coverage to expand
    - AI citation sources linking to competitors but not user

  TIER 2 - Medium-term moves
    - content formats competitors use that user lacks
    - backlink sources competitors share that user can target
    - SERP features competitors own that user can contest

  TIER 3 - Long-term positioning
    - high-difficulty keywords worth building toward
    - domain authority gap closure strategy
    - brand positioning differentiation

### STEP 10 - Validate before returning
Check that:
  - all competitor data is from actual tool output, not assumptions
  - keyword gaps are specific named terms, not categories
  - content gaps reference specific articles or topics
  - AI visibility comparison is included
  - action plan is tiered by impact and difficulty
  - user's actual competitive advantages are stated clearly

### STEP 11 - Build output
Present as:
  1. Competitor landscape overview
     (traffic, authority, keyword count per competitor)
  2. Keyword gap analysis
     (what they rank for that you do not, filtered by difficulty)
  3. Content gap analysis
     (topics and formats missing from your site)
  4. AI visibility gap
     (where competitors dominate AI responses vs you)
  5. Your competitive advantages
     (where you can realistically win)
  6. Prioritized action plan
     (Tier 1, 2, 3 moves)
  7. ONE recommended first move
     (highest ROI action this week)
  8. ONE follow-up question

### GUARDRAILS:
- Never use assumed competitor data without tool verification
- Never recommend targeting KD above user's realistic threshold
- Never skip AI visibility gap analysis
- Never present competitor strengths without pairing with user advantages
- Never give generic "create better content" advice without specific gap evidence
- Always tie every recommendation to a named keyword, page, or citation source
- Always separate traditional SEO gaps from AI visibility gaps
- Never treat all competitors as identical - flag where each one is strongest

## Purpose

Provide procedural guidance to compare competitor positioning, visibility, and content evidence.

## When to use

- Use when the authorized intent is `content_research` and the request is to compare competitor positioning, visibility, and content evidence.

## When NOT to use

- Do not use when the request belongs to `brand_safety`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `competitor_intelligence_report`

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
