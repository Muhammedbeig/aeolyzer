---
name: long-form-content-audit
description: |
  Audits a specific long-form page against its content, search intent, competitors, structure, depth, trust signals, and AI visibility. Use when the user provides a URL and wants page-specific improvements. Do NOT use without inspecting the page, for net-new drafts, or sitewide technical audits.
---

# Long Form Content Expert

## Trigger Conditions

- User provides a URL to analyze
- User wants page-specific content feedback
- User asks how to improve an existing article
- User wants competitive improvement recommendations
- User asks why a specific page is not ranking
- User wants content quality analysis for a URL
- User asks what is missing from an existing piece
- User wants to know how a page compares to top competitors

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Identify the target page
- Extract the URL from the user's request
- If no URL is provided, ask one clarifying question:
    "Which page URL do you want me to analyze?"
- Determine what the user wants to know:
    * overall content quality assessment
    * competitive gap analysis
    * specific section improvements
    * AI visibility optimization
    * ranking improvement recommendations

### STEP 2 - Inspect the page content
- Use scrapePage to extract full page content
- Identify:
    * primary topic and keyword focus
    * content structure (H1, H2, H3 hierarchy)
    * word count and depth
    * content format (guide, tutorial, list, etc.)
    * answer-ready blocks present or missing
    * structured data present or missing
    * internal and external links
    * images and media usage
    * author and E-E-A-T signals

### STEP 3 - Analyze the competitive landscape
- Use analyzeSERP for the primary keyword
- Identify top 3 ranking pages for the same query
- Compare against the user's page:
    * word count gap
    * structural differences
    * topics covered vs missing
    * format differences
    * E-E-A-T signal differences
    * schema markup differences
- Identify what top-ranking pages do that this page does not

### STEP 4 - Evaluate content quality
Apply this quality scoring logic:

  STRUCTURE:
    - Is there a clear H1?
    - Are H2s logically organized?
    - Does the intro answer the question quickly?
    - Is the content scannable?

  DEPTH:
    - Does it cover the topic comprehensively?
    - Are there missing subtopics vs competitors?
    - Are claims supported with specifics or data?
    - Is there a clear conclusion?

  SEARCH INTENT ALIGNMENT:
    - Does the format match what Google prefers?
    - Does the content satisfy the likely user goal?
    - Is the keyword used naturally and appropriately?

  E-E-A-T SIGNALS:
    - Is there a named author?
    - Are there credentials or expertise signals?
    - Are there citations or references?
    - Is the content original or generic?

  AI VISIBILITY POTENTIAL:
    - Are there answer-ready definition blocks?
    - Are there numbered steps for how-to queries?
    - Are there FAQ sections?
    - Is the content structured for extraction?

### STEP 5 - Identify specific improvement opportunities
Flag issues in order of impact:

  CRITICAL GAPS:
    - missing primary keyword in H1 or intro
    - content significantly shorter than top competitors
    - no structured data present
    - intent mismatch with current SERP format

  HIGH IMPACT IMPROVEMENTS:
    - missing subtopics covered by top competitors
    - no answer-ready blocks for AI citation
    - weak or missing introduction
    - no FAQ section when PAA boxes exist

  MEDIUM IMPACT IMPROVEMENTS:
    - thin sections that need expansion
    - outdated statistics or examples
    - missing internal links to related content
    - no author bio or E-E-A-T signals

  LOW IMPACT IMPROVEMENTS:
    - image alt text optimization
    - minor structural tweaks
    - additional external references

### STEP 6 - Build specific recommendations
For each gap identified, provide:
  - what is missing
  - why it matters for rankings or AI visibility
  - exactly how to fix it
  - where in the content to make the change

### STEP 7 - Validate mentally before returning
Check for:
  - page content actually inspected before advising
  - competitive comparison included
  - recommendations tied to specific content gaps
  - prioritization by impact
  - no generic advice without page-specific evidence

### STEP 8 - Build output
Present as:
  1. Page summary
     (topic, format, word count, current status)
  2. Competitive gap analysis
     (what top-ranking pages have that this page lacks)
  3. Content quality assessment
     (structure, depth, intent, E-E-A-T, AI visibility)
  4. Prioritized improvement recommendations
     (critical, high, medium, low impact)
  5. Quick wins
     (highest impact, lowest effort fixes)
  6. ONE follow-up question if needed

### GUARDRAILS:
- Never advise without inspecting the actual page first
- Never give generic content advice without competitive context
- Never ignore E-E-A-T signals in the assessment
- Never skip AI visibility evaluation
- Never recommend a full rewrite when targeted improvements will do
- Always tie recommendations to specific evidence from the page
- Always include competitive comparison as the benchmark
- Always prioritize improvements by impact not by ease
