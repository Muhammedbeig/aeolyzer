---
name: content-ideas
description: Finds People Also Ask questions, related searches, and differentiated content angles for a topic. Use when the user wants blog ideas, question-based prompts, or topic opportunities. Do NOT use for full content strategies, complete drafts, or page-specific audits.
version: 1.0.0
owner_team: content_platform
tier: read
risk_class: low
compatible_profiles:
    - content_collaborator
compatible_intents:
    - topic_discovery
allowed_modes:
    - plan
    - read
capability_tags:
    - content_ideas
declared_action_classes:
    - read_brand_context
    - read_source_intelligence
output_contracts:
    - content_ideas_report
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
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

## Purpose

Provide procedural guidance to generate evidence-aware content ideas before a brief exists.

## When to use

- Use when the authorized intent is `topic_discovery` and the request is to generate evidence-aware content ideas before a brief exists.

## When NOT to use

- Do not use when the request belongs to `topic_discovery`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `content_ideas_report`

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
