---
name: robots-txt-generation
description: Generates or reviews robots.txt rules for search and AI crawlers using verified public-path and access requirements. Use when the user wants crawler-access guidance or a robots.txt file. Do NOT use for sitemaps, llms.txt, authentication, or blocking private data by itself.
version: 1.0.0
owner_team: content_platform
tier: draft
risk_class: medium
compatible_profiles:
    - content_execution_guard
compatible_intents:
    - seo_planning
allowed_modes:
    - write
    - edit
    - optimize
capability_tags:
    - robots_txt_generation
declared_action_classes:
    - read_brand_context
    - canvas_write
output_contracts:
    - robots_txt_generation_draft
    - quality_summary
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Robots.txt and AI Crawler Access Expert

## Trigger Conditions

- User asks for robots.txt
- User asks about blocking or allowing crawlers
- User asks about AI bots like GPTBot, ClaudeBot, PerplexityBot
- User wants to control crawling behavior
- User wants to reference sitemaps in robots.txt
- User wants to improve discoverability without overexposing the site
- User asks how robots.txt affects search or AI systems

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Understand the site and goal
- Identify whether the user wants:
    * to allow search engine crawling
    * to block certain paths
    * to allow or block AI crawlers
    * to expose sitemap locations
    * to protect private or low-value sections
- If unclear, ask one clarifying question:
    "What are you trying to allow or block?"

### STEP 2 - Gather context
- Use quickContext for brand/domain context when available
- Inspect existing site structure if a URL is provided
- Determine if the site is:
    * public content site
    * SaaS
    * e-commerce
    * documentation hub
    * local business site
- Note any crawl-sensitive areas:
    * /admin/
    * /checkout/
    * /account/
    * /search/
    * staging or test environments

### STEP 3 - Decide crawler policy
- Separate search engine bots from AI bots
- Determine if the user wants:
    * full access
    * selective blocking
    * selective allowance
- Use conservative defaults:
    * allow important content
    * block private or duplicate areas
    * keep sitemap accessible
- Do not block core content unless the user explicitly wants that

### STEP 4 - Build robots.txt
- Output valid robots.txt syntax
- Include user-agent groups carefully
- Add Allow / Disallow rules as needed
- Include sitemap declarations when appropriate
- If AI crawler rules are requested, include them explicitly
- Keep it simple and readable
- Avoid conflicting directives

### STEP 5 - Optimize for SEO and AI discovery
- Make sure important content remains crawlable
- Avoid accidentally blocking CSS, JS, or key assets
- Ensure robots.txt supports:
    * search engine indexing
    * AI crawler access policy
    * sitemap discovery
    * clean crawl paths

### STEP 6 - Validate mentally before returning
Check for:
  - syntax correctness
  - no rule conflicts
  - no accidental blocking of valuable content
  - proper sitemap URL formatting
  - crawler names match intended bots

### STEP 7 - Build output
Present as:
  1. Recommended robots.txt content
  2. Short explanation of what it blocks or allows
  3. Notes on any risks or tradeoffs
  4. Optional next step for testing or deployment
  5. ONE follow-up question if needed

### GUARDRAILS:
- Never block core content by accident
- Never assume all bots should be treated the same
- Never create conflicting allow/disallow rules
- Never forget sitemap declarations when they matter
- Always explain the impact on SEO and AI crawling

## Purpose

Provide procedural guidance to prepare a safe robots.txt proposal without blocking required crawling.

## When to use

- Use when the authorized intent is `seo_planning` and the request is to prepare a safe robots.txt proposal without blocking required crawling.

## When NOT to use

- Do not use when the request belongs to `sitemap_generation`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `robots_txt_generation_draft`
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
