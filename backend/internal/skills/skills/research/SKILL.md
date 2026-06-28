---
name: research
description: |
  Researches competitor coverage, current news, first-party site overlap, and supporting sources for content. Use when a topic needs evidence, freshness, source URLs, or cannibalization checks. Do NOT use for drafting, editing, or direct publishing.
---

# Research

Research is what separates a strategist from a chatbot. Use it proactively and EARLY - especially BEFORE asking questions, so your options are informed by real data.

## Research depth

- Default for Blog Post: "light" research
- But ALWAYS let the user choose. The user may have their own research and not want you running background searches.
- Offer naturally: "Want me to dig into this? I can do a quick search for stats and quotes, or run deep research for a full competitive analysis. Or if you've got your own material, just share it."
- Respect the user's choice - if they say they have research, don't search anyway.

## TWO RESEARCH MODES

1. QUICK SEARCH (`webSearch`) - default for most situations:
   - Fast, inline results. Use for quick lookups, competitor checks, trending data.
   - Results come back immediately in the conversation.
   - Use this 90% of the time.

2. DEEP RESEARCH (`deepResearch`) - for substantial evidence gathering:
   - Runs in the background (several minutes). Recursively searches, evaluates sources, extracts statistics, quotes, and competitor insights.
   - Use when the content requires: extensive statistics/data, expert quotes, competitive analysis, or multi-faceted topic coverage.
   - Do NOT trigger deep research by default. Ask the user first: "This topic could benefit from deep research - I'd find stats, quotes, and analyze what competitors are saying. Takes several minutes. Want me to run it?"
   - After triggering, tell the user to hang tight: "Researching [topic]. I'll share what I find when it's done."
   - When results come back, present the highlights and use them to inform your plan/writing.
   - Deep research is especially valuable BEFORE `proposePlan` in Plan mode.

## RESEARCH TIMING (critical)

- BEFORE asking questions: Search the topic landscape so your question options reference real data, competitors, and angles - not generic categories.
- BEFORE suggesting topics: Search competitor blogs, industry news, and the user's own content to find gaps and timely opportunities.
- BEFORE proposing a plan: Gather specific sources, stats, and quotes so the plan has real evidence baked in.

## When to research

- User is unsure what to write about -> search competitor content, industry news, content gaps
- User picks a topic -> search for angles, data points, what competitors wrote
- Topic is time-sensitive -> search for latest developments
- User provides a URL -> use `scrapePage` to analyze it
- User mentions their own blog/site content -> search with `site:` prefix or scrape directly

## Searching the user's OWN site

- When user says "we have a blog about X" or "check our site for X":
  1. FIRST try `webSearch` with "site:[your-domain.com] [topic]"
  2. If that returns a URL, use `scrapePage` on it immediately
  3. NEVER ask the user for the URL if you can find it yourself
- When user provides a direct URL, use `scrapePage` immediately - don't search first
- The project domain is: [your-domain.com]

## Research queries that work well

- "[brand/company] recent news"
- "[industry] trends 2026"
- "[topic] statistics data"
- "[competitor] blog [topic]" (to see their angle)
- "[audience] common questions [topic]"
- "site:[your-domain.com] [topic]" (to find content on the user's own site)

## Rules

- Present findings briefly before writing - "Found some good angles..."
- Don't re-research what the user already provided
- Keep research summaries SHORT - bullet points, not essays
- Use research to inform your suggestions, not to dump info on the user
- When scraping fails on the first attempt, try `webSearch` for the URL instead
