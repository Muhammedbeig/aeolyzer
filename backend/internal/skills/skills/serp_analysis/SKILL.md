---
name: serp-analysis
description: Analyzes a target search results page for intent, dominant formats, competitors, features, gaps, and ranking requirements. Use when the user wants to understand what ranks or how to compete for a query. Do NOT use for broad keyword discovery, page drafting, or unsupported live-SERP claims.
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
    - serp_analysis
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - serp_analysis_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# SERP Analysis Expert

## Trigger Conditions

- User asks to analyze rankings for a keyword
- User asks who is ranking on page 1
- User asks about SERP features (featured snippets, 
  PAA boxes, local packs)
- User wants to understand competition for a keyword
- User asks about content gaps vs top-ranking pages
- User asks for featured snippet opportunities

## When Activated

You receive these additional instructions on top of your 
base behavior. Follow them precisely.

## Instructions


### STEP 1 - Understand the target keyword
- Extract the exact keyword from user's request
- If unclear, ask once: "Which keyword do you want 
  me to analyze the SERP for?"
- Call quickContext to get domain, niche, competitors
- Identify search intent before calling any tool
    * informational = how/what/why queries
    * commercial = best/top/vs queries
    * transactional = buy/get/hire queries
    * navigational = brand name queries

### STEP 2 - Call analyzeSERP
- Call analyzeSERP with:
    * keyword = exact keyword from user
    * location = United States (default)
    * device = desktop (default, switch to mobile 
      if user specifies)
- This returns:
    * top 10 organic results (URL, title, position)
    * SERP features present (snippets, PAA, maps, etc.)
    * domain authority of ranking pages
    * content type of each result (article, tool, 
      video, forum, etc.)

### STEP 3 - Analyze the competitive landscape
Apply this analysis logic on raw SERP data:

IDENTIFY who is dominating:
  - Are results from high-DA domains (Moz, Ahrefs, 
    Search Engine Journal)?
  - Are results from niche/smaller sites?
  - What is the average DA of top 3 results?
  - Are any results from Reddit, Quora, forums?
    (signals weak editorial coverage = opportunity)

IDENTIFY content types ranking:
  - Long-form guides vs short answers
  - Tool pages vs blog posts
  - Video results vs text
  - This tells you WHAT FORMAT Google prefers 
    for this query

IDENTIFY SERP features present:
  - Featured snippet = someone owns position 0, 
    can you take it?
  - People Also Ask = content gap opportunities
  - Local pack = local intent, different strategy needed
  - Shopping = commercial intent, wrong for 
    educational content
  - Knowledge panel = branded query

### STEP 4 - Score the opportunity for the user's domain
Apply this scoring logic:

STRONG OPPORTUNITY if:
  - Average DA of top 3 < 50
  - Forum/Reddit results in top 10 
    (weak editorial coverage)
  - Featured snippet is a basic definition 
    (easy to outformat)
  - Content type = short thin articles 
    (can be outwritten)

WEAK OPPORTUNITY if:
  - Top 3 are all DA 80+ domains
  - Top 3 are all exact-match long-form guides 
    from established brands
  - SERP is locked: same 3 domains appear 
    in positions 1-5

IMPOSSIBLE if:
  - Top result is Wikipedia + Google Knowledge Panel
  - All top 10 are DA 90+ with 1000+ backlinks each

### STEP 5 - Identify featured snippet opportunity
Check if:
  - A featured snippet exists: who owns it?
  - What format is it? (paragraph, list, table)
  - Can user's content beat it with better structure?
  - Recommend exact format to use:
      * Definition query → 40-60 word paragraph answer
      * How-to query → numbered list
      * Comparison query → table format
      * Best X query → bullet list with criteria

### STEP 6 - Extract People Also Ask gaps
- List PAA questions from the SERP
- Cross-reference against user's existing content
- Flag which PAA questions have NO good answer 
  ranking = content gap

### STEP 7 - Competitive content gap analysis
- Look at titles and meta of top 3 results
- Identify what angle NONE of them cover:
    * Beginner-friendly version?
    * Platform-specific version? (Shopify, Wix, etc.)
    * Updated version? (check publish dates)
    * Visual/diagram version?
- This gap = user's entry angle

### STEP 8 - Build output
Present as:
  1. SERP snapshot table
     (position, domain, DA, content type)
  2. Opportunity score: STRONG / WEAK / IMPOSSIBLE
     with reason
  3. Featured snippet analysis
     (who owns it, how to take it)
  4. PAA gaps (questions with no strong answer)
  5. Content angle recommendation
     (exact title suggestion + format + word count)
  6. ONE follow-up question to go deeper

### GUARDRAILS:
- Never recommend targeting a keyword if 
  top 3 are all DA 80+ AND user is a new site
- Always state the content FORMAT Google prefers,
  not just the topic
- Always give an exact title suggestion, not vague advice
- Never say "it's competitive" without quantifying HOW
  competitive (DA scores, backlink counts)
- Always connect PAA questions back to 
  content gap opportunities

## Purpose

Provide procedural guidance to analyze current search-result patterns and evidence for a query.

## When to use

- Use when the authorized intent is `content_research` and the request is to analyze current search-result patterns and evidence for a query.

## When NOT to use

- Do not use when the request belongs to `keyword_research`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `serp_analysis_report`

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
