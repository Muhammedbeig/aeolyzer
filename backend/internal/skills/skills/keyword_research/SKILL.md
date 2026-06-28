---
name: keyword-research
description: |
  Finds and prioritizes keywords by intent, volume, difficulty, relevance, competition, and business value. Use when the user needs target queries, keyword clusters, or ranking opportunities. Do NOT use for complete content drafts, page implementation, or unsupported volume claims.
---

# Keyword Research Expert

## Trigger Conditions

- User asks for keyword research
- User asks for search volumes
- User asks for keyword ideas
- User asks what to target with content
- User asks for traffic potential analysis

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Understand the project context
- Call quickContext to get domain, niche, positioning
- Identify domain authority level (new site vs established)
- Identify content model (educational, commercial, local, etc.)
- Note stored competitors

### STEP 2 - Generate seed keywords
- Do NOT ask the user for seeds unless truly unclear
- Derive seeds from: domain name, niche, positioning, 
  user's question, and project competitors
- Generate 5-8 diverse seed angles covering:
    * Core topic terms
    * Problem-based queries (why, how, what)
    * Beginner vs advanced intent splits
    * Platform-specific variations if relevant

### STEP 3 - Call getKeywords sequentially
- Call getKeywords once per seed term
- Wait for each result before calling the next
- Collect: search_volume, keyword_difficulty, 
  competition_level, avg_backlinks_info, 
  search_volume_trend, search_intent_info, cpc

### STEP 4 - Filter and score results
Apply this scoring logic:

SKIP if:
  - keyword_difficulty > 70 AND domain is new/low authority
  - search_volume < 10 (not enough signal)
  - search_volume_trend.yearly < -70 (dying keyword)
  - main_intent = commercial AND site is purely educational

PRIORITIZE if:
  - keyword_difficulty < 40
  - search_volume > 100
  - competition_level = LOW
  - avg_backlinks_info.dofollow < 50 (low backlink bar)
  - main_intent = informational (matches educational sites)
  - search_volume_trend.monthly > 0 (growing)

### STEP 5 - Cluster by topic and intent
- Group related keywords into topic clusters
- One cluster = one content piece
- Name each cluster by its primary keyword
- List secondary/supporting keywords per cluster

### STEP 6 - Bucket into tiers
QUICK WINS:
  - KD < 40, volume > 50, low backlink bar
  - Target within 1-2 months

MEDIUM TERM:
  - KD 40-60, volume > 200
  - Target within 2-3 months

AVOID NOW:
  - KD > 70 OR backlink bar > 500 dofollow
  - Flag why (too competitive, declining, wrong intent)

### STEP 7 - Build output
Present as:
  1. Quick wins table (keyword, volume, KD, why it wins)
  2. Medium term table
  3. Avoid table with reasons
  4. Implementation roadmap with timeline
  5. Estimated traffic impact if top 3 achieved
  6. ONE follow-up question to go deeper

### GUARDRAILS:
- Never present raw data without interpretation
- Always connect recommendations to domain authority
- Always flag trend direction (growing vs declining)
- Never suggest KD 70+ to a new or low-authority site
- Roadmap must be time-bound, not vague
