---
name: optimize-mode
description: |
  Use when running SEO optimization on an existing article or page. Covers a seven-step workflow: inspecting SEO settings, detecting search intent, generating metadata, generating FAQs, building schema JSON-LD, finding internal links, and validating all outputs. Trigger whenever the user wants to optimize, update SEO fields, generate schema, add FAQs, or build internal links for any piece of content. Do NOT use for new article drafting, topic discovery, or overwriting existing SEO values.
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
