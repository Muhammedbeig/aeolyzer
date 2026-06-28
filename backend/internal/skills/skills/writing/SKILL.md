---
name: writing
description: Guidelines for editorial voice, evidence depth, word count tracking, and anti-redundancy in article and blog post writing. Use when writing any article, blog post, or long-form content piece — even if the user doesn't say "skill" or mention structure explicitly. Do NOT use for research-only tasks, metadata-only work, or short coordination messages.
version: 1.0.0
owner_team: content_platform
tier: draft
risk_class: medium
compatible_profiles:
    - content_execution_guard
compatible_intents:
    - draft_article
allowed_modes:
    - write
    - edit
    - optimize
capability_tags:
    - writing
declared_action_classes:
    - read_brand_context
    - canvas_write
output_contracts:
    - writing_draft
    - quality_summary
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Writing

## EDITORIAL VOICE

- Take positions. Don't just present information - interpret it. Say what it MEANS.
- After presenting a fact, add "so what" analysis: what does this mean for the reader's business/decisions?
- Use direct analytical language:
  - "Here's what this actually means for [audience]:"
  - "The part most coverage misses:"
  - "The real risk is..."
  - "This matters because..."
  - "What no one is talking about:"
- NEVER hedge more than once per paragraph. Pick a position and defend it.
  BAD: "This could potentially maybe lead to some changes"
  GOOD: "This will force decision-makers to rethink their entire approach"
- Distinguish your piece from a press release. Press releases inform. Your article should ANALYZE.
  For every claim or feature described, ask: "So what? Why should the reader care?"

## DEPTH AND EVIDENCE (REQUIRED)

- Go beyond the press release. For every major topic, address:
  - What's missing from the standard narrative?
  - What second-order effects will this cause?
  - What's the competitive context? (Who else is doing this? What's their angle?)
  - What are the counter-arguments or risks?
- Reference SPECIFIC numbers, dates, companies, and quotes.
  BAD: "Many companies are adopting this approach"
  GOOD: "Over 1 million businesses, including [specific named examples], have adopted [specific thing] since [specific date]"
- Include counter-arguments. If you argue for X, acknowledge the strongest case against X.
  This builds credibility, not weakness.
- If the plan includes sources with key data, WEAVE that data into the prose.
  Don't list facts. Contextualize them: "[Stat] - which means [interpretation for reader]."

## PRE-WRITING CHECKS

- If no target word count was set during planning, ask: "Any length preference, or should I aim for ~500-2000 words?"
- Respect the user's answer. If they say "short", aim for the lower end.
- Blog Post: target 500-2000 words, long-form style
- Write from a first-person perspective (I/we). Conversational, personality-forward, relatable.

### VOICE

The reader should feel like they're hearing from a real person, not a brand.
Use personal anecdotes, specific experiences, and candid opinions.

### STRUCTURE

Shorter paragraphs than articles. Use H2 for sections, bullet lists for
tips/steps, bold for key takeaways. Break every 2-3 paragraphs with a formatting element.

### DEPTH

Blog posts can be lighter on research but MUST be specific. Replace 'it works great'
with 'it cut our deploy time from 12 minutes to 3'.

### AVOID

Corporate tone. Third-person detachment. Listing without commentary.
Padding a 500-word idea into 1500 words.

## WORD COUNT TRACKING

- Each writeSection call returns a "wordCount" field. This is the REAL word count for that section.
- Keep a RUNNING TOTAL by summing the wordCount from each writeSection result.
- Your context includes <canvas_state> with the current word count at session start.
- DO NOT estimate or guess the total. Trust the numbers from writeSection results.
- When your running total approaches the target, STOP writing. Do not overshoot.
- When reporting the final word count, add up the tool-returned counts. Do NOT fabricate a number.

## ANTI-REDUNDANCY (CRITICAL)

- NEVER make the same point in two different sections. Each section must advance the argument.
- If you established a distinction in section 2, do NOT re-explain it in sections 4, 5, and 6.
  Reference it briefly ("as discussed above") but add NEW analysis, not the same observation
  in different words.
- Before writing each section, mentally review what you've already said. If a point was covered, skip it.
- Common redundancy traps to avoid:
  - Restating the core thesis in every section introduction
  - Re-explaining what two things are after the explainer sections
  - Repeating the same shared characteristic of two things in 3+ sections
  - Summarizing the entire article in the conclusion instead of advancing it

## HIDDEN INTENT IN WRITING

- If you identified a hidden intent (emotional layer), the introduction MUST resolve it within
  the first 2-3 paragraphs. Don't bury the emotional payoff deep in the article.
- The golden thread should be visible in the intro, reinforced in the body, and resolved in
  the conclusion. It's not a single mention - it's the through-line of the entire piece.
- The title should already capture the hidden intent. If it doesn't, revisit it before writing.

## BRAND AWARENESS

- NEVER cite, link to, or reference any domain from the competitor blocklist.
- NEVER use a direct competitor as the hero case study.
- If the article topic overlaps with their product/service category, be especially careful
  about elevating competitors.
- When competitor data IS useful, use it one of these ways:
  1. Anonymize: "One major platform found that..."
  2. Frame as industry trend: "Across the industry, brands investing in this approach saw..."
  3. Use as contrast/cautionary: frame competitor data as evidence of what not to do
- For outbound links, ONLY use domains from the authoritySources list or well-known
  institutional sources (.gov, .edu, major publications).

## Purpose

Provide procedural guidance to write approved article sections with evidence, voice, and word-count controls.

## When to use

- Use when the authorized intent is `draft_article` and the request is to write approved article sections with evidence, voice, and word-count controls.

## When NOT to use

- Do not use when the request belongs to `content_craft`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `writing_draft`
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
