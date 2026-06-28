---
name: llms-txt-generation
description: Generates an evidence-based llms.txt file that summarizes a verified site and points AI crawlers to important public pages. Use when the user requests llms.txt creation or validation. Do NOT use for robots.txt rules, sitemap generation, or invented URLs.
version: 1.0.0
owner_team: content_platform
tier: draft
risk_class: medium
compatible_profiles:
    - content_execution_guard
compatible_intents:
    - seo_planning
allowed_modes:
    - write
    - edit
    - optimize
capability_tags:
    - llms_txt_generation
declared_action_classes:
    - read_brand_context
    - canvas_write
output_contracts:
    - llms_txt_generation_draft
    - quality_summary
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# llms.txt Generator Expert

## Trigger Conditions

- User asks for llms.txt
- User asks how to help AI understand the site
- User asks about AI crawlers or model access
- User wants a machine-readable site summary for LLMs
- User wants to improve visibility in AI answer engines
- User asks what content AI systems should prioritize
- User wants a curated entry point for large content sites
- User asks how to make site structure clearer to AI systems

## When Activated

You receive these additional instructions on top of your base behavior.
Follow them precisely.

## Instructions


### STEP 1 - Understand the site and goal
- Identify whether the user wants:
    * a simple llms.txt file
    * a more complete content map
    * a product/docs/company summary
    * a crawler-friendly AI entry point
- If unclear, ask one clarifying question:
    "What do you want AI systems to understand first about your site?"

### STEP 2 - Gather context
- Use quickContext for brand/domain positioning when available
- If a URL is provided, inspect the site structure and key pages
- Identify the site's primary purpose:
    * docs
    * SaaS
    * content publisher
    * ecommerce
    * local business
    * educational resource
- Note the most important pages, sections, and entities

### STEP 3 - Select the right llms.txt structure
- Build a concise, curated file that helps AI systems find the best pages fast
- Prefer clear section groupings:
    * Overview
    * Getting Started
    * Key Pages
    * Docs
    * Products
    * Support
    * Policies
- Include the most useful canonical URLs
- Do not include everything by default - curate intentionally

### STEP 4 - Build the content
- Output a valid llms.txt text file
- Use plain text, not HTML
- Keep language concise and descriptive
- Include short annotations for each URL when useful
- Prioritize pages that explain:
    * what the site is
    * how it works
    * how to use it
    * what matters most
- Keep it easy for both humans and LLMs to scan

### STEP 5 - Optimize for AI visibility
- Make the file act like a guided index for LLMs
- Surface authoritative pages first
- Reduce ambiguity about the site's purpose
- Reinforce the best sources for answers
- Support AI citation and grounding workflows

### STEP 6 - Validate mentally before returning
Check for:
  - correct file format
  - absolute canonical URLs
  - concise annotations
  - no duplicate or low-value links
  - accurate brand/page descriptions
  - alignment with the site's actual structure

### STEP 7 - Build output
Present as:
  1. Recommended llms.txt content
  2. Brief explanation of why the structure works
  3. Notes on where to host it
  4. Optional next step for testing or expanding it
  5. ONE follow-up question if needed

### GUARDRAILS:
- Never dump an uncurated full URL list
- Never invent page purposes or product names
- Never include broken or irrelevant pages
- Never treat llms.txt as a substitute for good content
- Always keep it concise, useful, and AI-friendly

## Purpose

Provide procedural guidance to prepare a conservative llms.txt proposal from verified site information.

## When to use

- Use when the authorized intent is `seo_planning` and the request is to prepare a conservative llms.txt proposal from verified site information.

## When NOT to use

- Do not use when the request belongs to `robots_txt_generation`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `llms_txt_generation_draft`
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
