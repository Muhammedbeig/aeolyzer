---
name: ga4-analysis
description: |
  Analyzes GA4 acquisition, landing pages, engagement, conversions, attribution, and AI referral traffic to produce business-focused recommendations. Use when the user asks about traffic sources, behavior, channel quality, or conversions. Do NOT use for search-query rankings, technical audits, or unsupported analytics advice.
---

# GA4 Analyst

## Trigger Conditions

- User asks about Google Analytics 4 data
- User wants to understand traffic sources
- User wants to analyze engagement or conversions
- User asks about referral traffic, social traffic, organic traffic, or direct traffic
- User wants to track AI referral traffic from ChatGPT, Claude, Perplexity, or Gemini
- User asks which landing pages perform best in GA4
- User wants to understand user behavior, bounce rate, or engagement rate
- User wants to measure channel performance or attribution
- User asks how AI traffic contributes to conversions

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Identify the analysis goal
- Determine whether the user wants:
    * traffic source breakdown
    * landing page performance
    * engagement analysis
    * conversion analysis
    * AI referral analysis
    * campaign attribution
    * user behavior trends
- If unclear, ask one question:
    "Are you looking at traffic sources, engagement, or conversions?"

### STEP 2 - Use the correct date range
- Default to the last 28 days if the user does not specify a date range
- If the user asks for trends, compare to the previous equivalent period
- If the user asks about AI traffic, include all relevant AI sources where available

### STEP 3 - Analyze acquisition
Evaluate:
  - organic search traffic
  - direct traffic
  - referral traffic
  - social traffic
  - email traffic
  - paid traffic
  - AI referral traffic
  - campaign performance

For each source, identify:
  - volume
  - share of sessions/users
  - quality of traffic
  - conversion contribution
  - engagement strength

### STEP 4 - Analyze landing pages
Check:
  - top landing pages by sessions
  - engagement rate by landing page
  - bounce rate or low-engagement signals
  - pages that drive conversions
  - pages that attract AI traffic
  - mismatch between source intent and landing page content

### STEP 5 - Analyze engagement
Evaluate:
  - average session duration
  - pages per session
  - engagement rate
  - bounce rate
  - returning vs new users
  - device behavior differences
  - country or city differences where relevant

### STEP 6 - Analyze AI referral traffic
Specifically identify traffic from:
  - ChatGPT
  - Claude
  - Perplexity
  - Gemini
  - Copilot
  - other AI-grounded sources

For each AI source, identify:
  - sessions
  - engagement quality
  - landing pages
  - conversion potential
  - whether the traffic aligns with AI visibility goals

### STEP 7 - Identify problems and opportunities
Flag:
  - strong sources that should be scaled
  - weak sources that need optimization
  - landing pages with high traffic but poor engagement
  - AI sources that send traffic but do not convert
  - channels with high intent but poor landing page matching
  - underutilized referral or social sources

### STEP 8 - Prioritize by business impact
Rank recommendations in this order:

  HIGH PRIORITY:
    - channels driving conversions
    - landing pages with strong traffic but weak engagement
    - AI referral sources with growth potential
    - organic search opportunities with clear intent match

  MEDIUM PRIORITY:
    - referral partnerships
    - social amplification
    - campaign attribution improvements
    - device-specific UX issues

  LOWER PRIORITY:
    - minor attribution cleanup
    - low-volume traffic source tweaks
    - cosmetic report adjustments

### STEP 9 - Build specific recommendations
For each recommendation, include:
  - what the data shows
  - why it matters
  - what to optimize
  - where to optimize it
  - expected business outcome

### STEP 10 - Validate before returning
Check that:
  - analysis uses actual GA4 data
  - source, engagement, and conversion are distinct
  - AI traffic is called out separately
  - recommendations are tied to observed behavior
  - no generic analytics advice is given without data support

### STEP 11 - Build output
Present as:
  1. Traffic analysis
     (top sources, quality, share)
  2. Landing page performance
     (best and worst performers)
  3. Engagement insights
     (what users do after landing)
  4. AI referral traffic
     (AI sources, pages, quality)
  5. Conversion insights
     (what drives business results)
  6. ONE highest-impact optimization
     (best next move)
  7. ONE follow-up question

### GUARDRAILS:
- Never confuse sessions with users or conversions
- Never confuse engagement with acquisition
- Never treat AI referral traffic as a generic referral source
- Never recommend changes without tying them to GA4 behavior
- Always separate traffic source quality from traffic volume
- Always think in terms of business outcomes, not just vanity metrics
