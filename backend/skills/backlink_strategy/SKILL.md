---
name: backlink-strategy
description: |
  Develops backlink and citation acquisition strategies from current link profiles, cited sources, competitors, and outreach targets. Use when the user wants links, mentions, listicle inclusion, or outreach planning. Do NOT use for internal linking, technical audits, or content drafting.
---

# Backlink Strategy Expert

## Trigger Conditions

- User asks how to get backlinks
- User wants to be featured in AI results
- User wants link acquisition strategy
- User asks what sites link to competitors
- User wants to find citation opportunities in AI platforms
- User asks how to get mentioned by ChatGPT, Claude, or Perplexity
- User wants to find articles or listicles to be included in
- User wants outreach strategy for link building

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Identify the goal
- Determine what the user wants:
    * traditional backlink acquisition (rankings)
    * AI citation acquisition (ChatGPT, Claude, Perplexity mentions)
    * both simultaneously
- If unclear, ask one question:
    "Are you focused on ranking higher in Google, appearing more in AI answers,
     or both?"

### STEP 2 - Identify the target brand and domain
- Use project context to confirm domain (searchenginebasics.io)
- Confirm primary topic or keyword focus if not already known
- Check memory for existing competitor list

### STEP 3 - Find what AI platforms are already citing
- Use getCitationSources to identify which domains are cited most
  in AI responses for this project
- Use getCitedUrls to see which specific pages get cited
- Note:
    * which domains appear repeatedly
    * what content types get cited (guides, lists, comparisons)
    * which AI platforms cite which domains

### STEP 4 - Identify link-worthy content targets
Use webSearch to find:
  a) Listicles and roundups
     - "best [topic] tools" articles
     - "top [topic] resources" lists
     - "beginner's guide to [topic]" compilations
     - "[topic] glossary" pages
     Search pattern: "[primary keyword] site:reddit.com OR site:medium.com OR
                      site:forbes.com OR site:[competitor domain]"

  b) Articles that AI platforms reference
     - Search for the exact phrases AI uses when citing competitors
     - Find the source articles those phrases come from
     - Identify whether those articles accept contributions or updates

  c) Outdated content with link gaps
     - Find articles older than 18 months covering the same topic
     - Identify missing tools, resources, or examples
     - Flag these as outreach targets for "you missed us" pitches

  d) Geographic or audience gaps
     - Find listicles that cover competitors but not the user's brand
     - Note which are missing beginner-focused or plain-language resources

### STEP 5 - Analyze competitor backlink profiles
- Use analyzeCompetitors if not already done for this session
- Use webSearch to find:
    * "[competitor domain] mentioned in" articles
    * "according to [competitor brand]" citations
    * "[competitor brand] review" roundups
- Identify patterns:
    * What content type earns them the most links
    * Which domains link to multiple competitors (link hubs)
    * Which publications cover the topic regularly

### STEP 6 - Find specific outreach targets
Apply link-opportunity-finder logic:

  TARGET TYPE 1 - Listicle inclusion targets
    - Find "best [topic]" articles ranking in top 10
    - Check if user's domain is already included
    - If not, flag as outreach target
    - Note author name, publication, contact method

  TARGET TYPE 2 - Resource page targets
    - Find "resources for [topic]" pages
    - Check if user's domain is linked
    - Flag missing inclusions as targets

  TARGET TYPE 3 - Broken link targets
    - Find pages linking to dead resources in the niche
    - Identify user's content that could replace the broken link
    - Flag for broken link outreach

  TARGET TYPE 4 - AI citation source targets
    - Find the specific articles ChatGPT, Claude, or Perplexity cite
      when answering questions in this topic area
    - Use webSearch with exact phrases from AI responses
    - Identify whether those articles link out to resources
    - Flag as highest priority (getting linked here = AI citation potential)

### STEP 7 - Prioritize targets by impact
Rank outreach targets using this logic:

  TIER 1 - AI citation sources
    (pages that AI platforms already cite in this topic area)
    Impact: getting linked here increases probability of AI citation
    Effort: medium - requires strong content and compelling pitch

  TIER 2 - High-DA listicles missing the brand
    (established articles ranking for key terms, not including user)
    Impact: direct referral traffic + ranking signal
    Effort: low to medium - "you missed us" pitch

  TIER 3 - Competitor link hubs
    (domains that link to 2+ competitors but not the user)
    Impact: closes competitive gap
    Effort: medium - requires demonstrating differentiation

  TIER 4 - Broken link opportunities
    (dead links replaceable with user's content)
    Impact: moderate ranking signal
    Effort: low - straightforward value proposition

### STEP 8 - Build outreach strategy
For each Tier 1 and Tier 2 target, provide:
  - publication name and URL
  - article title and URL
  - author name if findable
  - contact method (Twitter, LinkedIn, email pattern)
  - pitch angle:
      * what is currently missing from their article
      * what the user's content adds that competitors do not
      * why their readers would benefit from the inclusion
  - suggested subject line
  - suggested opening sentence

### STEP 9 - Content gap bridge
Identify whether the user has content worth linking to:
  - If yes: map existing content to outreach targets
  - If no: recommend creating one "link magnet" piece first
      * what format (data study, glossary, tool, guide)
      * what topic (based on what AI platforms cite most)
      * why this format earns links in this niche

### STEP 10 - Validate before returning
Check for:
  - specific named targets, not generic advice
  - AI citation sources identified separately from traditional links
  - outreach angles tied to actual content gaps
  - prioritization by impact, not ease
  - at least one actionable pitch template included

### STEP 11 - Build output
Present as:
  1. What is currently being cited
     (AI citation sources, competitor link patterns)
  2. Top outreach targets by tier
     (specific articles, publications, authors)
  3. Pitch strategy per target
     (angle, subject line, opening)
  4. Content gap assessment
     (do you have something worth linking to?)
  5. ONE recommended first action
     (highest impact move to make this week)
  6. ONE follow-up question

### GUARDRAILS:
- Never give generic "write great content and links will come" advice
- Never recommend paid links or link schemes
- Never skip AI citation source identification
- Always name specific targets, not categories
- Always tie outreach angle to a real content gap
- Always separate AI citation strategy from traditional link building
- Never recommend outreach without identifying what content to pitch
- Always prioritize Tier 1 (AI citation sources) above all else
