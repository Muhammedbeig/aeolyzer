---
name: link-opportunity-discovery
description: Discovers and prioritizes realistic backlink and citation prospects from competing links, cited domains, listicles, directories, and resource pages. Use when the user wants specific outreach targets. Do NOT use for internal links, generic backlink education, or sending outreach.
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
    - link_opportunity_discovery
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - link_opportunity_discovery_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Link Opportunity Finder

## Trigger Conditions

- User needs specific link targets
- User wants to find citation opportunities
- User asks where to get featured
- User wants to find articles missing their brand
- User asks what listicles or guides to target
- User wants broken link opportunities
- User asks which publications cover their topic
- User wants to find geographic or audience gaps in existing content

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Clarify scope
- Confirm domain from project context (user-domain.com)
- Identify primary topic focus:
    * use existing keyword/content context from conversation
    * if unclear, ask one question:
        "What topic or keyword should I focus the search on?"
- Identify target link type if specified:
    * listicle inclusion
    * resource page
    * broken link replacement
    * AI citation source
    * guest post opportunity
    * directory or glossary

### STEP 2 - Find listicles and roundups missing the brand
Use webSearch with these patterns:

  Pattern A - Best-of lists
    "best [topic] [resource type]" -site:[user domain]
    Examples:
      "best SEO learning resources for beginners"
      "best free SEO tools for beginners"
      "top SEO blogs for beginners"

  Pattern B - Competitor inclusion check
    "[competitor brand] OR [competitor brand] [topic] list"
    Examples:
      "moz ahrefs SEO beginner resources"
      "searchenginejournal beginner SEO guide list"

  Pattern C - Roundup articles
    "[topic] resources roundup [current year OR prior year]"
    Examples:
      "SEO resources roundup 2025"
      "learn SEO roundup 2024"

For each result found:
  - Check if user's domain is already included
  - Note publication name, article title, URL
  - Note author name if visible
  - Note how many resources are listed
  - Flag whether article appears active or outdated

### STEP 3 - Find resource pages
Use webSearch with these patterns:

  Pattern A - Explicit resource pages
    "[topic] resources" OR "[topic] learning resources" inurl:resources
    Examples:
      "SEO learning resources" inurl:resources
      "search engine optimization resources for beginners" inurl:resources

  Pattern B - Link pages
    "[topic] useful links" OR "[topic] recommended sites"
    Examples:
      "SEO useful links"
      "learn SEO recommended sites"

  Pattern C - University and educational resource pages
    "site:.edu [topic] resources"
    Examples:
      "site:.edu SEO resources"
      "site:.edu search engine optimization learning"

For each result:
  - Confirm it links out to external resources
  - Check if user's domain is listed
  - Note domain authority signals (university, government, major publication)
  - Flag as high priority if .edu or .gov

### STEP 4 - Find broken link opportunities
Use webSearch to find:

  Pattern A - Dead resource pages
    "[topic] [resource type]" then check for pages that reference
    tools or sites that may have shut down
    Examples:
      "SEO tools for beginners" + look for references to deprecated tools
      "free keyword research tools" + look for tools that no longer exist

  Pattern B - Wayback Machine signals
    Use webSearch for "[topic] tool discontinued" or
    "[topic] resource no longer available"

For each broken link opportunity:
  - Identify what content the dead link pointed to
  - Confirm user has equivalent or better content
  - Note the linking page URL and domain

### STEP 5 - Find AI citation source targets
This is the highest-priority target type.

Use webSearch to find what AI platforms reference:

  Pattern A - Cited source discovery
    Search for exact phrases AI systems use when explaining the topic
    Examples:
      "what is search engine optimization" - find pages AI cites
      "how does Google indexing work" - find pages AI cites
      "why is my website not showing on Google" - find pages AI cites

  Pattern B - Perplexity and ChatGPT source patterns
    Search for "[topic] explained" site:moz.com OR site:ahrefs.com
    OR site:searchenginejournal.com
    Identify which competitor pages dominate AI citations
    Then find whether those pages link out to additional resources

  Pattern C - "According to" citation patterns
    "[topic] according to" OR "as explained by [competitor brand]"
    Find which publications reference competitor content
    These publications are likely to cite user content too

For each AI citation source found:
  - Note the URL and domain
  - Note the specific topic it covers
  - Identify whether it links to external resources
  - Flag as Tier 1 target

### STEP 6 - Find guest post and contribution opportunities
Use webSearch with these patterns:

  Pattern A - Write for us pages
    "[topic] "write for us"" OR "[topic] "contribute""
    Examples:
      "SEO blog write for us"
      "digital marketing contribute guest post"

  Pattern B - Sites that have published guest posts before
    "[topic] guest post" OR "[topic] contributed by"
    Examples:
      "SEO guest post 2025"
      "search engine basics guest post"

For each opportunity:
  - Note domain authority signals
  - Identify whether they cover beginner-level content
  - Flag if they have published content similar to user's existing articles

### STEP 7 - Find geographic and audience gaps
Use webSearch to find:

  Pattern A - Audience-specific gaps
    "[topic] for beginners" - check if user is missing from results
    "[topic] for small businesses" - same check
    "[topic] explained simply" - same check

  Pattern B - Format gaps
    "[topic] checklist" - is user missing from checklist roundups?
    "[topic] glossary" - is user missing from glossary collections?
    "[topic] cheat sheet" - same check

For each gap:
  - Identify the specific angle not covered by user's content
  - Flag whether creating this content type would unlock multiple links

### STEP 8 - Score and rank all opportunities found
Apply this scoring matrix:

  SCORE each opportunity on:
    - Domain authority signal (high DA = higher score)
    - AI citation probability (does AI already reference this domain?)
    - Inclusion gap (is user clearly missing from a list they should be on?)
    - Recency (is the article recent and actively maintained?)
    - Contact accessibility (can author be reached?)

  TIER assignment:
    Tier 1: AI citation sources + high DA + clear inclusion gap
    Tier 2: High DA listicles + clear inclusion gap + reachable author
    Tier 3: Resource pages + moderate DA + user content matches
    Tier 4: Guest post opportunities + beginner audience match
    Tier 5: Broken link opportunities

### STEP 9 - Build contact and pitch information
For each Tier 1 and Tier 2 opportunity:

  Find author or editor:
    - Check article byline
    - Search "[author name] Twitter" or "[author name] LinkedIn"
    - Search "[publication] editor [topic]"

  Build pitch angle:
    - What is missing from their article that user's content covers?
    - What makes user's content different from what is already listed?
    - Why would their audience benefit from the addition?

  Generate outreach template:
    Subject: [Specific reference to their article] + [value proposition]
    Opening: Reference the specific article by name
    Body: State exactly what is missing + what user's content adds
    CTA: Single clear ask (add to list, include resource, etc.)

### STEP 10 - Validate before returning
Check for:
  - minimum 5 specific named targets identified
  - at least one Tier 1 AI citation source included
  - all targets have URL, publication name, and pitch angle
  - outreach templates are specific, not generic
  - broken link and resource page types both represented
  - no generic advice without named evidence

### STEP 11 - Build output
Present as:
  1. Tier 1 targets - AI citation sources
     (specific URLs, why they matter, pitch angle)
  2. Tier 2 targets - Listicles missing the brand
     (article title, URL, author, pitch angle)
  3. Tier 3 targets - Resource pages
     (URL, domain type, inclusion gap)
  4. Tier 4 targets - Guest post opportunities
     (publication, audience fit, topic angle)
  5. Tier 5 targets - Broken link opportunities
     (linking page, dead link topic, replacement content)
  6. ONE recommended first outreach
     (highest probability win this week)
  7. ONE follow-up question

### GUARDRAILS:
- Never return generic categories without specific named examples
- Never skip AI citation source identification
- Never recommend outreach without a specific pitch angle
- Always include at least one Tier 1 target
- Always provide author or contact method when findable
- Never confuse resource pages with listicles - treat them separately
- Always check whether user's domain is already included before flagging
- Never recommend a target without stating why the user belongs there

## Purpose

Provide procedural guidance to identify pages and relationships that create credible link opportunities.

## When to use

- Use when the authorized intent is `seo_planning` and the request is to identify pages and relationships that create credible link opportunities.

## When NOT to use

- Do not use when the request belongs to `backlink_strategy`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `link_opportunity_discovery_report`

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
