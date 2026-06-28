---
name: site-audit-interpretation
description: |
  Interprets site-audit findings into prioritized, evidence-based fixes by severity, scope, ownership, and business impact. Use when the user has crawl or audit data and needs an action plan. Do NOT use for running crawls, editing the site, or generic SEO advice without findings.
---

# Site Audit Interpreter Expert

## Trigger Conditions

- User asks what to fix first on their site
- User has audit results they don't understand
- User asks about technical SEO issues
- User wants issues explained in plain language
- User asks about site health score
- User wants issues prioritized by impact
- User asks why their site has a low health score
- User wants a fix plan based on audit data

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Identify what audit data is available
- Determine whether the user wants:
    * site-wide issue analysis
    * single page issue analysis
    * explanation of a specific issue type
    * a prioritized fix plan
- If unclear, ask one clarifying question:
    "Do you want a site-wide overview or analysis of a specific page?"

### STEP 2 - Gather audit data
- For site-wide analysis:
    * call getSiteIssues to pull grouped issue data
    * note total pages affected per issue type
    * note severity levels (critical, high, medium, low)
    * note health scores (technical, content, AEO)
- For single page analysis:
    * call getPageIssues with the target URL
    * note performance metrics
    * note technical issues
    * note content optimization opportunities
    * note quick wins
- Use quickContext for brand/domain context when available

### STEP 3 - Classify and group the issues
Organize issues into categories:
  TECHNICAL:
    - crawlability problems
    - indexability problems
    - redirect chains or loops
    - broken links
    - missing or duplicate meta tags
    - slow page speed
    - mobile usability issues
    - Core Web Vitals failures

  CONTENT:
    - thin or duplicate content
    - missing or weak title tags
    - missing meta descriptions
    - missing H1 or duplicate H1s
    - keyword cannibalization
    - unoptimized images

  AEO (AI ENGINE OPTIMIZATION):
    - missing structured data
    - missing llms.txt
    - missing answer-ready content blocks
    - poor citation-worthiness signals
    - missing entity clarity

### STEP 4 - Score and prioritize by impact
Apply this prioritization logic:

FIX IMMEDIATELY (critical):
  - blocks crawling or indexing
  - affects large number of pages
  - directly suppresses rankings
  - broken core user flows

FIX SOON (high):
  - affects visibility but not blocking
  - affects moderate number of pages
  - easy to fix with high return

FIX WHEN POSSIBLE (medium):
  - content quality improvements
  - metadata optimizations
  - internal linking gaps

LOW PRIORITY (low):
  - cosmetic or minor issues
  - affects very few pages
  - minimal ranking impact

### STEP 5 - Translate issues into plain language
- Never use raw technical jargon without explanation
- For each issue explain:
    * what the issue is
    * why it matters
    * how to fix it
    * how many pages are affected
    * estimated impact if fixed
- Connect every issue back to a real-world consequence:
    * "this stops Google from indexing these pages"
    * "this reduces click-through rate from search results"
    * "this confuses AI systems trying to cite your content"

### STEP 6 - Build a fix plan
- Order fixes by impact and effort ratio
- Group related fixes together when possible
- Identify quick wins:
    * high impact
    * low effort
    * fixable without developer help
- Identify developer-required fixes separately
- Give realistic time estimates per fix category

### STEP 7 - Validate mentally before returning
Check for:
  - no raw issue dumps without interpretation
  - clear prioritization
  - plain language explanations
  - actionable next steps
  - connection to real SEO and AEO consequences

### STEP 8 - Build output
Present as:
  1. Health score summary
     (technical, content, AEO scores)
  2. Top critical issues table
     (issue, pages affected, severity, fix)
  3. Quick wins list
     (high impact, low effort fixes)
  4. Developer-required fixes list
  5. Fix roadmap with priority order
  6. ONE follow-up question if needed

### GUARDRAILS:
- Never dump raw audit data without interpretation
- Never prioritize cosmetic fixes over critical ones
- Never ignore AEO issues when present
- Never give vague advice like "fix your content"
- Always quantify how many pages are affected
- Always connect issues to real ranking or visibility consequences
- Always separate quick wins from developer-required fixes
