---
name: optimize-mode
description: 'Use when running SEO optimization on an existing article or page. Covers a seven-step workflow: inspecting SEO settings, detecting search intent, generating metadata, generating FAQs, building schema JSON-LD, finding internal links, and validating all outputs. Trigger whenever the user wants to optimize, update SEO fields, generate schema, add FAQs, or build internal links for any piece of content. Do NOT use for new article drafting, topic discovery, or overwriting existing SEO values.'
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
    - optimize_mode
declared_action_classes:
    - read_brand_context
    - canvas_write
output_contracts:
    - optimize_mode_draft
    - quality_summary
token_budget:
    body_max_tokens: 3000
    references_max_tokens: 0
    assets_max_tokens: 0
    total_active_max_tokens: 3000
resource_manifest: resource-manifest.yaml
eval_manifest: eval-manifest.yaml
---

# Optimize Mode

## SEO OUTPUTS WORKFLOW (optimize mode)

### STEP 01: seo-inspect

- Read current seo_settings state
- Identify which fields are null or empty
- Identify which fields already have values
- Rule: never overwrite existing values

### STEP 02: intent-detection

- Read article title
- Read H2 and H3 headings
- Read body content for key entities and themes
- Classify search intent type:

| Intent | Description |
|---|---|
| `informational` | reader wants to learn |
| `commercial` | reader is comparing options |
| `navigational` | reader is looking for a specific page |
| `transactional` | reader wants to take an action |
| `news` | reader wants latest update |
| `opinion` | reader wants a point of view |

### STEP 03: metadata-generation

- **meta title:** clear, searchable, strong CTR, max 60 chars
- **meta description:** concise summary with click intent, max 155 chars
- **slug:** short, clean, lowercase, hyphenated, no stop words
- **og title:** social-friendly version of meta title
- **og description:** social-friendly summary, 1-2 sentences

### STEP 04: faq-generation

- Generate 3 to 5 real questions tied to article intent
- Questions must match what readers actually search
- Answers must be short, direct, and self-contained
- No fluff, no filler
- Format: Q + A pairs
- Each answer: 1-3 sentences max

### STEP 05: schema-generation

- Choose schema type based on page content:

| Schema Type | Use For |
|---|---|
| `Article` | general informational articles |
| `BlogPosting` | blog posts and opinion pieces |
| `NewsArticle` | match reports, breaking news, transfers |
| `FAQPage` | pages with question and answer content |
| `WebPage` | general landing or topic pages |
| `SportsEvent` | match previews or event pages |

- Build valid JSON-LD using: article title, meta description, page URL, publisher name and domain, date published / date modified, FAQ entities if FAQs are generated
- **Output format:** JSON-LD
- **Injection:** `<script type="application/ld+json">`

### STEP 06: internal-linking

- Find related pages on the same site via `getSitePages()`
- Match by topic, entity, or keyword overlap
- Suggest anchor text for each link
- Connect topically related pages only
- Output: page_title, url, anchor_text, relevance

### STEP 07: seo-validation

- Existing fields were NOT overwritten
- Meta title aligns with article title and intent
- Meta description matches article content
- Slug is clean, short, and keyword-aligned
- FAQs are answerable and non-duplicate
- Schema type matches content type
- Schema JSON-LD is valid and complete
- Internal links are relevant and safe
- No competitor domains in any links
- noIndex and noFollow are false unless explicitly set

## SUPPORTING TOOLS

| Tool | Purpose |
|---|---|
| `getSourcesInsights()` | competitor blocklist, authority domains |
| `getSitePages()` | find internal linking candidates |
| `scrapePage()` | read page content before linking |
| `readMemoryDoc()` | check brand/tone constraints |
| `quickContext()` | brand name, domain, positioning |

## Purpose

Provide procedural guidance to coordinate bounded improvements to existing selected content.

## When to use

- Use when the authorized intent is `optimize_content` and the request is to coordinate bounded improvements to existing selected content.

## When NOT to use

- Do not use when the request belongs to `post_write_checklist`.
- Do not use for direct publishing, policy bypass, or unapproved mutation.

## Inputs expected

- Sanitized project context
- Authorized intent and mode
- Evidence references or approved source summaries when required

## Procedure

Follow the skill-specific instructions above in order. Stop when required context, evidence, mode, or approval is absent.

## Output contract

- `optimize_mode_draft`
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
