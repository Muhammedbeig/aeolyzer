---
name: meta-optimization
description: Optimizes page titles and meta descriptions for relevance, intent, CTR, uniqueness, and length. Use when the user asks for metadata analysis or rewrites for known pages or keywords. Do NOT use for full content drafts, Open Graph packages, or sitewide technical audits.
version: 1.0.0
owner_team: content_platform
tier: draft
risk_class: medium
compatible_profiles:
    - content_execution_guard
compatible_intents:
    - optimize_content
allowed_modes:
    - write
    - edit
    - optimize
capability_tags:
    - meta_optimization
declared_action_classes:
    - read_brand_context
    - canvas_write
output_contracts:
    - meta_optimization_draft
    - quality_summary
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Meta Tag Optimization Expert

## Trigger Conditions

- User asks for title tags
- User asks for meta descriptions
- User wants better CTR from search results
- User wants optimized SEO metadata
- User wants to improve page snippets in Google
- User asks how to write better page titles or descriptions
- User wants metadata for a specific page type or URL

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Identify the page and goal
- Determine what page type the user is optimizing:
    * homepage
    * product page
    * blog post
    * category page
    * landing page
    * service page
- Identify the goal:
    * improve CTR
    * improve relevance
    * better keyword targeting
    * more compelling search snippets
- If unclear, ask one clarifying question:
    "What page are you optimizing and what keyword or theme should it target?"

### STEP 2 - Gather source context
- If a URL is provided, inspect the page content first
- Use quickContext when brand context matters
- Pull the main topic, audience, and value proposition from the source material
- Do not invent claims or benefits that are not supported by the page

### STEP 3 - Build the title tag
- Keep titles concise and readable
- Put the primary keyword near the front when natural
- Make the title specific to the page
- Differentiate from other pages on the site
- Avoid clickbait unless it truly fits the brand voice
- Avoid repeated or duplicated title structures

### STEP 4 - Build the meta description
- Summarize the page clearly
- Include the value proposition and a reason to click
- Use active language
- Match the search intent
- Include relevant entities or terms naturally
- Keep it persuasive but accurate

### STEP 5 - Optimize for snippets and CTR
- Write metadata that aligns with likely search intent
- Use benefit language when appropriate
- Make titles and descriptions distinct enough to avoid cannibalization
- Favor clarity over stuffing
- If helpful, provide multiple variants for testing

### STEP 6 - Validate mentally before returning
Check for:
  - correct page intent
  - natural keyword placement
  - no unsupported claims
  - no duplicate or overly generic metadata
  - concise, search-friendly wording
  - CTR potential without misleading wording

### STEP 7 - Build output
Present as:
  1. Recommended title tag
  2. Recommended meta description
  3. Short explanation of why they fit
  4. Optional alternate variants if useful
  5. ONE follow-up question if needed

### GUARDRAILS:
- Never invent features or benefits
- Never keyword-stuff titles or descriptions
- Never use the same metadata pattern on every page
- Never ignore search intent
- Always keep titles and descriptions aligned with visible page content

## Purpose

Provide procedural guidance to draft title and description recommendations within measured constraints.

## When to use

- Use when the authorized intent is `optimize_content` and the request is to draft title and description recommendations within measured constraints.

## When NOT to use

- Do not use when the request belongs to `title_generation`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `meta_optimization_draft`
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
