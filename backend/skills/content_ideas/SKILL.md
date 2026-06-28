---
name: content-ideas
description: |
  Finds People Also Ask questions, related searches, and differentiated content angles for a topic. Use when the user wants blog ideas, question-based prompts, or topic opportunities. Do NOT use for full content strategies, complete drafts, or page-specific audits.
---

# Content Ideas / PAA Research Expert

## Trigger Conditions

- User asks for content ideas
- User asks what to write about
- User wants People Also Ask questions
- User wants related searches
- User wants topic angles for an article
- User wants question-based content prompts
- User asks for blog topics around a keyword

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Understand the topic
- Extract the core topic, keyword, or theme from the user's request
- If unclear, ask once: "What topic are you thinking? (Want ideas? I can research what's trending)"
- Use quickContext to understand the brand niche and audience
- Determine whether the user wants:
    * educational content
    * industry trend content
    * commercial comparison content
    * beginner explainer content

### STEP 2 - Call getContentIdeas
- Call getContentIdeas with:
    * topic = exact topic from the user
    * location = United States (default)
- This returns:
    * People Also Ask questions
    * related search ideas
    * question variations
    * content angle suggestions

### STEP 3 - Interpret the idea set
Organize the results into:
  - beginner questions
  - problem/solution questions
  - comparison questions
  - deeper technical questions
  - commercial intent questions if present

### STEP 4 - Filter for best content opportunities
Prioritize ideas that:
  - are specific and question-based
  - match the user's niche and expertise
  - can support a full article or content cluster
  - reveal gaps in existing coverage
  - align with informational intent

Deprioritize ideas that:
  - are too generic
  - are purely navigational
  - are unrelated to the user's positioning
  - are too broad to own with one article

### STEP 5 - Build content recommendations
For each strong idea, provide:
  - exact question or angle
  - suggested content format
  - why it matters
  - how it fits into a cluster

Example:
  - "What is crawl budget?" -> definition guide
  - "How does Google crawl a website?" -> step-by-step explainer
  - "Why is my site not indexing?" -> troubleshooting guide

### STEP 6 - Build output
Present as:
  1. Top content ideas table
     (question, intent, recommended format)
  2. Best cluster opportunity
     (which questions can live in one article)
  3. Content angle recommendation
     (exact title or working headline)
  4. ONE follow-up question to go deeper

### GUARDRAILS:
- Never return raw question dumps without grouping
- Always turn questions into content strategy
- Always prioritize fit with the user's niche
- If the topic is vague, infer the most likely informational intent first
