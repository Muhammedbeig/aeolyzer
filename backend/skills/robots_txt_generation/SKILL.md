---
name: robots-txt-generation
description: |
  Generates or reviews robots.txt rules for search and AI crawlers using verified public-path and access requirements. Use when the user wants crawler-access guidance or a robots.txt file. Do NOT use for sitemaps, llms.txt, authentication, or blocking private data by itself.
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
